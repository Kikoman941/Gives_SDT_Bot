package publisher

import (
	"Gives_SDT_Bot/internal/data"
	"Gives_SDT_Bot/internal/give"
	"Gives_SDT_Bot/internal/images"
	"Gives_SDT_Bot/internal/member"
	"Gives_SDT_Bot/internal/user"
	"Gives_SDT_Bot/pkg/logging"
	"Gives_SDT_Bot/pkg/utils"
	"fmt"
	"github.com/go-pg/pg/v10"
	"gopkg.in/telebot.v3"
	"strconv"
	"time"
)

type Publisher struct {
	bot              *telebot.Bot
	userService      *user.Service
	giveService      *give.Service
	memberService    *member.Service
	imagesService    *images.Service
	publisherTimeout time.Duration
	location         *time.Location
	logger           *logging.Logger
}

func NewPublisher(
	bot *telebot.Bot,
	userService *user.Service,
	giveService *give.Service,
	memberService *member.Service,
	imagesService *images.Service,
	publisherTimeout time.Duration,
	location *time.Location,
	logger *logging.Logger,
) (*Publisher, error) {
	return &Publisher{
		bot:              bot,
		userService:      userService,
		giveService:      giveService,
		memberService:    memberService,
		imagesService:    imagesService,
		publisherTimeout: publisherTimeout,
		location:         location,
		logger:           logger,
	}, nil
}

func (p *Publisher) Run() {
	go func() {
		for {
			p.serveStartedGives()
			p.serveFinishedGives()
			time.Sleep(p.publisherTimeout)
		}
	}()
}

func (p *Publisher) serveStartedGives() {
	readyGives := p.giveService.GetStartedGives(p.location)
	if len(readyGives) != 0 {
		for _, g := range readyGives {
			recipient := p.getChatByUserId(g.Owner)
			if recipient == nil {
				p.logger.Errorf("cannot get recipient for userId=%d", g.Owner)
				continue
			}

			chanelId, err := utils.StringToInt64(g.Channel)
			if err != nil {
				p.logger.Errorf("cannot parse chatId=%d string to int64: %s", chanelId, err)
				continue
			}

			msg := &telebot.Photo{
				File: telebot.FromDisk(fmt.Sprintf("./.images/%s", g.Image)),
			}
			msg.Caption = data.ClearTextForMarkdownV2(
				fmt.Sprintf(
					data.GIVE_CONTENT_message,
					g.Title,
					g.Description,
				),
			)

			menu := data.CreateInlineMenu(
				telebot.Btn{
					Text: "Учавствовать",
					Data: strconv.Itoa(g.Id),
				},
			)

			message, err := p.bot.Send(
				&g,
				msg,
				menu,
				telebot.ModeMarkdownV2,
			)
			if err != nil {
				p.logger.Errorf("cannot publish giveId=%d in chanelId=%d: %s", g.Id, chanelId, err)

				_, err = p.bot.Send(recipient, fmt.Sprintf(data.CANNOT_PUBLISH_GIVE_message, g.Id))
				if err != nil {
					p.logger.Errorf("cannot send message to userId=%d: %s", g.Owner, err)
				}
				continue
			}

			err = p.giveService.UpdateGive(g.Id, `"messageId"=?`, message.ID)
			if err != nil {
				p.logger.Errorf("cannot update giveId=%d: %s", g.Id, err)

				_, err = p.bot.Send(recipient, fmt.Sprintf(data.CANNON_UPDATE_GIVE_ON_PUBLICATION_message, g.Id))
				if err != nil {
					p.logger.Errorf("cannot send message to userId=%d: %s", g.Owner, err)
				}
				continue
			}
		}
	}
}

