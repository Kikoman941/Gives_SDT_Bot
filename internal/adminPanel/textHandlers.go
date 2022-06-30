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

func (ad *AdminPanel) InitTextHandlers() {
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
					"workStatus": data.WORK_STATUS_EDIT,
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

				workStatus := userState.Data["workStatus"]
				state := ""
				replyMessage := ""
				menu := &telebot.ReplyMarkup{}
				if workStatus == data.WORK_STATUS_NEW {
					giveId, err := ad.giveService.CreateGive(channelStr, userId)
					if err != nil || giveId == 0 {
						return ctx.Reply(data.CANNOT_CREATE_GIVE_message, data.CANCEL_MENU)
					}

					userState.Data["giveId"] = strconv.Itoa(giveId)
					state = data.ENTER_GIVE_TITLE_state
					replyMessage = data.ENTER_GIVE_TITLE_message
					menu = data.CANCEL_MENU
				} else if workStatus == data.WORK_STATUS_EDIT {
					giveId, err := strconv.Atoi(userState.Data["giveId"])
					if err != nil {
						return ctx.Reply(data.CANNOT_GET_STATE_DATA_message, data.CANCEL_MENU)
					}
					err = ad.giveService.UpdateGive(giveId, `"channel"=?`, channelStr)
					if err != nil {
						return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
					}

					userState.Data["giveId"] = strconv.Itoa(giveId)
					state = data.SELECT_PROPERTY_TO_EDIT_state
					replyMessage = data.SELECT_PROPERTY_TO_EDIT_message
					menu = data.EDIT_GIVE_MENU
				}

				if err := ad.fsmService.Setstate(userId, state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(replyMessage, menu)
			// Ввод заголовка конкурса
			case data.ENTER_GIVE_TITLE_state:
				giveTitle := ctx.Message().Text

				giveId, err := strconv.Atoi(userState.Data["giveId"])
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_STATE_DATA_message, data.CANCEL_MENU)
				}

				err = ad.giveService.UpdateGive(giveId, `"title"=?`, giveTitle)
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
				}

				workStatus := userState.Data["workStatus"]
				state := ""
				replyMessage := ""
				menu := &telebot.ReplyMarkup{}
				if workStatus == data.WORK_STATUS_NEW {
					state = data.ENTER_GIVE_DESCRIPTION_state
					replyMessage = data.ENTER_GIVE_DESCRIPTION_message
					menu = data.CANCEL_MENU
				} else if workStatus == data.WORK_STATUS_EDIT {
					state = data.SELECT_PROPERTY_TO_EDIT_state
					replyMessage = data.SELECT_PROPERTY_TO_EDIT_message
					menu = data.EDIT_GIVE_MENU
				}

				if err := ad.fsmService.Setstate(userId, state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(replyMessage, menu)
			// Ввод описания конкурса
			case data.ENTER_GIVE_DESCRIPTION_state:
				giveDesc := ctx.Message().Text
				giveId, err := strconv.Atoi(userState.Data["giveId"])
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_STATE_DATA_message, data.CANCEL_MENU)
				}

				err = ad.giveService.UpdateGive(giveId, `"description"=?`, giveDesc)
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
				}

				workStatus := userState.Data["workStatus"]
				state := ""
				replyMessage := ""
				menu := &telebot.ReplyMarkup{}
				if workStatus == data.WORK_STATUS_NEW {
					state = data.UPLOAD_GIVE_IMAGE_state
					replyMessage = data.UPLOAD_GIVE_IMAGE_message
					menu = data.CANCEL_MENU
				} else if workStatus == data.WORK_STATUS_EDIT {
					state = data.SELECT_PROPERTY_TO_EDIT_state
					replyMessage = data.SELECT_PROPERTY_TO_EDIT_message
					menu = data.EDIT_GIVE_MENU
				}

				if err := ad.fsmService.Setstate(userId, state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(replyMessage, menu)
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

				workStatus := userState.Data["workStatus"]
				state := ""
				replyMessage := ""
				menu := &telebot.ReplyMarkup{}
				if workStatus == data.WORK_STATUS_NEW {
					state = data.ENTER_WINNERS_COUNT_state
					replyMessage = data.ENTER_WINNERS_COUNT_message
					menu = data.CANCEL_MENU
				} else if workStatus == data.WORK_STATUS_EDIT {
					state = data.SELECT_PROPERTY_TO_EDIT_state
					replyMessage = data.SELECT_PROPERTY_TO_EDIT_message
					menu = data.EDIT_GIVE_MENU
				}

				if err := ad.fsmService.Setstate(userId, state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(replyMessage, menu)
			// Ввод колличества победителей конкурса
			case data.ENTER_WINNERS_COUNT_state:
				winnersCount, err := strconv.Atoi(ctx.Message().Text)
				if err != nil || winnersCount <= 0 {
					return ctx.Reply(data.CANNOT_PARSE_WINNERS_COUNT_message, data.CANCEL_MENU)
				}

				giveId, err := strconv.Atoi(userState.Data["giveId"])
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_STATE_DATA_message, data.CANCEL_MENU)
				}
				err = ad.giveService.UpdateGive(giveId, `"winnersCount"=?`, winnersCount)
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
				}

				workStatus := userState.Data["workStatus"]
				state := ""
				replyMessage := ""
				menu := &telebot.ReplyMarkup{}
				if workStatus == data.WORK_STATUS_NEW {
					state = data.ENTER_SUBSCRIPTION_CHANNELS_state
					replyMessage = data.ENTER_SUBSCRIPTION_CHANNELS_message
					menu = data.CANCEL_MENU
				} else if workStatus == data.WORK_STATUS_EDIT {
					state = data.SELECT_PROPERTY_TO_EDIT_state
					replyMessage = data.SELECT_PROPERTY_TO_EDIT_message
					menu = data.EDIT_GIVE_MENU
				}

				if err := ad.fsmService.Setstate(userId, state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(replyMessage, menu)
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
					return ctx.Reply(data.CANNOT_GET_STATE_DATA_message, data.CANCEL_MENU)
				}
				err = ad.giveService.UpdateGive(giveId, `"targetChannels"=?`, pg.Array(channels))
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
				}

				workStatus := userState.Data["workStatus"]
				state := ""
				replyMessage := ""
				menu := &telebot.ReplyMarkup{}
				if workStatus == data.WORK_STATUS_NEW {
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
				} else if workStatus == data.WORK_STATUS_EDIT {
					state = data.SELECT_PROPERTY_TO_EDIT_state
					replyMessage = data.SELECT_PROPERTY_TO_EDIT_message
					menu = data.EDIT_GIVE_MENU
				}

				if err := ad.fsmService.Setstate(userId, state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(replyMessage, menu)
			default:
				return ctx.Reply(data.I_DONT_UNDERSTAND_message)
			}
		},
		middleware.Whitelist(ad.adminGroup...),
	)
}
