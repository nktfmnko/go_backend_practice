package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	core_logger "practice/internal/core/logger"
	core_redis2 "practice/internal/core/repository/cache/redis"
	"practice/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "practice/internal/core/transport/http/middleware"
	core_http_server "practice/internal/core/transport/http/server"
	statistics_postgres_repository "practice/internal/features/statistics/repository/postgres"
	statistics_service "practice/internal/features/statistics/service"
	statistics_cache_service "practice/internal/features/statistics/service/cache"
	statistics_transport_http "practice/internal/features/statistics/transport/http"
	tasks_postgres_repository "practice/internal/features/tasks/repository/postgres"
	tasks_redis_repository "practice/internal/features/tasks/repository/redis"
	tasks_service "practice/internal/features/tasks/service"
	tasks_transport_http "practice/internal/features/tasks/transport/http"
	users_postgres_repository "practice/internal/features/users/repository/postgres"
	users_redis_repository "practice/internal/features/users/repository/redis"
	users_service "practice/internal/features/users/service"
	users_transport_http "practice/internal/features/users/transport/http"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
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

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error: ", zap.Error(err))
	}
}