func (p *Publisher) serveFinishedGives() {
	finishedGives := p.giveService.GetFinishedGives(p.location)
	if len(finishedGives) != 0 {
		for _, g := range finishedGives {
			recipient := p.getChatByUserId(g.Owner)
			if recipient == nil {
				p.logger.Errorf("cannot get recipient for userId=%d", g.Owner)
				continue
			}

			winners, err := p.memberService.GetRandomMembersByGiveId(g.Id, g.WinnersCount)
			if err != nil || winners == nil || len(winners) == 0 {
				if _, err := p.bot.Send(recipient, fmt.Sprintf(data.CANNOT_GET_GIVE_WINNERS_message, g.Id)); err != nil {
					p.logger.Errorf("cannot send message to userId=%d: %s", g.Owner, err)
				}
				continue
			}

			winnersTgNicks := p.getTgLinkNicks(winners)
			if len(winnersTgNicks) != len(winners) {
				p.logger.Errorf("cannot get winners nicks for %s", winners)
				continue
			}

			err = p.giveService.UpdateGive(
				g.Id,
				`"isActive"=?, "winners"=?`,
				false,
				pg.Array(winners),
			)
			if err != nil {
				if _, err := p.bot.Send(recipient, fmt.Sprintf(data.CANNOT_UPDATE_FINISHED__GIVE_message, g.Id)); err != nil {
					p.logger.Errorf("cannot send message to userId=%d: %s", g.Owner, err)
				}
				continue
			}

			channelIdInt64, err := utils.StringToInt64(g.Channel)
			if err != nil {
				p.logger.Errorf("cannot parse channelId=%s string to int64: %s", g.Channel, err)
				continue
			}
			text := data.ClearTextForMarkdownV2(
				fmt.Sprintf(
					data.FINISHED_GIVE_CONTENT_message,
					g.Title,
					g.Description,
				),
			)
			_, err = p.bot.EditCaption(
				telebot.StoredMessage{
					MessageID: g.MessageId,
					ChatID:    channelIdInt64,
				},
				text,
				telebot.ModeMarkdownV2,
			)
			if err != nil {
				p.logger.Errorf(
					"cannot edit finished give message giveId=%d channelId=%s messageId=%s: %s",
					g.Id,
					g.Channel,
					g.MessageId,
					err,
				)
			}

			text = data.ClearTextForMarkdownV2(
				fmt.Sprintf(data.GIVE_SUCCESSFULLY_FINISHED_message, g.Title, winnersTgNicks),
			)
			if _, err := p.bot.Send(recipient, text, telebot.ModeMarkdownV2); err != nil {
				p.logger.Errorf("cannot send message to userId=%d: %s", g.Owner, err)
			}

			continue
		}
	}
}

func (p *Publisher) checkMemberSubscribe(memberId int64, channelId int64) bool {
	ch, err := p.bot.ChatByID(channelId)
	if err != nil {
		p.logger.Errorf("cannot get chat by channelId=%d: %s", channelId, err)
		return false
	}
	u, err := p.bot.ChatByID(memberId)
	if err != nil {
		p.logger.Errorf("cannot get user by memberId=%d: %s", memberId, err)
		return false
	}

	m, err := p.bot.ChatMemberOf(ch, u)
	if err != nil {
		p.logger.Errorf("cannot get userId=%d status in chatId=%d: %s", memberId, channelId, err)
		return false
	} else if m.Role == "left" || m.Role == "kicked" {
		return false
	}

	return true
}

func (p *Publisher) getChatByUserId(userId int) *telebot.Chat {
	userTgId, err := p.userService.GetTgIdByUserId(userId)
	if err != nil {
		p.logger.Errorf("cannot get user userId=%d: %s", userId, err)
		return nil
	}

	userTgIdInt64, err := utils.StringToInt64(userTgId)
	if err != nil {
		p.logger.Errorf("cannot parse userTgId=%s string to int64: %s", userTgId, err)
		return nil
	}

	chat, err := p.bot.ChatByID(userTgIdInt64)
	if err != nil {
		p.logger.Errorf("cannot get chat by userId userTgId=%d: %s", userTgIdInt64, err)
		return nil
	}

	return chat
}

func (p *Publisher) getTgLinkNicks(tgIds []string) []string {
	var tgNicks []string
	for _, tgId := range tgIds {
		tgIdInt64, err := utils.StringToInt64(tgId)
		if err != nil {
			p.logger.Errorf("cannot parse tgId=%s string to int64: %s", tgId, err)
			return nil
		}

		chat, err := p.bot.ChatByID(tgIdInt64)
		if err != nil {
			p.logger.Errorf("cannot get chat by userId userTgId=%d: %s", tgIdInt64, err)
			return nil
		}

		tgNicks = append(tgNicks, fmt.Sprintf("[%s](tg://user?id=%d)\n", chat.Username, tgIdInt64))
	}

	return tgNicks
}
