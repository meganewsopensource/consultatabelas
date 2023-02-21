package Web

import (
	"ConsultaTabelas/Banco"
	"ConsultaTabelas/Banco/NCM"
	"ConsultaTabelas/MockTestes"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewControllerSaude(t *testing.T) {
	conexao, err := MockTestes.GerarConexaoBanco()
	defer MockTestes.DeletarBanco(conexao)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão!")
	}
	repository := NCM.NewRepositorySaude(conexao)
	verificaSaude := Banco.NewVerificaSaudeBanco(repository)
	controller := NewControllerSaude(verificaSaude)
	if got := NewControllerSaude(verificaSaude); !reflect.DeepEqual(got, controller) {
		t.Errorf("ControllerSaude %v diferente de %v", got, controller)
	}
}

func Test_controllerSaude_VerificarSaude(t *testing.T) {
	conexao, err := MockTestes.GerarConexaoBanco()
	defer MockTestes.DeletarBanco(conexao)
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão! %v", err)
	}
	repository := NCM.NewRepositorySaude(conexao)
	verificaSaude := Banco.NewVerificaSaudeBanco(repository)
	controller := NewControllerSaude(verificaSaude)

	r := gin.Default()
	r.GET("saude", controller.VerificarSaude)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/saude", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Ocorreu um erro buscar o status de saude. %v %v", w.Code, w.Body)
	}
}

func Test_controllerSaude_VerificarSaude_Fail(t *testing.T) {
	conexao, err := MockTestes.GerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão! %v", err)
	}
	repository := NCM.NewRepositorySaude(conexao)
	verificaSaude := Banco.NewVerificaSaudeBanco(repository)
	controller := NewControllerSaude(verificaSaude)
	MockTestes.DeletarBanco(conexao)

	r := gin.Default()
	r.GET("saude", controller.VerificarSaude)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/saude", nil)
	r.ServeHTTP(w, req)
	if w.Code == 200 {
		t.Errorf("Não ocorreu o erro esperado ao buscar o status de saude. %v %v", w.Code, w.Body)
	}
}
