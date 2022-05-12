package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Config struct {
	IsProd            bool
	BotToken          string
	BotPollingTimeout time.Duration
	PostgreDSN        string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, errors.New("cannot load env")
	}

	duration, err := time.ParseDuration(os.Getenv("BOT_POLLING_TIMEOUT"))
	if err != nil {
		return nil, errors.New("cannot parse bot polling timeout")
	}

	return &Config{
		IsProd:            os.Getenv("IS_PROD") == "True",
		BotToken:          os.Getenv("BOT_TOKEN"),
		BotPollingTimeout: duration,
		PostgreDSN:        os.Getenv("POSTGRE_DSN"),
	}, nil
}
