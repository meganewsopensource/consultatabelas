package ConsultaNCM

type NcmReceita struct {
	DataUltimaAtualizacaoNcm string         `json:"Data_Ultima_Atualizacao_NCM"`
	Nomenclaturas            []Nomenclatura `json:"Nomenclaturas"`
}
