package adminPanel

import (
	"Gives_SDT_Bot/internal/adminPanel/data"
	"Gives_SDT_Bot/pkg/utils"
	"fmt"
	"github.com/go-pg/pg/v10"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"strconv"
	"strings"
	"time"
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

			if err := ad.fsmService.SetState(userID, data.START_MENU_STATE, nil); err != nil {
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

			if err := ad.fsmService.SetState(userId, data.START_MENU_STATE, nil); err != nil {
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

			if err := ad.fsmService.SetState(userId, data.ENTER_TARGET_CHANNEL_STATE, nil); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_STATE_MESSAGE, data.CANCEL_MENU)
			}

			return ctx.Reply(data.ENTER_TARGET_CHANNEL_MESSAGE, data.CANCEL_MENU)
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
			buttons := GivesToButtons(userGives)
			buttons = append(buttons, data.BACK_TO_START_BUTTON)

			givesMenu := data.CreateReplyMenu(buttons...)

			return ctx.Reply(data.SELECT_OWN_GIVE_MESSAGE, givesMenu)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка "Опубликовать ✅" активирует конкурс
	ad.bot.Handle(
		&data.ACTIVATE_GIVE_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_MESSAGE, data.CANCEL_MENU)
			}

			userState, err := ad.fsmService.GetState(userId)
			if err != nil || userState == nil {
				return ctx.Reply(data.CANNOT_GET_USER_STATE_MESSAGE, data.CANCEL_MENU)
			}

			giveId, err := strconv.Atoi(userState.Data["giveId"])
			if err != nil {
				return ctx.Reply(data.CANNOT_GET_STATE_DATA_MESSAGE, data.CANCEL_MENU)
			}

			give, err := ad.giveService.GetGiveById(giveId)
			if err != nil {
				return ctx.Reply(data.CANNOT_GET_GIVE_MESSAGE, data.CANCEL_MENU)
			}

			unfilledFields := ad.giveService.CheckFilling(&give)
			if len(unfilledFields) != 0 {
				return ctx.Reply(fmt.Sprintf(data.GIVE_FIELDS_MUST_BE_FILLED_MESSAGE, unfilledFields))
			}

			err = ad.giveService.UpdateGive(giveId, `"isActive"=?`, true)
			if err != nil {
				return ctx.Reply(data.CANNOT_UPDATE_GIVE_MESSAGE, data.CANCEL_MENU)
			}

			if err := ad.fsmService.SetState(userId, data.START_MENU_STATE, nil); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_STATE_MESSAGE, data.CANCEL_MENU)
			}

			return ctx.Reply(data.GIVE_SUCCESSFULLY_ACTIVATE_MESSAGE, data.START_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка "Редактировать 🅿️" меню редактирования конкурса
	ad.bot.Handle(
		&data.EDIT_GIVE_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_MESSAGE, data.CANCEL_MENU)
			}

			userState, err := ad.fsmService.GetState(userId)
			if err != nil || userState == nil {
				return ctx.Reply(data.CANNOT_GET_USER_STATE_MESSAGE, data.CANCEL_MENU)
			}

			giveId, err := strconv.Atoi(userState.Data["giveId"])
			if err != nil {
				return ctx.Reply(data.CANNOT_GET_STATE_DATA_MESSAGE, data.CANCEL_MENU)
			}

			d := map[string]string{
				"giveId": strconv.Itoa(giveId),
			}
			if err := ad.fsmService.SetState(userId, data.SELECT_PROPERTY_TO_EDIT_STATE, d); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_STATE_MESSAGE, data.CANCEL_MENU)
			}

			return ctx.Reply(data.SELECT_PROPERTY_TO_EDIT_MESSAGE, data.EDIT_GIVE_MENU)
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

			// Обслуживаем fsm
			switch userState.State {
			// Меню конкурса
			case data.SELECT_OWN_GIVE_STATE:
				giveTitle := ctx.Message().Text

				give, err := ad.giveService.GetGiveByTitle(giveTitle)
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_GIVE_MESSAGE, data.CANCEL_MENU)
				}

				if give.IsActive {
					return ctx.Reply("Снять с публикации")
				} else {
					return ctx.Reply("Полное редактированиее")
				}

			// Ввод канала на котором будет проходить конкурс
			case data.ENTER_TARGET_CHANNEL_STATE:
				channelStr := ctx.Message().Text
				channel, err := utils.StringToInt64(channelStr)
				if err != nil {
					ad.logger.Errorf("cannot parse chatId=%d string to int64: %s", channel, err)
					return ctx.Reply(fmt.Sprintf(data.CANNOT_PARSE_CHANNEL_MESSAGE, channel), data.CANCEL_MENU)
				}

				isAdmin, err := ad.checkBotIsAdmin(channel)
				if err != nil {
					return ctx.Reply(fmt.Sprintf(data.CANNOT_CHECK_BOT_IS_ADMIN_MESSAGE, channel), data.CANCEL_MENU)
				} else if isAdmin == false {
					return ctx.Reply(fmt.Sprintf(data.BOT_MUST_BE_ADMIN_MESSAGE, channel), data.CANCEL_MENU)
				}

				giveId, err := ad.giveService.CreateGive(channelStr, userId)
				if err != nil || giveId == 0 {
					return ctx.Reply(data.CANNOT_CREATE_GIVE_MESSAGE, data.CANCEL_MENU)
				}

				d := map[string]string{
					"giveId": strconv.Itoa(giveId),
				}
				if err := ad.fsmService.SetState(userId, data.ENTER_GIVE_TITLE_STATE, d); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_STATE_MESSAGE, data.CANCEL_MENU)
				}

				return ctx.Reply(data.ENTER_GIVE_TITLE_MESSAGE)
			// Ввод заголовка конкурса
			case data.ENTER_GIVE_TITLE_STATE:
				giveTitle := ctx.Message().Text

				giveId, err := strconv.Atoi(userState.Data["giveId"])
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_STATE_DATA_MESSAGE, data.CANCEL_MENU)
				}

				err = ad.giveService.UpdateGive(giveId, `"title"=?`, giveTitle)
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_MESSAGE, data.CANCEL_MENU)
				}

				d := map[string]string{
					"giveId": strconv.Itoa(giveId),
				}
				if err := ad.fsmService.SetState(userId, data.ENTER_GIVE_DESCRIPTION_STATE, d); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_STATE_MESSAGE, data.CANCEL_MENU)
				}

				return ctx.Reply(data.ENTER_GIVE_DESCRIPTION_MESSAGE)
			// Ввод описания конкурса
			case data.ENTER_GIVE_DESCRIPTION_STATE:
				giveDesc := ctx.Message().Text
				giveId, err := strconv.Atoi(userState.Data["giveId"])
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_STATE_DATA_MESSAGE, data.CANCEL_MENU)
				}

				err = ad.giveService.UpdateGive(giveId, `"description"=?`, giveDesc)
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_MESSAGE, data.CANCEL_MENU)
				}

				d := map[string]string{
					"giveId": userState.Data["giveId"],
				}
				if err := ad.fsmService.SetState(userId, data.UPLOAD_GIVE_IMAGE_STATE, d); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_STATE_MESSAGE, data.CANCEL_MENU)
				}

				return ctx.Reply(data.UPLOAD_GIVE_IMAGE_MESSAGE)
			// Ввод дат старта - финиша конкурса
			case data.ENTER_GIVE_START_FINISH_STATE:
				duration := strings.Split(ctx.Message().Text, " - ")
				startAt, err := StringToTimeMSK(duration[0], ad.logger)
				if err != nil || startAt.IsZero() {
					return ctx.Reply(fmt.Sprintf(data.CANNOT_PARSE_TIME_MESSAGE, duration[0]), data.CANCEL_MENU)
				}
				fmt.Println(startAt)
				finishAt, err := StringToTimeMSK(duration[1], ad.logger)
				if err != nil || finishAt.IsZero() {
					return ctx.Reply(fmt.Sprintf(data.CANNOT_PARSE_TIME_MESSAGE, duration[1]), data.CANCEL_MENU)
				}

				if finishAt.Before(time.Now()) {
					return ctx.Reply(data.FINISH_DATE_HAS_PASSED_MESSAGE)
				}
				if finishAt.Before(startAt) {
					return ctx.Reply(data.FINISH_DATE_BEFORE_START_MESSAGE)
				}

				giveId, err := strconv.Atoi(userState.Data["giveId"])
				err = ad.giveService.UpdateGive(
					giveId,
					`"startAt"=?`,
					startAt.Format(time.RFC3339),
				)
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_MESSAGE, data.CANCEL_MENU)
				}
				err = ad.giveService.UpdateGive(
					giveId,
					`"finishAt"=?`,
					finishAt.Format(time.RFC3339),
				)
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_MESSAGE, data.CANCEL_MENU)
				}

				d := map[string]string{
					"giveId": userState.Data["giveId"],
				}
				if err := ad.fsmService.SetState(userId, data.ENTER_WINNERS_COUNT_STATE, d); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_STATE_MESSAGE, data.CANCEL_MENU)
				}

				return ctx.Reply(data.ENTER_WINNERS_COUNT_MESSAGE)
			// Ввод колличества победителей конкурса
			case data.ENTER_WINNERS_COUNT_STATE:
				winnersCount, err := strconv.Atoi(ctx.Message().Text)
				if err != nil || winnersCount <= 0 {
					return ctx.Reply(data.CANNOT_PARSE_WINNERS_COUNT_MESSAGE, data.CANCEL_MENU)
				}

				giveId, err := strconv.Atoi(userState.Data["giveId"])
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_STATE_DATA_MESSAGE, data.CANCEL_MENU)
				}
				err = ad.giveService.UpdateGive(giveId, `"winnersCount"=?`, winnersCount)
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_MESSAGE, data.CANCEL_MENU)
				}

				d := map[string]string{
					"giveId": userState.Data["giveId"],
				}
				if err := ad.fsmService.SetState(userId, data.ENTER_SUBSCRIPTION_CHANNELS_STATE, d); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_STATE_MESSAGE, data.CANCEL_MENU)
				}

				return ctx.Reply(data.ENTER_SUBSCRIPTION_CHANNELS_MESSAGE)
			// Ввод каналов для проверки подписки
			case data.ENTER_SUBSCRIPTION_CHANNELS_STATE:
				channels := strings.Split(ctx.Message().Text, " ")
				for _, ch := range channels {
					channelId, err := utils.StringToInt64(ch)
					if err != nil {
						ad.logger.Errorf("cannot parse chatId=%d string to int64: %s", channelId, err)
						return ctx.Reply(fmt.Sprintf(data.CANNOT_PARSE_CHANNEL_MESSAGE, channelId), data.CANCEL_MENU)
					}

					isAdmin, err := ad.checkBotIsAdmin(channelId)
					if err != nil {
						return ctx.Reply(fmt.Sprintf(data.CANNOT_CHECK_BOT_IS_ADMIN_MESSAGE, channelId), data.CANCEL_MENU)
					} else if isAdmin == false {
						return ctx.Reply(fmt.Sprintf(data.BOT_MUST_BE_ADMIN_MESSAGE, channelId), data.CANCEL_MENU)
					}
				}

				giveId, err := strconv.Atoi(userState.Data["giveId"])
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_STATE_DATA_MESSAGE, data.CANCEL_MENU)
				}
				err = ad.giveService.UpdateGive(giveId, `"targetChannels"=?`, pg.Array(channels))
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_MESSAGE, data.CANCEL_MENU)
				}

				give, err := ad.giveService.GetGiveById(giveId)
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_GIVE_MESSAGE, data.CANCEL_MENU)
				}

				isActive := "Не активный"
				if give.IsActive {
					isActive = "Активный"
				}

				img := &telebot.Photo{
					File: telebot.FromDisk(fmt.Sprintf("./.images/%s", give.Image)),
				}
				img.Caption = fmt.Sprintf(
					data.GIVE_CONTENT_MESSAGE,
					give.Title,
					give.Description,
				)

				_ = data.CreateReplyMenu(
					data.ACTIVATE_GIVE_BUTTON,
					data.EDIT_GIVE_BUTTON,
				)

				d := map[string]string{
					"giveId": userState.Data["giveId"],
				}
				if err := ad.fsmService.SetState(userId, data.CHECK_AND_ACTIVATE_GIVE_STATE, d); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_STATE_MESSAGE, data.CANCEL_MENU)
				}

				_, err = ad.bot.Send(
					ctx.Recipient(),
					fmt.Sprintf(
						data.CHECK_AND_ACTIVATE_GIVE_MESSAGE,
						give.Channel,
						give.TargetChannels,
						give.StartAt,
						give.FinishAt,
						isActive,
					),
				)
				if err != nil {
					return ctx.Reply(fmt.Sprintf(data.CANNOT_SEND_MESSAGE, ctx.Recipient()))
				}

				return ctx.Reply(
					img,
				)
			default:
				return ctx.Reply(data.I_DONT_UNDERSTAND_MESSAGE)
			}
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Тригерится на любое фото, только фото, НЕ вложеный файл
	ad.bot.Handle(
		telebot.OnPhoto,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_MESSAGE, data.CANCEL_MENU)
			}

			userState, err := ad.fsmService.GetState(userId)
			if err != nil || userState == nil {
				return ctx.Reply(data.CANNOT_GET_USER_STATE_MESSAGE, data.CANCEL_MENU)
			}

			// Обслуживаем fsm
			// Загрузка обложки конкурса
			if userState.State == data.UPLOAD_GIVE_IMAGE_STATE {
				img := ctx.Message().Photo.File
				filename, err := ad.imagesService.SaveFile(&img, userState.Data["giveId"])
				if err != nil || filename == "" {
					return ctx.Reply(data.CANNOT_DOWNLOAD_IMAGE_MESSAGE, data.CANCEL_MENU)
				}

				giveId, err := strconv.Atoi(userState.Data["giveId"])
				err = ad.giveService.UpdateGive(giveId, fmt.Sprintf(`"image"='%s'`, filename))
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_MESSAGE, data.CANCEL_MENU)
				}

				d := map[string]string{
					"giveId": userState.Data["giveId"],
				}
				if err := ad.fsmService.SetState(userId, data.ENTER_GIVE_START_FINISH_STATE, d); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_STATE_MESSAGE, data.CANCEL_MENU)
				}

				return ctx.Reply(data.ENTER_GIVE_START_FINISH_MESSAGE)
			}

			return ctx.Reply(data.I_DONT_UNDERSTAND_MESSAGE)
		},
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
