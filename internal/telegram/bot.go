package telegram

import (
	"Gives_SDT_Bot/internal/config"
	"errors"
	"gopkg.in/telebot.v3"
)

type Bot struct {
	bot        *telebot.Bot
	adminGroup []int64
}

func NewBot(cfg *config.Config) (*Bot, error) {
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

	var adminGroup []int64
	adminGroup = append(adminGroup, cfg.SuperAdmin)

	return &Bot{
		bot:        b,
		adminGroup: adminGroup,
	}, nil
}

func (b *Bot) Start() {
	b.bot.Start()
}

func (b *Bot) refreshAdmins(admins []int64) {
	b.adminGroup = []int64{}
	b.adminGroup = append(b.adminGroup, admins...)
}
