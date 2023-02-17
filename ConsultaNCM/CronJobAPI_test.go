package ConsultaNCM

import (
	"reflect"
	"testing"
)

func TestNewCronJob(t *testing.T) {
	cron := NewCronJob("* * * * *")
	if got := NewCronJob("* * * * *"); !reflect.DeepEqual(got, cron) {
		t.Errorf("NewCronJob() = %v diferente de %v", got, cron)
	}
}

func Test_cronjob_ConfigurarCron(t *testing.T) {

}
