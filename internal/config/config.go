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
	BotUsername       string
	BotPollingTimeout time.Duration
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

	return &Config{
		IsProd:            os.Getenv("IS_PROD") == "True",
		SuperAdmin:        superAdmin,
		BotToken:          os.Getenv("BOT_TOKEN"),
		BotUsername:       os.Getenv("BOT_USERNAME"),
		BotPollingTimeout: duration,
		PostgresqlDSN:     os.Getenv("POSTGRESQL_DSN"),
	}, nil
}
