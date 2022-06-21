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
	adminGroup    []int64
	userService   *user.Service
	giveService   *give.Service
	fsmService    *fsm.Service
	imagesService *images.Service
	logger        *logging.Logger
}

func NewAdminPanel(
	bot *telebot.Bot,
	superAdmin int,
	userService *user.Service,
	giveService *give.Service,
	fsmService *fsm.Service,
	imagesService *images.Service,
	logger *logging.Logger,
) *AdminPanel {
	var adminGroup []int64
	adminGroup = append(adminGroup, int64(superAdmin))

	data.InitMenus()

	return &AdminPanel{
		bot:           bot,
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
