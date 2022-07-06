package adminPanel

import (
	"Gives_SDT_Bot/internal/data"
	"Gives_SDT_Bot/internal/fsm"
	"Gives_SDT_Bot/internal/give"
	"Gives_SDT_Bot/internal/images"
	"Gives_SDT_Bot/internal/user"
	"Gives_SDT_Bot/pkg/logging"
	"gopkg.in/telebot.v3"
	"time"
)

type AdminPanel struct {
	bot           *telebot.Bot
	adminGroup    []int64
	userService   *user.Service
	giveService   *give.Service
	fsmService    *fsm.Service
	imagesService *images.Service
	location      *time.Location
	logger        *logging.Logger
}

func NewAdminPanel(
	bot *telebot.Bot,
	superAdmin int,
	userService *user.Service,
	giveService *give.Service,
	fsmService *fsm.Service,
	imagesService *images.Service,
	location *time.Location,
	logger *logging.Logger,
) *AdminPanel {
	var adminGroup []int64
	adminGroup = append(adminGroup, int64(superAdmin))

	data.InitMenu()

	return &AdminPanel{
		bot:           bot,
		adminGroup:    adminGroup,
		userService:   userService,
		giveService:   giveService,
		fsmService:    fsmService,
		imagesService: imagesService,
		location:      location,
		logger:        logger,
	}
}

func (ad *AdminPanel) RefreshAdmins() error {
	admins, err := ad.userService.GetAdmins()
	if err != nil {
		return err
	} else if len(admins) == 0 {
		return data.ERROR_NO_ADMINS_FOR_REFRESH
	}
	ad.adminGroup = admins
	return nil
}

func (ad *AdminPanel) checkBotIsAdmin(channelId int64) (bool, error) {
	ch, err := ad.bot.ChatByID(channelId)
	if err != nil {
		ad.logger.Errorf("cannot get chat by channelId=%d: %s", channelId, err)
		return false, err
	}

	channelAdmins, err := ad.bot.AdminsOf(ch)
	if err != nil {
		ad.logger.Errorf("cannot get admins of channelId=%d: %s", channelId, err)
		return false, err
	}

	for _, admin := range channelAdmins {
		if admin.User.Username == ad.bot.Me.Username {
			return true, nil
		}
	}
	return false, nil
}
