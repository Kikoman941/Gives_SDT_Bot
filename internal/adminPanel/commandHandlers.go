package adminPanel

import (
	"Gives_SDT_Bot/internal/data"
	"errors"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func (ad *AdminPanel) InitCommandHandlers() {
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

			if err := ad.fsmService.SetState(userID, data.START_MENU_state, nil); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_state_message)
			}

			return reply
		},
	)

	// Комманда /RefreshAdmins. Обновляет список админов для мидлваря Whitelist
	ad.bot.Handle(
		data.COMMAND_REFRESH_ADMINS,
		func(ctx telebot.Context) error {
			err := ad.RefreshAdmins()
			if err != nil {
				if errors.Is(err, data.ERROR_NO_ADMINS_FOR_REFRESH) {
					return ctx.Reply(data.NO_ADMINS_message)
				}
				return ctx.Reply(data.CANNOT_GET_ADMINS_message)
			}

			ad.InitCommandHandlers()
			ad.InitButtonHandlers()
			ad.InitTextHandlers()
			ad.InitPhotoHandlers()

			return ctx.Reply(data.SUCCESS_REFRESH_ADMINS_message)
		},
		middleware.Whitelist(ad.adminGroup...),
	)
}
