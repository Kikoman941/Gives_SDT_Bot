package adminPanel

import (
	"Gives_SDT_Bot/internal/adminPanel/data"
	"Gives_SDT_Bot/internal/fsm"
	"Gives_SDT_Bot/internal/give"
	"Gives_SDT_Bot/internal/images"
	"Gives_SDT_Bot/internal/user"
	"Gives_SDT_Bot/pkg/logging"
	"gopkg.in/telebot.v3"
)

type AdminPanel struct {
	bot           *telebot.Bot
	botUsername   string
	adminGroup    []int64
	userService   *user.Service
	giveService   *give.Service
	fsmService    *fsm.Service
	imagesService *images.Service
	logger        *logging.Logger
}

func NewAdminPanel(
	bot *telebot.Bot,
	botUsername string,
	superAdmin int,
	userService *user.Service,
	giveService *give.Service,
	fsmService *fsm.Service,
	imagesService *images.Service,
	logger *logging.Logger,
) *AdminPanel {
	var adminGroup []int64
	adminGroup = append(adminGroup, int64(superAdmin))

	data.InitMenu()

	return &AdminPanel{
		bot:           bot,
		botUsername:   botUsername,
		adminGroup:    adminGroup,
		userService:   userService,
		giveService:   giveService,
		fsmService:    fsmService,
		imagesService: imagesService,
		logger:        logger,
	}
}

func (ad *AdminPanel) refreshAdmins(admins []int64) {
	ad.adminGroup = admins
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
		if admin.User.Username == ad.botUsername {
			return true, nil
		}
	}
	return false, nil
}
