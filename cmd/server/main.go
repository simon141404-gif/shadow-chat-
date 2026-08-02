package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/yourorg/shadowchat/backend/internal/config"
	"github.com/yourorg/shadowchat/backend/internal/db"
	"github.com/yourorg/shadowchat/backend/internal/repository"
	"github.com/yourorg/shadowchat/backend/internal/service"
	"github.com/yourorg/shadowchat/backend/internal/http/routes"
	"go.uber.org/zap"
)

type PostgresPool = db.PostgresPool

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration
	cfg := config.Load()
	logger.Info("loaded config", zap.Any("config", cfg))

	// Connect to PostgreSQL (optional - can work without)
	var pg *pgxpool.Pool
	pg, err := db.NewPostgres(cfg.PostgresURL)
	if err != nil {
		logger.Warn("failed to connect to postgres, continuing without DB", zap.Error(err))
	} else if pg != nil {
		defer pg.Close()
		logger.Info("connected to postgres")
		
		// Run migrations
		if err := runMigrations(pg); err != nil {
			logger.Warn("failed to run migrations", zap.Error(err))
		} else {
			logger.Info("migrations completed")
		}
	} else {
		logger.Info("running without PostgreSQL (in-memory mode)")
	}

	// Connect to Redis (optional - can work without)
	var redis *redis.Client
	redis, err = db.NewRedis(cfg.RedisURL)
	if err != nil {
		logger.Warn("failed to connect to redis, continuing without redis", zap.Error(err))
	} else if redis != nil {
		defer redis.Close()
		logger.Info("connected to redis")
	} else {
		logger.Info("running without Redis (in-memory mode)")
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

func runMigrations(pg *db.PostgresPool) error {
	if pg == nil {
		return nil // No DB, skip migrations
	}
	ctx := context.Background()
	migrations := []string{
		`CREATE EXTENSION IF NOT EXISTS pgcrypto;`,
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			public_id TEXT UNIQUE NOT NULL,
			display_name TEXT,
			avatar_url TEXT,
			bio TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			deleted_at TIMESTAMPTZ
		);`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			jti TEXT UNIQUE NOT NULL,
			refresh_token_hash TEXT NOT NULL,
			expires_at TIMESTAMPTZ NOT NULL,
			revoked_at TIMESTAMPTZ
		);`,
		`CREATE TABLE IF NOT EXISTS chats (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			type TEXT NOT NULL,
			name TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);`,
		`CREATE TABLE IF NOT EXISTS chat_members (
			chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			PRIMARY KEY (chat_id, user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS messages (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
			sender_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			client_msg_id TEXT,
			ciphertext BYTEA,
			message_type TEXT NOT NULL,
			reply_to_message_id UUID,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			edited_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ,
			expires_at TIMESTAMPTZ
		);`,
		`CREATE TABLE IF NOT EXISTS attachments (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
			file_name TEXT NOT NULL,
			file_size BIGINT,
			content_type TEXT,
			storage_path TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);`,
		`CREATE TABLE IF NOT EXISTS contacts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			contact_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			UNIQUE(user_id, contact_user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS notifications (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			body TEXT,
			data JSONB,
			read_at TIMESTAMPTZ,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);`,
	}

	for _, m := range migrations {
		if _, err := pg.Exec(ctx, m); err != nil {
			return err
		}
	}
	return nil
}
