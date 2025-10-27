package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"voteweb/internal/config"
	"voteweb/internal/domain"
	"voteweb/internal/repo"
	"voteweb/internal/util"
)

// App represents the application
type App struct {
	Config  *config.Config
	Pool    *pgxpool.Pool
	Service domain.VoteService
	Logger  *slog.Logger
}

// New creates and initializes a new App
func New(ctx context.Context) (*App, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	logger := initLogger()

	// Connect to database
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Ping database to verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Connected to database")

	// Initialize repository
	repository := repo.NewPostgresRepository(pool)

	// Initialize IP hasher
	ipHasher := util.NewIPHasher(cfg.IPHashSalt)

	// Initialize service
	service := domain.NewVoteService(repository, ipHasher, logger)

	return &App{
		Config:  cfg,
		Pool:    pool,
		Service: service,
		Logger:  logger,
	}, nil
}

// Close closes the application resources
func (a *App) Close() {
	if a.Pool != nil {
		a.Pool.Close()
		a.Logger.Info("Database connection closed")
	}
}

func initLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	// Use JSON handler for structured logging
	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}


