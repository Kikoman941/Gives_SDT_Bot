package data

import "gopkg.in/telebot.v3"

var (
	CREATE_GIVE_BUTTON   = telebot.Btn{Text: "Новый конкурс 🎁"}
	MY_GIVES_BUTTON      = telebot.Btn{Text: "Мои конкурсы 📋"}
	BACK_TO_START_BUTTON = telebot.Btn{Text: "Назад в главное меню 🔙"}
	ACTIVATE_GIVE_BUTTON = telebot.Btn{Text: "Опубликовать ✅"}
	EDIT_GIVE_BUTTON     = telebot.Btn{Text: "Редактировать 🅿️"}
)
