package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Server  ServerConfig
	PG      PGConfig
	Redis   RedisConfig
	JWT     JWTConfig
	Hasher  HasherConfig
	Decoder DecoderConfig
	Logger  LoggerConfig
	Speller SpellerConfig
}

type (
	ServerConfig struct {
		Host    string `env:"HTTP_HOST"`
		Port    string `env:"HTTP_PORT"`
		Address string
	}

	PGConfig struct {
		URL      string
		User     string `env:"POSTGRES_USER"`
		Password string `env:"POSTGRES_PASSWORD"`
		Host     string `env:"POSTGRES_HOST"`
		Port     uint16 `env:"POSTGRES_PORT"`
		DB       string `env:"POSTGRES_DB"`
		SSLMode  string `env:"POSTGRES_SSL_MODE"`
	}

	RedisConfig struct {
		Host     string `env:"REDIS_HOST"`
		Port     string `env:"REDIS_PORT"`
		Password string `env:"REDIS_PASSWORD"`
		DB       int    `env:"REDIS_DB"`
	}

	JWTConfig struct {
		SignKey string `env:"JWT_SIGN_KEY"`
	}

	HasherConfig struct {
		Salt string `env:"HASHER_SALT"`
	}

	DecoderConfig struct {
		SecretKey string `env:"SECRET_KEY"`
	}

	LoggerConfig struct {
		LogFilePath string `env:"LOG_FILE_PATH"`
		Level       string `env:"LOG_LVL"`
	}
	SpellerConfig struct {
		URL string `env:"SPELLER_URL"`
	}
)

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}
	cfg := &Config{}

	cfg.Server.Host = os.Getenv("HTTP_HOST")
	cfg.Server.Port = os.Getenv("HTTP_PORT")
	cfg.Server.Address = fmt.Sprintf("%s%s", cfg.Server.Host, cfg.Server.Port)
	cfg.PG.User = os.Getenv("POSTGRES_USER")
	cfg.PG.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.PG.Host = os.Getenv("POSTGRES_HOST")
	cfg.PG.Port = parseUint16Env("POSTGRES_PORT")
	cfg.PG.DB = os.Getenv("POSTGRES_DB")
	cfg.PG.SSLMode = os.Getenv("POSTGRES_SSL_MODE")
	cfg.PG.URL = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.PG.User, cfg.PG.Password, cfg.PG.Host, cfg.PG.Port, cfg.PG.DB, cfg.PG.SSLMode)
	cfg.Redis.Host = os.Getenv("REDIS_HOST")
	cfg.Redis.Port = os.Getenv("REDIS_PORT")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
	cfg.Redis.DB, _ = strconv.Atoi("REDIS_DB")
	cfg.JWT.SignKey = os.Getenv("JWT_SIGN_KEY")
	cfg.Hasher.Salt = os.Getenv("HASHER_SALT")
	cfg.Decoder.SecretKey = os.Getenv("SECRET_KEY")
	cfg.Logger.LogFilePath = os.Getenv("LOG_FILE_PATH")
	cfg.Logger.Level = os.Getenv("LOG_LVL")
	cfg.Speller.URL = os.Getenv("SPELLER_URL")
	return cfg, err
}

func parseUint16Env(key string) uint16 {
	val, err := strconv.ParseUint(os.Getenv(key), 10, 16)
	if err != nil {
		log.Printf("failed to parse %s: %v\n", key, err)
	}
	return uint16(val)
}
