package data

import (
	"Gives_SDT_Bot/internal/give"
	"Gives_SDT_Bot/pkg/logging"
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
