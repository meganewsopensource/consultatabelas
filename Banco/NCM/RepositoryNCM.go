package NCM

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	transacao := repository.db.Begin()

	err := transacao.Create(ncm).Error
	if err != nil {
		transacao.Rollback()
		return err
	}

	err = transacao.CreateInBatches(ncm.Nomenclaturas, 1000).Error
	if err != nil {
		transacao.Rollback()
		return err
	}

	transacao.Commit()

	return nil
}

func (repository *repositoryNCM) GetAll() (lista []*NcmBanco, err error) {
	err = repository.db.Find(&lista).Error
	return
}

func (repository *repositoryNCM) GetByID(id uint) (ncmSelecionado *NcmBanco, err error) {
	err = repository.db.Find(&ncmSelecionado, id).Error
	return
}

func (repository *repositoryNCM) Update(ncm *NcmBanco) error {
	transacao := repository.db.Begin()
	err := transacao.Updates(ncm).Error
	if err != nil {
		transacao.Rollback()
		return err
	}

	err = transacao.Table("nomenclatura_bancos").Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(ncm.Nomenclaturas, 1000).Error
	if err != nil {
		transacao.Rollback()
		return err
	}

	transacao.Commit()

	return nil
}

func (repository *repositoryNCM) Delete(ncm *NcmBanco) error {
	return repository.db.Delete(ncm).Error
}
