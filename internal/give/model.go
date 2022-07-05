package give

import "time"

type Give struct {
	tableName      struct{}  `pg:"public.gives"`
	Id             int       `pg:"id,pk"`
	Owner          int       `pg:"owner"`
	IsActive       bool      `pg:"isActive"`
	IsDeleted      bool      `pg:"isDeleted"`
	StartAt        time.Time `pg:"startAt"`
	FinishAt       time.Time `pg:"finishAt"`
	Title          string    `pg:"title"`
	Description    string    `pg:"description"`
	Image          string    `pg:"image"`
	WinnersCount   int       `pg:"winnersCount"`
	Channel        string    `pg:"channel"`
	MessageId      string    `pg:"messageId"`
	TargetChannels []string  `pg:"targetChannels,array"`
}

func (g *Give) Recipient() string {
	return g.Channel
}
