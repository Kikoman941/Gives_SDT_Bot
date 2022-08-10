package publisher

import (
	"Gives_SDT_Bot/internal/data"
	"Gives_SDT_Bot/pkg/utils"
	"fmt"
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
			memberId := ctx.Sender().ID

			give, err := p.giveService.GetGiveById(giveId)
			if err != nil {
				return ctx.Respond(
					&telebot.CallbackResponse{
						Text:      data.CANNOT_GET_GIVE_message,
						ShowAlert: true,
					},
				)
			}

			for _, channel := range give.TargetChannels {
				channelId, err := utils.StringToInt64(channel)
				if err != nil {
					return ctx.Respond(
						&telebot.CallbackResponse{
							Text:      fmt.Sprintf(data.CANNOT_PARSE_CHANNEL_message, channelId),
							ShowAlert: true,
						},
					)
				}
				if !p.checkMemberSubscribe(memberId, channelId) {
					return ctx.Respond(
						&telebot.CallbackResponse{
							Text:      data.MUST_SUBSCRIBE_message,
							ShowAlert: true,
						},
					)
				}
			}

			if err := p.memberService.SaveGiveMember(giveId, utils.Int64ToString(ctx.Sender().ID)); err != nil {
				if err.Error() == data.ERROR_MEMBER_ALREADY_EXIST.Error() {
					return ctx.Respond(
						&telebot.CallbackResponse{
							Text:      data.MEMBER_ALREADY_EXIST_message,
							ShowAlert: true,
						},
					)
				}
				return ctx.Respond(
					&telebot.CallbackResponse{
						Text:      fmt.Sprintf(data.CANNOT_SAVE_MEMBER_message, giveId, ctx.Sender().ID),
						ShowAlert: true,
					},
				)
			}

			return ctx.Respond(
				&telebot.CallbackResponse{
					Text:      data.MEMBER_SAVE_SUCCESSFULLY_message,
					ShowAlert: true,
				},
			)
		},
	)
}
