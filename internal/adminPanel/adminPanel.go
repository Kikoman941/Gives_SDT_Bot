package adminPanel

import (
	"Gives_SDT_Bot/internal/fsm"
	"Gives_SDT_Bot/pkg/logging"
	"gopkg.in/telebot.v3"
)

type AdminPanel struct {
	bot    *telebot.Bot
	logger *logging.Logger
	fsm    *fsm.FSM
}

func NewAdminPanel(bot *telebot.Bot, logger *logging.Logger) (*AdminPanel, error) {
	return &AdminPanel{
		bot:    bot,
		logger: logger,
	}, nil
}
