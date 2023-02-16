package ConsultaNCM

import (
	"ConsultaTabelas/Banco/NCM"
	"ConsultaTabelas/ConsultaHTTP"
	"ConsultaTabelas/ConsultaNCMSefaz"
	"encoding/json"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

func gerarConexaoBanco() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	migracao := NCM.NewMigration(db)
	_ = migracao.Executar()
	return db, err
}

func criarServidor() *httptest.Server {
	ncm := preencheNcmReceita()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonNCM, _ := json.Marshal(ncm)
		fmt.Fprintf(w, string(jsonNCM[:]))
	}))
	return server
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

func preencheNcmReceita() ConsultaNCMSefaz.NcmReceita {
	return ConsultaNCMSefaz.NcmReceita{
		DataUltimaAtualizacaoNcm: "01/01/2023",
		Nomenclaturas:            preencheNomenclaturas(),
	}
}

func preencheNomenclaturas() []ConsultaNCMSefaz.Nomenclatura {
	lista := []ConsultaNCMSefaz.Nomenclatura{}
	lista = append(lista, ConsultaNCMSefaz.Nomenclatura{
		Codigo:     "01",
		Descricao:  "Teste 01",
		DataInicio: "01/01/2023",
		DataFim:    "31/12/2023",
		TipoAto:    "Regex",
		NumeroAto:  "20",
		AnoAto:     "2021",
	},
		ConsultaNCMSefaz.Nomenclatura{
			Codigo:     "02",
			Descricao:  "Teste 02",
			DataInicio: "01/02/2023",
			DataFim:    "28/02/2023",
			TipoAto:    "Regex 2",
			NumeroAto:  "202",
			AnoAto:     "2022",
		})
	return lista
}

func TestNewConsultaNCM(t *testing.T) {
	server := criarServidor()
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)

	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)
	consulta := NewConsultaNCM(consultaSefaz, repositoryNCM, repositoryNomenclatura)

	if got := NewConsultaNCM(consultaSefaz, repositoryNCM, repositoryNomenclatura); !reflect.DeepEqual(got, consulta) {
		t.Errorf("ConsultaNCM() = %v, diferente de %v", got, consulta)
	}
}

func Test_consultaNCM_AtualizarNCM(t *testing.T) {
	server := criarServidor()
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}

	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)

	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)
	consulta := NewConsultaNCM(consultaSefaz, repositoryNCM, repositoryNomenclatura)

	err = consulta.AtualizarNCM()
	if err != nil {
		t.Errorf("AtualizarNCM() error : %v", err)
	}
}

func Test_consultaNCM_AtualizarNCM_ConsultaSefazFail(t *testing.T) {
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}

	defer func() {
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New("teste123")
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)

	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)
	consulta := NewConsultaNCM(consultaSefaz, repositoryNCM, repositoryNomenclatura)

	err = consulta.AtualizarNCM()
	if err == nil {
		t.Errorf("O erro esperado ao tentar atualizar os NCMs não ocorreu. ")
	}
}

func Test_consultaNCM_AtualizarNCM_BuscaNCMBancoFail(t *testing.T) {
	server := criarServidor()
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}

	mig := db.Migrator()
	err = mig.DropTable(&NCM.NcmBanco{})
	if err != nil {
		t.Errorf("Ocorreu um erro ao apagar a tabela : %v", err)
	}

	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)

	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)
	consulta := NewConsultaNCM(consultaSefaz, repositoryNCM, repositoryNomenclatura)

	err = consulta.AtualizarNCM()
	if err == nil {
		t.Errorf("O erro esperado ao tentar atualizar os NCMs não ocorreu. ")
	}
}

func Test_consultaNCM_AtualizarNCMDataParseFail(t *testing.T) {
	ncm := preencheNcmReceita()
	ncm.DataUltimaAtualizacaoNcm = "2000/31/02"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonNCM, _ := json.Marshal(ncm)
		fmt.Fprintf(w, string(jsonNCM[:]))
	}))

	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}

	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)

	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)
	consulta := NewConsultaNCM(consultaSefaz, repositoryNCM, repositoryNomenclatura)

	err = consulta.AtualizarNCM()
	if err == nil {
		t.Errorf("O erro esperado ao tentar atualizar os NCMs não ocorreu. ")
	}
}

