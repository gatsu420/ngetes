package jobs

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gatsu420/ngetes/models"
)

type WeatherForecastOperations interface {
	CreateForecast(f *models.WeatherForecast) (id int, createdAt time.Time, err error)
}

type WeatherForecastJobs struct {
	Operations WeatherForecastOperations
}

func NewWeatherForecastJobs(ops WeatherForecastOperations) *WeatherForecastJobs {
	return &WeatherForecastJobs{
		Operations: ops,
	}
}

func (j *WeatherForecastJobs) CreateJob() error {
	httpResp, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&hourly=temperature_2m")
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	respPrettyBody, respStatusCode := string(respBody), httpResp.StatusCode

	forecast := &models.WeatherForecast{
		StatusCode: respStatusCode,
		Payload:    respPrettyBody,
	}
	forecastID, forecastCreatedAt, err := j.Operations.CreateForecast(forecast)
	if err != nil {
		return err
	}
	fmt.Println(forecastID)
	fmt.Println(forecastCreatedAt)

	return nil
}
