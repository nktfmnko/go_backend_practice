package core_http_middleware

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	core_error "practice/internal/core/errors"
	core_logger "practice/internal/core/logger"
	core_auth "practice/internal/core/transport/auth"
	core_http_response "practice/internal/core/transport/http/response"
	core_transport "practice/internal/core/transport/limiter"
	"practice/metrics"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			requestID := request.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			request.Header.Set(requestIDHeader, requestID)
			writer.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(writer, request)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			requestID := request.Header.Get(requestIDHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", request.URL.String()),
			)

			ctx := core_logger.ToContext(request.Context(), l)

			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_response.NewResponseWriter(writer)

			before := time.Now()
			log.Debug(
				">>> incoming HTTP request",
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, request)

			log.Debug(
				"<<< done HTTP request",
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("latency", time.Now().Sub(before)),
			)
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHttpResponseHandler(log, writer)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(
						p,
						"during handle HTTP request got unexpected panic",
					)
				}
			}()

			next.ServeHTTP(writer, request)
		})
	}
}

func Auth(secret string, appID int) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHttpResponseHandler(log, writer)

			header := request.Header.Get("Authorization")
			tokenStr := strings.TrimPrefix(header, "Bearer ")
			if tokenStr == header || tokenStr == "" {
				responseHandler.ErrorResponse(core_error.ErrMissingBearerToken, "missing bearer token")
			}

			claims, err := core_auth.ParseToken(tokenStr, secret, appID)
			if err != nil {
				handleAuthError(responseHandler, err)
				return
			}
			ctx = core_auth.WithUser(ctx, claims)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

func handleAuthError(h *core_http_response.HTTPResponseHandler, err error) {
	switch {
	case errors.Is(err, core_auth.ErrWrongApp):
		h.ErrorResponse(fmt.Errorf("%w: %w", core_error.ErrInvalidToken, err), "token issued for another app")
	case errors.Is(err, core_auth.ErrInvalidToken):
		h.ErrorResponse(fmt.Errorf("%w: %w", core_error.ErrInvalidToken, err), "invalid token")
	default:
		h.ErrorResponse(err, "failed to parse token")
	}
}

func RateLimit(il *core_transport.IPLimiter) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHttpResponseHandler(log, writer)

			ip, _, err := net.SplitHostPort(request.RemoteAddr)
			if err != nil {
				ip = request.RemoteAddr
			}

			if !il.Get(ip).Allow() {
				responseHandler.ErrorResponse(core_error.ErrTooManyRequests, "too many requests")
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}

func Metrics() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			rw := core_http_response.NewResponseWriter(writer)
			start := time.Now()

			metrics.IncInFlight()
			defer metrics.DecInFlight()

			next.ServeHTTP(rw, request)

			duration := time.Since(start).Seconds()
			status := strconv.Itoa(rw.GetStatusCode())

			metrics.IncRequests(request.Method, request.URL.Path, status)
			metrics.ObserveDuration(request.Method, request.URL.Path, duration)
		})
	}
}
