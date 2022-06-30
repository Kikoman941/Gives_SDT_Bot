package main

import (
	"Gives_SDT_Bot/internal/app"
	"Gives_SDT_Bot/internal/config"
	"Gives_SDT_Bot/pkg/logging"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logging.Init(cfg.IsProd)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Creating application")
	a := app.NewApp(cfg, logger)

	if err := a.Run(); err != nil {
		logger.Fatal(err)
	}
}
