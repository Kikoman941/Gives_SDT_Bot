package data

import "gopkg.in/telebot.v3"

var (
	CREATE_GIVE_BUTTON          = telebot.Btn{Text: "Новый конкурс 🎁"}
	MY_GIVES_BUTTON             = telebot.Btn{Text: "Мои конкурсы 📋"}
	BACK_TO_START_BUTTON        = telebot.Btn{Text: "Назад в главное меню 🔙"}
	ACTIVATE_GIVE_BUTTON        = telebot.Btn{Text: "Опубликовать ✅"}
	DEACTIVATE_GIVE_BUTTON      = telebot.Btn{Text: "Снять с публикации ⛔"}
	EDIT_GIVE_BUTTON            = telebot.Btn{Text: "Редактировать 🅿️"}
	EDIT_TITLE_BUTTON           = telebot.Btn{Text: "Заголовок"}
	EDIT_DESCRIPTION_BUTTON     = telebot.Btn{Text: "Описание"}
	EDIT_START_FINISH_BUTTON    = telebot.Btn{Text: "Дата старта - финиша"}
	EDIT_IMAGE_BUTTON           = telebot.Btn{Text: "Картинка"}
	EDIT_WINNERS_COUNT_BUTTON   = telebot.Btn{Text: "Колличество победителей"}
	EDIT_CHANNEL_BUTTON         = telebot.Btn{Text: "Канал конкурса"}
	EDIT_TARGET_CHANNELS_BUTTON = telebot.Btn{Text: "Каналы проверки подписки"}
	DELETE_GIVE_BUTTON          = telebot.Btn{Text: "Удалить ❌"}
)
