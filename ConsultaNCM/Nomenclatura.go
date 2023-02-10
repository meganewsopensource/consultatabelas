package ConsultaNCM

type Nomenclatura struct {
	Codigo     string `json:"Codigo"`
	Descricao  string `json:"Descricao"`
	DataInicio string `json:"Data_Inicio"`
	DataFim    string `json:"Data_Fim"`
	TipoAto    string `json:"Tipo_Ato"`
	NumeroAto  string `json:"Numero_Ato"`
	AnoAto     string `json:"Ano_Ato"`
}
