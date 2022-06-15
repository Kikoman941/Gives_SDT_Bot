package publisher

import (
	"Gives_SDT_Bot/pkg/logging"
	"fmt"
	"gopkg.in/telebot.v3"
	"time"
)

type Publisher struct {
	bot    *telebot.Bot
	logger *logging.Logger
}

func NewPublisher(bot *telebot.Bot, logger *logging.Logger) (*Publisher, error) {
	return &Publisher{
		bot:    bot,
		logger: logger,
	}, nil
}

func (p *Publisher) Run() {
	go func() {
		for {
			fmt.Println("test publisger")
			time.Sleep(time.Second * 5)
		}
	}()
}
