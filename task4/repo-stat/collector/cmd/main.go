package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/artem-smola/golang-course/task4/collector/config"
	"github.com/artem-smola/golang-course/task4/collector/internal/adapter"
	"github.com/artem-smola/golang-course/task4/collector/internal/controller"
	"github.com/artem-smola/golang-course/task4/collector/internal/usecase"
	"github.com/artem-smola/golang-course/task4/platform/grpcserver"
	"github.com/artem-smola/golang-course/task4/platform/logger"
	collectorpb "github.com/artem-smola/golang-course/task4/proto/collector"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting collector server...")
	log.Debug("debug messages are enabled")

	githubClient := &adapter.GitHubClient{}
	subscriberClient, err := adapter.NewSubscriberClient(cfg.Services.Subscriber, log)
	if err != nil {
		return fmt.Errorf("create subscriber grpc client: %w", err)
	}
	defer func() {
		if err := subscriberClient.Close(); err != nil {
			log.Error("cannot close subscriber adapter", "error", err)
		}
	}()

	collectorUsecase := usecase.NewUsecase(githubClient, subscriberClient)
	server := controller.NewServer(log, collectorUsecase)

	grpcServer, err := grpcserver.NewServer(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}
	collectorpb.RegisterCollectorServer(grpcServer.GRPC(), server)

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
