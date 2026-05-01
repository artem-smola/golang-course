package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	"repo-stat/processor/config"
	"repo-stat/processor/internal/adapter"
	"repo-stat/processor/internal/controller"
	"repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/processor"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting processor server...")
	log.Debug("debug messages are enabled")

	collectorClient, err := adapter.NewGRPCClient(cfg.Services.Collector, log)
	if err != nil {
		return fmt.Errorf("create collector grpc client: %w", err)
	}

	getRepoInfo := usecase.NewGetRepoInfo(collectorClient)
	ping := usecase.NewPing()
	repoServer := controller.NewGRPCServer(log, getRepoInfo, ping)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	processorpb.RegisterProcessorServer(srv.GRPC(), repoServer)

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)

	if err := run(ctx); err != nil {
		if _, errPrint := fmt.Fprintln(os.Stderr, err); errPrint != nil {
			fmt.Printf("launching server error: %s\n", errPrint)
		}
		cancel()
		os.Exit(1)
	}

	cancel()
}
