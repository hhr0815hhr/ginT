package cron

import (
	cron2 "github.com/hhr0815hhr/gint/internal/cron"
)

type CronFunc func()
type CronItem struct {
	Time string
	Func CronFunc
}

var (
	// 存储所有的定时任务
	cronItems = []CronItem{
		{Time: cron2.CronDayly, Func: func() {}},   // 每天0点执行
		{Time: cron2.Every("1s"), Func: func() {}}, // 每s执行
	}
)
