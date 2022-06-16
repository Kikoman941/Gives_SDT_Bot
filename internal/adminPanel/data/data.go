package data

import "gopkg.in/telebot.v3"

func InitMenus() {
	START_MENU.Reply(
		[]telebot.Btn{
			CREATE_GIVE_BUTTON,
			MY_GIVES_BUTTON,
		},
	)
}
