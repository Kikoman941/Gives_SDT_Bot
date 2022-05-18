package telegram

import (
	"errors"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

const (
	commandStart         = "/start"
	commandRefreshAdmins = "/refreshAdmins"
	commandTestDB        = "/testDB"
)

var (
	errorUserNotFound = errors.New("telegram: Bad Request: user not found (400)")
)

func (b *Bot) InitHandlers() {
	// Комманда /start
	b.bot.Handle(
		commandStart,
		func(ctx telebot.Context) error {
			userID, err := b.addUser(ctx.Chat().ID, false)
			if err != nil {
				return err
			}

			if err := b.fsm.setState(userID, MAIN_MENU); err != nil {
				return err
			}

			return ctx.Reply(b.telegramData.Messages["startMessage"], b.telegramData.Menus.StartMenu)
		},
	)

	// Комманда /refreshAdmins. Обновляет список админов для мидлваря Whitelist
	b.bot.Handle(
		commandRefreshAdmins,
		func(ctx telebot.Context) error {
			admins, err := b.getAdmins()
			if err != nil {
				return err
			}

			b.refreshAdmins(admins)
			b.InitHandlers()

			return ctx.Reply(b.telegramData.Messages["successRefreshAdmins"])
		},
		middleware.Whitelist(b.adminGroup...),
	)

	b.bot.Handle(
		&b.telegramData.Buttons.CreateGiveButton,
		func(ctx telebot.Context) error {
			return ctx.Reply(b.telegramData.Messages["addTargetChanel"])
		},
		middleware.Whitelist(b.adminGroup...),
	)

	b.bot.Handle(
		telebot.OnText,
		func(ctx telebot.Context) error {
			_, err := b.bot.ChatMemberOf(getChatIdFromInt(1001422240135), getChatIdFromInt(10784525812))
			if errors.As(err, &errorUserNotFound) {
				return ctx.Reply("Пользователь не подписан")
			} else {
				return err
			}
		},
		middleware.Whitelist(b.adminGroup...),
	)

	b.bot.Handle(
		commandTestDB,
		func(ctx telebot.Context) error {
			//var users []User
			//err := b.storage.Select(&users, "users", "")
			//if err != nil {
			//	return err
			//}
			//fmt.Println(users[0].IsAdmin)
			userId, err := b.addUser(ctx.Chat().ID, true)
			if err != nil {
				return err
			}

			if err := b.fsm.setState(userId, MAIN_MENU); err != nil {
				return err
			}

			return ctx.Reply("You push button")
		},
		middleware.Whitelist(b.adminGroup...),
	)
}
