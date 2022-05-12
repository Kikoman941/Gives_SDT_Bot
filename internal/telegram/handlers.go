package telegram

import (
	"errors"
	"gopkg.in/telebot.v3"
	"strconv"
)

const (
	commandStart  = "/start"
	commandTest   = "/test"
	commandTestDB = "/testDB"
)

var (
	errorUserNotFound = errors.New("telegram: Bad Request: user not found (400)")
)

func (b *Bot) InitHandlers() {
	b.bot.Handle(commandStart, func(ctx telebot.Context) error {
		return ctx.Reply(b.telegramData.Messages["startMessage"], b.telegramData.Menus.StartMenu)
	})

	b.bot.Handle(&b.telegramData.Buttons.CreateGiveButton, func(ctx telebot.Context) error {
		return ctx.Reply("You push button")
	})

	b.bot.Handle(commandTest, func(ctx telebot.Context) error {
		_, err := b.bot.ChatMemberOf(getChatId("1001422240135"), getChatId("10784525812"))
		if errors.As(err, &errorUserNotFound) {
			return ctx.Reply("Пользователь не подписан")
		} else {
			return err
		}
	})

	b.bot.Handle(commandTestDB, func(ctx telebot.Context) error {
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
	})
}

func (b *Bot) addUser(telegramID int64, isAdmin bool) (int, error) {
	user := &User{
		TgID:    strconv.FormatInt(telegramID, 10),
		IsAdmin: isAdmin,
	}
	result, err := b.storage.Insert(user)
	if err != nil {
		return 0, err
	}
	return result.(*User).ID, nil
}
