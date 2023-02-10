package NCM

import (
	"gorm.io/gorm"
)

type repositoryNomenclatura struct {
	db *gorm.DB
}

type IRepositoryNomenclatura interface {
	Create(ncm *NomenclaturaBanco) error
	GetAll() ([]*NomenclaturaBanco, error)
	GetByCodigo(codigo string) (*NomenclaturaBanco, error)
	Update(ncm *NomenclaturaBanco) error
	Delete(ncm *NomenclaturaBanco) error
}

func NewRepositoryNomenclatura(db *gorm.DB) IRepositoryNomenclatura {
	return &repositoryNomenclatura{db}
}

func (repository *repositoryNomenclatura) Create(ncm *NomenclaturaBanco) error {
	return repository.db.Create(ncm).Error
}

func (repository *repositoryNomenclatura) GetAll() ([]*NomenclaturaBanco, error) {
	var listaNomenclaturas []*NomenclaturaBanco
	err := repository.db.Find(&listaNomenclaturas).Error
	return listaNomenclaturas, err
}

func (repository *repositoryNomenclatura) GetByCodigo(codigo string) (*NomenclaturaBanco, error) {
	var nomenclatura NomenclaturaBanco
	err := repository.db.Find(&nomenclatura, codigo).Error
	return &nomenclatura, err
}

func (repository *repositoryNomenclatura) Update(ncm *NomenclaturaBanco) error {
	return repository.db.Updates(ncm).Error
}

func (repository *repositoryNomenclatura) Delete(ncm *NomenclaturaBanco) error {
	return repository.db.Delete(ncm).Error
}
