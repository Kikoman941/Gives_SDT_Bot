package images

import "gopkg.in/telebot.v3"

type Repository interface {
	Download(file *telebot.File, filePath string) error
}
