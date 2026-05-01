package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"repo-stat/collector/config"
	"repo-stat/collector/internal/adapter"
	"repo-stat/collector/internal/controller"
	"repo-stat/collector/internal/usecase"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	gen "repo-stat/proto/collector"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting collector server...")
	log.Debug("debug messages are enabled")

	client := &adapter.HTTPClient{}
	usecase := usecase.NewUsecase(client)
	server := controller.NewGRPCServer(log, usecase)

	grpcServer, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}
	gen.RegisterCollectorServer(grpcServer.GRPC(), server)

	if err := grpcServer.Run(ctx); err != nil {
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
