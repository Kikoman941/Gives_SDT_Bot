package adminPanel

import (
	"Gives_SDT_Bot/internal/adminPanel/data"
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
