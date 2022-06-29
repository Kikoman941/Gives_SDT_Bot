package data

import "gopkg.in/telebot.v3"

var (
	START_MENU         = &telebot.ReplyMarkup{ResizeKeyboard: true}
	CANCEL_MENU        = &telebot.ReplyMarkup{ResizeKeyboard: true}
	ACTIVATE_GIVE_MENU = &telebot.ReplyMarkup{ResizeKeyboard: true}
	EDIT_GIVE_MENU     = &telebot.ReplyMarkup{ResizeKeyboard: true}
	ACTIVE_GIVE_MENU   = &telebot.ReplyMarkup{ResizeKeyboard: true}
)

func CreateReplyMenu(buttons ...telebot.Btn) *telebot.ReplyMarkup {
	var rows []telebot.Row
	menu := &telebot.ReplyMarkup{ResizeKeyboard: true}

	for _, button := range buttons {
		rows = append(rows, telebot.Row{button})
	}

	menu.Reply(rows...)

	return menu
}

func CreateInlineMenu() *telebot.ReplyMarkup {
	menu := &telebot.ReplyMarkup{}
	return menu
}
