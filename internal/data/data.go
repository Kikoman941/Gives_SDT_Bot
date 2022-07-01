package data

import (
	"gopkg.in/telebot.v3"
	"time"
)

func InitMenu() {
	START_MENU.Reply(
		[]telebot.Btn{
			CREATE_GIVE_BUTTON,
		},
		[]telebot.Btn{
			MY_GIVES_BUTTON,
		},
	)

	CANCEL_MENU.Reply(
		[]telebot.Btn{
			BACK_TO_START_BUTTON,
		},
	)

	ACTIVATE_GIVE_MENU.Reply(
		[]telebot.Btn{
			ACTIVATE_GIVE_BUTTON,
		},
		[]telebot.Btn{
			EDIT_GIVE_BUTTON,
		},
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
			EDIT_START_FINISH_BUTTON,
		},
		[]telebot.Btn{
			EDIT_IMAGE_BUTTON,
			EDIT_WINNERS_COUNT_BUTTON,
		},
		[]telebot.Btn{
			EDIT_CHANNEL_BUTTON,
			EDIT_TARGET_CHANNELS_BUTTON,
		},
		[]telebot.Btn{
			DELETE_GIVE_BUTTON,
			BACK_TO_START_BUTTON,
		},
	)

	ACTIVE_GIVE_MENU.Reply(
		[]telebot.Btn{
			DEACTIVATE_GIVE_BUTTON,
		},
		[]telebot.Btn{
			BACK_TO_START_BUTTON,
		},
	)
}

func LoadLocation(loc string) error {
	var err error
	LOCATION, err = time.LoadLocation(loc)
	if err != nil {
		return err
	}
	return nil
}
