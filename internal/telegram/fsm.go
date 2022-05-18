package telegram

import (
	"Gives_SDT_Bot/internal/storage"
)

type FSM struct {
	storage storage.BotStorage
}

func NewFSM(storage storage.BotStorage) *FSM {
	return &FSM{
		storage: storage,
	}
}

func (fsm *FSM) setState(userID int64, state string) error {
	userState := &UserState{
		UserID: userID,
		State:  state,
	}
	_, err := fsm.storage.Insert(userState)
	if err != nil {
		return err
	}
	return nil
}
