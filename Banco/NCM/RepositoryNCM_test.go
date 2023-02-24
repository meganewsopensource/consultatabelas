package NCM

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"reflect"
	"testing"
	"time"
)

func gerarConexaoBanco() (*gorm.DB, error) {
	db, err := gerarConexaoBancoSemMigration()
	migracao := NewMigration(db)
	migracao.Executar()
	return db, err
}

func gerarConexaoBancoSemMigration() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
}

func deletarBanco(db *gorm.DB) {
	sqlDB, err := db.DB()
	err = sqlDB.Close()
	if err != nil {
		panic(err)
	}
	err = os.Remove("gorm.db")
	if err != nil {
		panic(err)
	}
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
	conexao, _ := gerarConexaoBanco()
	repository := NewRepositoryNCM(conexao)
	if got := NewRepositoryNCM(conexao); !reflect.DeepEqual(got, repository) {
		t.Errorf("Esperado %v, recebi %v", got, repository)
	}
	deletarBanco(conexao)
}

func Test_repositoryNCM_Create(t *testing.T) {
	conexao, _ := gerarConexaoBanco()
	repository := NewRepositoryNCM(conexao)
	ncmGravar := gerarNCMBanco()
	err := repository.Create(&ncmGravar)
	if err != nil {
		t.Errorf("Ocorreu um erro ao tentar gravar %v", err)
	}

	deletarBanco(conexao)
}

func Test_repositoryNCM_Create_FailNCM(t *testing.T) {
	conexao, _ := gerarConexaoBancoSemMigration()
	repository := NewRepositoryNCM(conexao)
	var ncmGravar NcmBanco
	ncmGravar = NcmBanco{}
	err := repository.Create(&ncmGravar)
	if err == nil {
		t.Errorf("Não ocorreu um erro esperado ao tentar gravar %v", err)
	}

	deletarBanco(conexao)
}

func Test_repositoryNCM_Create_FailNomenclatura(t *testing.T) {
	conexao, _ := gerarConexaoBanco()
	mig := conexao.Migrator()
	err := mig.DropTable(&NomenclaturaBanco{})
	if err != nil {
		t.Errorf("Ocorreu um erro esperado ao tentar migrar o banco %v", err)
	}

	repository := NewRepositoryNCM(conexao)
	ncmGravar := gerarNCMBanco()
	err = repository.Create(&ncmGravar)
	if err == nil {
		t.Errorf("Não ocorreu um erro esperado ao tentar gravar %v", err)
	}

	deletarBanco(conexao)
}

func Test_repositoryNCM_Delete(t *testing.T) {
	conexao, _ := gerarConexaoBanco()
	repository := NewRepositoryNCM(conexao)
	ncmGravar := gerarNCMBanco()
	repository.Create(&ncmGravar)
	err := repository.Delete(&ncmGravar)
	if err != nil {
		t.Errorf("Ocorreu um erro ao tentar deletar %v", err)
	}
	deletarBanco(conexao)
}

func Test_repositoryNCM_GetAll(t *testing.T) {
	conexao, _ := gerarConexaoBanco()
	repository := NewRepositoryNCM(conexao)
	ncmGravar := gerarNCMBanco()
	repository.Create(&ncmGravar)
	lista, err := repository.GetAll()
	if err != nil {
		t.Errorf("Ocorreu um erro ao tentar listar %v", err)
	}
	if len(lista) == 0 {
		t.Errorf("Erro, a lista está vazia")
	}
	deletarBanco(conexao)
}

func Test_repositoryNCM_GetByCodigo(t *testing.T) {
	conexao, _ := gerarConexaoBanco()
	repository := NewRepositoryNCM(conexao)
	ncmGravar := gerarNCMBanco()
	repository.Create(&ncmGravar)
	ncm, err := repository.GetByID(ncmGravar.ID)
	if err != nil {
		t.Errorf("Ocorreu um erro ao tentar listar %v", err)
	}

	if !reflect.DeepEqual(ncmGravar.ID, ncm.ID) && !reflect.DeepEqual(ncmGravar.DataUltimaAtualizacaoNcm, ncm.DataUltimaAtualizacaoNcm) {
		t.Errorf("Esperado %v, recebi %v", ncmGravar, ncm)
	}
	deletarBanco(conexao)
}

func Test_repositoryNCM_Update(t *testing.T) {
	conexao, _ := gerarConexaoBanco()
	repository := NewRepositoryNCM(conexao)
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
	deletarBanco(conexao)
}

func Test_repositoryNCM_Update_FailNCM(t *testing.T) {
	conexao, _ := gerarConexaoBanco()
	repository := NewRepositoryNCM(conexao)
	ncmGravar := gerarNCMBanco()
	repository.Create(&ncmGravar)
	data := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	ncmUpdate := ncmGravar
	ncmUpdate.DataUltimaAtualizacaoNcm = data
	mig := conexao.Migrator()
	err := mig.DropTable(&NcmBanco{})
	if err != nil {
		t.Errorf("Ocorreu um erro esperado ao tentar migrar o banco %v", err)
	}

	err = repository.Update(&ncmUpdate)
	if err == nil {
		t.Errorf("Ocorreu um erro ao tentar atualizar NcmBanco %v", err)
	}

	deletarBanco(conexao)
}

func Test_repositoryNCM_Update_FailNomenclaturas(t *testing.T) {
	conexao, _ := gerarConexaoBanco()
	repository := NewRepositoryNCM(conexao)
	ncmGravar := gerarNCMBanco()
	repository.Create(&ncmGravar)
	data := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	ncmUpdate := ncmGravar
	ncmUpdate.DataUltimaAtualizacaoNcm = data
	mig := conexao.Migrator()
	err := mig.DropTable(&NomenclaturaBanco{})
	if err != nil {
		t.Errorf("Ocorreu um erro esperado ao tentar migrar o banco %v", err)
	}

	err = repository.Update(&ncmUpdate)
	if err == nil {
		t.Errorf("Não ocorreu o erro esperado ao tentar atualizar Nomenclaturas %v", err)
	}

	deletarBanco(conexao)
}
