package data

import "gopkg.in/telebot.v3"

type Buttons struct {
	CreateGiveButton telebot.Btn
	MyGivesButton    telebot.Btn
	MainMenuButton   telebot.Btn
}

func newButtons(menus *Menus) *Buttons {
	return &Buttons{
		CreateGiveButton: telebot.Btn{Text: "Новый конкурс"},
		MyGivesButton:    telebot.Btn{Text: "Мои конкурсы"},
		MainMenuButton:   telebot.Btn{Text: "Назад в главное меню"},
	}
}
