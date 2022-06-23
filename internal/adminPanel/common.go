package adminPanel

import (
	"Gives_SDT_Bot/internal/give"
	"gopkg.in/telebot.v3"
	"strconv"
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
