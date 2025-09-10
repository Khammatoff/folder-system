package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Logging  LoggingConfig  `yaml:"logging"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type JWTConfig struct {
	AccessSecret  string `yaml:"access_secret"`
	RefreshSecret string `yaml:"refresh_secret"`
	AccessTTL     int    `yaml:"access_ttl"`
	RefreshTTL    int    `yaml:"refresh_ttl"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

func LoadConfig(path string) (*Config, error) {
	// Открываем файл конфигурации
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Читаем содержимое файла
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Парсим YAML
	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	// Устанавливаем значения по умолчанию, если нужно
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}

	return cfg, nil
}

func (c *Config) SaveConfig(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
