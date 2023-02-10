package NCM

import (
	"ConsultaTabelas/Banco"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func gerarConexaoBanco() Banco.IConexaoBanco {
	conexao := Banco.NewBanco("postgres://admin:admin@localhost:5432/tabelas")
	conexao.Conectar()
	return conexao
}

func gerarNCMBanco() NcmBanco {
	nomenclaturas := []NomenclaturaBanco{}
	nomenclaturas = append(nomenclaturas,
		NomenclaturaBanco{
			Codigo:     "01",
			Descricao:  "teste de gravação",
			DataInicio: time.Time{},
			DataFim:    time.Time{},
			TipoAto:    "teste de tipo",
			NumeroAto:  "1",
			AnoAto:     "2023",
		})

	ncmGravar := NcmBanco{
		DataUltimaAtualizacaoNcm: time.Time{},
		Nomenclaturas:            nomenclaturas,
	}
	return ncmGravar
}

func TestNew(t *testing.T) {
	conexao := gerarConexaoBanco()
	repository, err := NewRepositoryNCM(conexao)
	if err != nil {
		t.Errorf("Ocorreu um erro ao criar RepositoryNCM : %v", err)
	}
	if got, _ := NewRepositoryNCM(conexao); !reflect.DeepEqual(got, repository) {
		t.Errorf("Esperado %v, recebi %v", got, repository)
	}
}

func TestNewConexaoFail(t *testing.T) {
	conexao := Banco.NewBanco("")
	_, err := NewRepositoryNCM(conexao)
	if err == nil {
		t.Errorf("Problema, não ocorreu um erro ao criar RepositoryNCM.")
	}
}

func Test_repositoryNCM_Create(t *testing.T) {
	conexao := gerarConexaoBanco()
	repository, _ := NewRepositoryNCM(conexao)
	ncmGravar := gerarNCMBanco()
	err := repository.Create(&ncmGravar)
	if err != nil {
		t.Errorf("Ocorreu um erro ao tentar gravar %v", err)
	}
}

func Test_repositoryNCM_Delete(t *testing.T) {
	conexao := gerarConexaoBanco()
	repository, _ := NewRepositoryNCM(conexao)
	ncmGravar := gerarNCMBanco()
	repository.Create(&ncmGravar)
	err := repository.Delete(&ncmGravar)
	if err != nil {
		t.Errorf("Ocorreu um erro ao tentar deletar %v", err)
	}
}

func Test_repositoryNCM_GetAll(t *testing.T) {
	conexao := gerarConexaoBanco()
	repository, _ := NewRepositoryNCM(conexao)
	ncmGravar := gerarNCMBanco()
	repository.Create(&ncmGravar)
	lista, err := repository.GetAll()
	if err != nil {
		t.Errorf("Ocorreu um erro ao tentar listar %v", err)
	}
	if len(lista) == 0 {
		t.Errorf("Erro, a lista está vazia")
	}
}

func Test_repositoryNCM_GetByCodigo(t *testing.T) {
	conexao := gerarConexaoBanco()
	repository, _ := NewRepositoryNCM(conexao)
	ncmGravar := gerarNCMBanco()
	repository.Create(&ncmGravar)
	ncm, err := repository.GetByID(ncmGravar.ID)
	if err != nil {
		t.Errorf("Ocorreu um erro ao tentar listar %v", err)
	}

	if !reflect.DeepEqual(ncmGravar.ID, ncm.ID) && !reflect.DeepEqual(ncmGravar.DataUltimaAtualizacaoNcm, ncm.DataUltimaAtualizacaoNcm) {
		t.Errorf("Esperado %v, recebi %v", ncmGravar, ncm)
	}
}

func Test_repositoryNCM_Update(t *testing.T) {
	conexao := gerarConexaoBanco()
	repository, _ := NewRepositoryNCM(conexao)
	ncmGravar := gerarNCMBanco()
	repository.Create(&ncmGravar)
	data := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	ncmUpdate := ncmGravar
	ncmUpdate.DataUltimaAtualizacaoNcm = data
	err := repository.Update(&ncmUpdate)
	if err != nil {
		t.Errorf("Ocorreu um erro ao tentar atualizar NcmBanco %v", err)
	}
	ncm, _ := repository.GetByID(ncmUpdate.ID)
	fmt.Println(ncm)
	if !reflect.DeepEqual(data, ncm.DataUltimaAtualizacaoNcm) {
		t.Errorf("Ocorreu um erro, a data esperada %s, data recebida %s", data, ncm.DataUltimaAtualizacaoNcm)
	}
}
