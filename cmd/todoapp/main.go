package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	core_logger "practice/internal/core/logger"
	core_postgres_pool "practice/internal/core/repository/postgres/pool"
	core_http_middleware "practice/internal/core/transport/http/middleware"
	core_http_server "practice/internal/core/transport/http/server"
	users_postgres_repository "practice/internal/features/users/repository/postgres"
	users_service "practice/internal/features/users/service"
	users_transport_http "practice/internal/features/users/transport/http"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init app logger: ", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("init connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(core_postgres_pool.NewConfigMust(), ctx)
	if err != nil {
		logger.Fatal("failed to init connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("init feature", zap.String("feature", "users"))

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("init HTTP Server", zap.String("feature", "users"))
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)
	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.APIVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)

	httpServer.RegisterAPIRoutes(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error: ", zap.Error(err))
	}
}
