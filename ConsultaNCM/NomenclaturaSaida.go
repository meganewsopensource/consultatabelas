package ConsultaNCM

type NomenclaturaSaida struct {
	Codigo                   string `json:"codigo"`
	DataInicio               string `json:"dataInicio"`
	DataFim                  string `json:"dataFim"`
	Descricao                string `json:"descricao"`
	TipoAto                  string `json:"tipoAto"`
	NumeroAto                string `json:"numeroAto"`
	AnoAto                   string `json:"anoAto"`
	DataUltimaAtualizacaoNcm string `json:"dataUltimaAtualizacaoNcm"`
}
