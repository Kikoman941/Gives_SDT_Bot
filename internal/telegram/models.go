package telegram

import "time"

// Users states
const (
	MAIN_MENU = "main_menu"
)

type User struct {
	tableName struct{} `pg:"public.users"`
	ID        int      `pg:"id,pk"`
	TgID      string   `pg:"tg_id,unique"`
	IsAdmin   bool     `pg:"is_admin"`
}

type UserState struct {
	tableName struct{} `pg:"public.users_states"`
	UserID    int      `pg:"user_id,fk,unique"`
	State     string   `pg:"state"`
}

type Give struct {
	tableName  struct{}  `pg:"public.gives"`
	Id         int       `pg:"id,pk"`
	Owner      int       `pg:"owner"`
	StartTime  time.Time `pg:"start_time"`
	FinishTime time.Time `pg:"finish_time"`
	Text       string    `pg:"text"`
	Image      string    `pg:"image"`
}
