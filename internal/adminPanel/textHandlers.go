package adminPanel

import (
	"Gives_SDT_Bot/internal/data"
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
				selectedGiveId, err := strconv.Atoi(ctx.Callback().Data)
				if err != nil {
					return err
				}

				give, err := ad.giveService.GetGiveById(selectedGiveId)
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_GIVE_message, data.CANCEL_MENU)
				}

				d := map[string]string{
					"giveId":     strconv.Itoa(give.Id),
					"workStatus": data.WORK_STATUS_EDIT,
				}
				if err := ad.fsmService.SetState(userId, data.OWN_GIVE_MENU_state, d); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				isActive := "Не активный"
				if give.IsActive {
					isActive = "Активный"
				}

				msg := &telebot.Photo{
					File: telebot.FromDisk(fmt.Sprintf("./.images/%s", give.Image)),
				}

				msg.Caption = data.ClearTextForMarkdownV2(
					fmt.Sprintf(
						data.GIVE_CONTENT_message,
						give.Title,
						give.Description,
					),
				)

				text := data.ClearTextForMarkdownV2(
					fmt.Sprintf(
						data.GIVE_OUTPUT_message,
						give.Channel,
						give.TargetChannels,
						give.WinnersCount,
						give.StartAt.In(ad.location).Format(time.RFC822),
						give.FinishAt.In(ad.location).Format(time.RFC822),
						isActive,
					),
				)
				_, err = ad.bot.Send(
					ctx.Recipient(),
					text,
					telebot.ModeMarkdownV2,
				)
				if err != nil {
					return ctx.Reply(fmt.Sprintf(data.CANNOT_SEND_message, err))
				}

				if give.IsActive {
					return ctx.Reply(
						msg,
						data.ACTIVE_GIVE_MENU,
						telebot.ModeMarkdownV2,
					)
				} else {
					return ctx.Reply(
						msg,
						data.ACTIVATE_GIVE_MENU,
						telebot.ModeMarkdownV2,
					)
				}
			// Ввод заголовка конкурса
			case data.ENTER_GIVE_TITLE_state:
				giveTitle := ctx.Message().Text

				workStatus := userState.Data["workStatus"]
				state := ""
				replyMessage := ""
				menu := &telebot.ReplyMarkup{}
				if workStatus == data.WORK_STATUS_NEW {
					giveId, err := ad.giveService.CreateGive(giveTitle, userId)
					if err != nil || giveId == 0 {
						return ctx.Reply(data.CANNOT_CREATE_GIVE_message, data.CANCEL_MENU)
					}

					userState.Data["giveId"] = strconv.Itoa(giveId)

					state = data.ENTER_GIVE_DESCRIPTION_state
					replyMessage = data.ENTER_GIVE_DESCRIPTION_message
					menu = data.CANCEL_MENU
				} else if workStatus == data.WORK_STATUS_EDIT {
					giveId, err := strconv.Atoi(userState.Data["giveId"])
					if err != nil {
						return ctx.Reply(data.CANNOT_GET_STATE_DATA_message, data.CANCEL_MENU)
					}

					err = ad.giveService.UpdateGive(giveId, `"title"=?`, giveTitle)
					if err != nil {
						return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
					}

					state = data.SELECT_PROPERTY_TO_EDIT_state
					replyMessage = data.SELECT_PROPERTY_TO_EDIT_message
					menu = data.EDIT_GIVE_MENU
				}

				if err := ad.fsmService.SetState(userId, state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				if err := ctx.Reply(data.GIVE_TITLE_OK_message); err != nil {
					return err
				}
				return ctx.Send(replyMessage, menu)
			// Ввод описания конкурса
			case data.ENTER_GIVE_DESCRIPTION_state:
				giveDesc := ctx.Message().Text
				giveId, err := strconv.Atoi(userState.Data["giveId"])
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_STATE_DATA_message, data.CANCEL_MENU)
				}

				if err = ad.giveService.UpdateGive(giveId, `"description"=?`, giveDesc); err != nil {
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

				if err := ad.fsmService.SetState(userId, state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				if err := ctx.Reply(data.GIVE_DESCRIPTION_OK_message); err != nil {
					return err
				}
				return ctx.Send(replyMessage, menu)
			// Ввод дат старта - финиша конкурса
			case data.ENTER_GIVE_START_FINISH_state:
				duration := strings.Split(ctx.Message().Text, " - ")
				startAt, err := data.StringToTimeLocation(duration[0], ad.logger, ad.location)
				if err != nil || startAt.IsZero() {
					return ctx.Reply(fmt.Sprintf(data.CANNOT_PARSE_TIME_message, duration[0]), data.CANCEL_MENU)
				}

				finishAt, err := data.StringToTimeLocation(duration[1], ad.logger, ad.location)
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
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_STATE_DATA_message, data.CANCEL_MENU)
				}
				if err = ad.giveService.UpdateGive(giveId, `"startAt"=?`, startAt.Format(time.RFC3339)); err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
				}

				if err = ad.giveService.UpdateGive(giveId, `"finishAt"=?`, finishAt.Format(time.RFC3339)); err != nil {
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

				if err := ad.fsmService.SetState(userId, state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				if err := ctx.Reply(data.GIVE_START_FINISH_OK_message); err != nil {
					return err
				}
				return ctx.Send(replyMessage, menu)
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

				if err = ad.giveService.UpdateGive(giveId, `"winnersCount"=?`, winnersCount); err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
				}

				workStatus := userState.Data["workStatus"]
				state := ""
				replyMessage := ""
				menu := &telebot.ReplyMarkup{}
				if workStatus == data.WORK_STATUS_NEW {
					state = data.ENTER_TARGET_CHANNEL_state
					replyMessage = data.ENTER_TARGET_CHANNEL_message
					menu = data.CANCEL_MENU
				} else if workStatus == data.WORK_STATUS_EDIT {
					state = data.SELECT_PROPERTY_TO_EDIT_state
					replyMessage = data.SELECT_PROPERTY_TO_EDIT_message
					menu = data.EDIT_GIVE_MENU
				}

				if err := ad.fsmService.SetState(userId, state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				if err := ctx.Reply(data.WINNERS_COUNT_OK_message); err != nil {
					return err
				}
				return ctx.Send(replyMessage, menu)
			// Ввод канала на котором будет проходить конкурс
			case data.ENTER_TARGET_CHANNEL_state:
				channelStr := ctx.Message().Text
				channel, err := utils.StringToInt64(ctx.Message().Text)
				if err != nil {
					ad.logger.Errorf("cannot parse chatId=%d string to int64: %s", channel, err)
					return ctx.Reply(fmt.Sprintf(data.CANNOT_PARSE_CHANNEL_message, channel), data.CANCEL_MENU)
				}

				isAdmin, err := ad.checkBotIsAdmin(channel)
				if err != nil {
					return ctx.Reply(fmt.Sprintf(data.CANNOT_CHECK_BOT_IS_ADMIN_message, channel), data.CANCEL_MENU)
				} else if !isAdmin {
					return ctx.Reply(fmt.Sprintf(data.BOT_MUST_BE_ADMIN_message, channel), data.CANCEL_MENU)
				}

				giveId, err := strconv.Atoi(userState.Data["giveId"])
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_STATE_DATA_message, data.CANCEL_MENU)
				}

				if err = ad.giveService.UpdateGive(giveId, `"channel"=?`, channelStr); err != nil {
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

				if err := ad.fsmService.SetState(userId, state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				if err := ctx.Reply(data.TARGET_CHANNEL_OK_message); err != nil {
					return nil
				}
				return ctx.Send(replyMessage, menu)
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
					} else if !isAdmin {
						return ctx.Reply(fmt.Sprintf(data.BOT_MUST_BE_ADMIN_message, channelId), data.CANCEL_MENU)
					}
				}

				giveId, err := strconv.Atoi(userState.Data["giveId"])
				if err != nil {
					return ctx.Reply(data.CANNOT_GET_STATE_DATA_message, data.CANCEL_MENU)
				}

				if err = ad.giveService.UpdateGive(giveId, `"targetChannels"=?`, pg.Array(channels)); err != nil {
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

					msg := &telebot.Photo{
						File: telebot.FromDisk(fmt.Sprintf("./.images/%s", give.Image)),
					}
					msg.Caption = data.ClearTextForMarkdownV2(
						fmt.Sprintf(
							data.GIVE_CONTENT_message,
							give.Title,
							give.Description,
						),
					)

					d := map[string]string{
						"giveId":     userState.Data["giveId"],
						"workStatus": data.WORK_STATUS_EDIT,
					}
					if err := ad.fsmService.SetState(userId, data.EDIT_GIVE_state, d); err != nil {
						return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
					}

					giveChannelId, err := utils.StringToInt64(give.Channel)
					if err != nil {
						return ctx.Reply(data.CANNOT_PARSE_CHANNEL_message, give.Channel, err)
					}
					giveChannelName := data.GetChatNameByChatID(ad.bot, giveChannelId, ad.logger)
					if giveChannelName == "" {
						return ctx.Reply(fmt.Sprintf(data.CANNOT_GET_CHAT_NAME_message, giveChannelId), data.CANCEL_MENU)
					}
					var targetChannelsNames []string
					for _, targetChannel := range give.TargetChannels {
						targetChannelId, err := utils.StringToInt64(targetChannel)
						if err != nil {
							return ctx.Reply(data.CANNOT_PARSE_CHANNEL_message, give.Channel, err)
						}
						targetChannelName := data.GetChatNameByChatID(ad.bot, targetChannelId, ad.logger)
						if targetChannelName == "" {
							return ctx.Reply(fmt.Sprintf(data.CANNOT_GET_CHAT_NAME_message, targetChannelId), data.CANCEL_MENU)
						}
						targetChannelsNames = append(targetChannelsNames, fmt.Sprintf("@%s", targetChannelName))
					}

					text := data.ClearTextForMarkdownV2(
						fmt.Sprintf(
							data.GIVE_OUTPUT_message,
							fmt.Sprintf("@%s", giveChannelName),
							targetChannelsNames,
							give.WinnersCount,
							give.StartAt.In(ad.location).Format(time.RFC822),
							give.FinishAt.In(ad.location).Format(time.RFC822),
							isActive,
						),
					)
					_, err = ad.bot.Send(
						ctx.Recipient(),
						text,
						telebot.ModeMarkdownV2,
					)
					if err != nil {
						return ctx.Reply(fmt.Sprintf(data.CANNOT_SEND_message, err))
					}

					return ctx.Reply(
						msg,
						&telebot.SendOptions{
							ReplyMarkup: data.ACTIVATE_GIVE_MENU,
							ParseMode:   telebot.ModeMarkdownV2,
						},
					)
				} else if workStatus == data.WORK_STATUS_EDIT {
					state = data.SELECT_PROPERTY_TO_EDIT_state
					replyMessage = data.SELECT_PROPERTY_TO_EDIT_message
					menu = data.EDIT_GIVE_MENU
				}

				if err := ad.fsmService.SetState(userId, state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				if err := ctx.Reply(data.SUBSCRIPTION_CHANNELS_OK_message); err != nil {
					return err
				}
				return ctx.Send(replyMessage, menu)
			default:
				return ctx.Reply(data.I_DONT_UNDERSTAND_message)
			}
		},
		middleware.Whitelist(ad.adminGroup...),
	)
}