func Test_consultaNCM_AtualizarNomenclaturaDataInicioParseFail(t *testing.T) {
	ncm := preencheNcmReceita()
	ncm.Nomenclaturas[0].DataInicio = "2023/35/01"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonNCM, _ := json.Marshal(ncm)
		fmt.Fprintf(w, string(jsonNCM[:]))
	}))

	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}

	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)

	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)
	consulta := NewConsultaNCM(consultaSefaz, repositoryNCM, repositoryNomenclatura)

	err = consulta.AtualizarNCM()
	if err == nil {
		t.Errorf("O erro esperado ao tentar atualizar os NCMs não ocorreu. ")
	}
}

func Test_consultaNCM_AtualizarNomenclaturaDataFimParseFail(t *testing.T) {
	ncm := preencheNcmReceita()
	ncm.Nomenclaturas[0].DataFim = "2023/35/01"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonNCM, _ := json.Marshal(ncm)
		fmt.Fprintf(w, string(jsonNCM[:]))
	}))

	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}

	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)

	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)
	consulta := NewConsultaNCM(consultaSefaz, repositoryNCM, repositoryNomenclatura)

	err = consulta.AtualizarNCM()
	if err == nil {
		t.Errorf("O erro esperado ao tentar atualizar os NCMs não ocorreu. ")
	}
}

func Test_consultaNCM_gravarNCM(t *testing.T) {
	server := criarServidor()
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)
	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)

	consulta := &consultaNCM{
		consultaSefaz:          consultaSefaz,
		respotoryNCM:           repositoryNCM,
		repositoryNomenclatura: repositoryNomenclatura,
		modeloData:             "02/01/2006",
	}
	ncm := preencheNcmReceita()
	lista, err := consulta.listaNomenclatura(ncm)
	if err != nil {
		t.Errorf("Ocorreu um erro ao preencher lista de nomeclaturas : %v", err)
	}
	data, _ := time.Parse(consulta.modeloData, "01/01/2023")
	ncmBanco := NCM.NcmBanco{
		ID:                       1,
		DataUltimaAtualizacaoNcm: data,
		Nomenclaturas:            lista,
	}

	err = consulta.gravarNCM(ncmBanco)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gravar o NCM %v", err)
	}
}

func Test_consultaNCM_gravarNCM_NCMFail(t *testing.T) {
	server := criarServidor()
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)
	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)

	consulta := &consultaNCM{
		consultaSefaz:          consultaSefaz,
		respotoryNCM:           repositoryNCM,
		repositoryNomenclatura: repositoryNomenclatura,
		modeloData:             "02/01/2006",
	}
	ncm := preencheNcmReceita()
	lista, err := consulta.listaNomenclatura(ncm)
	if err != nil {
		t.Errorf("Ocorreu um erro ao preencher lista de nomeclaturas : %v", err)
	}
	data, _ := time.Parse(consulta.modeloData, "01/01/2023")
	ncmBanco := NCM.NcmBanco{
		ID:                       1,
		DataUltimaAtualizacaoNcm: data,
		Nomenclaturas:            lista,
	}

	mig := db.Migrator()
	err = mig.DropTable(&NCM.NcmBanco{})
	if err != nil {
		t.Errorf("Ocorreu um erro ao apagar tabela %v", err)
	}

	err = consulta.gravarNCM(ncmBanco)
	if err == nil {
		t.Errorf("Não ocorreu o erro esperado ao gravar o NCM %v", err)
	}
}

func Test_consultaNCM_gravarNCM_NomenclaturaUpdateFail(t *testing.T) {
	server := criarServidor()
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)
	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)

	consulta := &consultaNCM{
		consultaSefaz:          consultaSefaz,
		respotoryNCM:           repositoryNCM,
		repositoryNomenclatura: repositoryNomenclatura,
		modeloData:             "02/01/2006",
	}
	ncm := preencheNcmReceita()
	lista, err := consulta.listaNomenclatura(ncm)
	if err != nil {
		t.Errorf("Ocorreu um erro ao preencher lista de nomeclaturas : %v", err)
	}
	data, _ := time.Parse(consulta.modeloData, "01/01/2023")
	ncmBanco := NCM.NcmBanco{
		ID:                       1,
		DataUltimaAtualizacaoNcm: data,
		Nomenclaturas:            lista,
	}

	_ = consulta.gravarNCM(ncmBanco)

	mig := db.Migrator()
	err = mig.DropTable(&NCM.NomenclaturaBanco{})
	if err != nil {
		t.Errorf("Ocorreu um erro ao apagar tabela %v", err)
	}

	err = consulta.gravarNCM(ncmBanco)
	if err == nil {
		t.Errorf("Não ocorreu o erro esperado ao gravar o NCM %v", err)
	}
}

