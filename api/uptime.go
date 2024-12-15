package api

import (
	"time"

	"github.com/gatsu420/ngetes/logger"
	"go.uber.org/zap"
)

func (rs *uptimeResource) Worker() {
	for {
		err := rs.workers.CreateUptimeWorker()
		logger.Logger.Info("uptime recorded")
		if err != nil {
			logger.Logger.Error("uptime not recorded", zap.Error(err))
		}

		time.Sleep(15 * time.Second)
	}
}
