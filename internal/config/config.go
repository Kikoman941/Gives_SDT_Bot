package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

type Config struct {
	IsProd            bool
	SuperAdmin        int64
	BotToken          string
	BotPollingTimeout time.Duration
	PostgreDSN        string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, errors.New("cannot load env")
	}

	superadmin, err := strconv.ParseInt(os.Getenv("SUPERADMIN"), 10, 64)
	if err != nil {
		return nil, errors.New("cannot parse superadmin id")
	}

	duration, err := time.ParseDuration(os.Getenv("BOT_POLLING_TIMEOUT"))
	if err != nil {
		return nil, errors.New("cannot parse bot polling timeout")
	}

	return &Config{
		IsProd:            os.Getenv("IS_PROD") == "True",
		SuperAdmin:        superadmin,
		BotToken:          os.Getenv("BOT_TOKEN"),
		BotPollingTimeout: duration,
		PostgreDSN:        os.Getenv("POSTGRESQL_DSN"),
	}, nil
}
