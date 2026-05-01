package config

import (
	"github.com/artem-smola/golang-course/task4/platform/env"
	"github.com/artem-smola/golang-course/task4/platform/grpcserver"
	"github.com/artem-smola/golang-course/task4/platform/logger"
	"time"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-subscriber"`
}

type Postgres struct {
	DSN      string `yaml:"dsn" env:"POSTGRES_DSN" env-default:"postgres://postgres:postgres@localhost:5432/repo_stat?sslmode=disable"`
	MinConns int32  `yaml:"min_conns" env:"POSTGRES_MIN_CONNS" env-default:"1"`
	MaxConns int32  `yaml:"max_conns" env:"POSTGRES_MAX_CONNS" env-default:"10"`
}

type GitHub struct {
	BaseURL string        `yaml:"base_url" env:"GITHUB_BASE_URL" env-default:"https://api.github.com"`
	Token   string        `yaml:"token" env:"GITHUB_TOKEN"`
	Timeout time.Duration `yaml:"timeout" env:"GITHUB_TIMEOUT" env-default:"10s"`
}

type Migrations struct {
	Path string `yaml:"path" env:"MIGRATIONS_PATH" env-default:"file://subscriber/migrations"`
}

type Config struct {
	App        App               `yaml:"app"`
	GRPC       grpcserver.Config `yaml:"grpc"`
	Logger     logger.Config     `yaml:"logger"`
	Postgres   Postgres          `yaml:"postgres"`
	GitHub     GitHub            `yaml:"github"`
	Migrations Migrations        `yaml:"migrations"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
