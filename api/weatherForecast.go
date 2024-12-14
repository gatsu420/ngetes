package api

import (
	"fmt"
	"time"
)

func (rs *weatherForecastResource) Worker() {
	for {
		err := rs.workers.CreateForecastWorker()
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(10 * time.Second)
	}
}
