package Web

import (
	"ConsultaTabelas/Banco"
	"github.com/gin-gonic/gin"
	"reflect"
	"testing"
)

func TestNewControllerSaude(t *testing.T) {
	type args struct {
		verificaSaude Banco.IVerificaSaudeBanco
	}
	tests := []struct {
		name string
		args args
		want IControllerSaude
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewControllerSaude(tt.args.verificaSaude); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewControllerSaude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_controllerSaude_VerificarSaude(t *testing.T) {
	type fields struct {
		verificaSaude Banco.IVerificaSaudeBanco
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
			controller := &controllerSaude{
				verificaSaude: tt.fields.verificaSaude,
			}
			controller.VerificarSaude(tt.args.context)
		})
	}
}
