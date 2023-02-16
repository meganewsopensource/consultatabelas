package Banco

import (
	"ConsultaTabelas/Banco/NCM"
)

type verificaSaudeBanco struct {
	repository NCM.IRepositorySaude
}

type IVerificaSaudeBanco interface {
	VerificarSaude() (string, bool)
}

func NewVerificaSaudeBanco(repository NCM.IRepositorySaude) IVerificaSaudeBanco {
	return &verificaSaudeBanco{repository}
}

func (verifica *verificaSaudeBanco) VerificarSaude() (string, bool) {
	if verifica.repository.Saudavel() {
		return "Healthy", true
	} else {
		return "Unhealthy", false
	}
}
