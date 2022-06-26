package fsm

type UserState struct {
	tableName struct{}          `pg:"public.users_states"`
	UserID    int               `pg:"userId,fk,unique"`
	State     string            `pg:"state"`
	Data      map[string]string `pg:"data"`
}
