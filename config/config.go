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
}

type (
	ServerConfig struct {
		Host string `env:"SERVER_HOST"`
		Port string `env:"SERVER_PORT"`
	}

	PGConfig struct {
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
)

//func NewConfig() (*Config, error) {
//	cfg := Config{}
//	err := env.Parse(&cfg)
//	if err != nil {
//		return nil, fmt.Errorf("failed to parse config from environment variables: %w", err)
//	}
//
//	return &cfg, nil
//}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	cfg := Config{
		Server: ServerConfig{
			Host: getKey("SERVER_HOST"),
			Port: getKey("SERVER_PORT"),
		},
		PG: PGConfig{
			User:     getKey("POSTGRES_USER"),
			Password: getKey("POSTGRES_PASSWORD"),
			Host:     getKey("POSTGRES_HOST"),
			Port:     parseUint16Env("POSTGRES_PORT"),
			DB:       getKey("POSTGRES_DB"),
			SSLMode:  getKey("POSTGRES_SSL_MODE"),
		},
		Redis: RedisConfig{
			Host:     getKey("REDIS_HOST"),
			Port:     getKey("REDIS_PORT"),
			Password: getKey("REDIS_PASSWORD"),
			DB:       parseIntEnv("REDIS_DB"),
		},
		JWT: JWTConfig{
			SignKey: getKey("JWT_SIGN_KEY"),
		},
		Hasher: HasherConfig{
			Salt: getKey("HASHER_SALT"),
		},
		Decoder: DecoderConfig{
			SecretKey: getKey("SECRET_KEY"),
		},
		Logger: LoggerConfig{
			LogFilePath: getKey("LOG_FILE_PATH"),
			Level:       getKey("LOG_LVL"),
		},
	}

	return &cfg, nil
}

func getKey(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
		return ""
	}

	return os.Getenv(key)
}

func parseUint16Env(key string) uint16 {
	val, err := strconv.ParseUint(os.Getenv(key), 10, 16)
	if err != nil {
		log.Printf("failed to parse %s: %v\n", key, err)
	}
	return uint16(val)
}

func parseIntEnv(key string) int {
	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		log.Printf("failed to parse %s: %v/n", key, err)
	}
	return val
}
