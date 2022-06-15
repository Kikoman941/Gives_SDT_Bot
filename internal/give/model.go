package give

import "time"

type Give struct {
	tableName  struct{}  `pg:"public.gives"`
	Id         int       `pg:"id,pk"`
	Owner      int       `pg:"owner"`
	StartTime  time.Time `pg:"start_time"`
	FinishTime time.Time `pg:"finish_time"`
	Text       string    `pg:"text"`
	Image      string    `pg:"image"`
}
