package server

import "github.com/go-kratos/kratos/v2/log"

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
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

type JWTConfig struct {
	Secret   string `yaml:"secret"`
	Expiry   int    `yaml:"expiry"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Name:     "auth-service",
			HttpPort: 8006,
			GrpcPort: 9006,
		},
		Database: DatabaseConfig{
			Url: "root:password@tcp(localhost:3306)/auth_db?charset=utf8mb4&parseTime=True&loc=Local",
		},
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
		JWT: JWTConfig{
			Secret: "your-secret-key",
			Expiry: 86400,
		},
		Log: LogConfig{
			Level: "info",
		},
	}
}

func NewLogger(cfg *Config) log.Logger {
	return log.NewStdLogger(log.NewHelper(log.DefaultLogger).Logger())
}