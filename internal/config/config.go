package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	errConfigEnvEmpty    = errors.New("config path env is not set")
	errConfigFileInvalid = errors.New("config file has invalid format")
)

type EnvType string

const (
	EnvLocal   EnvType = "local"
	EnvDev     EnvType = "dev"
	EnvTesting EnvType = "testing"

	EnvConfigPath string = "APP_CONFIG_PATH"
)

type Config struct {
	Env             EnvType         `yaml:"env" env-required:"true"`
	ApplicationName string          `yaml:"application_name" env-required:"true"`
	Storage         Storage         `yaml:"storage" env-required:"true"`
	GrpcServer      GrpcServer      `yaml:"grpc_server" env-required:"true"`
	Migration       MigrationConfig `yaml:"migrations" env-required:"true"`
	Timeouts        Timeouts        `yaml:"timeouts" env-required:"true"`
	Tracing         Tracing         `yaml:"tracing" env-required:"true"`
	Redis           Redis           `yaml:"redis" env-required:"true"`
	Rabbit          Rabbit          `yaml:"rabbit" env-required:"true"`
}

type Storage struct {
	Host         string `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	Port         string `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	DatabaseName string `yaml:"database_name" env:"POSTGRES_DATABASE_NAME" env-required:"true"`
	User         string `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
	Password     string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
}

type GrpcServer struct {
	Port    string `yaml:"port" env:"GRPC_SERVER_PORT" env-required:"true"`
	Timeout string `yaml:"timeout" env:"GRPC_SERVER_TIMEOUT" env-required:"true"`
}

type MigrationConfig struct {
	MigrationsPath      string `yaml:"migrations_path" env:"MIGRATIONS_PATH" env-required:"true"`
	MigrationsTableName string `yaml:"migrations_table_name" env:"MIGRATIONS_TABLE_NAME" env-default:"migrations"`
}

type Timeouts struct {
	EndpointExecutionTimeoutMS string `yaml:"endpoint_execution_timeout_ms" env:"ENDPOINT_EXECUTION_TIMEOUT_MS" env-required:"true"`
}

type Tracing struct {
	Host                  string `yaml:"host" env:"JAEGER_HOST" env-required:"true"`
	Port                  string `yaml:"port" env:"JAEGER_PORT" env-required:"true"`
	CredentialsTracerName string `yaml:"credentials_tracer_name" env:"CREDENTIALS_TRACER_NAME" env-required:"true"`
	TokenTracerName       string `yaml:"token_tracer_name" env:"TOKEN_TRACER_NAME" env-required:"true"`
}

type Redis struct {
	User           string `yaml:"user" env:"REDIS_USER" env-required:"true"`
	Password       string `yaml:"password" env:"REDIS_PASSWORD" env-required:"true"`
	Host           string `yaml:"host" env:"REDIS_HOST" env-required:"true"`
	Port           string `yaml:"port" env:"REDIS_PORT" env-required:"true"`
	DatabaseNumber string `yaml:"database_number" env:"REDIS_DATABASE_NUMBER" env-required:"true"`
	EntityTTL      string `yaml:"entity_ttl" env:"REDIS_ENTITY_TTL" env-required:"true"`
}

type Rabbit struct {
	User     string         `yaml:"user" env:"RABBIT_USER" env-required:"true"`
	Password string         `yaml:"password" env:"RABBIT_PASSWORD" env-required:"true"`
	Host     string         `yaml:"host" env:"RABBIT_HOST" env-required:"true"`
	Port     string         `yaml:"port" env:"RABBIT_PORT" env-required:"true"`
	Exchange RabbitExchange `yaml:"exchange" env-required:"true"`
	Queues   Queues         `yaml:"queues" env-required:"true"`
}

type Queues struct {
	CredentialsCreate CredentialsCreateQueue `yaml:"credentials_create" env-required:"true"`
}

type RabbitExchange struct {
	Name string `yaml:"name" env:"RABBIT_EXCHANGE_NAME" env-required:"true"`
	Type string `yaml:"type" env:"RABBIT_EXCHANGE_TYPE" env-required:"true"`
}

type CredentialsCreateQueue struct {
	Name       string `yaml:"name" env:"RABBIT_CREDENTIALS_CREATE_QUEUE_NAME" env-required:"true"`
	RoutingKey string `yaml:"routing_key" env:"RABBIT_CREDENTIALS_CREATE_QUEUE_ROUTING_KEY" env-required:"true"`
}

func Load() (*Config, error) {
	op := runtime.FuncForPC(reflect.ValueOf(Load).Pointer()).Name()

	pathToConfig := os.Getenv(EnvConfigPath)
	if pathToConfig == "" {
		return nil, fmt.Errorf("%s -> %w", op, errConfigEnvEmpty)
	}

	return loadByPath(pathToConfig)
}

func loadByPath(pathToConfig string) (*Config, error) {
	op := runtime.FuncForPC(reflect.ValueOf(loadByPath).Pointer()).Name()

	var config Config

	if err := cleanenv.ReadConfig(pathToConfig, &config); err != nil {
		return nil, fmt.Errorf("%s -> %w; %w", op, errConfigFileInvalid, err)
	}

	return &config, nil
}

func (cfg *Config) GetEndpointExecutionTimeout() time.Duration {
	op := runtime.FuncForPC(reflect.ValueOf(cfg.GetEndpointExecutionTimeout).Pointer()).Name()

	milliseconds, err := strconv.Atoi(cfg.Timeouts.EndpointExecutionTimeoutMS)
	if err != nil {
		log.Fatal(fmt.Sprintf("%s, op=%s, err=%s", "unable to convert endpoint timeout duration", op, err.Error()))
	}

	return time.Duration(milliseconds) * time.Millisecond
}

func (cfg *Config) GetPostgresConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.User,
		cfg.Storage.Password,
		cfg.Storage.DatabaseName,
	)
}

func (cfg *Config) GetMigrationPostgresConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table=%s",
		cfg.Storage.User,
		cfg.Storage.Password,
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.DatabaseName,
		cfg.Migration.MigrationsTableName,
	)
}

func (cfg *Config) GetRedisEntityTTL() time.Duration {
	op := runtime.FuncForPC(reflect.ValueOf(cfg.GetRedisEntityTTL).Pointer()).Name()

	parsedTTL, err := time.ParseDuration(cfg.Redis.EntityTTL)
	if err != nil {
		log.Fatal(fmt.Sprintf("%s, op=%s, err=%s", "unable to parse entity ttl into time.Duration", op, err.Error()))
	}

	return parsedTTL
}

func (cfg *Config) GetRabbitConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.Rabbit.User,
		cfg.Rabbit.Password,
		cfg.Rabbit.Host,
		cfg.Rabbit.Port,
	)
}
