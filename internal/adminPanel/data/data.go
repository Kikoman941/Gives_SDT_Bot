package data

import "gopkg.in/telebot.v3"

func InitMenu() {
	START_MENU.Reply(
		[]telebot.Btn{
			CREATE_GIVE_BUTTON,
			MY_GIVES_BUTTON,
		},
	)

	CANCEL_MENU.Reply(
		[]telebot.Btn{
			BACK_TO_START_BUTTON,
		},
	)
}
