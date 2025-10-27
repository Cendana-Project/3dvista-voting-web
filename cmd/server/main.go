package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appPkg "voteweb/internal/app"
	httpPkg "voteweb/internal/http"
	"voteweb/seed"
)

func main() {
	ctx := context.Background()

	// Initialize application
	app, err := appPkg.New(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer app.Close()

	app.Logger.Info("Application initialized successfully")

	// Run seeds if SEED environment variable is set
	if os.Getenv("SEED") == "true" {
		app.Logger.Info("Running seed data...")
		if err := seed.SeedInnovations(ctx, app.Pool); err != nil {
			log.Fatalf("Failed to seed innovations: %v", err)
		}
		app.Logger.Info("Seed data completed successfully")
	}

	// Setup router
	router := httpPkg.SetupRouter(app.Config, app.Pool, app.Service, app.Logger)

	// Create HTTP server
	addr := fmt.Sprintf(":%s", app.Config.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		app.Logger.Info("Starting HTTP server", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	app.Logger.Info("Server started successfully", "port", app.Config.Port)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.Logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	app.Logger.Info("Server exited gracefully")
}


