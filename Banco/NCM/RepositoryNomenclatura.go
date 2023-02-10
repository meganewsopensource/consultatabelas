package NCM

import "ConsultaTabelas/Banco"

type repositoryNomenclatura struct {
	db Banco.IConexaoBanco
}

type IRepositoryNomenclatura interface {
	Create(ncm *NomenclaturaBanco) error
	GetAll() (*[]NomenclaturaBanco, error)
	GetByCodigo(codigo string) (*NomenclaturaBanco, error)
	Update(ncm *NomenclaturaBanco) error
	Delete(ncm *NcmBanco) error
}

func NewRepositoryNomenclatura(db Banco.IConexaoBanco) IRepositoryNomenclatura {
	return &repositoryNomenclatura{db}
}

func (repository *repositoryNomenclatura) Create(ncm *NomenclaturaBanco) error {
	return repository.db.Conexao().Create(ncm).Error
}

func (repository *repositoryNomenclatura) GetAll() (*[]NomenclaturaBanco, error) {
	var listaNomenclaturas []NomenclaturaBanco
	err := repository.db.Conexao().Find(&listaNomenclaturas).Error
	if err != nil {
		var lista *[]NomenclaturaBanco
		return lista, err
	}
	return &listaNomenclaturas, nil
}

func (repository *repositoryNomenclatura) GetByCodigo(codigo string) (*NomenclaturaBanco, error) {
	var nomenclatura NomenclaturaBanco
	err := repository.db.Conexao().Find(&nomenclatura, codigo).Error
	if err != nil {
		return &NomenclaturaBanco{}, err
	}
	return &nomenclatura, nil
}

func (repository *repositoryNomenclatura) Update(ncm *NomenclaturaBanco) error {
	return repository.db.Conexao().Updates(ncm).Error
}

func (repository *repositoryNomenclatura) Delete(ncm *NcmBanco) error {
	return repository.db.Conexao().Delete(ncm).Error
}
