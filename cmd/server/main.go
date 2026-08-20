package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/ultimate-ci-cd-pipeline/internal/config"
	"github.com/example/ultimate-ci-cd-pipeline/internal/handler"
	"github.com/example/ultimate-ci-cd-pipeline/internal/logger"
	"github.com/example/ultimate-ci-cd-pipeline/internal/middleware"
	"github.com/example/ultimate-ci-cd-pipeline/internal/service"
)

// Version is injected during compilation via -ldflags.
var Version = "dev"

func main() {
	cfg, err := config.LoadConfig(Version)
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	log := logger.SetupLogger(cfg.LogLevel)
	log.Info("Starting Enterprise Microservice",
		slog.String("app_name", cfg.AppName),
		slog.String("environment", cfg.Environment),
		slog.String("version", cfg.Version),
		slog.Int("port", cfg.Port),
	)

	// Initialize services and handlers
	orderService := service.NewInMemoryOrderService()
	orderHandler := handler.NewOrderHandler(orderService)
	healthHandler := handler.NewHealthHandler(cfg.Version)
	uiHandler := handler.NewUIHandler(cfg.Version, cfg.Environment, cfg.AppName)

	// Setup Router
	mux := http.NewServeMux()
	mux.Handle("/", uiHandler)
	mux.HandleFunc("/healthz", healthHandler.Healthz)
	mux.HandleFunc("/livez", healthHandler.Livez)
	mux.HandleFunc("/readyz", healthHandler.Readyz)
	mux.Handle("/api/v1/orders", orderHandler)
	mux.Handle("/api/v1/orders/", orderHandler)

	// Apply Middlewares
	handlerStack := middleware.PanicRecovery(log)(
		middleware.RequestLogger(log)(mux),
	)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      handlerStack,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Server shutdown listener context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		s := <-sig
		log.Info("Received OS shutdown signal", slog.String("signal", s.String()))

		shutdownCtx, cancel := context.WithTimeout(serverCtx, cfg.ShutdownTimeout)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Error("Graceful shutdown timed out, forcing exit")
				os.Exit(1)
			}
		}()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error("Failed to gracefully shutdown server", slog.Any("error", err))
		}
		serverStopCtx()
	}()

	log.Info("HTTP Server listening", slog.String("addr", server.Addr))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("HTTP Server fatal error", slog.Any("error", err))
		os.Exit(1)
	}

	<-serverCtx.Done()
	log.Info("Server stopped cleanly")
}
