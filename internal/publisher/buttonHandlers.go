package publisher

import (
	"Gives_SDT_Bot/pkg/utils"
	"gopkg.in/telebot.v3"
	"strconv"
)

func (p *Publisher) InitButtonHandlers() {
	// Тригерится на Inline кнопку
	p.bot.Handle(
		telebot.OnCallback,
		func(ctx telebot.Context) error {
			giveId, err := strconv.Atoi(ctx.Callback().Data)
			if err != nil {
				p.logger.Errorf("cannot parse giveId=%s to int: %s", ctx.Callback().Data, err)
				return nil
			}

			if err := p.memberService.SaveGiveMember(giveId, utils.Int64ToString(ctx.Sender().ID)); err != nil {
				p.logger.Errorf("cannot save give member giveId=%d memderTgId=%d: %s", giveId, ctx.Sender().ID, err)
				return nil
			}

			return ctx.Respond(
				&telebot.CallbackResponse{
					Text:      "YES",
					ShowAlert: false,
				},
			)
		},
	)
}
