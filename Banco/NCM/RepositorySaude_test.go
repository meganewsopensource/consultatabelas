package NCM

import (
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestNewRepositorySaude(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want IRepositorySaude
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRepositorySaude(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRepositorySaude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_saudeBanco_Saudavel(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saude := &saudeBanco{
				db: tt.fields.db,
			}
			if got := saude.Saudavel(); got != tt.want {
				t.Errorf("Saudavel() = %v, want %v", got, tt.want)
			}
		})
	}
}
