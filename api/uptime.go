package api

import (
	"time"

	"go.uber.org/zap"
)

func (rs *uptimeResource) Worker() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	for {
		err := rs.workers.CreateUptimeWorker()
		logger.Info("uptime recorded")
		if err != nil {
			logger.Error("uptime not recorded")
		}

		time.Sleep(15 * time.Second)
	}
}
