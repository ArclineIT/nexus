package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"git.arcline.it/ArclineIT/nexus/internal/config"
	"git.arcline.it/ArclineIT/nexus/internal/handler"
	"git.arcline.it/ArclineIT/nexus/internal/middleware"
)

// Run starts the Nexus HTTP server and blocks until shutdown.
func Run() error {
	cfg := config.Load()

	logger := setupLogger(cfg)

	authHandler := handler.NewAuthHandler(cfg)

	uiHandler, err := handler.NewUIHandler(cfg)
	if err != nil {
		return fmt.Errorf("initializing UI handler: %w", err)
	}

	mux := http.NewServeMux()

	// Public API routes
	mux.HandleFunc("GET /health", handler.Ready())
	mux.HandleFunc("GET /ready", handler.Ready())
	mux.HandleFunc("POST /auth/login", authHandler.Login())
	mux.HandleFunc("POST /auth/refresh", authHandler.Refresh())

	// Protected API routes (Bearer token)
	mux.Handle("GET /auth/me", middleware.Authenticate(cfg)(authHandler.Me()))

	// Public web UI routes
	mux.HandleFunc("GET /{$}", uiHandler.Root())
	mux.HandleFunc("GET /login", uiHandler.LoginPage())
	mux.HandleFunc("POST /login", uiHandler.LoginSubmit())
	mux.HandleFunc("GET /signup", uiHandler.SignupPage())
	mux.HandleFunc("POST /signup", uiHandler.SignupSubmit())
	mux.HandleFunc("GET /forgot-password", uiHandler.ForgotPasswordPage())
	mux.HandleFunc("POST /forgot-password", uiHandler.ForgotPasswordSubmit())
	mux.HandleFunc("GET /reset-password", uiHandler.ResetPasswordPage())
	mux.HandleFunc("POST /reset-password", uiHandler.ResetPasswordSubmit())

	// Protected web UI routes (cookie or Bearer)
	mux.Handle("GET /dashboard", middleware.WebAuth(cfg)(uiHandler.DashboardPage()))
	mux.Handle("POST /logout", middleware.WebAuth(cfg)(uiHandler.Logout()))

	// Apply global middleware — outermost first
	var h http.Handler = mux
	h = middleware.CORS(h)
	h = middleware.RequestID(h)
	h = middleware.Logger(logger)(h)
	h = middleware.Recoverer(logger)(h)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Graceful shutdown
	errCh := make(chan error, 1)
	go func() {
		logger.Info("nexus control panel starting", slog.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		logger.Info("shutting down", slog.String("signal", sig.String()))
	case err := <-errCh:
		logger.Error("server error", slog.Any("error", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutting down server: %w", err)
	}

	logger.Info("server stopped gracefully")
	return nil
}

func setupLogger(cfg *config.Config) *slog.Logger {
	var level slog.Level
	switch cfg.Logging.Level {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}

	var h slog.Handler
	if cfg.Logging.Format == "json" {
		h = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		h = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(h)
}
