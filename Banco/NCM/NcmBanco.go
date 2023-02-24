package NCM

import (
	"time"
)

type NcmBanco struct {
	ID                       uint `gorm:"primaryKey"`
	DataUltimaAtualizacaoNcm time.Time
	Nomenclaturas            []NomenclaturaBanco `gorm:"-"`
}
