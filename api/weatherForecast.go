package api

import (
	"fmt"
	"time"
)

func (rs *weatherForecastResource) Job() {
	for {
		err := rs.jobs.CreateJob()
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(30 * time.Second)
	}
}
