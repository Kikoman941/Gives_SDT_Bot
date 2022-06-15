package app

import (
	"Gives_SDT_Bot/internal/adminPanel"
	"Gives_SDT_Bot/internal/config"
	"Gives_SDT_Bot/internal/publisher"
	"Gives_SDT_Bot/pkg/errors"
	"Gives_SDT_Bot/pkg/localImages"
	"Gives_SDT_Bot/pkg/logging"
	"gopkg.in/telebot.v3"
)

type App struct {
	config      *config.Config
	logger      *logging.Logger
	bot         *telebot.Bot
	publisher   *publisher.Publisher
	adminPanel  *adminPanel.AdminPanel
	localImages *localImages.LocalImages
}

func NewApp(config *config.Config, logger *logging.Logger) (*App, error) {
	logger.Info("Creating telegram bot")
	bot, err := telebot.NewBot(
		telebot.Settings{
			Token: config.BotToken,
			Poller: &telebot.LongPoller{
				Timeout: config.BotPollingTimeout,
			},
		},
	)
	if err != nil {
		return nil, errors.FormatError("cannot create telegram bot", err)
	}

	logger.Info("Initialization local image service")
	lImages, err := localImages.NewLocalImage(".images", logger)
	if err != nil {
		return nil, err
	}

	logger.Info("Initialization publisher")
	pub, err := publisher.NewPublisher(bot, logger)
	if err != nil {
		return nil, err
	}

	logger.Info("Initialization admin panel")
	ap, err := adminPanel.NewAdminPanel(bot, logger)
	if err != nil {
		return nil, err
	}

	return &App{
		config:      config,
		logger:      logger,
		bot:         bot,
		publisher:   pub,
		adminPanel:  ap,
		localImages: lImages,
	}, nil
}

func (a *App) Run() {
	a.publisher.Run()
	a.bot.Start()
	a.logger.Info("App successfully started")
}