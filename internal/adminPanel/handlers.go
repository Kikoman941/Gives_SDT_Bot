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
			reply := ctx.Reply(data.START_message, data.START_MENU)

			userID, err := ad.userService.AddUser(ctx.Chat().ID, false)
			if err != nil {
				return ctx.Reply(data.CANNOT_CREATE_USER_message)
			} else if userID == 0 {
				ad.logger.Infof("User with tgId=%d, already exists", ctx.Chat().ID)
				return reply
			}

			if err := ad.fsmService.Setstate(userID, data.START_MENU_state, nil); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_state_message)
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
				return ctx.Reply(data.NO_ADMINS_message)
			}

			ad.refreshAdmins(admins)
			ad.InitHandlers()

			return ctx.Reply(data.SUCCESS_REFRESH_ADMINS_message)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка "Назад в главное меню", отмена любого состояния до старта
	ad.bot.Handle(
		&data.BACK_TO_START_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_message, data.START_MENU)
			}

			if err := ad.fsmService.Setstate(userId, data.START_MENU_state, nil); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_state_message, data.START_MENU)
			}

			return ctx.Reply(data.START_message, data.START_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка "Новый конкурс 🎁", заускает цепочку создания конкурса
	ad.bot.Handle(
		&data.CREATE_GIVE_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_message, data.CANCEL_MENU)
			}

			if err := ad.fsmService.Setstate(userId, data.ENTER_TARGET_CHANNEL_state, nil); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
			}

			return ctx.Reply(data.ENTER_TARGET_CHANNEL_message, data.CANCEL_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)
	// Кнопка "Мои конкурсы 📋", выводит список конкурсов пользователя в виде кнопок
	ad.bot.Handle(
		&data.MY_GIVES_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_message, data.CANCEL_MENU)
			}

			if err := ad.fsmService.Setstate(userId, data.SELECT_OWN_GIVE_state, nil); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
			}

			userGives, err := ad.giveService.GetAllUserGives(userId)
			if err != nil {
				return ctx.Reply(data.CANNOT_GET_USER_GIVES_message, data.START_MENU)
			}
			buttons := GivesToButtons(userGives)
			buttons = append(buttons, data.BACK_TO_START_BUTTON)

			givesMenu := data.CreateReplyMenu(buttons...)

			return ctx.Reply(data.SELECT_OWN_GIVE_message, givesMenu)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка "Опубликовать ✅" активирует конкурс
	ad.bot.Handle(
		&data.ACTIVATE_GIVE_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_message, data.CANCEL_MENU)
			}

			userState, err := ad.fsmService.GetState(userId)
			if err != nil || userState == nil {
				return ctx.Reply(data.CANNOT_GET_USER_state_message, data.CANCEL_MENU)
			}

			giveId, err := strconv.Atoi(userState.Data["giveId"])
			if err != nil {
				return ctx.Reply(data.CANNOT_GET_state_DATA_message, data.CANCEL_MENU)
			}

			give, err := ad.giveService.GetGiveById(giveId)
			if err != nil {
				return ctx.Reply(data.CANNOT_GET_GIVE_message, data.CANCEL_MENU)
			}

			unfilledFields := ad.giveService.CheckFilling(&give)
			if len(unfilledFields) != 0 {
				return ctx.Reply(fmt.Sprintf(data.GIVE_FIELDS_MUST_BE_FILLED_message, unfilledFields))
			}

			err = ad.giveService.UpdateGive(giveId, `"isActive"=?`, true)
			if err != nil {
				return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
			}

			if err := ad.fsmService.Setstate(userId, data.START_MENU_state, nil); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
			}

			return ctx.Reply(data.GIVE_SUCCESSFULLY_ACTIVATE_message, data.START_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка "Редактировать 🅿️" меню редактирования конкурса
	ad.bot.Handle(
		&data.EDIT_GIVE_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_message, data.CANCEL_MENU)
			}

			userState, err := ad.fsmService.GetState(userId)
			if err != nil || userState == nil {
				return ctx.Reply(data.CANNOT_GET_USER_state_message, data.CANCEL_MENU)
			}

			giveId, err := strconv.Atoi(userState.Data["giveId"])
			if err != nil {
				return ctx.Reply(data.CANNOT_GET_state_DATA_message, data.CANCEL_MENU)
			}

			d := map[string]string{
				"giveId": strconv.Itoa(giveId),
			}
			if err := ad.fsmService.Setstate(userId, data.SELECT_PROPERTY_TO_EDIT_state, d); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
			}

			return ctx.Reply(data.SELECT_PROPERTY_TO_EDIT_message, data.EDIT_GIVE_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка "Снять с публикации ⛔", выключает активный конкурс
	ad.bot.Handle(
		&data.DEACTIVATE_GIVE_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_message, data.CANCEL_MENU)
			}

			userState, err := ad.fsmService.GetState(userId)
			if err != nil || userState == nil {
				return ctx.Reply(data.CANNOT_GET_USER_state_message, data.CANCEL_MENU)
			}

			giveId, err := strconv.Atoi(userState.Data["giveId"])
			if err != nil {
				return ctx.Reply(data.CANNOT_GET_state_DATA_message, data.CANCEL_MENU)
			}

			err = ad.giveService.UpdateGive(giveId, `"isActive"=?`, false)
			if err != nil {
				return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
			}

			if err := ad.fsmService.Setstate(userId, data.START_MENU_state, nil); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
			}

			return ctx.Reply(data.GIVE_SUCCESSFULL_DEACTIVATE_message, data.START_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// кнопка редактировани "Заголовок"
	ad.bot.Handle(
		&data.EDIT_TITLE_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_message, data.CANCEL_MENU)
			}

			userState, err := ad.fsmService.GetState(userId)
			if err != nil || userState == nil {
				return ctx.Reply(data.CANNOT_GET_USER_state_message, data.CANCEL_MENU)
			}

			if userState.State == data.SELECT_PROPERTY_TO_EDIT_state {
				giveId, err := strconv.Atoi(userState.Data["giveId"])
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_state_DATA_message, data.CANCEL_MENU)
				}

				d := map[string]string{
					"giveId":      strconv.Itoa(giveId),
					"work_status": data.WORK_STATUS_EDIT,
				}

				if err := ad.fsmService.Setstate(userId, data.ENTER_GIVE_TITLE_state, d); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(data.ENTER_GIVE_TITLE_message, data.CANCEL_MENU)
			}
			return ctx.Reply(data.I_DONT_UNDERSTAND_message, data.CANCEL_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Тригерится на любой текст
	ad.bot.Handle(
		telebot.OnText,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_message, data.CANCEL_MENU)
			}

			userState, err := ad.fsmService.GetState(userId)
			if err != nil || userState == nil {
				return ctx.Reply(data.CANNOT_GET_USER_state_message, data.CANCEL_MENU)
			}

			// Обслуживаем fsm
			switch userState.State {
			// Меню конкурса
			case data.SELECT_OWN_GIVE_state:
				giveTitle := ctx.Message().Text

				give, err := ad.giveService.GetGiveByTitle(giveTitle)
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_GIVE_message, data.CANCEL_MENU)
				}

				d := map[string]string{
					"giveId":     strconv.Itoa(give.Id),
					"workStatus": "edit",
				}
				if err := ad.fsmService.Setstate(userId, data.OWN_GIVE_MENU_state, d); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				isActive := "Не активный"
				if give.IsActive {
					isActive = "Активный"
				}

				img := &telebot.Photo{
					File: telebot.FromDisk(fmt.Sprintf("./.images/%s", give.Image)),
				}
				img.Caption = fmt.Sprintf(
					data.GIVE_CONTENT_message,
					give.Title,
					give.Description,
				)

				_, err = ad.bot.Send(
					ctx.Recipient(),
					fmt.Sprintf(
						data.GIVE_OUTPUT_message,
						give.Channel,
						give.TargetChannels,
						give.StartAt,
						give.FinishAt,
						isActive,
					),
				)
				if err != nil {
					return ctx.Reply(fmt.Sprintf(data.CANNOT_SEND_message, ctx.Recipient()))
				}

				if give.IsActive {
					return ctx.Reply(
						img,
						data.ACTIVE_GIVE_MENU,
					)
				} else {
					return ctx.Reply(
						img,
						data.ACTIVATE_GIVE_MENU,
					)
				}

			// Ввод канала на котором будет проходить конкурс
			case data.ENTER_TARGET_CHANNEL_state:
				channelStr := ctx.Message().Text
				channel, err := utils.StringToInt64(channelStr)
				if err != nil {
					ad.logger.Errorf("cannot parse chatId=%d string to int64: %s", channel, err)
					return ctx.Reply(fmt.Sprintf(data.CANNOT_PARSE_CHANNEL_message, channel), data.CANCEL_MENU)
				}

				isAdmin, err := ad.checkBotIsAdmin(channel)
				if err != nil {
					return ctx.Reply(fmt.Sprintf(data.CANNOT_CHECK_BOT_IS_ADMIN_message, channel), data.CANCEL_MENU)
				} else if isAdmin == false {
					return ctx.Reply(fmt.Sprintf(data.BOT_MUST_BE_ADMIN_message, channel), data.CANCEL_MENU)
				}

				giveId, err := ad.giveService.CreateGive(channelStr, userId)
				if err != nil || giveId == 0 {
					return ctx.Reply(data.CANNOT_CREATE_GIVE_message, data.CANCEL_MENU)
				}

				d := map[string]string{
					"giveId": strconv.Itoa(giveId),
				}
				if err := ad.fsmService.Setstate(userId, data.ENTER_GIVE_TITLE_state, d); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(data.ENTER_GIVE_TITLE_message)
			// Ввод заголовка конкурса
			case data.ENTER_GIVE_TITLE_state:
				giveTitle := ctx.Message().Text

				giveId, err := strconv.Atoi(userState.Data["giveId"])
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_state_DATA_message, data.CANCEL_MENU)
				}

				err = ad.giveService.UpdateGive(giveId, `"title"=?`, giveTitle)
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
				}

				if userState.Data["workStatus"] == data.WORK_STATUS_NEW {
					d := map[string]string{
						"giveId":     strconv.Itoa(giveId),
						"workStatus": data.WORK_STATUS_NEW,
					}
					if err := ad.fsmService.Setstate(userId, data.ENTER_GIVE_DESCRIPTION_state, d); err != nil {
						return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
					}

					return ctx.Reply(data.ENTER_GIVE_DESCRIPTION_message)
				} else {
					d := map[string]string{
						"giveId":     strconv.Itoa(giveId),
						"workStatus": "edit",
					}
					if err := ad.fsmService.Setstate(userId, data.SELECT_PROPERTY_TO_EDIT_state, d); err != nil {
						return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
					}

					return ctx.Reply(data.SELECT_PROPERTY_TO_EDIT_message, data.EDIT_GIVE_MENU)
				}
			// Ввод описания конкурса
			case data.ENTER_GIVE_DESCRIPTION_state:
				giveDesc := ctx.Message().Text
				giveId, err := strconv.Atoi(userState.Data["giveId"])
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_state_DATA_message, data.CANCEL_MENU)
				}

				err = ad.giveService.UpdateGive(giveId, `"description"=?`, giveDesc)
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
				}

				d := map[string]string{
					"giveId": userState.Data["giveId"],
				}
				if err := ad.fsmService.Setstate(userId, data.UPLOAD_GIVE_IMAGE_state, d); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(data.UPLOAD_GIVE_IMAGE_message)
			// Ввод дат старта - финиша конкурса
			case data.ENTER_GIVE_START_FINISH_state:
				duration := strings.Split(ctx.Message().Text, " - ")
				startAt, err := StringToTimeMSK(duration[0], ad.logger)
				if err != nil || startAt.IsZero() {
					return ctx.Reply(fmt.Sprintf(data.CANNOT_PARSE_TIME_message, duration[0]), data.CANCEL_MENU)
				}
				fmt.Println(startAt)
				finishAt, err := StringToTimeMSK(duration[1], ad.logger)
				if err != nil || finishAt.IsZero() {
					return ctx.Reply(fmt.Sprintf(data.CANNOT_PARSE_TIME_message, duration[1]), data.CANCEL_MENU)
				}

				if finishAt.Before(time.Now()) {
					return ctx.Reply(data.FINISH_DATE_HAS_PASSED_message)
				}
				if finishAt.Before(startAt) {
					return ctx.Reply(data.FINISH_DATE_BEFORE_START_message)
				}

				giveId, err := strconv.Atoi(userState.Data["giveId"])
				err = ad.giveService.UpdateGive(
					giveId,
					`"startAt"=?`,
					startAt.Format(time.RFC3339),
				)
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
				}
				err = ad.giveService.UpdateGive(
					giveId,
					`"finishAt"=?`,
					finishAt.Format(time.RFC3339),
				)
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
				}

				d := map[string]string{
					"giveId": userState.Data["giveId"],
				}
				if err := ad.fsmService.Setstate(userId, data.ENTER_WINNERS_COUNT_state, d); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(data.ENTER_WINNERS_COUNT_message)
			// Ввод колличества победителей конкурса
			case data.ENTER_WINNERS_COUNT_state:
				winnersCount, err := strconv.Atoi(ctx.Message().Text)
				if err != nil || winnersCount <= 0 {
					return ctx.Reply(data.CANNOT_PARSE_WINNERS_COUNT_message, data.CANCEL_MENU)
				}

				giveId, err := strconv.Atoi(userState.Data["giveId"])
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_state_DATA_message, data.CANCEL_MENU)
				}
				err = ad.giveService.UpdateGive(giveId, `"winnersCount"=?`, winnersCount)
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
				}

				d := map[string]string{
					"giveId": userState.Data["giveId"],
				}
				if err := ad.fsmService.Setstate(userId, data.ENTER_SUBSCRIPTION_CHANNELS_state, d); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(data.ENTER_SUBSCRIPTION_CHANNELS_message)
			// Ввод каналов для проверки подписки
			case data.ENTER_SUBSCRIPTION_CHANNELS_state:
				channels := strings.Split(ctx.Message().Text, " ")
				for _, ch := range channels {
					channelId, err := utils.StringToInt64(ch)
					if err != nil {
						ad.logger.Errorf("cannot parse chatId=%d string to int64: %s", channelId, err)
						return ctx.Reply(fmt.Sprintf(data.CANNOT_PARSE_CHANNEL_message, channelId), data.CANCEL_MENU)
					}

					isAdmin, err := ad.checkBotIsAdmin(channelId)
					if err != nil {
						return ctx.Reply(fmt.Sprintf(data.CANNOT_CHECK_BOT_IS_ADMIN_message, channelId), data.CANCEL_MENU)
					} else if isAdmin == false {
						return ctx.Reply(fmt.Sprintf(data.BOT_MUST_BE_ADMIN_message, channelId), data.CANCEL_MENU)
					}
				}

				giveId, err := strconv.Atoi(userState.Data["giveId"])
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_state_DATA_message, data.CANCEL_MENU)
				}
				err = ad.giveService.UpdateGive(giveId, `"targetChannels"=?`, pg.Array(channels))
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
				}

				give, err := ad.giveService.GetGiveById(giveId)
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_GIVE_message, data.CANCEL_MENU)
				}

				isActive := "Не активный"
				if give.IsActive {
					isActive = "Активный"
				}

				img := &telebot.Photo{
					File: telebot.FromDisk(fmt.Sprintf("./.images/%s", give.Image)),
				}
				img.Caption = fmt.Sprintf(
					data.GIVE_CONTENT_message,
					give.Title,
					give.Description,
				)

				d := map[string]string{
					"giveId": userState.Data["giveId"],
				}
				if err := ad.fsmService.Setstate(userId, data.EDIT_GIVE_state, d); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				_, err = ad.bot.Send(
					ctx.Recipient(),
					fmt.Sprintf(
						data.GIVE_OUTPUT_message,
						give.Channel,
						give.TargetChannels,
						give.StartAt,
						give.FinishAt,
						isActive,
					),
				)
				if err != nil {
					return ctx.Reply(fmt.Sprintf(data.CANNOT_SEND_message, ctx.Recipient()))
				}

				return ctx.Reply(
					img,
					data.ACTIVATE_GIVE_MENU,
				)
			default:
				return ctx.Reply(data.I_DONT_UNDERSTAND_message)
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
				return ctx.Reply(data.CANNOT_FIND_USER_message, data.CANCEL_MENU)
			}

			userState, err := ad.fsmService.GetState(userId)
			if err != nil || userState == nil {
				return ctx.Reply(data.CANNOT_GET_USER_state_message, data.CANCEL_MENU)
			}

			// Обслуживаем fsm
			// Загрузка обложки конкурса
			if userState.State == data.UPLOAD_GIVE_IMAGE_state {
				img := ctx.Message().Photo.File
				filename, err := ad.imagesService.SaveFile(&img, userState.Data["giveId"])
				if err != nil || filename == "" {
					return ctx.Reply(data.CANNOT_DOWNLOAD_IMAGE_message, data.CANCEL_MENU)
				}

				giveId, err := strconv.Atoi(userState.Data["giveId"])
				err = ad.giveService.UpdateGive(giveId, fmt.Sprintf(`"image"='%s'`, filename))
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
				}

				d := map[string]string{
					"giveId": userState.Data["giveId"],
				}
				if err := ad.fsmService.Setstate(userId, data.ENTER_GIVE_START_FINISH_state, d); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(data.ENTER_GIVE_START_FINISH_message)
			}

			return ctx.Reply(data.I_DONT_UNDERSTAND_message)
		},
	)
}
