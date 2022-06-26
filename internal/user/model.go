package user

type User struct {
	tableName struct{} `pg:"public.users"`
	ID        int      `pg:"id,pk"`
	TgID      string   `pg:"tgId,unique"`
	IsAdmin   bool     `pg:"isAdmin"`
}
