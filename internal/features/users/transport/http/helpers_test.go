package users_transport_http

import (
	"context"
	core_logger "practice/internal/core/logger"

	"go.uber.org/zap"
)

func testContext() context.Context {
	nopZapLogger := zap.NewNop()
	testLogger := &core_logger.Logger{
		Logger: nopZapLogger,
		File:   nil,
	}
	return core_logger.ToContext(context.Background(), testLogger)
}