func Test_consultaNCM_gravarNCM_NomenclaturaCreateFail(t *testing.T) {
	server := criarServidor()
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)
	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)

	consulta := &consultaNCM{
		consultaSefaz:          consultaSefaz,
		respotoryNCM:           repositoryNCM,
		repositoryNomenclatura: repositoryNomenclatura,
		modeloData:             "02/01/2006",
	}
	ncm := preencheNcmReceita()
	lista, err := consulta.listaNomenclatura(ncm)
	if err != nil {
		t.Errorf("Ocorreu um erro ao preencher lista de nomeclaturas : %v", err)
	}
	data, _ := time.Parse(consulta.modeloData, "01/01/2023")
	ncmBanco := NCM.NcmBanco{
		ID:                       0,
		DataUltimaAtualizacaoNcm: data,
		Nomenclaturas:            lista,
	}

	mig := db.Migrator()
	err = mig.DropTable(&NCM.NomenclaturaBanco{})
	if err != nil {
		t.Errorf("Ocorreu um erro ao apagar tabela %v", err)
	}

	err = consulta.gravarNCM(ncmBanco)
	if err == nil {
		t.Errorf("Não ocorreu o erro esperado ao gravar o NCM %v", err)
	}
}

func Test_consultaNCM_listaNomenclaturaDataInicialErrada(t *testing.T) {
	server := criarServidor()
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)

	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)

	consulta := &consultaNCM{
		consultaSefaz:          consultaSefaz,
		respotoryNCM:           repositoryNCM,
		repositoryNomenclatura: repositoryNomenclatura,
		modeloData:             "01/02/2006",
	}

	listaNcm := []ConsultaNCMSefaz.Nomenclatura{}
	listaNcm = append(listaNcm, ConsultaNCMSefaz.Nomenclatura{
		Codigo:     "01",
		Descricao:  "teste 01",
		DataInicio: "2023/01/02",
		DataFim:    "01/02/2023",
		TipoAto:    "Teste de ato",
		NumeroAto:  "25",
		AnoAto:     "2021",
	})

	ncm := ConsultaNCMSefaz.NcmReceita{
		DataUltimaAtualizacaoNcm: "01/01/2023",
		Nomenclaturas:            listaNcm,
	}

	lista, err := consulta.listaNomenclatura(ncm)
	if err == nil {
		t.Errorf("Não ocorreu um erro ao preencher a lista")
	}

	if len(lista) > 0 {
		t.Errorf("Valor da lista retornado não é válido : %v", err)
	}
}

func Test_consultaNCM_listaNomenclaturaDataFinalErrada(t *testing.T) {
	server := criarServidor()
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)

	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)

	consulta := &consultaNCM{
		consultaSefaz:          consultaSefaz,
		respotoryNCM:           repositoryNCM,
		repositoryNomenclatura: repositoryNomenclatura,
		modeloData:             "01/02/2006",
	}

	listaNcm := []ConsultaNCMSefaz.Nomenclatura{}
	listaNcm = append(listaNcm, ConsultaNCMSefaz.Nomenclatura{
		Codigo:     "01",
		Descricao:  "teste 01",
		DataInicio: "01/02/2023",
		DataFim:    "2026/02/01",
		TipoAto:    "Teste de ato",
		NumeroAto:  "25",
		AnoAto:     "2021",
	})

	ncm := ConsultaNCMSefaz.NcmReceita{
		DataUltimaAtualizacaoNcm: "01/01/2023",
		Nomenclaturas:            listaNcm,
	}

	_, err = consulta.listaNomenclatura(ncm)
	if err == nil {
		t.Errorf("Não ocorreu um erro ao preencher a lista")
	}
}

func Test_consultaNCM_listaNomenclatura(t *testing.T) {
	server := criarServidor()
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)

	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)

	consulta := &consultaNCM{
		consultaSefaz:          consultaSefaz,
		respotoryNCM:           repositoryNCM,
		repositoryNomenclatura: repositoryNomenclatura,
		modeloData:             "02/01/2006",
	}

	listaNcm := []ConsultaNCMSefaz.Nomenclatura{}
	listaNcm = append(listaNcm, ConsultaNCMSefaz.Nomenclatura{
		Codigo:     "01",
		Descricao:  "teste 01",
		DataInicio: "01/02/2023",
		DataFim:    "31/12/2023",
		TipoAto:    "Teste de ato",
		NumeroAto:  "25",
		AnoAto:     "2021",
	})

	ncm := ConsultaNCMSefaz.NcmReceita{
		DataUltimaAtualizacaoNcm: "01/01/2023",
		Nomenclaturas:            listaNcm,
	}
	lista, err := consulta.listaNomenclatura(ncm)
	if err != nil {
		t.Errorf("Ocorreu um erro ao listar %v", err)
	}

	if len(lista) != 1 {
		t.Errorf("A quantidade de registros retornados não é igual a um. ")
	}
}

