package telegram

import (
	"gopkg.in/telebot.v3"
)

func getChatIdFromInt(i int) telebot.ChatID {
	i64 := int64(i)
	return telebot.ChatID(i64)
}
