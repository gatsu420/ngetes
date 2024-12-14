package workers

import "github.com/gatsu420/ngetes/models"

type UptimeOperations interface {
	CreateUptime(u *models.Uptime) error
}

type UptimeWorkers struct {
	Operations UptimeOperations
}

func NewUptimeWorkers(operations UptimeOperations) *UptimeWorkers {
	return &UptimeWorkers{
		Operations: operations,
	}
}

func (w *UptimeWorkers) CreateUptimeWorker() error {
	uptime := &models.Uptime{}
	err := w.Operations.CreateUptime(uptime)
	if err != nil {
		return err
	}

	return nil
}
