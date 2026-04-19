package config

import (
	"github.com/artem-smola/golang-course/task4/platform/env"
	"github.com/artem-smola/golang-course/task4/platform/grpcserver"
	"github.com/artem-smola/golang-course/task4/platform/logger"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-processor"`
}

type Services struct {
	Collector string `yaml:"collector" env:"COLLECTOR_ADDRESS" env-default:"localhost:8080"`
}

type Config struct {
	App      App               `yaml:"app"`
	Services Services          `yaml:"services"`
	GRPC     grpcserver.Config `yaml:"grpc"`
	Logger   logger.Config     `yaml:"logger"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
