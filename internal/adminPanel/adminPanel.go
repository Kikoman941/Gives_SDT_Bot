package adminPanel

import (
	"Gives_SDT_Bot/internal/fsm"
	"Gives_SDT_Bot/pkg/logging"
	"gopkg.in/telebot.v3"
)

type AdminPanel struct {
	bot        *telebot.Bot
	adminGroup []int
	fsm        *fsm.Service
	logger     *logging.Logger
}

func NewAdminPanel(bot *telebot.Bot, fsm *fsm.Service, superAdmin int, logger *logging.Logger) (*AdminPanel, error) {
	var adminGroup []int
	adminGroup = append(adminGroup, superAdmin)

	return &AdminPanel{
		bot:        bot,
		adminGroup: adminGroup,
		fsm:        fsm,
		logger:     logger,
	}, nil
}
