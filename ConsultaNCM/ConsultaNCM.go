package ConsultaNCM

import (
	"ConsultaTabelas/Banco/NCM"
	"ConsultaTabelas/ConsultaNCMSefaz"
	"time"
)

type IConsultaNCM interface {
	AtualizarNCM() error
}

type consultaNCM struct {
	consultaSefaz          ConsultaNCMSefaz.IConsultaSefaz
	respotoryNCM           NCM.IRepositoryNCM
	repositoryNomenclatura NCM.IRepositoryNomenclatura
	modeloData             string
}

func NewConsultaNCM(consultaSefaz ConsultaNCMSefaz.IConsultaSefaz,
	respotoryNCM NCM.IRepositoryNCM,
	repositoryNomenclatura NCM.IRepositoryNomenclatura) IConsultaNCM {
	return &consultaNCM{
		consultaSefaz:          consultaSefaz,
		respotoryNCM:           respotoryNCM,
		repositoryNomenclatura: repositoryNomenclatura,
		modeloData:             "02/01/2006",
	}
}

func (consulta *consultaNCM) AtualizarNCM() error {
	dadosConsulta, err := consulta.consultaSefaz.ConsultarNCM()
	if err != nil {
		return err
	}
	ncm, err := consulta.respotoryNCM.GetByID(1)
	if err != nil {
		return err
	}
	data, err := time.Parse(consulta.modeloData, dadosConsulta.DataUltimaAtualizacaoNcm)
	if data.After(ncm.DataUltimaAtualizacaoNcm) {
		lista, err := consulta.listaNomenclatura(dadosConsulta.Nomenclaturas)
		if err != nil {
			return err
		}
		ncmBanco := NCM.NcmBanco{
			ID:                       ncm.ID,
			DataUltimaAtualizacaoNcm: data,
			Nomenclaturas:            lista,
		}
		err = consulta.gravarNCM(ncmBanco)
	}
	return err
}

func (consulta *consultaNCM) gravarNCM(ncm NCM.NcmBanco) error {
	var err error
	err = nil
	resposta, err := consulta.respotoryNCM.GetByID(1)
	if err != nil {
		return err
	}

	if resposta.ID == 1 {
		resposta.DataUltimaAtualizacaoNcm = ncm.DataUltimaAtualizacaoNcm
		err = consulta.respotoryNCM.Update(resposta)
	} else {
		err = consulta.respotoryNCM.Create(&ncm)
	}

	for _, ncm := range ncm.Nomenclaturas {
		err = consulta.repositoryNomenclatura.Create(&ncm)
		if err != nil {
			break
		}
	}
	return err
}

func (consulta *consultaNCM) listaNomenclatura(listaNCM []ConsultaNCMSefaz.Nomenclatura) ([]NCM.NomenclaturaBanco, error) {
	var lista []NCM.NomenclaturaBanco
	for _, nomenclatura := range listaNCM {
		dataInicial, err := time.Parse(consulta.modeloData, nomenclatura.DataInicio)
		if err != nil {
			return nil, err
		}
		dataFinal, err := time.Parse(consulta.modeloData, nomenclatura.DataFim)
		if err != nil {
			return nil, err
		}
		lista = append(lista, NCM.NomenclaturaBanco{
			Codigo:     nomenclatura.Codigo,
			Descricao:  nomenclatura.Descricao,
			DataInicio: dataInicial,
			DataFim:    dataFinal,
			TipoAto:    nomenclatura.TipoAto,
			NumeroAto:  nomenclatura.NumeroAto,
			AnoAto:     nomenclatura.AnoAto,
		})
	}
	return lista, nil
}
