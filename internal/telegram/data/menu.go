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
	m.StartMenu.Inline(
		m.StartMenu.Row(buttons.CreateGiveButton, buttons.MyGivesButton),
	)
}

func (m *Menus) CreateReplyMenu(buttons []telebot.Btn) *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{}
}

func (m *Menus) CreateInlineMenu() *telebot.ReplyMarkup {
	menu := &telebot.ReplyMarkup{}
	btn := menu.Data("TEST INLINE", "haha", "test_inline")
	menu.Inline(menu.Row(btn))
	return menu
}
