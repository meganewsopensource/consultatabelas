package NCM

import (
	"gorm.io/gorm"
)

type IMigration interface {
	Executar() error
}

type migration struct {
	db *gorm.DB
}

func NewMigration(db *gorm.DB) IMigration {
	return &migration{db}
}

func (migra *migration) Executar() error {
	err := migra.db.AutoMigrate(&NomenclaturaBanco{})
	if err != nil {
		return err
	}

	err = migra.db.AutoMigrate(&NcmBanco{})
	if err != nil {
		return err
	}
	return nil
}
