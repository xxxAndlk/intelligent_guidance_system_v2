package server

import (
	"time"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Qdrant   QdrantConfig   `yaml:"qdrant"`
	LLM      LLMConfig      `yaml:"llm"`
	Redis    RedisConfig    `yaml:"redis"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
	Consul   ConsulConfig   `yaml:"consul"`
}

type ServerConfig struct {
	GRPC GRPCConfig `yaml:"grpc"`
	HTTP HTTPConfig `yaml:"http"`
}

type GRPCConfig struct {
	Addr    string        `yaml:"addr"`
	Timeout time.Duration `yaml:"timeout"`
}

type HTTPConfig struct {
	Addr    string        `yaml:"addr"`
	Timeout time.Duration `yaml:"timeout"`
}

type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func (c DatabaseConfig) DSN() string {
	return c.User + ":" + c.Password + "@tcp(" + c.Host + ":" + string(rune(c.Port)) + ")/" + c.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
}

type QdrantConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	DiseaseCollection string `yaml:"disease_collection"`
	DoctorCollection  string `yaml:"doctor_collection"`
	DrugCollection    string `yaml:"drug_collection"`
}

type LLMConfig struct {
	BaseURL     string        `yaml:"base_url"`
	APIKey      string        `yaml:"api_key"`
	Model       string        `yaml:"model"`
	EmbedModel  string        `yaml:"embed_model"`
	MaxTokens   int           `yaml:"max_tokens"`
	Temperature float32       `yaml:"temperature"`
	Timeout     time.Duration `yaml:"timeout"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type RabbitMQConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	VHost    string `yaml:"vhost"`
}

type ConsulConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Key  string `yaml:"key"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			GRPC: GRPCConfig{
				Addr:    ":9000",
				Timeout: 30 * time.Second,
			},
			HTTP: HTTPConfig{
				Addr:    ":8000",
				Timeout: 30 * time.Second,
			},
		},
		Database: DatabaseConfig{
			Driver:   "mysql",
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "root",
			Database: "ai_diagnosis",
		},
		Qdrant: QdrantConfig{
			Host:              "localhost",
			Port:              6333,
			DiseaseCollection: "diseases",
			DoctorCollection:  "doctors",
			DrugCollection:    "drugs",
		},
		LLM: LLMConfig{
			BaseURL:     "https://api.deepseek.com",
			Model:       "deepseek-chat",
			EmbedModel:  "deepseek-embedding",
			MaxTokens:   4096,
			Temperature: 0.7,
			Timeout:     60 * time.Second,
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
		RabbitMQ: RabbitMQConfig{
			Host:     "localhost",
			Port:     5672,
			User:     "guest",
			Password: "guest",
			VHost:    "/",
		},
		Consul: ConsulConfig{
			Host: "localhost",
			Port: 8500,
			Key:  "ai-service/config",
		},
	}

	return cfg, nil
}