package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	// "github.com/joho/godotenv"
)

type Config struct {
	Env              string `env:"ENV" env-default:"local"`
	PGConfig         PostgreConfig
	ServerConfig     ServerConfig
	OtelConfig       OtelConfig
	GrinexConfig     GrinexConfig
	PrometheusConfig PrometheusConfig
}

type PostgreConfig struct {
	Username string `env:"PG_USERNAME" env-default:"user"`
	Host     string `env:"PG_HOST" env-default:"db"`
	Port     string `env:"PG_PORT" env-default:"5432"`
	DBName   string `env:"PG_DBNAME" env-default:"rates"`
	SSLMode  string `env:"PG_SSLMODE" env-default:"disable"`
	Password string `env:"PG_PASSWORD"`
}

type ServerConfig struct {
	GRPCPort int `env:"GRPC_PORT" env-default:"50051"`
}

type OtelConfig struct {
	URL         string `env:"OTEL_EXPORTER" env-default:"jaeger:4317"`
	ServiceName string `env:"OTEL_SERVICE_NAME" env-default:"usdt-rates-service"`
}

type PrometheusConfig struct {
	Port string `env:"METRICS_PORT" env-default:":9090"`
}

type GrinexConfig struct {
	URL string `env:"GRINEX_URL" env-default:"https://grinex.io"`
}

func MustLoad() *Config {
	// If we need to load .env file
	// if err := godotenv.Load(); err != nil {
	// }

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(fmt.Sprintf("cannot read config: %s", err))
	}

	dbHostFlag := flag.String("pg-host", "", "Database host")
	dbPortFlag := flag.String("pg-port", "", "Database host")
	dbUserFlag := flag.String("pg-user", "", "gRPC server port")
	dbNameFlag := flag.String("pg-dbname", "", "Database host")
	dbPasswordFlag := flag.String("pg-password", "", "gRPC server port")
	dbSSLModeFlag := flag.String("pg-sslmode", "", "Postgres SSL mode")

	flag.Parse()

	if *dbHostFlag != "" {
		cfg.PGConfig.Host = *dbHostFlag
	}
	if *dbPortFlag != "" {
		cfg.PGConfig.Port = *dbPortFlag
	}
	if *dbUserFlag != "" {
		cfg.PGConfig.Username = *dbUserFlag
	}
	if *dbNameFlag != "" {
		cfg.PGConfig.DBName = *dbNameFlag
	}
	if *dbPasswordFlag != "" {
		cfg.PGConfig.Password = *dbPasswordFlag
	}
	if *dbSSLModeFlag != "" {
		cfg.PGConfig.SSLMode = *dbSSLModeFlag
	}

	return &cfg
}

func DTO(cfg *Config) string {
	storage := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.PGConfig.Host, cfg.PGConfig.Port, cfg.PGConfig.Username, cfg.PGConfig.DBName, os.Getenv("PG_PASSWORD"), cfg.PGConfig.SSLMode)

	return storage
}
