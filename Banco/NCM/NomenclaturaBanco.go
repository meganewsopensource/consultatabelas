package NCM

import "time"

type NomenclaturaBanco struct {
	Codigo     string    `gorm:"primaryKey;autoIncrement:false"`
	DataInicio time.Time `gorm:"primaryKey;autoIncrement:false"`
	DataFim    time.Time `gorm:"primaryKey;autoIncrement:false"`
	Descricao  string
	TipoAto    string
	NumeroAto  string
	AnoAto     string
}
