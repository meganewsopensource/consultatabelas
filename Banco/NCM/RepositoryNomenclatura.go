package NCM

import (
	"gorm.io/gorm"
	"time"
)

type repositoryNomenclatura struct {
	db *gorm.DB
}

type IRepositoryNomenclatura interface {
	Create(ncm *NomenclaturaBanco) error
	GetAll() ([]*NomenclaturaBanco, error)
	GetByData(data time.Time) ([]*NomenclaturaBanco, error)
	GetByCodigo(codigo string) (NomenclaturaBanco, error)
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
	err := repository.db.Where("data_fim > ? and  ? >= data_inicio ", time.Now(), time.Now()).Find(&listaNomenclaturas).Error
	return listaNomenclaturas, err
}

func (repository *repositoryNomenclatura) GetByData(data time.Time) ([]*NomenclaturaBanco, error) {
	var nomenclatura []*NomenclaturaBanco
	err := repository.db.Where("data_ultima_atualizacao_ncm = ?", data).Find(&nomenclatura).Error
	return nomenclatura, err
}

func (repository *repositoryNomenclatura) Update(ncm *NomenclaturaBanco) error {
	return repository.db.Updates(ncm).Error
}

func (repository *repositoryNomenclatura) Delete(ncm *NomenclaturaBanco) error {
	return repository.db.Delete(ncm).Error
}

func (repository *repositoryNomenclatura) GetByCodigo(codigo string) (NomenclaturaBanco, error) {
	var ncm NomenclaturaBanco
	err := repository.db.Find(&ncm, codigo).Error
	return ncm, err
}
