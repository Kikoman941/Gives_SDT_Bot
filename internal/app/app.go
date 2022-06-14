package app

import (
	"Gives_SDT_Bot/internal/config"
	"Gives_SDT_Bot/internal/service/adminPanel"
	"Gives_SDT_Bot/internal/service/publisher"
	"Gives_SDT_Bot/pkg/logging"
	"errors"
	"gopkg.in/telebot.v3"
)

type App struct {
	config     *config.Config
	logger     *logging.Logger
	bot        *telebot.Bot
	publisher  *publisher.Publisher
	adminPanel *adminPanel.AdminPanel
}

func NewApp(config *config.Config, logger *logging.Logger) (*App, error) {
	bot, err := telebot.NewBot(
		telebot.Settings{
			Token: config.BotToken,
			Poller: &telebot.LongPoller{
				Timeout: config.BotPollingTimeout,
			},
		},
	)
	if err != nil {
		return nil, errors.New("cannot create bot")
	}

	pub, err := publisher.NewPublisher(bot, logger)
	if err != nil {
		return nil, errors.New("cannot create publisher")
	}

	ap, err := adminPanel.NewAdminPanel(bot, logger)
	if err != nil {
		return nil, errors.New("cannot create admin panel")
	}

	return &App{
		config:     config,
		logger:     logger,
		bot:        bot,
		publisher:  pub,
		adminPanel: ap,
	}, nil
}

func (a *App) Run() {
	a.publisher.Run()
	a.bot.Start()
	a.logger.Info("App successfully started")
}
