package fsm

type Userstate struct {
	tableName struct{}          `pg:"public.users_states"`
	UserID    int               `pg:"userId,fk,unique"`
	state     string            `pg:"state"`
	Data      map[string]string `pg:"data"`
}
