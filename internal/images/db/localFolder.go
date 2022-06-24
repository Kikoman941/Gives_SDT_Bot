package db

import (
	"Gives_SDT_Bot/internal/images"
	"Gives_SDT_Bot/pkg/logging"
	"gopkg.in/telebot.v3"
)

type repository struct {
	bot    *telebot.Bot
	logger *logging.Logger
}

func (r *repository) Download(file *telebot.File, filePath string) error {
	if err := r.bot.Download(file, filePath); err != nil {
		return err
	}
	return nil
}

func NewRepository(bot *telebot.Bot, logger *logging.Logger) images.Repository {
	return &repository{
		bot:    bot,
		logger: logger,
	}
}
