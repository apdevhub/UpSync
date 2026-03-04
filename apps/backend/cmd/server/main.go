package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/upsync/backend/internal/config"
	"github.com/upsync/backend/internal/database"
	"github.com/upsync/backend/internal/handlers"
	"github.com/upsync/backend/internal/middleware"
	"github.com/upsync/backend/internal/scheduler"
	"github.com/upsync/backend/internal/services"
)

func main() {
	// ── Structured logging ────────────────────────────────────
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	// ── Load .env (ignore error in production) ────────────────
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found — reading from environment")
	}

	// ── Config ────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config error", slog.Any("error", err))
		os.Exit(1)
	}

	// ── Supabase DB client ────────────────────────────────────
	db := database.New(cfg)

	// ── Services ──────────────────────────────────────────────
	fileSvc := services.New(db, cfg)

	// ── Handlers ──────────────────────────────────────────────
	fileHandler := handlers.NewFileHandler(fileSvc, cfg)

	// ── Gin router ────────────────────────────────────────────
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Middleware stack
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(cfg.FrontendURL))
	r.Use(middleware.MaxBodySize(cfg.MaxFileSizeBytes + 1024)) // +1 KB for form fields

	// Routes
	r.GET("/health", handlers.Health)

	api := r.Group("/api")
	{
		files := api.Group("/files")
		{
			files.POST("/upload", fileHandler.Upload)
			files.GET("/:id", fileHandler.GetMeta)
			files.GET("/:id/download", fileHandler.GetDownloadURL)
		}
	}

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found", "code": "NOT_FOUND"})
	})

	// ── Scheduler ─────────────────────────────────────────────
	sched := scheduler.New(fileSvc)
	sched.Start()
	defer sched.Stop()

	// ── HTTP server with graceful shutdown ────────────────────
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		slog.Info(fmt.Sprintf("🚀 UpSync API listening on http://localhost:%s", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("graceful shutdown initiated...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", slog.Any("error", err))
	}
	slog.Info("server stopped")
}
