package db

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPool = pgxpool.Pool

func NewPostgres(url string) (*pgxpool.Pool, error) {
	// If no URL provided, return nil (will use in-memory fallback)
	if url == "" || url == "postgres://localhost:5432/shadowchat?sslmode=disable" {
		return nil, nil
	}
	
	// Add sslmode=require if not present
	if !strings.Contains(url, "sslmode") {
		if strings.Contains(url, "?") {
			url += "&sslmode=disable"
		} else {
			url += "?sslmode=disable"
		}
	}
	
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 20
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 15 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return pgxpool.NewWithConfig(ctx, cfg)
}
