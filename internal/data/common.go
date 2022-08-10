package data

import (
	"Gives_SDT_Bot/internal/give"
	"Gives_SDT_Bot/pkg/logging"
	"Gives_SDT_Bot/pkg/utils"
	"fmt"
	"gopkg.in/telebot.v3"
	"strconv"
	"strings"
	"time"
)

func GivesToButtons(gives []give.Give) []telebot.Btn {
	var buttons []telebot.Btn
	for _, g := range gives {
		buttons = append(
			buttons,
			telebot.Btn{
				Text: g.Title,
				Data: strconv.Itoa(g.Id),
			},
		)
	}
	return buttons
}

func StringToTimeLocation(str string, logger *logging.Logger, location *time.Location) (time.Time, error) {
	t, err := time.ParseInLocation("02.01.2006 15:04", str, location)
	if err != nil {
		logger.Error(err)
		return time.Time{}, err
	}

	return t, nil
}

func ClearTextForMarkdownV2(text string) string {
	badCharacters := []string{"`", "<", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, ch := range badCharacters {
		text = strings.ReplaceAll(text, ch, fmt.Sprintf("\\%s", ch))
	}
	return text
}

func GetTgLinkNicks(b *telebot.Bot, logger *logging.Logger, tgIds []string) []string {
	var tgNicks []string
	for _, tgId := range tgIds {
		tgIdInt64, err := utils.StringToInt64(tgId)
		if err != nil {
			logger.Errorf("cannot parse tgId=%s string to int64: %s", tgId, err)
			return nil
		}

		chat, err := b.ChatByID(tgIdInt64)
		if err != nil {
			logger.Errorf("cannot get chat by userId userTgId=%d: %s", tgIdInt64, err)
			return nil
		}

		tgNicks = append(tgNicks, fmt.Sprintf("[%s](tg://user?id=%d)\n", chat.Username, tgIdInt64))
	}

	return tgNicks
}

func GetChatNameByChatID(b *telebot.Bot, chatId int64, logger *logging.Logger) string {
	chat, err := b.ChatByID(chatId)
	if err != nil {
		logger.Errorf("cannot get chatName by chatID=%d: %s", chatId, err)
		return ""
	}
	return chat.Title
}

func CreateMessageFromGive(t string, give *give.Give) interface{} {
	var msg interface{}

	switch t {
	case "text":
		msg = ClearTextForMarkdownV2(
			fmt.Sprintf(
				GIVE_CONTENT_message,
				give.Title,
				give.Description,
			),
		)
	case "photo":
		msg = &telebot.Photo{
			File: telebot.FromDisk(fmt.Sprintf("./.images/%s", give.Image)),
		}
		if msg, ok := msg.(telebot.Photo); ok {
			msg.Caption = ClearTextForMarkdownV2(
				fmt.Sprintf(
					GIVE_CONTENT_message,
					give.Title,
					give.Description,
				),
			)
		}
	}

	return msg
}
