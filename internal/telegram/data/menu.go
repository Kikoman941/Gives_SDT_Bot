package data

import "gopkg.in/telebot.v3"

type Menus struct {
	StartMenu *telebot.ReplyMarkup
}

func newMenus() *Menus {
	return &Menus{
		StartMenu: &telebot.ReplyMarkup{ResizeKeyboard: true},
	}
}

func (m *Menus) init(buttons *Buttons) {
	m.StartMenu.Reply(
		m.StartMenu.Row(buttons.CreateGiveButton, buttons.MyGivesButton),
	)
}
