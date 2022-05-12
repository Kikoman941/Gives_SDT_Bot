package telegram

import (
	"gopkg.in/telebot.v3"
	"log"
	"strconv"
)

func getChatId(str string) *telebot.ChatID {
	var chatId telebot.ChatID
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	chatId = telebot.ChatID(i)
	return &chatId
}
