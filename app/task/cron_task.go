package task

import (
	"go.uber.org/zap"
	"time"

	"github.com/robfig/cron/v3"
)

func RunCorn() {
	nyc, _ := time.LoadLocation("local")
	cronder := cron.New(cron.WithSeconds(), cron.WithLocation(nyc))
	cronder.AddFunc("0 0 0 * * *", ResetDepartMp)
}

func ResetDepartMp() {
	defer func() {
		if err := recover(); err != nil {
			zap.S().Error(err)
			return
		}
	}()
	TaskStorage.DepartMp = make(map[int64]int64)
}
