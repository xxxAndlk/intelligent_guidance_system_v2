package server

import "github.com/go-kratos/kratos/v2/log"

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Log      LogConfig      `yaml:"log"`
}

type ServerConfig struct {
	Name      string `yaml:"name"`
	HttpPort  int    `yaml:"http_port"`
	GrpcPort  int    `yaml:"grpc_port"`
}

type DatabaseConfig struct {
	Url string `yaml:"url"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Name:     "payment-service",
			HttpPort: 8010,
			GrpcPort: 9010,
		},
		Database: DatabaseConfig{
			Url: "root:password@tcp(localhost:3306)/payment_db?charset=utf8mb4&parseTime=True&loc=Local",
		},
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
		Log: LogConfig{
			Level: "info",
		},
	}
}

func NewLogger(cfg *Config) log.Logger {
	return log.NewStdLogger(log.NewHelper(log.DefaultLogger).Logger())
}