package main

import (
	"fmt"
	"time"
)

func main() {
	loc, _ := time.LoadLocation("Europe/Moscow")
	t, err := time.ParseInLocation("02.01.2006 15:04", "01.06.2022 13:00", loc)
	if err != nil {
		fmt.Errorf("%s", err)
	}
	fmt.Println(t.Format(time.RFC3339))
}
