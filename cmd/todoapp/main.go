package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	core_logger "practice/internal/core/logger"
	core_redis2 "practice/internal/core/repository/cache/redis"
	"practice/internal/core/repository/postgres/pool/pgx"
	core_grpc_interceptors "practice/internal/core/transport/grpc/interceptors"
	core_grpc_server "practice/internal/core/transport/grpc/server"
	core_http_middleware "practice/internal/core/transport/http/middleware"
	core_http_server "practice/internal/core/transport/http/server"
	statistics_postgres_repository "practice/internal/features/statistics/repository/postgres"
	statistics_service "practice/internal/features/statistics/service"
	statistics_cache_service "practice/internal/features/statistics/service/cache"
	statistics_transport_http "practice/internal/features/statistics/transport/http"
	tasks_postgres_repository "practice/internal/features/tasks/repository/postgres"
	tasks_redis_repository "practice/internal/features/tasks/repository/redis"
	tasks_service "practice/internal/features/tasks/service"
	tasks_transport_grpc "practice/internal/features/tasks/transport/grpc"
	tasks_transport_http "practice/internal/features/tasks/transport/http"
	users_postgres_repository "practice/internal/features/users/repository/postgres"
	users_redis_repository "practice/internal/features/users/repository/redis"
	users_service "practice/internal/features/users/service"
	users_transport_grpc "practice/internal/features/users/transport/grpc"
	users_transport_http "practice/internal/features/users/transport/http"
	v1 "practice/pkg/pb/v1"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var timeZone = time.UTC

const (
	usersCacheTTL      = 5 * time.Minute
	tasksCacheTTL      = 1 * time.Minute
	statisticsCacheTTL = 30 * time.Second
)

func main() {
	_ = godotenv.Load()

	time.Local = timeZone

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
	logger.Debug("application time zone", zap.Any("zone", timeZone))

	logger.Debug("init connection pool")
	pool, err := core_pgx_pool.NewPool(core_pgx_pool.NewConfigMust(), ctx)
	if err != nil {
		logger.Fatal("failed to init connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("init redis")
	redisCache, err := core_redis2.NewClient(context.Background(), core_redis2.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init redis ", zap.Error(err))
	}
	defer func(redisCache *core_redis2.RedisClient) {
		err := redisCache.Close()
		if err != nil {

		}
	}(redisCache)

	logger.Debug("init feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersCachedRepository := users_redis_repository.NewCachedUsersRepository(usersRepository, redisCache, usersCacheTTL)
	usersService := users_service.NewUsersService(usersCachedRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("init feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksCachedRepository := tasks_redis_repository.NewCachedTasksRepository(tasksRepository, redisCache, tasksCacheTTL)
	tasksService := tasks_service.NewTasksService(tasksCachedRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("init feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsCacheService := statistics_cache_service.NewStatisticsCacheService(statisticsService, redisCache, statisticsCacheTTL)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsCacheService)

	logger.Debug("init HTTP Server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.APIVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	httpServer.RegisterAPIRoutes(apiVersionRouter)

	usersGRPC := users_transport_grpc.NewUsersGRPC(usersService)
	tasksGRPC := tasks_transport_grpc.NewTasksGRPC(tasksService)
	grpcServer := core_grpc_server.NewGRPCServer(
		core_grpc_server.NewConfigMust(),
		grpc.ChainUnaryInterceptor(
			core_grpc_interceptors.RequestIDUnaryServerInterceptor(),
			core_grpc_interceptors.LoggerUnaryServerInterceptor(logger),
			core_grpc_interceptors.TraceUnaryServerInterceptor(),
			core_grpc_interceptors.PanicUnaryServerInterceptor(),
		))
	grpcServer.Register(func(server *grpc.Server) {
		v1.RegisterUsersServiceV1Server(server, usersGRPC)
	}, func(server *grpc.Server) {
		v1.RegisterTasksServiceV1Server(server, tasksGRPC)
	})

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		logger.Info("HTTP server starting...")
		if err := httpServer.Run(gCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server error", zap.Error(err))
			return err
		}
		return nil
	})

	g.Go(func() error {
		logger.Info("gRPC server starting...")
		if err := grpcServer.Start(gCtx); err != nil {
			logger.Error("gRPC server error", zap.Error(err))
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		logger.Error("application stopped with error", zap.Error(err))
	} else {
		logger.Info("application stopped gracefully")
	}
}
