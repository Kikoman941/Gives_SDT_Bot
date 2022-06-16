package data

import "gopkg.in/telebot.v3"

var (
	CREATE_GIVE_BUTTON = telebot.Btn{Text: "Новый конкурс 🎁", Data: "test_inline_callback_CREATE", InlineQuery: "query"}
	MY_GIVES_BUTTON    = telebot.Btn{Text: "Мои конкурсы 📋", Data: "test_inline_callback_GIVES", InlineQuery: "query"}
	MAIN_MENU_BUTTON   = telebot.Btn{Text: "Назад в главное меню 🔙"}
)
