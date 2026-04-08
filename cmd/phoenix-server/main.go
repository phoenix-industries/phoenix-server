package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/phoenix-industries/phoenix-server/internal/apiservice/v1"
	"github.com/phoenix-industries/phoenix-server/internal/authservice"
	"github.com/phoenix-industries/phoenix-server/internal/buildinfo"
	"github.com/phoenix-industries/phoenix-server/pkg/auth"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
	"github.com/phoenix-industries/phoenix-server/pkg/httputil"
	"github.com/phoenix-industries/phoenix-server/pkg/kernel"
)

const defaultPort = ":5000"

var port = flag.String("port", defaultPort, "")

func main() {
	flag.Parse()

	// TODO: load env files explicitly
	if err := godotenv.Load(); err != nil {
		slog.Error("failed to load environment variables", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done() // wait for first signal
		forceCtx, forceStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer forceStop()
		<-forceCtx.Done()
		slog.Warn("second shutdown signal received - forcing immediate exit")
		os.Exit(130)
	}()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("application panic", "panic", r)
			os.Exit(1)
		}
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := run(ctx); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()

	select {
	case <-done:
		return
	case <-ctx.Done():
		slog.Info("shutdown signal received, waiting for server to exit gracefully")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	select {
	case <-done:
		slog.Info("successfully shutdown server gracefully")
	case <-shutdownCtx.Done():
		slog.Info("shutdown timeout reached, forcing server to exit")
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	if *port == "" {
		*port = os.Getenv("PORT")
		if *port == "" {
			*port = defaultPort
		}
	} else if (*port)[0] != ':' {
		*port = ":" + *port
	}

	logger := slog.Default()

	logger.Info(
		"build info",
		"build_tag", buildinfo.BuildTag,
		"go_version", buildinfo.GoVersion,
		"system_tag", buildinfo.SystemTag,
	)
	if buildinfo.DevMode() {
		logger.Info("dev mode enabled")
	}

	dbConfig, err := database.ConfigFromEnv()
	if err != nil {
		return err
	}
	if buildinfo.DevMode() {
		logger.Info("database: connecting", "connection_string", dbConfig.ConnectionString())
	}

	db, err := database.ConnectWithLogger(ctx, dbConfig, logger)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Migrate(ctx); err != nil {
		return err
	}

	auth, err := auth.NewFromEnv()
	if err != nil {
		return err
	}

	kernel := kernel.NewKernel(db, auth, logger)
	if err := kernel.Run(authservice.New(), apiservice.New()); err != nil {
		return err
	}

	mux := kernel.Mux()
	mux.HandleFunc("/", httputil.NotFoundHandler(logger))

	handler := httputil.ChainMiddlewares(
		httputil.RecoveryMiddleware(logger),
		httputil.LoggingMiddleware(logger),
	)(mux)

	srv := http.Server{
		Addr:         *port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("server started", "port", *port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutting down http server")
		shutdownError := srv.Shutdown(ctx)
		if shutdownError != nil {
			logger.Error("server shutdown failed", "error", err)
		}
		select {
		case err := <-serverErr:
			return err
		default:
			return shutdownError
		}
	case err := <-serverErr:
		return err
	}
}
