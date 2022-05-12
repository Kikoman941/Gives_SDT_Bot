package telegram

import (
	"Gives_SDT_Bot/internal/config"
	"Gives_SDT_Bot/internal/storage"
	"Gives_SDT_Bot/internal/telegram/data"
	"errors"
	"gopkg.in/telebot.v3"
)

type Bot struct {
	bot          *telebot.Bot
	fsm          *FSM
	telegramData *data.TelegramData
	storage      storage.BotStorage
}

func NewBot(cfg *config.Config, storage storage.BotStorage) (*Bot, error) {
	b, err := telebot.NewBot(
		telebot.Settings{
			Token: cfg.BotToken,
			Poller: &telebot.LongPoller{
				Timeout: cfg.BotPollingTimeout,
			},
		},
	)
	if err != nil {
		return nil, errors.New("cannot create bot")
	}

	fsm := NewFSM(storage)
	td := data.NewTelegramData()

	return &Bot{
		bot:          b,
		fsm:          fsm,
		telegramData: td,
		storage:      storage,
	}, nil
}

func (b *Bot) Start() {
	b.bot.Start()
}
