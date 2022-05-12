package data

type TelegramData struct {
	Buttons  *Buttons
	Menus    *Menus
	Messages map[string]string
}

func NewTelegramData() *TelegramData {
	menus := newMenus()
	buttons := newButtons(menus)
	menus.init(buttons)

	return &TelegramData{
		Menus:    menus,
		Buttons:  buttons,
		Messages: getMessages(),
	}
}
