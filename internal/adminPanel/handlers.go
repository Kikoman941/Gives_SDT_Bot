package adminPanel

import (
	"Gives_SDT_Bot/internal/adminPanel/data"
	"fmt"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"strconv"
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

			if err := ad.fsmService.SetState(userID, data.MAIN_MENU_STATE, nil); err != nil {
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
		&data.BACK_TO_START_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_MESSAGE, data.START_MENU)
			}

			if err := ad.fsmService.SetState(userId, data.MAIN_MENU_STATE, nil); err != nil {
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
				return ctx.Reply(data.CANNOT_FIND_USER_MESSAGE, data.CANCEL_MENU)
			}

			if err := ad.fsmService.SetState(userId, data.ENTER_GIVE_TITLE_STATE, nil); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_STATE_MESSAGE, data.CANCEL_MENU)
			}

			return ctx.Reply(data.ENTER_GIVE_TITLE_MESSAGE, data.CANCEL_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)
	// Кнопка "Мои конкурсы 📋", выводит список конкурсов пользователя в виде кнопок
	ad.bot.Handle(
		&data.MY_GIVES_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_MESSAGE, data.CANCEL_MENU)
			}

			if err := ad.fsmService.SetState(userId, data.SELECT_OWN_GIVE_STATE, nil); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_STATE_MESSAGE, data.CANCEL_MENU)
			}

			userGives, err := ad.giveService.GetAllUserGives(userId)
			if err != nil {
				return ctx.Reply(data.CANNOT_GET_USER_GIVES_MESSAGE, data.START_MENU)
			}

			givesMenu := data.CreateReplyMenu(
				GivesToButtons(userGives)...,
			)

			return ctx.Reply(data.SELECT_OWN_GIVE_MESSAGE, givesMenu)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Тригерится на любой текст
	ad.bot.Handle(
		telebot.OnText,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_MESSAGE, data.CANCEL_MENU)
			}

			userState, err := ad.fsmService.GetState(userId)
			if err != nil || userState == nil {
				return ctx.Reply(data.CANNOT_GET_USER_STATE_MESSAGE, data.CANCEL_MENU)
			}

			switch userState.State {
			case data.ENTER_GIVE_TITLE_STATE:
				giveTitle := ctx.Message().Text
				giveId, err := ad.giveService.CreateGive(giveTitle, userId)
				if err != nil || giveId == 0 {
					return ctx.Reply(data.CANNOT_CREATE_GIVE_MESSAGE, data.CANCEL_MENU)
				}

				d := map[string]string{
					"giveId": strconv.Itoa(giveId),
				}
				if err := ad.fsmService.SetState(userId, data.ENTER_GIVE_DESCRIPTION_STATE, d); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_STATE_MESSAGE, data.CANCEL_MENU)
				}

				return ctx.Reply(data.ENTER_GIVE_DESCRIPTION_MESSAGE)
			case data.ENTER_GIVE_DESCRIPTION_STATE:
				giveDesc := ctx.Message().Text
				giveId, err := strconv.Atoi(userState.Data["giveId"])
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_STATE_DATA_MESSAGE, data.CANCEL_MENU)
				}

				err = ad.giveService.UpdateGive(giveId, fmt.Sprintf("description='%s'", giveDesc))
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_MESSAGE, data.CANCEL_MENU)
				}

				if err := ad.fsmService.SetState(userId, data.UPLOAD_GIVE_IMAGE_STATE, nil); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_STATE_MESSAGE, data.CANCEL_MENU)
				}

				return ctx.Reply(data.UPLOAD_GIVE_IMAGE_MESSAGE)
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
