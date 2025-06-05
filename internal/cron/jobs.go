package cron

var (
	// 存储所有的定时任务
	cronItems = []CronItem{
		{Time: CronDayly, Func: func() {}},   // 每天0点执行
		{Time: Every("1s"), Func: func() {}}, // 每1秒执行
	}
)
