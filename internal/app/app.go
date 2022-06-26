package app

import (
	"Gives_SDT_Bot/internal/adminPanel"
	"Gives_SDT_Bot/internal/config"
	"Gives_SDT_Bot/internal/fsm"
	fsmDB "Gives_SDT_Bot/internal/fsm/db"
	"Gives_SDT_Bot/internal/give"
	giveDB "Gives_SDT_Bot/internal/give/db"
	"Gives_SDT_Bot/internal/images"
	imagesDB "Gives_SDT_Bot/internal/images/db"
	"Gives_SDT_Bot/internal/publisher"
	"Gives_SDT_Bot/internal/user"
	userDB "Gives_SDT_Bot/internal/user/db"
	"Gives_SDT_Bot/pkg/client/postgresql"
	"Gives_SDT_Bot/pkg/logging"
	"context"
	"gopkg.in/telebot.v3"
)

type App struct {
	config     *config.Config
	logger     *logging.Logger
	bot        *telebot.Bot
	publisher  *publisher.Publisher
	adminPanel *adminPanel.AdminPanel
}

func NewApp(config *config.Config, logger *logging.Logger) *App {
	ctx := context.TODO()

	logger.Info("Initialization postgresql client")
	postgresqlClient, err := postgresql.NewClient(ctx, config.PostgreDSN)
	if err != nil {
		logger.Fatalf("cannot get postgresql client: %s", err)
		return nil
	}

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
		logger.Fatalf("cannot create telegram bot: %s", err)
		return nil
	}

	logger.Info("Initialization local image service")

	userRepo := userDB.NewRepository(postgresqlClient, logger)
	fsmRepo := fsmDB.NewRepository(postgresqlClient, logger)
	giveRepo := giveDB.NewRepository(postgresqlClient, logger)
	imagesRepo := imagesDB.NewRepository(bot, logger)

	userService := user.NewUserService(userRepo, logger)
	fsmService := fsm.NewFSMService(fsmRepo, logger)
	giveService := give.NewGiveService(giveRepo, logger)
	imagesService := images.NewImagesService(imagesRepo, logger)

	logger.Info("Initialization publisher")
	pub, err := publisher.NewPublisher(bot, logger)
	if err != nil {
		logger.Fatalf("cannot init publisher: %s", err)
		return nil
	}

	logger.Info("Initialization admin panel")
	ap := adminPanel.NewAdminPanel(
		bot,
		config.BotUsername,
		config.SuperAdmin,
		userService,
		giveService,
		fsmService,
		imagesService,
		logger,
	)

	return &App{
		config:     config,
		logger:     logger,
		bot:        bot,
		publisher:  pub,
		adminPanel: ap,
	}
}

func (a *App) Run() {
	a.adminPanel.InitHandlers()
	a.publisher.Run()
	a.bot.Start()
	a.logger.Info("App successfully started")
}
