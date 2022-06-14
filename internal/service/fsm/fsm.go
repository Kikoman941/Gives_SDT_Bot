package fsm

import (
	"Gives_SDT_Bot/pkg/logging"
	"gopkg.in/telebot.v3"
)

type FSM struct {
	bot    *telebot.Bot
	logger *logging.Logger
}

func NewFSM(bot *telebot.Bot, logger *logging.Logger) (*FSM, error) {
	return &FSM{
		bot:    bot,
		logger: logger,
	}, nil
}
