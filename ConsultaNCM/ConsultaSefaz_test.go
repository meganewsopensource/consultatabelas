package ConsultaNCM

import (
	"ConsultaTabelas/ConsultaHTTP"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	consulta := ConsultaHTTP.New("teste.com")
	consultaServer := New(consulta)
	if got := New(consulta); !reflect.DeepEqual(got, consultaServer) {
		t.Errorf("Esperado %v, recebi %v", got, consulta)
	}
}

func Test_consultaSefaz_ConsultarNCM(t *testing.T) {
	ncm := preencheNcmReceita()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonNCM, _ := json.Marshal(ncm)
		fmt.Fprintf(w, string(jsonNCM[:]))
	}))
	defer server.Close()

	consulta := consultaSefaz{
		consultaHttp: ConsultaHTTP.New(server.URL),
	}
	got, err := consulta.ConsultarNCM()
	if err != nil {
		t.Errorf("ConsultarNCM() error = %v, wantErr nil", err)
		return
	}
	if !reflect.DeepEqual(got, ncm) {
		t.Errorf("ConsultarNCM() got = %v, want %v", got, ncm)
	}
}

func Test_consultaSefaz_ConsultarNCMFailHttp(t *testing.T) {
	consulta := consultaSefaz{
		consultaHttp: preencheMockConsulta(),
	}
	_, err := consulta.ConsultarNCM()
	if err == nil {
		t.Errorf("Esperado mensagem de erro preenchida, recebi %v", err)
		return
	}
}

func Test_consultaSefaz_ConsultarNCMDesserealizarFail(t *testing.T) {
	ncm := "teste de NCM"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonNCM, _ := json.Marshal(ncm)
		fmt.Fprintf(w, string(jsonNCM[:]))
	}))
	defer server.Close()

	consulta := consultaSefaz{
		consultaHttp: ConsultaHTTP.New(server.URL),
	}
	_, err := consulta.ConsultarNCM()
	if err == nil {
		t.Errorf("Esperado mensagem de erro preenchida, recebi %v", err)
		return
	}
}

func Test_consultaSefaz_desserealizarNCM(t *testing.T) {
	ncm := preencheNcmReceita()
	jsonNCM, _ := json.Marshal(ncm)
	consulta := consultaSefaz{
		consultaHttp: ConsultaHTTP.New("teste"),
	}
	got, err := consulta.desserealizarNCM(jsonNCM)
	if err != nil {
		t.Errorf("desserealizarNCM() error = %v, wantErr nil", err)
		return
	}
	if !reflect.DeepEqual(got, ncm) {
		t.Errorf("desserealizarNCM() got = %v, want %v", got, ncm)
	}
}

func preencheNcmReceita() NcmReceita {
	nomenclaturas := []Nomenclatura{}
	nomenclaturas = append(nomenclaturas, Nomenclatura{
		Codigo:     "01",
		Descricao:  "teste de NCM",
		DataInicio: "01/01/2023",
		DataFim:    "31/12/2023",
		TipoAto:    "teste de ato",
		NumeroAto:  "2",
		AnoAto:     "2023",
	})

	ncm := NcmReceita{
		DataUltimaAtualizacaoNcm: "09/02/2023",
		Nomenclaturas:            nomenclaturas,
	}
	return ncm
}

func preencheMockConsulta() ConsultaHTTP.IConsultaHttp {
	mensagem := "ocorreu um erro"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(mensagem)
	}))
	defer server.Close()
	return ConsultaHTTP.New(server.URL)
}
