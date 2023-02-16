package Web

import (
	"ConsultaTabelas/ConsultaNCM"
	"github.com/gin-gonic/gin"
	"reflect"
	"testing"
)

func TestNewControllerNCM(t *testing.T) {
	type args struct {
		consulta ConsultaNCM.IConsultaNCM
	}
	tests := []struct {
		name string
		args args
		want IControllerNCM
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewControllerNCM(tt.args.consulta); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewControllerNCM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_controllerNCM_AtualizarNCM(t *testing.T) {
	type fields struct {
		consulta ConsultaNCM.IConsultaNCM
	}
	type args struct {
		context *gin.Context
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
			controller := &controllerNCM{
				consulta: tt.fields.consulta,
			}
			controller.AtualizarNCM(tt.args.context)
		})
	}
}

func Test_controllerNCM_DataUltimaAtualizacao(t *testing.T) {
	type fields struct {
		consulta ConsultaNCM.IConsultaNCM
	}
	type args struct {
		context *gin.Context
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
			controller := &controllerNCM{
				consulta: tt.fields.consulta,
			}
			controller.DataUltimaAtualizacao(tt.args.context)
		})
	}
}

func Test_controllerNCM_ListarNCMPorData(t *testing.T) {
	type fields struct {
		consulta ConsultaNCM.IConsultaNCM
	}
	type args struct {
		context *gin.Context
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
			controller := &controllerNCM{
				consulta: tt.fields.consulta,
			}
			controller.ListarNCMPorData(tt.args.context)
		})
	}
}

func Test_controllerNCM_ListarNCMS(t *testing.T) {
	type fields struct {
		consulta ConsultaNCM.IConsultaNCM
	}
	type args struct {
		context *gin.Context
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
			controller := &controllerNCM{
				consulta: tt.fields.consulta,
			}
			controller.ListarNCMS(tt.args.context)
		})
	}
}
