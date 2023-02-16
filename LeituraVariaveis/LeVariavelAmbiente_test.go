package LeituraVariaveis

import (
	"reflect"
	"testing"
)

func TestNewLeVariavelAmbiente(t *testing.T) {
	tests := []struct {
		name string
		want IVariavelAmbiente
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLeVariavelAmbiente(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLeVariavelAmbiente() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_leVariavelAmbiente_ConnectionHttp(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variavel := &leVariavelAmbiente{}
			if got := variavel.ConnectionHttp(); got != tt.want {
				t.Errorf("ConnectionHttp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_leVariavelAmbiente_ConnectionString(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variavel := &leVariavelAmbiente{}
			if got := variavel.ConnectionString(); got != tt.want {
				t.Errorf("ConnectionString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_leVariavelAmbiente_CronExpression(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variavel := &leVariavelAmbiente{}
			if got := variavel.CronExpression(); got != tt.want {
				t.Errorf("CronExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}
