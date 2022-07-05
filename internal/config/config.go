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
	SuperAdmin        int
	BotToken          string
	BotPollingTimeout time.Duration
	PublisherTimeout  time.Duration
	PostgresqlDSN     string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, errors.New("cannot load env")
	}

	superAdmin, err := strconv.Atoi(os.Getenv("SUPERADMIN"))
	if err != nil {
		return nil, errors.New("cannot parse superadmin id")
	}

	duration, err := time.ParseDuration(os.Getenv("BOT_POLLING_TIMEOUT"))
	if err != nil {
		return nil, errors.New("cannot parse bot polling timeout")
	}

	publisherTimeout, err := time.ParseDuration(os.Getenv("PUBLISHER_TIMEOUT"))
	if err != nil {
		return nil, errors.New("cannot parse publisher timeout")
	}

	return &Config{
		IsProd:            os.Getenv("IS_PROD") == "True",
		SuperAdmin:        superAdmin,
		BotToken:          os.Getenv("BOT_TOKEN"),
		BotPollingTimeout: duration,
		PublisherTimeout:  publisherTimeout,
		PostgresqlDSN:     os.Getenv("POSTGRESQL_DSN"),
	}, nil
}
