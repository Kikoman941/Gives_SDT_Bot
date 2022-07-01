package adminPanel

import (
	data2 "Gives_SDT_Bot/internal/data"
	"fmt"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"strconv"
)

func (ad *AdminPanel) InitButtonHandlers() {
	// Кнопка "Назад в главное меню", отмена любого состояния до старта
	ad.bot.Handle(
		&data2.BACK_TO_START_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data2.CANNOT_FIND_USER_message, data2.START_MENU)
			}

			if err := ad.fsmService.SetState(userId, data2.START_MENU_state, nil); err != nil {
				return ctx.Reply(data2.CANNOT_SET_USER_state_message, data2.START_MENU)
			}

			return ctx.Reply(data2.START_message, data2.START_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка "Новый конкурс 🎁", заускает цепочку создания конкурса
	ad.bot.Handle(
		&data2.CREATE_GIVE_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data2.CANNOT_FIND_USER_message, data2.CANCEL_MENU)
			}

			d := map[string]string{
				"workStatus": data2.WORK_STATUS_NEW,
			}
			if err := ad.fsmService.SetState(userId, data2.ENTER_TARGET_CHANNEL_state, d); err != nil {
				return ctx.Reply(data2.CANNOT_SET_USER_state_message, data2.CANCEL_MENU)
			}

			return ctx.Reply(data2.ENTER_TARGET_CHANNEL_message, data2.CANCEL_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)
	// Кнопка "Мои конкурсы 📋", выводит список конкурсов пользователя в виде кнопок
	ad.bot.Handle(
		&data2.MY_GIVES_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data2.CANNOT_FIND_USER_message, data2.CANCEL_MENU)
			}

			userGives, err := ad.giveService.GetAllUserGives(userId)

			if err != nil {
				return ctx.Reply(data2.CANNOT_GET_USER_GIVES_message, data2.START_MENU)
			} else if len(userGives) == 0 {
				return ctx.Reply(data2.NO_GIVES_message, data2.START_MENU)
			}

			buttons := data2.GivesToButtons(userGives)
			buttons = append(buttons, data2.BACK_TO_START_BUTTON)

			givesMenu := data2.CreateReplyMenu(buttons...)

			if err := ad.fsmService.SetState(userId, data2.SELECT_OWN_GIVE_state, nil); err != nil {
				return ctx.Reply(data2.CANNOT_SET_USER_state_message, data2.CANCEL_MENU)
			}

			return ctx.Reply(data2.SELECT_OWN_GIVE_message, givesMenu)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка "Опубликовать ✅" активирует конкурс
	ad.bot.Handle(
		&data2.ACTIVATE_GIVE_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data2.CANNOT_FIND_USER_message, data2.CANCEL_MENU)
			}

			userState, err := ad.fsmService.GetState(userId)
			if err != nil || userState == nil {
				return ctx.Reply(data2.CANNOT_GET_USER_state_message, data2.CANCEL_MENU)
			}

			giveId, err := strconv.Atoi(userState.Data["giveId"])
			if err != nil {
				return ctx.Reply(data2.CANNOT_GET_STATE_DATA_message, data2.CANCEL_MENU)
			}

			give, err := ad.giveService.GetGiveById(giveId)
			if err != nil {
				return ctx.Reply(data2.CANNOT_GET_GIVE_message, data2.CANCEL_MENU)
			}

			unfilledFields := ad.giveService.CheckFilling(&give)
			if len(unfilledFields) != 0 {
				return ctx.Reply(fmt.Sprintf(data2.GIVE_FIELDS_MUST_BE_FILLED_message, unfilledFields))
			}

			err = ad.giveService.UpdateGive(giveId, `"isActive"=?`, true)
			if err != nil {
				return ctx.Reply(data2.CANNOT_UPDATE_GIVE_message, data2.CANCEL_MENU)
			}

			if err := ad.fsmService.SetState(userId, data2.START_MENU_state, nil); err != nil {
				return ctx.Reply(data2.CANNOT_SET_USER_state_message, data2.CANCEL_MENU)
			}

			return ctx.Reply(data2.GIVE_SUCCESSFULLY_ACTIVATE_message, data2.START_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка "Редактировать 🅿️" меню редактирования конкурса
	ad.bot.Handle(
		&data2.EDIT_GIVE_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data2.CANNOT_FIND_USER_message, data2.CANCEL_MENU)
			}

			userState, err := ad.fsmService.GetState(userId)
			if err != nil || userState == nil {
				return ctx.Reply(data2.CANNOT_GET_USER_state_message, data2.CANCEL_MENU)
			}

			if err := ad.fsmService.SetState(userId, data2.SELECT_PROPERTY_TO_EDIT_state, userState.Data); err != nil {
				return ctx.Reply(data2.CANNOT_SET_USER_state_message, data2.CANCEL_MENU)
			}

			return ctx.Reply(data2.SELECT_PROPERTY_TO_EDIT_message, data2.EDIT_GIVE_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка "Снять с публикации ⛔", выключает активный конкурс
	ad.bot.Handle(
		&data2.DEACTIVATE_GIVE_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data2.CANNOT_FIND_USER_message, data2.CANCEL_MENU)
			}

			userState, err := ad.fsmService.GetState(userId)
			if err != nil || userState == nil {
				return ctx.Reply(data2.CANNOT_GET_USER_state_message, data2.CANCEL_MENU)
			}

			giveId, err := strconv.Atoi(userState.Data["giveId"])
			if err != nil {
				return ctx.Reply(data2.CANNOT_GET_STATE_DATA_message, data2.CANCEL_MENU)
			}

			err = ad.giveService.UpdateGive(giveId, `"isActive"=?`, false)
			if err != nil {
				return ctx.Reply(data2.CANNOT_UPDATE_GIVE_message, data2.CANCEL_MENU)
			}

			if err := ad.fsmService.SetState(userId, data2.START_MENU_state, nil); err != nil {
				return ctx.Reply(data2.CANNOT_SET_USER_state_message, data2.CANCEL_MENU)
			}

			return ctx.Reply(data2.GIVE_SUCCESSFULL_DEACTIVATE_message, data2.START_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка редактировани "Заголовок"
	ad.bot.Handle(
		&data2.EDIT_TITLE_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data2.CANNOT_FIND_USER_message, data2.CANCEL_MENU)
			}

			userState, err := ad.fsmService.GetState(userId)
			if err != nil || userState == nil {
				return ctx.Reply(data2.CANNOT_GET_USER_state_message, data2.CANCEL_MENU)
			}

			if userState.State == data2.SELECT_PROPERTY_TO_EDIT_state {
				if err := ad.fsmService.SetState(userId, data2.ENTER_GIVE_TITLE_state, userState.Data); err != nil {
					return ctx.Reply(data2.CANNOT_SET_USER_state_message, data2.CANCEL_MENU)
				}

				return ctx.Reply(data2.ENTER_GIVE_TITLE_message, data2.CANCEL_MENU)
			}
			return ctx.Reply(data2.I_DONT_UNDERSTAND_message, data2.CANCEL_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка редактировани "Описание"
	ad.bot.Handle(
		&data2.EDIT_DESCRIPTION_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data2.CANNOT_FIND_USER_message, data2.CANCEL_MENU)
			}

			userState, err := ad.fsmService.GetState(userId)
			if err != nil || userState == nil {
				return ctx.Reply(data2.CANNOT_GET_USER_state_message, data2.CANCEL_MENU)
			}

			if userState.State == data2.SELECT_PROPERTY_TO_EDIT_state {
				if err := ad.fsmService.SetState(userId, data2.ENTER_GIVE_DESCRIPTION_state, userState.Data); err != nil {
					return ctx.Reply(data2.CANNOT_SET_USER_state_message, data2.CANCEL_MENU)
				}

				return ctx.Reply(data2.ENTER_GIVE_DESCRIPTION_message, data2.CANCEL_MENU)
			}
			return ctx.Reply(data2.I_DONT_UNDERSTAND_message, data2.CANCEL_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)
}
