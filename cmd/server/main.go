package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourorg/shadowchat/backend/internal/config"
	"github.com/yourorg/shadowchat/backend/internal/db"
	"github.com/yourorg/shadowchat/backend/internal/repository"
	"github.com/yourorg/shadowchat/backend/internal/service"
	"github.com/yourorg/shadowchat/backend/internal/http/routes"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration
	cfg := config.Load()
	logger.Info("loaded config", zap.Any("config", cfg))

	// Connect to PostgreSQL
	pg, err := db.NewPostgres(cfg.PostgresURL)
	if err != nil {
		logger.Fatal("failed to connect to postgres", zap.Error(err))
	}
	defer pg.Close()
	logger.Info("connected to postgres")

	// Connect to Redis
	redis, err := db.NewRedis(cfg.RedisURL)
	if err != nil {
		logger.Warn("failed to connect to redis, continuing without redis", zap.Error(err))
	} else {
		defer redis.Close()
		logger.Info("connected to redis")
	}

	// Initialize repositories
	repos := repository.NewBundle(pg, redis)

	// Initialize services
	svcs := service.NewBundle(cfg, repos, logger)

	// Create HTTP router
	router := routes.NewRouter(cfg, svcs, logger)

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		logger.Info("starting server", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server exited")
}
