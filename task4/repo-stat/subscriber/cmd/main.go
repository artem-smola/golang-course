package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/artem-smola/golang-course/task4/platform/grpcserver"
	"github.com/artem-smola/golang-course/task4/platform/logger"
	subscriberpb "github.com/artem-smola/golang-course/task4/proto/subscriber"
	"github.com/artem-smola/golang-course/task4/subscriber/config"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/artem-smola/golang-course/task4/subscriber/internal/adapter"
	"github.com/artem-smola/golang-course/task4/subscriber/internal/controller"
	migrator "github.com/artem-smola/golang-course/task4/subscriber/internal/storage/migrate"
	"github.com/artem-smola/golang-course/task4/subscriber/internal/storage/postgres"
	db "github.com/artem-smola/golang-course/task4/subscriber/internal/storage/postgres/sqlc"
	"github.com/artem-smola/golang-course/task4/subscriber/internal/usecase"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting subscriber server...")
	log.Debug("debug messages are enabled")

	if err := migrator.Up(cfg.Migrations.Path, cfg.Postgres.DSN); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.Postgres.DSN)
	if err != nil {
		return fmt.Errorf("parse postgres config: %w", err)
	}
	poolConfig.MinConns = cfg.Postgres.MinConns
	poolConfig.MaxConns = cfg.Postgres.MaxConns

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("create postgres pool: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping postgres: %w", err)
	}

	queries := db.New(pool)
	repository := postgres.NewRepository(queries)
	githubClient := adapter.NewClient(cfg.GitHub.BaseURL, cfg.GitHub.Token, cfg.GitHub.Timeout)

	pingUseCase := usecase.NewPing()
	subscriptionUseCase := usecase.NewSubscription(repository, githubClient)
	pingServer := controller.NewServer(log, pingUseCase, subscriptionUseCase)

	srv, err := grpcserver.NewServer(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	subscriberpb.RegisterSubscriberServer(srv.GRPC(), pingServer)

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	if err := run(ctx); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Printf("launching server error: %s\n", err)
		}
		cancel()
		os.Exit(1)
	}
}
