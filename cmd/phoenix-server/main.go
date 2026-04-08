package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		done()
		if r := recover(); r != nil {
			slog.Error("application panic", "panic", r)
			os.Exit(1)
		}
	}()

	err := run(ctx)
	done()

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("successfully exited")
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

	logger.Info("build info", "build_tag", buildinfo.BuildTag, "go_version", buildinfo.GoVersion, "system_tag", buildinfo.SystemTag)
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

	logger.Info("server started", "port", *port)

	mux := kernel.Mux()
	mux.HandleFunc("/", httputil.NotFoundHandler(logger))

	handler := httputil.ChainMiddlewares(
		httputil.RecoveryMiddleware(logger),
		httputil.LoggingMiddleware(logger),
	)(mux)

	return http.ListenAndServe(*port, handler)
}
