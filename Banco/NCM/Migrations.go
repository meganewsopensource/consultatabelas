package NCM

import (
	"gorm.io/gorm"
)

type IMigration interface {
	Executar() (err error)
}

type migration struct {
	db *gorm.DB
}

func NewMigration(db *gorm.DB) IMigration {
	return &migration{db}
}

func (migra *migration) Executar() (err error) {
	err = migra.db.AutoMigrate(&NomenclaturaBanco{})
	err = migra.db.AutoMigrate(&NcmBanco{})
	return
}
