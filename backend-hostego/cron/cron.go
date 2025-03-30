package cron

import (
	"sync"

	"github.com/robfig/cron/v3"
)

var (
	cronScheduler *cron.Cron
	once          sync.Once
)

func InitCronJobs() {

}
