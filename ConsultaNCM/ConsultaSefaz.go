package ConsultaNCM

import (
	"ConsultaTabelas/ConsultaHTTP"
	"encoding/json"
)

type IConsultaSefaz interface {
	ConsultarNCM() (NcmReceita, error)
}

type consultaSefaz struct {
	consultaHttp ConsultaHTTP.IConsultaHttp
}

func New(consulta ConsultaHTTP.IConsultaHttp) IConsultaSefaz {
	return &consultaSefaz{consulta}
}

func (consulta consultaSefaz) ConsultarNCM() (NcmReceita, error) {
	resposta, err := consulta.consultaHttp.Consultar()
	if err != nil {
		return NcmReceita{}, err
	}
	ncm, err := consulta.desserealizarNCM(resposta)
	if err != nil {
		return NcmReceita{}, err
	}
	return ncm, nil
}

func (consulta consultaSefaz) desserealizarNCM(resposta []byte) (NcmReceita, error) {
	var ncm NcmReceita
	err := json.Unmarshal(resposta, &ncm)
	if err != nil {
		return NcmReceita{}, err
	}
	return ncm, nil
}
