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

	EDIT_GIVE_MENU.Reply(
		[]telebot.Btn{
			EDIT_TITLE_BUTTON,
			EDIT_DESCRIPTION_BUTTON,
		},
		[]telebot.Btn{
			EDIT_START_BUTTON,
			EDIT_FINISH_BUTTON,
		},
		[]telebot.Btn{
			EDIT_IMAGE_BUTTON,
			EDIT_WINNERS_COUNT_BUTTON,
		},
		[]telebot.Btn{
			EDIT_CHANNEL_BUTTON,
			EDIT_TARGET_CHANNELS_BUTTON,
		},
	)
}
