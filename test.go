package main

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"time"
)

type Give struct {
	tableName      struct{}  `pg:"public.gives"`
	Id             int       `pg:"id,pk"`
	Owner          int       `pg:"owner"`
	IsActive       bool      `pg:"isActive"`
	StartAt        time.Time `pg:"startAt"`
	FinishAt       time.Time `pg:"finishAt"`
	Title          string    `pg:"title"`
	Description    string    `pg:"description"`
	Image          string    `pg:"image"`
	WinnersCount   int       `pg:"winnersCount"`
	Channel        string    `pg:"channel"`
	TargetChannels []string  `pg:"targetChannels,array"`
}

func main() {
	opt, _ := pg.ParseURL("postgresql://superadmin:123123@192.168.57.81:8432/telegram?sslmode=disable")
	db := pg.Connect(opt)

	var give []Give
	err := db.Model(&give).
		Where(`"id"=?`, 3).
		Select()
	if err != nil {
		fmt.Println(err)
	}
}
