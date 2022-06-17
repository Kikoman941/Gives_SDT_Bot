package adminPanel

import (
	"Gives_SDT_Bot/internal/adminPanel/data"
	"gopkg.in/telebot.v3"
)

func (ad *AdminPanel) InitHandlers() {
	// Комманда /start
	ad.bot.Handle(
		data.COMMAND_START,
		func(ctx telebot.Context) error {
			userID, err := ad.userService.AddUser(ctx.Chat().ID, false)
			if err != nil {
				return err
			}

			if err := ad.fsmService.SetState(userID, data.MAIN_MENU_STATE); err != nil {
				return err
			}

			return ctx.Reply(data.START_MESSAGE, data.START_MENU)
		},
	)

	// Комманда /refreshAdmins. Обновляет список админов для мидлваря Whitelist
	//ad.bot.Handle(
	//	data.COMMAND_REFRESH_ADMINS,
	//	func(ctx telebot.Context) error {
	//		admins, err := ad.getAdmins()
	//		if err != nil {
	//			return err
	//		}
	//
	//		ad.refreshAdmins(admins)
	//		ad.InitHandlers()
	//
	//		return ctx.Reply(ad.telegramData.Messages["successRefreshAdmins"])
	//	},
	//	middleware.Whitelist(ad.adminGroup...),
	//)
	//
	//// Кнопка "Назад в главное меню", отмена любого состояния до старта
	//ad.bot.Handle(
	//	&ad.telegramData.Buttons.MainMenuButton,
	//	func(ctx telebot.Context) error {
	//		if err := ad.fsm.setState(ctx.Chat().ID, telegram.MAIN_MENU); err != nil {
	//			return err
	//		}
	//
	//		return ctx.Reply(ad.telegramData.Messages["startMessage"], ad.telegramData.Menus.StartMenu)
	//	},
	//	middleware.Whitelist(ad.adminGroup...),
	//)
	//
	//ad.bot.Handle(
	//	&ad.telegramData.Buttons.CreateGiveButton,
	//	func(ctx telebot.Context) error {
	//		if err := ad.fsm.setState(ctx.Chat().ID, telegram.ADD_TARGET_CHANNELS); err != nil {
	//			return err
	//		}
	//
	//		return ctx.Reply(ad.telegramData.Messages["addTargetChannels"])
	//	},
	//	middleware.Whitelist(ad.adminGroup...),
	//)
	//
	//ad.bot.Handle(
	//	telebot.OnText,
	//	func(ctx telebot.Context) error {
	//		//_, err := ad.bot.ChatMemberOf(getChatIdFromInt(1001422240135), getChatIdFromInt(10784525812))
	//		//if errors.As(err, &ERROR_USER_NOT_FOUND) {
	//		//	return ctx.Reply("Пользователь не подписан")
	//		//} else {
	//		//	return err
	//		//}
	//		return ctx.Reply("DONE", ad.telegramData.Menus.CreateInlineMenu())
	//	},
	//	middleware.Whitelist(ad.adminGroup...),
	//)
	//
	//ad.bot.Handle(
	//	&ad.telegramData.Buttons.MyGivesButton,
	//	func(ctx telebot.Context) error {
	//		return ctx.Respond()
	//	},
	//	middleware.Whitelist(ad.adminGroup...),
	//)
	//
	//ad.bot.Handle(
	//	telebot.OnCallback,
	//	func(ctx telebot.Context) error {
	//		return ctx.Reply(ctx.Callback().Data)
	//	},
	//	middleware.Whitelist(ad.adminGroup...),
	//)
	//
	//ad.bot.Handle(
	//	data.COMMAND_TEST_DB,
	//	func(ctx telebot.Context) error {
	//		//var users []User
	//		//err := ad.storage.Select(&users, "users", "")
	//		//if err != nil {
	//		//	return err
	//		//}
	//		//fmt.Println(users[0].IsAdmin)
	//		userId, err := ad.addUser(ctx.Chat().ID, true)
	//		if err != nil {
	//			return err
	//		}
	//
	//		if err := ad.fsm.setState(userId, telegram.MAIN_MENU); err != nil {
	//			return err
	//		}
	//
	//		return ctx.Reply("You push button")
	//	},
	//	middleware.Whitelist(ad.adminGroup...),
	//)
}
