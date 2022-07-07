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
			time.Sleep(p.publisherTimeout)
		}
	}()
}

func (p *Publisher) serveStartedGives() {
	readyGives := p.giveService.GetStartedGive(p.location)
	if len(readyGives) != 0 {
		for _, g := range readyGives {
			ownerTgId, err := p.userService.GetTgIdByUserId(g.Owner)
			if err != nil {
				p.logger.Errorf("cannot get user userId=%d: %s", g.Owner, err)
				continue
			}

			ownerTgIdInt64, err := utils.StringToInt64(ownerTgId)
			if err != nil {
				p.logger.Errorf("cannot parse ownerTgId=%s string to int64: %s", ownerTgId, err)
				continue
			}

			recipient, err := p.bot.ChatByID(ownerTgIdInt64)
			if err != nil {
				p.logger.Errorf("cannot get chat by id ownerTgId=%d: %s", ownerTgIdInt64, err)
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
					continue
				}
				continue
			}

			err = p.giveService.UpdateGive(g.Id, `"messageId"=?`, message.ID)
			if err != nil {
				p.logger.Errorf("cannot update giveId=%d: %s", g.Id, err)

				_, err = p.bot.Send(recipient, fmt.Sprintf(data.CANNON_UPDATE_GIVE_ON_PUBLICATION_message, g.Id))
				if err != nil {
					p.logger.Errorf("cannot send message to userId=%d: %s", g.Owner, err)
					continue
				}
				continue
			}
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
	fmt.Println(m.Rights, m.Role)
	if err != nil {
		return false
	}

	return true
}
