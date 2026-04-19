package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/artem-smola/golang-course/task4/api/config"
	"github.com/artem-smola/golang-course/task4/api/internal/adapter/processor"
	"github.com/artem-smola/golang-course/task4/api/internal/adapter/subscriber"
	"github.com/artem-smola/golang-course/task4/api/internal/controller"
	"github.com/artem-smola/golang-course/task4/api/internal/usecase"
	"github.com/artem-smola/golang-course/task4/platform/httpserver"
	"github.com/artem-smola/golang-course/task4/platform/logger"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/artem-smola/golang-course/task4/api/docs"
)

// @title Repo Stat API
// @version 1.0
// @description REST gateway for service health and repository info.
// @host localhost:8080
// @BasePath /
func run(ctx context.Context) error {
	// config
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	// logger

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	log.Info("starting server...")
	log.Debug("debug messages are enabled")

	subscriberClient, err := subscriber.NewClient(cfg.Services.Subscriber, log)
	if err != nil {
		return fmt.Errorf("cannot init subscriber adapter: %w", err)
	}
	defer func() {
		if err := subscriberClient.Close(); err != nil {
			log.Error("cannot close subscriber adapter", "error", err)
		}
	}()

	processorClient, err := processor.NewClient(cfg.Services.Processor, log)
	if err != nil {
		return fmt.Errorf("cannot init processor adapter: %w", err)
	}
	defer func() {
		if err := processorClient.Close(); err != nil {
			log.Error("cannot close processor adapter", "error", err)
		}
	}()

	pingUseCase := usecase.NewPing(subscriberClient, processorClient)
	getRepoInfoUseCase := usecase.NewGetRepoInfoUsecase(processorClient)
	subscriptionUseCase := usecase.NewSubscription(subscriberClient)
	getSubscriptionsRepoInfoUseCase := usecase.NewGetSubscriptionsRepoInfoUsecase(processorClient)

	server := controller.NewServer(
		log,
		pingUseCase,
		getRepoInfoUseCase,
		subscriptionUseCase,
		getSubscriptionsRepoInfoUseCase,
	)
	handler := server.Handler()
	router, ok := handler.(*gin.Engine)
	if !ok {
		return fmt.Errorf("unexpected handler type: %T", handler)
	}

	docs.SwaggerInfo.Host = cfg.HTTP.Address
	if host := os.Getenv("SWAGGER_HOST"); host != "" {
		docs.SwaggerInfo.Host = host
	}
	scheme := os.Getenv("SWAGGER_SCHEME")
	if scheme == "" {
		scheme = "http"
	}
	docs.SwaggerInfo.Schemes = []string{scheme}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// server
	srv := httpserver.NewServer(cfg.HTTP, router)
	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run http server: %w", err)
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
	cancel()
}
