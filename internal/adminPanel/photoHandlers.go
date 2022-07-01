package adminPanel

import (
	data2 "Gives_SDT_Bot/internal/data"
	"fmt"
	"gopkg.in/telebot.v3"
	"strconv"
)

func (ad *AdminPanel) InitPhotoHandlers() {
	// Тригерится на любое фото, только фото, НЕ вложеный файл
	ad.bot.Handle(
		telebot.OnPhoto,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data2.CANNOT_FIND_USER_message, data2.CANCEL_MENU)
			}

			userState, err := ad.fsmService.GetState(userId)
			if err != nil || userState == nil {
				return ctx.Reply(data2.CANNOT_GET_USER_state_message, data2.CANCEL_MENU)
			}

			// Обслуживаем fsm
			// Загрузка обложки конкурса
			if userState.State == data2.UPLOAD_GIVE_IMAGE_state {
				img := ctx.Message().Photo.File
				filename, err := ad.imagesService.SaveFile(&img, userState.Data["giveId"])
				if err != nil || filename == "" {
					return ctx.Reply(data2.CANNOT_DOWNLOAD_IMAGE_message, data2.CANCEL_MENU)
				}

				giveId, err := strconv.Atoi(userState.Data["giveId"])
				err = ad.giveService.UpdateGive(giveId, fmt.Sprintf(`"image"='%s'`, filename))
				if err != nil {
					return ctx.Reply(data2.CANNOT_UPDATE_GIVE_message, data2.CANCEL_MENU)
				}

				workStatus := userState.Data["workStatus"]
				state := ""
				replyMessage := ""
				menu := &telebot.ReplyMarkup{}
				if workStatus == data2.WORK_STATUS_NEW {
					state = data2.ENTER_GIVE_START_FINISH_state
					replyMessage = data2.ENTER_GIVE_START_FINISH_message
					menu = data2.CANCEL_MENU
				} else if workStatus == data2.WORK_STATUS_EDIT {
					state = data2.SELECT_PROPERTY_TO_EDIT_state
					replyMessage = data2.SELECT_PROPERTY_TO_EDIT_message
					menu = data2.EDIT_GIVE_MENU
				}

				if err := ad.fsmService.SetState(userId, state, userState.Data); err != nil {
					return ctx.Reply(data2.CANNOT_SET_USER_state_message, data2.CANCEL_MENU)
				}

				return ctx.Reply(replyMessage, menu)
			}

			return ctx.Reply(data2.I_DONT_UNDERSTAND_message)
		},
	)
}
