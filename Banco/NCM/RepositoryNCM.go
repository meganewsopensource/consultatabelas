package NCM

import (
	"gorm.io/gorm"
)

type IRepositoryNCM interface {
	Create(ncm *NcmBanco) error
	GetAll() ([]*NcmBanco, error)
	GetByID(id uint) (*NcmBanco, error)
	Update(ncm *NcmBanco) error
	Delete(ncm *NcmBanco) error
}

type repositoryNCM struct {
	db *gorm.DB
}

func NewRepositoryNCM(db *gorm.DB) IRepositoryNCM {
	return &repositoryNCM{db}
}

func (repository *repositoryNCM) Create(ncm *NcmBanco) error {
	return repository.db.Create(ncm).Error
}

func (repository *repositoryNCM) GetAll() ([]*NcmBanco, error) {
	var listaNcm []*NcmBanco
	err := repository.db.Find(&listaNcm).Error
	return listaNcm, err
}

func (repository *repositoryNCM) GetByID(id uint) (*NcmBanco, error) {
	var ncmSelecionado NcmBanco
	err := repository.db.Find(&ncmSelecionado, id).Error
	return &ncmSelecionado, err
}

func (repository *repositoryNCM) Update(ncm *NcmBanco) error {
	return repository.db.Updates(ncm).Error
}

func (repository *repositoryNCM) Delete(ncm *NcmBanco) error {
	return repository.db.Delete(ncm).Error
}
