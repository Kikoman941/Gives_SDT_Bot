package telegram

import (
	"strconv"
)

func (b *Bot) addUser(telegramID int64, isAdmin bool) (int64, error) {
	user := &User{
		TgID: strconv.FormatInt(telegramID, 10),
	}

	if isAdmin {
		user.IsAdmin = isAdmin
	}

	result, err := b.storage.Insert(user)
	if err != nil {
		return 0, err
	}
	return result.(*User).ID, nil
}

func (b *Bot) getAdmins() ([]int64, error) {
	var users []User
	var usersIDs []int64
	if err := b.storage.Select(&users, "is_admin = true"); err != nil {
		return nil, err
	}
	for _, user := range users {
		userID, err := strconv.ParseInt(user.TgID, 10, 64)
		if err != nil {
			return nil, err
		}
		usersIDs = append(usersIDs, userID)
	}

	return usersIDs, nil
}
