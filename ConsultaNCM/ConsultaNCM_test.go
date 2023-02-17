package ConsultaNCM

import (
	"ConsultaTabelas/Banco/NCM"
	"ConsultaTabelas/ConsultaHTTP"
	"ConsultaTabelas/ConsultaNCMSefaz"
	"ConsultaTabelas/MockTestes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestNewConsultaNCM(t *testing.T) {
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}

	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}

	defer func() {
		MockTestes.DeletarBanco(db)
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
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
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
		MockTestes.DeletarBanco(db)
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

func Test_consultaNCM_AtualizarNCMgravaNCMFail(t *testing.T) {
	ncm := MockTestes.PreencheNcmReceita()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonNCM, _ := json.Marshal(ncm)
		fmt.Fprintf(w, string(jsonNCM[:]))
	}))

	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}

	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
	}()

	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)

	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)
	consulta := NewConsultaNCM(consultaSefaz, repositoryNCM, repositoryNomenclatura)

	mig := db.Migrator()
	err = mig.DropTable(&NCM.NomenclaturaBanco{})
	if err != nil {
		t.Errorf("Ocorreu um erro ao deletar tabela.")
	}

	err = consulta.AtualizarNCM()
	if err == nil {
		t.Errorf("O erro esperado ao tentar atualizar os NCMs não ocorreu. ")
	}
}

func Test_consultaNCM_AtualizarNCMDataParseFail(t *testing.T) {
	ncm := MockTestes.PreencheNcmReceita()
	ncm.DataUltimaAtualizacaoNcm = "2000/31/02"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonNCM, _ := json.Marshal(ncm)
		fmt.Fprintf(w, string(jsonNCM[:]))
	}))

	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}

	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	ncm := MockTestes.PreencheNcmReceita()
	ncm.Nomenclaturas[0].DataInicio = "2023/35/01"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonNCM, _ := json.Marshal(ncm)
		fmt.Fprintf(w, string(jsonNCM[:]))
	}))

	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}

	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	ncm := MockTestes.PreencheNcmReceita()
	ncm.Nomenclaturas[0].DataFim = "2023/35/01"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonNCM, _ := json.Marshal(ncm)
		fmt.Fprintf(w, string(jsonNCM[:]))
	}))

	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}

	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	ncm := MockTestes.PreencheNcmReceita()
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
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	ncm := MockTestes.PreencheNcmReceita()
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
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	ncm := MockTestes.PreencheNcmReceita()
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
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	ncm := MockTestes.PreencheNcmReceita()
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
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	ncm := MockTestes.PreencheNcmReceita()
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
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	ncm := MockTestes.PreencheNcmReceita()
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
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	ncm := MockTestes.PreencheNcmReceita()
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
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	ncm := MockTestes.PreencheNcmReceita()
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
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	ncm := MockTestes.PreencheNcmReceita()
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
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	ncm := MockTestes.PreencheNcmReceita()
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
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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
	ncm := MockTestes.PreencheNcmReceita()
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

func Test_paraData(t *testing.T) {
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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

	data, err := consulta.paraData("01-01-2023")
	if err != nil {
		t.Errorf("Ocorreu um erro ao converter a data. %v", err)
	}
	dataComparacao, err := time.Parse(consulta.modeloData, "01/01/2023")
	if err != nil {
		t.Errorf("Erro no parse de comparação do teste. %v", err)
	}

	if !data.Equal(dataComparacao) {
		t.Errorf("Ocorreu um erro, as datas são divergentes %v <> %v", data, dataComparacao)
	}
}

func Test_paraDataModeloSoNumero(t *testing.T) {
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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

	data, err := consulta.paraData("01012023")
	if err != nil {
		t.Errorf("Ocorreu um erro ao converter a data. %v", err)
	}
	dataComparacao, err := time.Parse(consulta.modeloData, "01/01/2023")
	if err != nil {
		t.Errorf("Erro no parse de comparação do teste. %v", err)
	}

	if !data.Equal(dataComparacao) {
		t.Errorf("Ocorreu um erro, as datas são divergentes %v <> %v", data, dataComparacao)
	}
}

