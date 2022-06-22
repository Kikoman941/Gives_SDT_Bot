package data

import "gopkg.in/telebot.v3"

var (
	START_MENU  = &telebot.ReplyMarkup{ResizeKeyboard: true}
	CANCEL_MENU = &telebot.ReplyMarkup{ResizeKeyboard: true}
)

func CreateReplyMenu() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{}
}

func CreateInlineMenu() *telebot.ReplyMarkup {
	menu := &telebot.ReplyMarkup{}
	return menu
}
