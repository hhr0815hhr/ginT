package cron

import "fmt"

const (
	CronDayly  = "0 0 * * *"
	CronHourly = "0 * * * *"
	CronMinute = "* * * * *"
)

func Every(str string) string {
	return fmt.Sprintf("@every %s", str)
}

type CronFunc func()
type CronItem struct {
	Time string
	Func CronFunc
}


