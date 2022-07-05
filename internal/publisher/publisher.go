package publisher

import (
	"Gives_SDT_Bot/internal/data"
	"Gives_SDT_Bot/internal/give"
	"Gives_SDT_Bot/internal/images"
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
	imagesService    *images.Service
	publisherTimeout time.Duration
	logger           *logging.Logger
}

func NewPublisher(
	bot *telebot.Bot,
	userService *user.Service,
	giveService *give.Service,
	imagesService *images.Service,
	publisherTimeout time.Duration,
	logger *logging.Logger,
) (*Publisher, error) {
	return &Publisher{
		bot:              bot,
		userService:      userService,
		giveService:      giveService,
		imagesService:    imagesService,
		publisherTimeout: publisherTimeout,
		logger:           logger,
	}, nil
}

func (p *Publisher) Run() {
	go func() {
		for {
			readyGives := p.giveService.GetStartedGive()
			if len(readyGives) == 0 {
				continue
			}

			for _, g := range readyGives {
				tgId, err := p.userService.GetTgIdByUserId(g.Owner)
				if err != nil {
					p.logger.Errorf("cannot get user userId=%d: %s", g.Owner, err)
					continue
				}

				recipient, err := p.bot.ChatByID(int64(tgId))
				if err != nil {
					p.logger.Errorf("cannot get chat by id tgId=%d: %s", tgId, err)
					continue
				}

				chanelId, err := utils.StringToInt64(g.Channel)
				if err != nil {
					p.logger.Errorf("cannot parse chatId=%d string to int64: %s", chanelId, err)
					continue
				}

				menu := data.CreateReplyMenu(
					telebot.Btn{
						Text:   "Учавствовать",
						Unique: strconv.Itoa(g.Id),
					},
				)

				message, err := p.bot.Send(
					&g,
					fmt.Sprintf(data.GIVE_CONTENT_message, g.Title, g.Description),
					menu,
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

			time.Sleep(p.publisherTimeout)
		}
	}()
}
