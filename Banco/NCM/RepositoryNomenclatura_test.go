package NCM

import (
	"reflect"
	"testing"
	"time"
)

func TestNewRepositoryNomenclatura(t *testing.T) {
	conexao, _ := gerarConexaoBanco()
	repository := NewRepositoryNomenclatura(conexao)

	if got := NewRepositoryNomenclatura(conexao); !reflect.DeepEqual(got, repository) {
		t.Errorf("NewRepositoryNomenclatura() = %v, want %v", got, repository)
	}
	deletarBanco()
}

func CriarNomenclatura() NomenclaturaBanco {
	return NomenclaturaBanco{
		Codigo:     "01",
		DataInicio: time.Time{},
		DataFim:    time.Time{},
		Descricao:  "teste 01",
		TipoAto:    "2023",
		NumeroAto:  "1",
		AnoAto:     "2021",
	}
}

func Test_repositoryNomenclatura_Create(t *testing.T) {
	conexao, _ := gerarConexaoBanco()
	repository := &repositoryNomenclatura{
		conexao,
	}
	nomenclatura := CriarNomenclatura()
	err := repository.Create(&nomenclatura)
	if err != nil {
		t.Errorf("Create() error = %v ", err)
	}
	deletarBanco()
}

func Test_repositoryNomenclatura_Delete(t *testing.T) {
	conexao, _ := gerarConexaoBanco()
	nomenclatura := CriarNomenclatura()
	repository := &repositoryNomenclatura{
		conexao,
	}
	_ = repository.Create(&nomenclatura)
	err := repository.Delete(&nomenclatura)
	if err != nil {
		t.Errorf("Delete() error = %v ", err)
	}
	deletarBanco()
}

func Test_repositoryNomenclatura_GetAll(t *testing.T) {
	conexao, _ := gerarConexaoBanco()
	nomenclatura1 := CriarNomenclatura()
	nomenclatura2 := CriarNomenclatura()
	nomenclatura2.Codigo = "02"
	repository := &repositoryNomenclatura{conexao}
	_ = repository.Create(&nomenclatura1)
	_ = repository.Create(&nomenclatura2)
	lista, err := repository.GetAll()
	if err != nil {
		t.Errorf("GetAll() error = %v ", err)
	}
	if len(lista) != 2 {
		t.Errorf("A quantidade retornada não é  = %v ", err)
	}
	deletarBanco()
}

func Test_repositoryNomenclatura_GetByCodigo(t *testing.T) {
	conexao, _ := gerarConexaoBanco()
	nomenclatura := CriarNomenclatura()
	repository := &repositoryNomenclatura{
		conexao,
	}
	_ = repository.Create(&nomenclatura)

	got, err := repository.GetByCodigo(nomenclatura.Codigo)
	if err != nil {
		t.Errorf("GetByCodigo() error = %v ", err)
	}
	if (got.Codigo == nomenclatura.Codigo) && (got.AnoAto == nomenclatura.AnoAto) {
		t.Errorf("Nomenclatura encontrada %v diferente da gravada %v", got, nomenclatura)
	}
	deletarBanco()
}

func Test_repositoryNomenclatura_Update(t *testing.T) {
	conexao, _ := gerarConexaoBanco()
	nomenclatura := CriarNomenclatura()
	repository := &repositoryNomenclatura{
		conexao,
	}
	_ = repository.Create(&nomenclatura)
	nomenclatura.AnoAto = "2024"
	err := repository.Update(&nomenclatura)
	if err != nil {
		t.Errorf("Update() error = %v ", err)
	}
	nomenclaturaReturno, _ := repository.GetByCodigo(nomenclatura.Codigo)

	if nomenclaturaReturno.AnoAto != nomenclatura.AnoAto {
		t.Errorf("Valor esperado %v, recebido %v", nomenclatura.AnoAto, nomenclaturaReturno.AnoAto)
	}
	deletarBanco()
}
