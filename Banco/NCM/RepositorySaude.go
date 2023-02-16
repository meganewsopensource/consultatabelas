package NCM

import "gorm.io/gorm"

type saudeBanco struct {
	db *gorm.DB
}

type IRepositorySaude interface {
	Saudavel() bool
}

func NewRepositorySaude(db *gorm.DB) IRepositorySaude {
	return &saudeBanco{db}
}

func (saude *saudeBanco) Saudavel() bool {
	sqlDB, err := saude.db.DB()
	if err != nil {
		return false
	}
	err = sqlDB.Ping()
	if err != nil {
		return false
	}
	return true
}
