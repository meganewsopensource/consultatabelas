package Web

import (
	"ConsultaTabelas/Banco/NCM"
	"ConsultaTabelas/ConsultaHTTP"
	"ConsultaTabelas/ConsultaNCM"
	"ConsultaTabelas/ConsultaNCMSefaz"
	"ConsultaTabelas/MockTestes"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewControllerNCM(t *testing.T) {
	var consultaNCM ConsultaNCM.IConsultaNCM
	conexao, err := MockTestes.GerarConexaoBanco()
	defer MockTestes.DeletarBanco(conexao)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão! %v", err)
	}
	consultaNCM = ConfiguraCasoDeUso(conexao)
	controller := NewControllerNCM(consultaNCM)

	if got := NewControllerNCM(consultaNCM); !reflect.DeepEqual(got, controller) {
		t.Errorf("NewControllerNCM() = %v, diferente de %v", got, controller)
	}
}

func ConfiguraCasoDeUso(conexao *gorm.DB) ConsultaNCM.IConsultaNCM {
	server := MockTestes.CriarServidor()
	consultaHttp := ConsultaHTTP.New(server.URL)
	repositoryNCM := NCM.NewRepositoryNCM(conexao)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(conexao)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)
	consultaNCM := ConsultaNCM.NewConsultaNCM(consultaSefaz, repositoryNCM, repositoryNomenclatura)
	return consultaNCM
}

func Test_controllerNCM_AtualizarNCM(t *testing.T) {
	var consultaNCM ConsultaNCM.IConsultaNCM
	conexao, err := MockTestes.GerarConexaoBanco()
	defer MockTestes.DeletarBanco(conexao)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão! %v", err)
	}
	consultaNCM = ConfiguraCasoDeUso(conexao)
	controller := NewControllerNCM(consultaNCM)

	r := gin.Default()
	r.GET("/AtualizaNCM", controller.AtualizarNCM)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/AtualizaNCM", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Ocorreu um erro ao fazer a requisição de atualização. %v %v", w.Code, w.Body)
	}
}

func Test_controllerNCM_AtualizarNCM_Fail(t *testing.T) {
	var consultaNCM ConsultaNCM.IConsultaNCM
	conexao, err := MockTestes.GerarConexaoBanco()
	defer MockTestes.DeletarBanco(conexao)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão! %v", err)
	}
	consultaNCM = ConfiguraCasoDeUso(conexao)
	controller := NewControllerNCM(consultaNCM)

	mig := conexao.Migrator()
	err = mig.DropTable(&NCM.NomenclaturaBanco{})
	if err != nil {
		t.Errorf("Ocorreu um erro ao apagar tabela! %v", err)
	}

	r := gin.Default()
	r.GET("/AtualizaNCM", controller.AtualizarNCM)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/AtualizaNCM", nil)
	r.ServeHTTP(w, req)
	if w.Code == 200 {
		t.Errorf("Não ocorreu um erro esperado ao fazer a requisição de atualização. %v %v", w.Code, w.Body)
	}
}

func Test_controllerNCM_DataUltimaAtualizacao(t *testing.T) {
	var consultaNCM ConsultaNCM.IConsultaNCM
	conexao, err := MockTestes.GerarConexaoBanco()
	defer MockTestes.DeletarBanco(conexao)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão! %v", err)
	}
	consultaNCM = ConfiguraCasoDeUso(conexao)
	controller := NewControllerNCM(consultaNCM)

	r := gin.Default()
	r.GET("/atualizacao/ultima", controller.DataUltimaAtualizacao)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/atualizacao/ultima", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Ocorreu um erro ao fazer a requisição da data de atualização. %v %v", w.Code, w.Body)
	}
}

func Test_controllerNCM_DataUltimaAtualizacao_fail(t *testing.T) {
	var consultaNCM ConsultaNCM.IConsultaNCM
	conexao, err := MockTestes.GerarConexaoBanco()
	defer MockTestes.DeletarBanco(conexao)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão! %v", err)
	}
	consultaNCM = ConfiguraCasoDeUso(conexao)
	controller := NewControllerNCM(consultaNCM)

	mig := conexao.Migrator()
	err = mig.DropTable(&NCM.NcmBanco{})
	if err != nil {
		t.Errorf("Ocorreu um erro ao apagar tabela! %v", err)
	}

	r := gin.Default()
	r.GET("/atualizacao/ultima", controller.DataUltimaAtualizacao)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/atualizacao/ultima", nil)
	r.ServeHTTP(w, req)
	if w.Code == 200 {
		t.Errorf("Ocorreu um erro ao fazer a requisição da data de atualização. %v %v", w.Code, w.Body)
	}
}

func Test_controllerNCM_ListarNCMPorData(t *testing.T) {
	var consultaNCM ConsultaNCM.IConsultaNCM
	conexao, err := MockTestes.GerarConexaoBanco()
	defer MockTestes.DeletarBanco(conexao)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão! %v", err)
	}
	consultaNCM = ConfiguraCasoDeUso(conexao)
	controller := NewControllerNCM(consultaNCM)

	r := gin.Default()
	r.GET("ncms/:data", controller.ListarNCMPorData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ncms/01-01-2023", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Ocorreu um erro ao fazer a requisição da data de atualização. %v %v", w.Code, w.Body)
	}
}

func Test_controllerNCM_ListarNCMPorData_Fail(t *testing.T) {
	var consultaNCM ConsultaNCM.IConsultaNCM
	conexao, err := MockTestes.GerarConexaoBanco()
	defer MockTestes.DeletarBanco(conexao)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão! %v", err)
	}
	consultaNCM = ConfiguraCasoDeUso(conexao)
	controller := NewControllerNCM(consultaNCM)

	r := gin.Default()
	r.GET("ncms/:data", controller.ListarNCMPorData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ncms/2023-01-01", nil)
	r.ServeHTTP(w, req)
	if w.Code == 200 {
		t.Errorf("Não ocorreu o erro esperado ao fazer a requisição da data de atualização. %v %v", w.Code, w.Body)
	}
}

func Test_controllerNCM_ListarNCMS(t *testing.T) {
	var consultaNCM ConsultaNCM.IConsultaNCM
	conexao, err := MockTestes.GerarConexaoBanco()
	defer MockTestes.DeletarBanco(conexao)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão! %v", err)
	}
	consultaNCM = ConfiguraCasoDeUso(conexao)
	controller := NewControllerNCM(consultaNCM)

	r := gin.Default()
	r.GET("ncms", controller.ListarNCMS)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ncms", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Ocorreu um erro ao fazer a requisição de listagem de NCMs. %v %v", w.Code, w.Body)
	}
}

func Test_controllerNCM_ListarNCMS_Fail(t *testing.T) {
	var consultaNCM ConsultaNCM.IConsultaNCM
	conexao, err := MockTestes.GerarConexaoBanco()
	defer MockTestes.DeletarBanco(conexao)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão! %v", err)
	}
	consultaNCM = ConfiguraCasoDeUso(conexao)
	controller := NewControllerNCM(consultaNCM)

	mig := conexao.Migrator()
	err = mig.DropTable(NCM.NomenclaturaBanco{})
	if err != nil {
		t.Errorf("Ocorreu um erro ao deletar a tabela! %v", err)
	}

	r := gin.Default()
	r.GET("ncms", controller.ListarNCMS)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ncms", nil)
	r.ServeHTTP(w, req)
	if w.Code == 200 {
		t.Errorf("Não ocorreu o erro esperado ao fazer a requisição de listagem de NCMs. %v %v", w.Code, w.Body)
	}
}
