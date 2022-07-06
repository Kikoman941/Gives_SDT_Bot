package member

type Member struct {
	tableName  struct{} `pg:"public.gives_members"`
	GiveId     int      `pg:"giveId"`
	memberTgId string   `pg:"memberTgId"`
}
