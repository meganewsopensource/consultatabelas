package NCM

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestNewRepositorySaude(t *testing.T) {
	conexao, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar a conexão : %v", err)
	}
	repository := NewRepositorySaude(conexao)

	if got := NewRepositorySaude(conexao); !reflect.DeepEqual(got, repository) {
		t.Errorf("NewRepositorySaude() %v diferente de %v", got, repository)
	}

	deletarBanco(conexao)
}

func Test_saudeBanco_Saudavel(t *testing.T) {
	conexao, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar a conexão : %v", err)
	}
	repository := NewRepositorySaude(conexao)

	if !repository.Saudavel() {
		t.Errorf("A conexão com o banco de dados não é saudável : %v", err)
	}
	deletarBanco(conexao)
}

func Test_saudeBanco_Saudavel_FailPing(t *testing.T) {
	conexao, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar a conexão : %v", err)
	}
	repository := NewRepositorySaude(conexao)

	sqlDB, err := conexao.DB()
	err = sqlDB.Close()
	if err != nil {
		t.Errorf("Ocorreu um erro ao desconectar com o banco : %v", err)
	}

	if repository.Saudavel() {
		t.Errorf("Não ocorreu um erro esperado com a conexão com o banco de dados")
	}
	deletarBanco(conexao)
}

func Test_saudeBanco_Saudavel_FailConexao(t *testing.T) {
	conexao, _ := gorm.Open(postgres.Open("teste.db"), &gorm.Config{})
	repository := NewRepositorySaude(conexao)
	if repository.Saudavel() {
		t.Errorf("Não ocorreu um erro esperado com a conexão com o banco de dados. ")
	}
}
