package adminPanel

import (
	"Gives_SDT_Bot/internal/adminPanel/data"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func (ad *AdminPanel) InitHandlers() {
	// Комманда /start
	ad.bot.Handle(
		data.COMMAND_START,
		func(ctx telebot.Context) error {
			reply := ctx.Reply(data.START_MESSAGE, data.START_MENU)

			userID, err := ad.userService.AddUser(ctx.Chat().ID, false)
			if err != nil {
				return ctx.Reply(data.CANNOT_CREATE_USER_MESSAGE)
			} else if userID == 0 {
				ad.logger.Infof("User with tgId=%d, already exists", ctx.Chat().ID)
				return reply
			}

			if err := ad.fsmService.SetState(userID, data.MAIN_MENU_STATE); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_STATE_MESSAGE)
			}

			return reply
		},
	)

	// Комманда /refreshAdmins. Обновляет список админов для мидлваря Whitelist
	ad.bot.Handle(
		data.COMMAND_REFRESH_ADMINS,
		func(ctx telebot.Context) error {
			admins, err := ad.userService.GetAdmins()
			if err != nil || len(admins) == 0 {
				return ctx.Reply(data.NO_ADMINS_MESSAGE)
			}

			ad.refreshAdmins(admins)
			ad.InitHandlers()

			return ctx.Reply(data.SUCCESS_REFRESH_ADMINS_MESSAGE)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка "Назад в главное меню", отмена любого состояния до старта
	ad.bot.Handle(
		&data.MAIN_MENU_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_MESSAGE, data.START_MENU)
			}

			if err := ad.fsmService.SetState(userId, data.MAIN_MENU_STATE); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_STATE_MESSAGE, data.START_MENU)
			}

			return ctx.Reply(data.START_MESSAGE, data.START_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка "Новый конкурс 🎁", заускает цепочку создания конкурса
	ad.bot.Handle(
		&data.CREATE_GIVE_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_MESSAGE, data.START_MENU)
			}

			if err := ad.fsmService.SetState(userId, data.ENTER_GIVE_TITLE_STATE); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_STATE_MESSAGE, data.START_MENU)
			}

			return ctx.Reply(data.ENTER_GIVE_TITLE_MESSAGE, data.CANCEL_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Тригерится на любой текст
	ad.bot.Handle(
		telebot.OnText,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_MESSAGE, data.START_MENU)
			}

			userState, err := ad.fsmService.GetState(userId)
			if err != nil || userState == "" {
				return ctx.Reply(data.CANNOT_GET_USER_STATE_MESSAGE, data.START_MENU)
			}

			switch userState {
			case data.ENTER_GIVE_TITLE_STATE:
				giveTitle := ctx.Message().Text
				giveId, err := ad.giveService.CreateGive(giveTitle, userId)
				ad.logger.Info("туту")
				if err != nil || giveId == 0 {
					return ctx.Reply(data.CANNOT_CREATE_GIVE_MESSAGE, data.START_MENU)
				}
				return ctx.Reply("svfs")
			default:
				return ctx.Reply(data.I_DONT_UNDERSTAND_MESSAGE)
			}
		},
		middleware.Whitelist(ad.adminGroup...),
	)

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
