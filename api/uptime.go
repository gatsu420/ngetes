package api

import (
	"log"
	"time"
)

func (rs *uptimeResource) Worker() {
	for {
		err := rs.workers.CreateUptimeWorker()
		log.Println("uptime recorded")
		if err != nil {
			log.Println("uptime not recorded")
		}

		time.Sleep(15 * time.Second)
	}
}
