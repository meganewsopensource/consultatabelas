package ConsultaNCM

import (
	"ConsultaTabelas/Banco/NCM"
	"ConsultaTabelas/ConsultaNCMSefaz"
	"regexp"
	"time"
)

type IConsultaNCM interface {
	AtualizarNCM() error
	ListarNCMs() ([]*NCM.NomenclaturaBanco, error)
	UltimaAtualizacao() (time.Time, error)
	ListarNCMPorData(data string) ([]*NCM.NomenclaturaBanco, error)
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
	if err != nil {
		return err
	}

	if data.After(ncm.DataUltimaAtualizacaoNcm) {
		lista, erro := consulta.listaNomenclatura(dadosConsulta)
		if erro != nil {
			return erro
		}
		ncmBanco := NCM.NcmBanco{
			ID:                       ncm.ID,
			DataUltimaAtualizacaoNcm: data,
			Nomenclaturas:            lista,
		}
		erro = consulta.gravarNCM(ncmBanco)
		if erro != nil {
			return erro
		}
	}
	return nil
}

func (consulta *consultaNCM) gravarNCM(ncm NCM.NcmBanco) error {
	var err error
	resposta, err := consulta.respotoryNCM.GetByID(1)
	if err != nil {
		return err
	}

	if resposta.ID == 1 {
		resposta.DataUltimaAtualizacaoNcm = ncm.DataUltimaAtualizacaoNcm
		resposta.Nomenclaturas = ncm.Nomenclaturas
		err = consulta.respotoryNCM.Update(resposta)
		if err != nil {
			return err
		}
	} else {
		ncm.ID = 1
		err = consulta.respotoryNCM.Create(&ncm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (consulta *consultaNCM) listaNomenclatura(listaNCM ConsultaNCMSefaz.NcmReceita) ([]NCM.NomenclaturaBanco, error) {
	var lista []NCM.NomenclaturaBanco
	for _, nomenclatura := range listaNCM.Nomenclaturas {
		dataInicial, err := time.Parse(consulta.modeloData, nomenclatura.DataInicio)
		if err != nil {
			return nil, err
		}
		dataFinal, err := time.Parse(consulta.modeloData, nomenclatura.DataFim)
		if err != nil {
			return nil, err
		}
		dataAtualizacao, _ := time.Parse(consulta.modeloData, listaNCM.DataUltimaAtualizacaoNcm)
		lista = append(lista, NCM.NomenclaturaBanco{
			Codigo:                   nomenclatura.Codigo,
			Descricao:                nomenclatura.Descricao,
			DataInicio:               dataInicial,
			DataFim:                  dataFinal,
			TipoAto:                  nomenclatura.TipoAto,
			NumeroAto:                nomenclatura.NumeroAto,
			AnoAto:                   nomenclatura.AnoAto,
			DataUltimaAtualizacaoNcm: dataAtualizacao,
		})
	}
	return lista, nil
}

func (consulta *consultaNCM) ListarNCMs() ([]*NCM.NomenclaturaBanco, error) {
	lista, err := consulta.repositoryNomenclatura.GetAll()
	if err != nil {
		return nil, err
	}
	return lista, nil
}

func (consulta *consultaNCM) UltimaAtualizacao() (time.Time, error) {
	ncm, err := consulta.respotoryNCM.GetByID(1)
	if err != nil {
		return time.Now(), err
	}
	return ncm.DataUltimaAtualizacaoNcm, nil
}

func (consulta *consultaNCM) ListarNCMPorData(data string) ([]*NCM.NomenclaturaBanco, error) {
	dataConvertida, err := consulta.paraData(data)
	if err != nil {
		return nil, err
	}
	lista, err := consulta.repositoryNomenclatura.GetByData(dataConvertida)
	if err != nil {
		return nil, err
	}
	return lista, nil
}

func (consulta *consultaNCM) paraData(data string) (time.Time, error) {
	r, _ := regexp.Compile("\\D+")
	entrada := []byte(data)
	saida := r.ReplaceAll(entrada, []byte(""))
	dataConvertida, err := time.Parse("02012006", string(saida))
	if err != nil {
		return time.Time{}, err
	}
	return dataConvertida, nil
}
