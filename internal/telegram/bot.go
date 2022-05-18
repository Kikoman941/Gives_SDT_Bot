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
	adminGroup   []int64
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

	var adminGroup []int64
	adminGroup = append(adminGroup, cfg.Superadmin)

	return &Bot{
		bot:          b,
		fsm:          fsm,
		telegramData: td,
		adminGroup:   adminGroup,
		storage:      storage,
	}, nil
}

func (b *Bot) Start() {
	b.bot.Start()
}

func (b *Bot) refreshAdmins(admins []int64) {
	b.adminGroup = []int64{}
	b.adminGroup = append(b.adminGroup, admins...)
}
