package user

type User struct {
	tableName struct{} `pg:"public.users"`
	ID        int      `pg:"id,pk"`
	TgID      string   `pg:"tg_id,unique"`
	IsAdmin   bool     `pg:"is_admin"`
}