func Test_paraNomenclaturaSaida(t *testing.T) {
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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

	nomenclatura := NCM.NomenclaturaBanco{
		Codigo:                   "01",
		DataInicio:               time.Now(),
		DataFim:                  time.Now(),
		Descricao:                "teste",
		TipoAto:                  "02",
		NumeroAto:                "03",
		AnoAto:                   "04",
		DataUltimaAtualizacaoNcm: time.Now(),
	}

	nomeConvertido := consulta.paraNomenclaturaSaida(nomenclatura)
	if nomeConvertido.Codigo != nomenclatura.Codigo {
		t.Errorf("Ocorreu um erro, os campos AnoAto são divergentes %v <> %v", nomeConvertido.AnoAto, nomenclatura.AnoAto)
	}
	if nomeConvertido.DataInicio != nomenclatura.DataInicio.Format(consulta.modeloData) {
		t.Errorf("Ocorreu um erro, os campos AnoAto são divergentes %v <> %v", nomeConvertido.AnoAto, nomenclatura.AnoAto)
	}
	if nomeConvertido.DataFim != nomenclatura.DataFim.Format(consulta.modeloData) {
		t.Errorf("Ocorreu um erro, os campos AnoAto são divergentes %v <> %v", nomeConvertido.AnoAto, nomenclatura.AnoAto)
	}
	if nomeConvertido.Descricao != nomenclatura.Descricao {
		t.Errorf("Ocorreu um erro, os campos AnoAto são divergentes %v <> %v", nomeConvertido.AnoAto, nomenclatura.AnoAto)
	}
	if nomeConvertido.TipoAto != nomenclatura.TipoAto {
		t.Errorf("Ocorreu um erro, os campos AnoAto são divergentes %v <> %v", nomeConvertido.AnoAto, nomenclatura.AnoAto)
	}
	if nomeConvertido.NumeroAto != nomenclatura.NumeroAto {
		t.Errorf("Ocorreu um erro, os campos AnoAto são divergentes %v <> %v", nomeConvertido.AnoAto, nomenclatura.AnoAto)
	}
	if nomeConvertido.AnoAto != nomenclatura.AnoAto {
		t.Errorf("Ocorreu um erro, os campos AnoAto são divergentes %v <> %v", nomeConvertido.AnoAto, nomenclatura.AnoAto)
	}
	if nomeConvertido.DataUltimaAtualizacaoNcm != nomenclatura.DataUltimaAtualizacaoNcm.Format(consulta.modeloData) {
		t.Errorf("Ocorreu um erro, os campos AnoAto são divergentes %v <> %v", nomeConvertido.AnoAto, nomenclatura.AnoAto)
	}
}

func Test_paraNcmSaida(t *testing.T) {
	server := MockTestes.CriarServidor()
	db, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	defer func() {
		server.Close()
		MockTestes.DeletarBanco(db)
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

	ncm := NCM.NcmBanco{
		ID:                       1,
		DataUltimaAtualizacaoNcm: time.Now(),
		Nomenclaturas:            nil,
	}

	ncmConvertido := consulta.paraNcmSaida(ncm)
	if ncmConvertido.ID != ncm.ID {
		t.Errorf("Ocorreu um erro, os campos ID são divergentes %v <> %v", ncmConvertido.ID, ncm.ID)
	}
	if ncmConvertido.DataUltimaAtualizacaoNcm != ncm.DataUltimaAtualizacaoNcm.Format(consulta.modeloData) {
		t.Errorf("Ocorreu um erro, os campos DataUltimaAtualizacaoNCM são divergentes %v <> %v", ncmConvertido.DataUltimaAtualizacaoNcm, ncm.DataUltimaAtualizacaoNcm)
	}
}
