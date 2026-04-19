package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/artem-smola/golang-course/task4/platform/grpcserver"
	"github.com/artem-smola/golang-course/task4/platform/logger"
	"github.com/artem-smola/golang-course/task4/processor/config"
	"github.com/artem-smola/golang-course/task4/processor/internal/adapter"
	"github.com/artem-smola/golang-course/task4/processor/internal/controller"
	"github.com/artem-smola/golang-course/task4/processor/internal/usecase"
	processorpb "github.com/artem-smola/golang-course/task4/proto/processor"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)
	log.Info("starting processor server...")
	log.Debug("debug messages are enabled")

	collectorClient, err := adapter.NewClient(cfg.Services.Collector, log)
	if err != nil {
		return fmt.Errorf("create collector grpc client: %w", err)
	}
	defer func() {
		if err := collectorClient.Close(); err != nil {
			log.Error("cannot close collector adapter", "error", err)
		}
	}()

	getRepoInfo := usecase.NewGetRepoInfoUsecase(collectorClient)
	getSubscriptionsRepoInfo := usecase.NewGetSubscriptionsRepoInfoUsecase(collectorClient)
	ping := usecase.NewPing()
	repoServer := controller.NewServer(log, getRepoInfo, getSubscriptionsRepoInfo, ping)

	srv, err := grpcserver.NewServer(cfg.GRPC.Address)
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
