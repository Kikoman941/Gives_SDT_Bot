package data

import "gopkg.in/telebot.v3"

type Buttons struct {
	CreateGiveButton telebot.Btn
	MyGivesButton    telebot.Btn
	MainMenuButton   telebot.Btn
}

func newButtons(menus *Menus) *Buttons {
	return &Buttons{
		CreateGiveButton: telebot.Btn{Text: "Новый конкурс 🎁", Data: "test_inline_callback_CREATE", InlineQuery: "query"},
		MyGivesButton:    telebot.Btn{Text: "Мои конкурсы 📋", Data: "test_inline_callback_GIVES", InlineQuery: "query"},
		MainMenuButton:   telebot.Btn{Text: "Назад в главное меню 🔙"},
	}
}
