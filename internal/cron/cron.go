package cron

import (
	"github.com/hhr0815hhr/gint/internal/log"
	c3 "github.com/robfig/cron/v3"
)

type CronJob struct {
	c *c3.Cron
}

func New() *CronJob {
	return &CronJob{
		c: c3.New(),
	}
}

func (c *CronJob) AddFunc(spec string, fn func()) {
	fd, err := c.c.AddFunc(spec, fn)
	if err != nil {
		log.Logger.Errorf("CronJob add func %s error: %s", spec, err.Error())
		return
	}
	log.Logger.Infof("CronJob add func %s, id: %d", spec, fd)
}

func (c *CronJob) AddFuncs() {
	for _, item := range cronItems {
		c.AddFunc(item.Time, item.Func)
	}
}

func (c *CronJob) Start() {
	c.c.Start()
	log.Logger.Info("CronJob started")
}
