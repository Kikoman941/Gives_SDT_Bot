package adminPanel

import (
	"Gives_SDT_Bot/internal/give"
	"Gives_SDT_Bot/pkg/logging"
	"gopkg.in/telebot.v3"
	"strconv"
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

func StringToTime(str string, logger *logging.Logger) (time.Time, error) {
	location, _ := time.LoadLocation("Europe/Moscow")
	t, err := time.ParseInLocation("02.01.2006 15:04", str, location)
	if err != nil {
		logger.Error(err)
		return time.Time{}, err
	}

	return t, nil
}
