package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/phoenix-industries/phoenix-server/internal/authservice"
	"github.com/phoenix-industries/phoenix-server/internal/buildinfo"
	"github.com/phoenix-industries/phoenix-server/pkg/auth"
	"github.com/phoenix-industries/phoenix-server/pkg/database"
	"github.com/phoenix-industries/phoenix-server/pkg/kernel"
)

const defaultPort = ":5000"

var port = flag.String("port", defaultPort, "")

func main() {
	flag.Parse()
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	if *port == "" {
		*port = os.Getenv("PORT")
		if *port == "" {
			*port = defaultPort
		}
	} else if (*port)[0] != ':' {
		*port = ":" + *port
	}

	slog.Info("build info", "build_tag", buildinfo.BuildTag, "go_version", buildinfo.GoVersion, "system_tag", buildinfo.SystemTag)
	if buildinfo.DevMode() {
		slog.Info("dev mode enabled")
	}

	ctx := context.Background()
	logger := slog.Default()

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

	jwtSecret, err := auth.JWTSecretFromEnv()
	if err != nil {
		return err
	}

	auth := auth.New(jwtSecret)

	kernel := kernel.NewKernel(auth, logger)
	kernel.Run(
		authservice.New(),
	)

	if err := http.ListenAndServe(*port, kernel.Mux()); err != nil {
		return err
	}

	return nil
}