func Test_ListarNCMs(t *testing.T) {
	server := criarServidor()
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)
	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)

	consulta := &consultaNCM{
		consultaSefaz:          consultaSefaz,
		respotoryNCM:           repositoryNCM,
		repositoryNomenclatura: repositoryNomenclatura,
		modeloData:             "02/01/2006",
	}
	ncm := preencheNcmReceita()
	lista, err := consulta.listaNomenclatura(ncm)
	if err != nil {
		t.Errorf("Ocorreu um erro ao preencher lista de nomeclaturas : %v", err)
	}
	data, _ := time.Parse(consulta.modeloData, "01/01/2023")
	ncmBanco := NCM.NcmBanco{
		ID:                       1,
		DataUltimaAtualizacaoNcm: data,
		Nomenclaturas:            lista,
	}

	err = consulta.gravarNCM(ncmBanco)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gravar o NCM %v", err)
	}

	ncmLista, err := consulta.ListarNCMs()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gravar o NCM %v", err)
	}

	if len(ncmLista) != len(lista) {
		t.Errorf("O valor de retorno da listagem não é igual a %v!", len(lista))
	}
}

func Test_ListarNCMs_Fail(t *testing.T) {
	server := criarServidor()
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)
	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)

	consulta := &consultaNCM{
		consultaSefaz:          consultaSefaz,
		respotoryNCM:           repositoryNCM,
		repositoryNomenclatura: repositoryNomenclatura,
		modeloData:             "02/01/2006",
	}
	ncm := preencheNcmReceita()
	lista, err := consulta.listaNomenclatura(ncm)
	if err != nil {
		t.Errorf("Ocorreu um erro ao preencher lista de nomeclaturas : %v", err)
	}
	data, _ := time.Parse(consulta.modeloData, "01/01/2023")
	ncmBanco := NCM.NcmBanco{
		ID:                       1,
		DataUltimaAtualizacaoNcm: data,
		Nomenclaturas:            lista,
	}

	err = consulta.gravarNCM(ncmBanco)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gravar o NCM %v", err)
	}

	mig := db.Migrator()
	err = mig.DropTable(&NCM.NomenclaturaBanco{})
	if err != nil {
		t.Errorf("Ocorreu um erro tentar apagar a tabela Nomenclatura : %v", err)
	}

	_, err = consulta.ListarNCMs()
	if err == nil {
		t.Errorf("Não ocorreu um erro ao buscar a lista")
	}
}

func Test_UltimaAtualizacao(t *testing.T) {
	server := criarServidor()
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)
	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)

	consulta := &consultaNCM{
		consultaSefaz:          consultaSefaz,
		respotoryNCM:           repositoryNCM,
		repositoryNomenclatura: repositoryNomenclatura,
		modeloData:             "02/01/2006",
	}
	ncm := preencheNcmReceita()
	lista, err := consulta.listaNomenclatura(ncm)
	if err != nil {
		t.Errorf("Ocorreu um erro ao preencher lista de nomeclaturas : %v", err)
	}
	data, _ := time.Parse(consulta.modeloData, "01/01/2023")
	ncmBanco := NCM.NcmBanco{
		ID:                       1,
		DataUltimaAtualizacaoNcm: data,
		Nomenclaturas:            lista,
	}

	err = consulta.gravarNCM(ncmBanco)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gravar o NCM %v", err)
	}

	_, err = consulta.UltimaAtualizacao()
	if err != nil {
		t.Errorf("Ocorreu um erro ao buscar a última atualização! %v", err)
	}
}

