package give

import "time"

type Give struct {
	tableName      struct{}  `pg:"public.gives"`
	Id             int       `pg:"id,pk"`
	Owner          int       `pg:"owner"`
	IsActive       bool      `pg:"is_active"`
	StartTime      time.Time `pg:"start_at"`
	FinishTime     time.Time `pg:"finish_at"`
	Title          string    `pg:"title"`
	Description    string    `pg:"description"`
	Images         []string  `pg:"images"`
	Channel        int       `pg:"channel"`
	TargetChannels []int     `pg:"target_channels"`
}
