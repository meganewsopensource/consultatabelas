package ConsultaNCM

import (
	"reflect"
	"testing"
)

func TestNewCronJob(t *testing.T) {
	type args struct {
		expression string
	}
	tests := []struct {
		name string
		args args
		want ICronJob
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCronJob(tt.args.expression); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCronJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cronjob_ConfigurarCron(t *testing.T) {
	type fields struct {
		expression string
	}
	type args struct {
		consulta func() error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cron := &cronjob{
				expression: tt.fields.expression,
			}
			cron.ConfigurarCron(tt.args.consulta)
		})
	}
}
