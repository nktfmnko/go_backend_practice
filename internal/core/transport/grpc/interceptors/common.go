package core_grpc_interceptors

import (
	"context"
	"net"
	core_logger "practice/internal/core/logger"
	core_transport "practice/internal/core/transport/limiter"
	"practice/metrics"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const requestIDHeader = "x-request-id"

type contextKey string

const requestIDCtxKey contextKey = "request_id"

func SetRequestIDToContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDCtxKey, requestID)
}

func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDCtxKey).(string); ok {
		return id
	}
	return ""
}

func RequestIDUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		var requestID string
		if ids := md.Get(requestIDHeader); len(ids) > 0 && ids[0] != "" {
			requestID = ids[0]
		} else {
			requestID = uuid.NewString()
		}

		header := metadata.Pairs(requestIDHeader, requestID)
		_ = grpc.SetHeader(ctx, header)
		ctx = SetRequestIDToContext(ctx, requestID)

		return handler(ctx, req)
	}
}

func LoggerUnaryServerInterceptor(log *core_logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		requestID := RequestIDFromContext(ctx)

		l := log.With(
			zap.String("request_id", requestID),
			zap.String("grpc_method", info.FullMethod),
		)

		ctx = core_logger.ToContext(ctx, l)

		return handler(ctx, req)
	}
}

func TraceUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		logger := core_logger.FromContext(ctx)

		before := time.Now()
		logger.Debug(
			">>> incoming gRPC request",
			zap.Time("time", before.UTC()),
			zap.String("method", info.FullMethod),
		)
		res, err := handler(ctx, req)

		stCode := status.Code(err)

		logger.Debug(
			"<<< done gRPC request",
			zap.String("status_code", stCode.String()),
			zap.Duration("latency", time.Since(before.UTC())),
		)

		return res, err
	}
}

func PanicUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		logger := core_logger.FromContext(ctx)
		defer func() {
			if p := recover(); p != nil {
				logger.Error(
					"during handle gRPC request got unexpected panic",
					zap.Any("panic", p),
					zap.String("stack", string(debug.Stack())),
				)

				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

type validator interface {
	Validate() error
}

func ValidationUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if v, ok := req.(validator); ok {
			if err := v.Validate(); err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "validation err: %v", err)
			}
		}

		return handler(ctx, req)
	}
}

func RateLimitUnaryServerInterceptor(il *core_transport.IPLimiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		p, ok := peer.FromContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Internal, "failed to get peer info")
		}

		ip, _, err := net.SplitHostPort(p.Addr.String())
		if err != nil {
			ip = p.Addr.String()
		}

		if !il.Get(ip).Allow() {
			return nil, status.Errorf(codes.ResourceExhausted, "too many requests")
		}
		return handler(ctx, req)
	}
}

func MetricsUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()

		metrics.IncGRPCInFlight()
		defer metrics.DecGRPCInFlight()

		res, err := handler(ctx, req)

		duration := time.Since(start).Seconds()

		stCode := status.Code(err)

		metrics.IncGRPCRequests(info.FullMethod, stCode.String())
		metrics.ObserveGRPCDuration(info.FullMethod, duration)

		return res, err
	}
}
