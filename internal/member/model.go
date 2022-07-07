package member

type Member struct {
	tableName  struct{} `pg:"public.gives_members"`
	GiveId     int      `pg:"giveId"`
	MemberTgId string   `pg:"memberTgId"`
}
