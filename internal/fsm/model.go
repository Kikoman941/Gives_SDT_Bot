package fsm

type UserState struct {
	tableName struct{} `pg:"public.users_states"`
	UserID    int      `pg:"user_id,fk,unique"`
	State     string   `pg:"state"`
}
