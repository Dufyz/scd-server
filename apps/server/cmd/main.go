package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	_ "time/tzdata"

	db "github.com/Dufyz/scd-server/infra/database"
	"github.com/Dufyz/scd-server/internal/env"
	queue "github.com/Dufyz/scd-server/internal/queues"
	"github.com/Dufyz/scd-server/internal/rest/middlewares"
	"github.com/Dufyz/scd-server/internal/rest/routes"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	core := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(os.Stdout), zap.InfoLevel)
	logger = zap.New(core)
	zap.ReplaceGlobals(logger)

	envPath := filepath.Join("..", ".env")
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			logger.Fatal("Error loading .env file", zap.Error(err))
		}
		logger.Info(".env file loaded successfully")
	} else {
		logger.Warn(".env file not found, skipping loading")
	}

	requiredEnvVars := []string{
		"GO_ENV",
		"PORT",
		"POSTGRES_DB",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"DATABASE_URL",
		"REDIS_PORT",
		"REDIS_URL",
		"REDIS_TTL_SECONDS",
		"WEB_URL",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			logger.Fatal("Missing required environment variable", zap.String("envVar", envVar))
		}
	}
}

func main() {
	fmt.Printf("Starting scd-server \n")

	time.Local, _ = time.LoadLocation("America/Sao_Paulo")

	connection, err := db.NewDBConnectionWithRetries(20)
	if err != nil {
		logger.Fatal("Could not connect to database", zap.Error(err))
	}

	defer connection.Close()

	redisQueueAddr := env.GetString("REDIS_URL", "server-redis:6379")

	go func() {
		err = queue.NewAsynqServerQueue(redisQueueAddr, connection)
		if err != nil {
			logger.Fatal("Could not start Redis Asynq Server", zap.Error(err))
		}
	}()

	queueClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisQueueAddr})
	defer queueClient.Close()

	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 10 * time.Minute,
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*.vercel.app", "http://localhost:*"},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		ExposeHeaders:    []string{echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	e.Use(middlewares.LoggerMiddleware(logger))

	routes.SetupRoutes(e, connection, queueClient)

	port := env.GetInt("PORT", 3000)

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", port)); err != nil {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

}
