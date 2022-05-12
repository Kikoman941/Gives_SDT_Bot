package main

import (
	"Gives_SDT_Bot/internal/config"
	"Gives_SDT_Bot/internal/storage/postgre"
	"Gives_SDT_Bot/internal/telegram"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgre.NewPostgresDB(cfg.PostgreDSN)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := telegram.NewBot(cfg, db)
	if err != nil {
		log.Fatal(err)
	}

	bot.InitHandlers()
	bot.Start()
}
