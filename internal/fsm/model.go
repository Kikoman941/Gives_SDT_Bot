package fsm

import (
	"Gives_SDT_Bot/pkg/logging"
)

type FSM struct {
	repository Repository
	logger     *logging.Logger
}

type UserState struct {
	tableName struct{} `pg:"public.user_state"`
	UserID    int      `pg:"user_id,fk,unique"`
	State     string   `pg:"state"`
}
