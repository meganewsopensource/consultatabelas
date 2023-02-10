package ConsultaHTTP

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestConsultaHttp_Consultar(t *testing.T) {
	esperado := "qualquer dado"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, esperado)
	}))
	defer server.Close()
	consulta := New(server.URL)
	resposta, erro := consulta.Consultar()
	if erro != nil {
		t.Errorf("Esperado erro igual a nil, recebi %v", erro)
	}
	textoResposta := string(resposta[:])
	if textoResposta != textoResposta {
		t.Errorf("Esperava resposta igual a %s, recebi %s", esperado, textoResposta)
	}
}

func TestConsultaHttp_ConsultarFail(t *testing.T) {
	mensagem := "ocorreu um erro"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(mensagem)
	}))
	defer server.Close()
	consulta := New(server.URL)
	_, erro := consulta.Consultar()
	if erro == nil {
		t.Errorf("Esperado mensagem %s, recebi %v", mensagem, erro)
	}
}

func TestNew(t *testing.T) {
	esperado := "url_qualquer"
	consulta := New(esperado)
	if got := New(esperado); !reflect.DeepEqual(got, consulta) {
		t.Errorf("Esperado %v, recebi %v", got, consulta)
	}
}
