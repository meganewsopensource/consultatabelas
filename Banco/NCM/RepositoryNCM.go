package NCM

import (
	"ConsultaTabelas/Banco"
)

type IRepositoryNCM interface {
	Create(ncm *NcmBanco) error
	GetAll() ([]*NcmBanco, error)
	GetByID(id uint) (*NcmBanco, error)
	Update(ncm *NcmBanco) error
	Delete(ncm *NcmBanco) error
}

type repositoryNCM struct {
	db Banco.IConexaoBanco
}

func NewRepositoryNCM(db Banco.IConexaoBanco) (IRepositoryNCM, error) {
	err := db.Conexao().AutoMigrate(&NomenclaturaBanco{})
	if err != nil {
		return nil, err
	}

	err = db.Conexao().AutoMigrate(&NcmBanco{})
	if err != nil {
		return nil, err
	}
	return &repositoryNCM{db}, nil
}

func (repository *repositoryNCM) Create(ncm *NcmBanco) error {
	return repository.db.Conexao().Create(ncm).Error
}

func (repository *repositoryNCM) GetAll() ([]*NcmBanco, error) {
	var listaNcm []*NcmBanco
	err := repository.db.Conexao().Find(&listaNcm).Error
	if err != nil {
		var lista []*NcmBanco
		return lista, err
	}
	return listaNcm, nil
}

func (repository *repositoryNCM) GetByID(id uint) (*NcmBanco, error) {
	var ncmSelecionado NcmBanco
	err := repository.db.Conexao().Find(&ncmSelecionado, id).Error
	if err != nil {
		ncm := NcmBanco{}
		return &ncm, err
	}
	return &ncmSelecionado, nil
}

func (repository *repositoryNCM) Update(ncm *NcmBanco) error {
	return repository.db.Conexao().Updates(ncm).Error
}

func (repository *repositoryNCM) Delete(ncm *NcmBanco) error {
	return repository.db.Conexao().Delete(ncm).Error
}
