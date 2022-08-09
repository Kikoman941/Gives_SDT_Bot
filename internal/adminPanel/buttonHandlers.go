package adminPanel

import (
	"Gives_SDT_Bot/internal/data"
	"Gives_SDT_Bot/pkg/utils"
	"fmt"
	"github.com/go-pg/pg/v10"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"strconv"
)

func (ad *AdminPanel) InitButtonHandlers() {
	// Кнопка "Назад в главное меню", отмена любого состояния до старта
	ad.bot.Handle(
		&data.BACK_TO_START_BUTTON,
		func(ctx telebot.Context) error {
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_message, data.START_MENU)
			}

			if err := ad.fsmService.SetState(userId, data.START_MENU_state, nil); err != nil {
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

			d := map[string]string{
				"workStatus": data.WORK_STATUS_NEW,
			}
			if err := ad.fsmService.SetState(userId, data.ENTER_GIVE_TITLE_state, d); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
			}

			if err := ctx.Reply(data.NEW_GIVE_message); err != nil {
				return err
			}
			return ctx.Send(data.ENTER_GIVE_TITLE_message, data.CANCEL_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)
	// Кнопка "Мои конкурсы 📋", выводит список конкурсов пользователя в виде кнопок
	ad.bot.Handle(
		&data.MY_GIVES_BUTTON,
		func(ctx telebot.Context) error {
			fmt.Println()
			userId, err := ad.userService.GetUserIdByTgId(ctx.Chat().ID)
			if err != nil || userId == 0 {
				return ctx.Reply(data.CANNOT_FIND_USER_message, data.CANCEL_MENU)
			}

			userGives, err := ad.giveService.GetAllUserGives(userId)

			if err != nil {
				return ctx.Reply(data.CANNOT_GET_USER_GIVES_message, data.START_MENU)
			} else if len(userGives) == 0 {
				return ctx.Reply(data.NO_GIVES_message, data.START_MENU)
			}

			buttons := data.GivesToButtons(userGives)
			buttons = append(buttons, data.BACK_TO_START_BUTTON)

			givesMenu := data.CreateReplyMenu(buttons...)

			if err := ad.fsmService.SetState(userId, data.SELECT_OWN_GIVE_state, nil); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
			}

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
				return ctx.Reply(data.CANNOT_GET_STATE_DATA_message, data.CANCEL_MENU)
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

			if err := ad.fsmService.SetState(userId, data.START_MENU_state, nil); err != nil {
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

			if err := ad.fsmService.SetState(userId, data.SELECT_PROPERTY_TO_EDIT_state, userState.Data); err != nil {
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
				return ctx.Reply(data.CANNOT_GET_STATE_DATA_message, data.CANCEL_MENU)
			}

			give, err := ad.giveService.GetGiveById(giveId)
			if err != nil {
				return ctx.Reply(data.CANNOT_GET_GIVE_message, data.CANCEL_MENU)
			}

			winners, err := ad.memberService.GetRandomMembersByGiveId(give.Id, give.WinnersCount)
			if err != nil || winners == nil || len(winners) == 0 {
				return ctx.Reply(fmt.Sprintf(data.CANNOT_GET_GIVE_WINNERS_message, give.Id))
			}

			winnersTgNicks := data.GetTgLinkNicks(ad.bot, ad.logger, winners)
			if len(winnersTgNicks) != len(winners) {
				return ctx.Reply(fmt.Sprintf(data.CANNOT_GET_WINNERS_DATA_message, winners), data.CANCEL_MENU)
			}

			err = ad.giveService.UpdateGive(
				give.Id,
				`"isActive"=?, "winners"=?`,
				false,
				pg.Array(winners),
			)
			if err != nil {
				if _, err := ad.bot.Send(ctx.Recipient(), fmt.Sprintf(data.CANNOT_UPDATE_FINISHED_GIVE_message, give.Id)); err != nil {
					return ctx.Reply(fmt.Sprintf(data.CANNOT_SEND_message, err))
				}
			}

			if err := ad.fsmService.SetState(userId, data.START_MENU_state, nil); err != nil {
				return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
			}

			channelIdInt64, err := utils.StringToInt64(give.Channel)
			if err != nil {
				return ctx.Reply(data.CANNOT_PARSE_CHANNEL_message, give.Channel, err)
			}
			text := data.ClearTextForMarkdownV2(
				fmt.Sprintf(
					data.FINISHED_GIVE_CONTENT_message,
					give.Title,
					give.Description,
				),
			)
			_, err = ad.bot.EditCaption(
				telebot.StoredMessage{
					MessageID: give.MessageId,
					ChatID:    channelIdInt64,
				},
				text,
				telebot.ModeMarkdownV2,
			)
			if err != nil {
				return ctx.Reply(fmt.Sprintf(data.CANNOT_UPDATE_FINISHED_GIVE_message, give.Id))
			}

			text = data.ClearTextForMarkdownV2(
				fmt.Sprintf(data.GIVE_SUCCESSFULLY_FINISHED_message, give.Title, winnersTgNicks),
			)
			if _, err := ad.bot.Send(ctx.Recipient(), text, telebot.ModeMarkdownV2); err != nil {
				return ctx.Reply(fmt.Sprintf(data.CANNOT_SEND_message, err))
			}

			return ctx.Reply(data.GIVE_SUCCESSFULL_DEACTIVATE_message, data.START_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка редактировани "Заголовок"
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
				if err := ad.fsmService.SetState(userId, data.ENTER_GIVE_TITLE_state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(data.ENTER_GIVE_TITLE_message, data.CANCEL_MENU)
			}
			return ctx.Reply(data.I_DONT_UNDERSTAND_message, data.CANCEL_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка редактирования "Описание"
	ad.bot.Handle(
		&data.EDIT_DESCRIPTION_BUTTON,
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
				if err := ad.fsmService.SetState(userId, data.ENTER_GIVE_DESCRIPTION_state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(data.ENTER_GIVE_DESCRIPTION_message, data.CANCEL_MENU)
			}
			return ctx.Reply(data.I_DONT_UNDERSTAND_message, data.CANCEL_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка редактирования "Дата старта - финиша"
	ad.bot.Handle(
		&data.EDIT_START_FINISH_BUTTON,
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
				if err := ad.fsmService.SetState(userId, data.ENTER_GIVE_START_FINISH_state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(data.ENTER_GIVE_START_FINISH_message, data.CANCEL_MENU)
			}
			return ctx.Reply(data.I_DONT_UNDERSTAND_message, data.CANCEL_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка редактирования "Картинка"
	ad.bot.Handle(
		&data.EDIT_IMAGE_BUTTON,
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
				if err := ad.fsmService.SetState(userId, data.UPLOAD_GIVE_IMAGE_state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(data.UPLOAD_GIVE_IMAGE_message, data.CANCEL_MENU)
			}
			return ctx.Reply(data.I_DONT_UNDERSTAND_message, data.CANCEL_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка редактирования "Колличество победителей"
	ad.bot.Handle(
		&data.EDIT_WINNERS_COUNT_BUTTON,
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
				if err := ad.fsmService.SetState(userId, data.ENTER_WINNERS_COUNT_state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(data.ENTER_WINNERS_COUNT_message, data.CANCEL_MENU)
			}
			return ctx.Reply(data.I_DONT_UNDERSTAND_message, data.CANCEL_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка редактирования "Канал конкурса"
	ad.bot.Handle(
		&data.EDIT_CHANNEL_BUTTON,
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
				if err := ad.fsmService.SetState(userId, data.ENTER_TARGET_CHANNEL_state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(data.ENTER_TARGET_CHANNEL_message, data.CANCEL_MENU)
			}
			return ctx.Reply(data.I_DONT_UNDERSTAND_message, data.CANCEL_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка редактирования "Каналы проверки подписки"
	ad.bot.Handle(
		&data.EDIT_TARGET_CHANNELS_BUTTON,
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
				if err := ad.fsmService.SetState(userId, data.ENTER_SUBSCRIPTION_CHANNELS_state, userState.Data); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(data.ENTER_SUBSCRIPTION_CHANNELS_message, data.CANCEL_MENU)
			}
			return ctx.Reply(data.I_DONT_UNDERSTAND_message, data.CANCEL_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)

	// Кнопка редактирования "Удалить ❌"
	ad.bot.Handle(
		&data.DELETE_GIVE_BUTTON,
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
					return ctx.Reply(data.CANNOT_GET_STATE_DATA_message, data.CANCEL_MENU)
				}
				err = ad.giveService.UpdateGive(giveId, `"isDeleted"=?`, true)
				if err != nil {
					return ctx.Reply(data.CANNOT_UPDATE_GIVE_message, data.CANCEL_MENU)
				}

				if err := ad.memberService.ClearGiveMembers(giveId); err != nil {
					return ctx.Reply(fmt.Sprintf(data.CANNOT_CLEAR_GIVE_MEMBERS_message, giveId))
				}

				if err := ad.fsmService.SetState(userId, data.START_MENU_state, nil); err != nil {
					return ctx.Reply(data.CANNOT_SET_USER_state_message, data.CANCEL_MENU)
				}

				return ctx.Reply(data.GIVE_SUCCESSFULLY_DELETE_message, data.START_MENU)
			}
			return ctx.Reply(data.I_DONT_UNDERSTAND_message, data.CANCEL_MENU)
		},
		middleware.Whitelist(ad.adminGroup...),
	)
}
