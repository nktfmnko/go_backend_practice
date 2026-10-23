package core_http_middleware

import (
	"net"
	"net/http"
	core_error "practice/internal/core/errors"
	core_logger "practice/internal/core/logger"
	core_http_response "practice/internal/core/transport/http/response"
	core_transport "practice/internal/core/transport/limiter"
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
