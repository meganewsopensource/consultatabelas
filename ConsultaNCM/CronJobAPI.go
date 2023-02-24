package ConsultaNCM

import (
	"github.com/go-co-op/gocron"
	"log"
	"time"
)

type cronjob struct {
	expression string
}

func NewCronJob(expression string) ICronJob {
	return &cronjob{expression: expression}
}

type ICronJob interface {
	ConfigurarCron(consulta func() error)
}

func (cron *cronjob) ConfigurarCron(consulta func() error) {
	s := gocron.NewScheduler(time.Local)
	s.Cron(cron.expression).Do(
		func() {
			err := consulta()
			if err != nil {
				log.Fatal(err)
			}
		})
	s.StartAsync()
}
