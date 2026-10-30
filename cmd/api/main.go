package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/kisalto/Feel-The-Night-Swagger/docs"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/app"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/database"
)

// @title           Feel The Night API
// @version         1.0
// @description     API for managing guides, events and characters.
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	if err := godotenv.Load(); err != nil {
		zap.L().Warn(".env file not found, reading environment variables from the system.")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		zap.L().Fatal("JWT_SECRET is not set")
	}

	if len(secret) < 32 {
		zap.L().Fatal("JWT_SECRET must be at least 32 characters long")
	}

	// 1. Port
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("APP_PORT")
	}
	if port == "" {
		port = "8080"
	}

	// 2. Dynamic Swagger configuration via docs.SwaggerInfo
	renderHost := os.Getenv("RENDER_EXTERNAL_HOSTNAME")

	if renderHost != "" {
		// Production (Render)
		docs.SwaggerInfo.Host = renderHost
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		// Local development
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", port)
		docs.SwaggerInfo.Schemes = []string{"http"}
	}

	// 3. Database connection
	zap.L().Info("Connecting to the database...")
	db, err := database.Connect()
	if err != nil {
		zap.L().Fatal("Failed to initialize the database", zap.Error(err))
	}

	// 4. Routes
	routes := app.NewRouter(db)

	// 5. HTTP server
	serverPort := fmt.Sprintf(":%s", port)

	srv := &http.Server{
		Addr:              serverPort,
		Handler:           routes,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		zap.L().Info("Server listening", zap.String("port", serverPort))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("Failed to start the server", zap.Error(err))
		}
	}()

	<-ctx.Done()
	zap.L().Info("Shutdown signal received, stopping the server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		zap.L().Error("Forced shutdown: in-flight requests did not finish in time", zap.Error(err))
	}

	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}

	zap.L().Info("Server stopped")

}
