package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	// Окружение: local, dev, prod
	Env string `env:"ENV" env-default:"local"`

	// PostgreSQL
	DB struct {
		Host     string `env:"DB_HOST" env-default:"localhost"`
		Port     int    `env:"DB_PORT" env-default:"5432"`
		User     string `env:"DB_USER" env-default:"postgres"`
		Password string `env:"DB_PASSWORD" env-default:"postgres"`
		Name     string `env:"DB_NAME" env-default:"notes_service"`
		SSLMode  string `env:"DB_SSLMODE" env-default:"disable"`
	}

	// HTTP Server
	HTTPServer struct {
		Address     string        `env:"HTTP_ADDRESS" env-default:"localhost:8083"`
		Timeout     time.Duration `env:"HTTP_TIMEOUT" env-default:"4s"`
		IdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
		User        string        `env:"HTTP_USER" env-default:"user"`
		Password    string        `env:"HTTP_PASSWORD" env-default:"user"`
	}
}

func MustLoad() *Config {
	// 1. Загружаем .env файл (если существует)
	// Игнорируем ошибку — переменные могут быть установлены напрямую
	_ = godotenv.Load()

	// 2. Создаём структуру конфигурации
	cfg := &Config{}

	// 3. Загружаем переменные окружения в структуру
	// Теперь НЕ читаем файл конфигурации — только переменные окружения!
	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatalf("Cannot read environment variables: %s", err)
	}

	// 4. Валидация
	validate(cfg)

	return cfg
}
func validate(cfg *Config) {
	// Проверка окружения
	allowedEnvs := map[string]bool{
		"local": true,
		"dev":   true,
		"prod":  true,
	}
	if !allowedEnvs[cfg.Env] {
		log.Fatalf("Invalid ENV: %s (allowed: local, dev, prod)", cfg.Env)
	}

	// Проверка порта БД
	if cfg.DB.Port < 1 || cfg.DB.Port > 65535 {
		log.Fatalf("Invalid DB_PORT: %d (must be 1-65535)", cfg.DB.Port)
	}

	// Проверка таймаутов
	if cfg.HTTPServer.Timeout <= 0 {
		log.Fatal("HTTP_TIMEOUT must be positive")
	}
	if cfg.HTTPServer.IdleTimeout <= 0 {
		log.Fatal("HTTP_IDLE_TIMEOUT must be positive")
	}
}

func (c *Config) StoragePath() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DB.Host,
		c.DB.Port,
		c.DB.User,
		c.DB.Password,
		c.DB.Name,
		c.DB.SSLMode,
	)
}