func Test_UltimaAtualizacao_Fail(t *testing.T) {
	server := criarServidor()
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)
	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)

	consulta := &consultaNCM{
		consultaSefaz:          consultaSefaz,
		respotoryNCM:           repositoryNCM,
		repositoryNomenclatura: repositoryNomenclatura,
		modeloData:             "02/01/2006",
	}
	ncm := preencheNcmReceita()
	lista, err := consulta.listaNomenclatura(ncm)
	if err != nil {
		t.Errorf("Ocorreu um erro ao preencher lista de nomeclaturas : %v", err)
	}
	data, _ := time.Parse(consulta.modeloData, "01/01/2023")
	ncmBanco := NCM.NcmBanco{
		ID:                       1,
		DataUltimaAtualizacaoNcm: data,
		Nomenclaturas:            lista,
	}

	err = consulta.gravarNCM(ncmBanco)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gravar o NCM %v", err)
	}

	mig := db.Migrator()
	err = mig.DropTable(&NCM.NcmBanco{})
	if err != nil {
		t.Errorf("Ocorreu um erro tentar apagar a tabela Nomenclatura : %v", err)
	}

	_, err = consulta.UltimaAtualizacao()
	if err == nil {
		t.Errorf("Não ocorreu um erro ao buscar a lista")
	}
}

func Test_ListarNCMPorData(t *testing.T) {
	server := criarServidor()
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)
	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)

	consulta := &consultaNCM{
		consultaSefaz:          consultaSefaz,
		respotoryNCM:           repositoryNCM,
		repositoryNomenclatura: repositoryNomenclatura,
		modeloData:             "02/01/2006",
	}
	ncm := preencheNcmReceita()
	lista, err := consulta.listaNomenclatura(ncm)
	if err != nil {
		t.Errorf("Ocorreu um erro ao preencher lista de nomeclaturas : %v", err)
	}
	data, _ := time.Parse(consulta.modeloData, "01/01/2023")
	ncmBanco := NCM.NcmBanco{
		ID:                       1,
		DataUltimaAtualizacaoNcm: data,
		Nomenclaturas:            lista,
	}

	err = consulta.gravarNCM(ncmBanco)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gravar o NCM %v", err)
	}

	listaNCM, err := consulta.ListarNCMPorData("01/01/2023")
	if err != nil {
		t.Errorf("Ocorreu um erro ao buscar a lista")
	}

	if len(listaNCM) == 0 {
		t.Errorf("O valor da listagem não deve ser zero!")
	}
}

func Test_ListarNCMPorData_ParseFail(t *testing.T) {
	server := criarServidor()
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)
	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)

	consulta := &consultaNCM{
		consultaSefaz:          consultaSefaz,
		respotoryNCM:           repositoryNCM,
		repositoryNomenclatura: repositoryNomenclatura,
		modeloData:             "02/01/2006",
	}
	ncm := preencheNcmReceita()
	lista, err := consulta.listaNomenclatura(ncm)
	if err != nil {
		t.Errorf("Ocorreu um erro ao preencher lista de nomeclaturas : %v", err)
	}
	data, _ := time.Parse(consulta.modeloData, "01/01/2023")
	ncmBanco := NCM.NcmBanco{
		ID:                       1,
		DataUltimaAtualizacaoNcm: data,
		Nomenclaturas:            lista,
	}

	err = consulta.gravarNCM(ncmBanco)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gravar o NCM %v", err)
	}

	_, err = consulta.ListarNCMPorData("01/31/2023")
	if err == nil {
		t.Errorf("Não ocorreu um erro ao buscar a lista relacionado a formação de data")
	}
}

func Test_ListarNCMPorData_BuscaFail(t *testing.T) {
	server := criarServidor()
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		deletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)
	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)

	consulta := &consultaNCM{
		consultaSefaz:          consultaSefaz,
		respotoryNCM:           repositoryNCM,
		repositoryNomenclatura: repositoryNomenclatura,
		modeloData:             "02/01/2006",
	}
	ncm := preencheNcmReceita()
	lista, err := consulta.listaNomenclatura(ncm)
	if err != nil {
		t.Errorf("Ocorreu um erro ao preencher lista de nomeclaturas : %v", err)
	}
	data, _ := time.Parse(consulta.modeloData, "01/01/2023")
	ncmBanco := NCM.NcmBanco{
		ID:                       1,
		DataUltimaAtualizacaoNcm: data,
		Nomenclaturas:            lista,
	}

	err = consulta.gravarNCM(ncmBanco)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gravar o NCM %v", err)
	}

	mig := db.Migrator()
	err = mig.DropTable(&NCM.NomenclaturaBanco{})

	_, err = consulta.ListarNCMPorData("01/01/2023")
	if err == nil {
		t.Errorf("Não ocorreu um erro ao buscar a lista de NCM")
	}
}
