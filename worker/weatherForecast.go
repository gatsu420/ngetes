package worker

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

type WeatherForecastWorkers struct {
	Operations WeatherForecastOperations
}

func NewWeatherForecastWorkers(operations WeatherForecastOperations) *WeatherForecastWorkers {
	return &WeatherForecastWorkers{
		Operations: operations,
	}
}

func (w *WeatherForecastWorkers) CreateForecastWorker() error {
	client := &http.Client{
		Timeout: 3000 * time.Millisecond,
	}

	forecastURL := "https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&hourly=temperature_2m"
	forecast, err := client.Get(forecastURL)
	if err != nil {
		return err
	}
	defer forecast.Body.Close()

	forecastBody, err := io.ReadAll(forecast.Body)
	if err != nil {
		return err
	}

	forecastModel := &models.WeatherForecast{
		StatusCode: forecast.StatusCode,
		URL:        forecastURL,
		Payload:    string(forecastBody),
	}
	forecastID, forecastCreatedAt, err := w.Operations.CreateForecast(forecastModel)
	if err != nil {
		return err
	}

	fmt.Println(forecastID, forecastCreatedAt)
	return nil
}
