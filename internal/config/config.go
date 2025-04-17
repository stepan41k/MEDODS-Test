package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Env         string     `yaml:"env" env-default:"local"`
	Server      HTTPServer `yaml:"http_server"`
	Storage     DataBase   `yaml:"db"`
}

type HTTPServer struct {
	Port         string        `yaml:"server_port" env-default:"localhost:8020"`
	Timeout      time.Duration `yaml:"timeout" env-default:"4s"`
	Idle_timeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DataBase struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}


func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		logrus.Fatalf("config path is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		logrus.Fatalf("config file does not exist: %s", err.Error())
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		logrus.Fatalf("cannot read config: %s", err.Error())
	}

	return &cfg
}
