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
	migracao.Executar()
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
	defer server.Close()
	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)
	consulta := NewConsultaNCM(consultaSefaz, repositoryNCM, repositoryNomenclatura)

	if got := NewConsultaNCM(consultaSefaz, repositoryNCM, repositoryNomenclatura); !reflect.DeepEqual(got, consulta) {
		t.Errorf("ConsultaNCM() = %v, diferente de %v", got, consulta)
	}

	deletarBanco(db)
}

func Test_consultaNCM_AtualizarNCM(t *testing.T) {
	server := criarServidor()
	defer server.Close()
	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)
	consulta := NewConsultaNCM(consultaSefaz, repositoryNCM, repositoryNomenclatura)

	err = consulta.AtualizarNCM()
	if err != nil {
		t.Errorf("AtualizarNCM() error : %v", err)
	}
}

func Test_consultaNCM_gravarNCM(t *testing.T) {
	server := criarServidor()
	defer server.Close()
	consultaHttp := ConsultaHTTP.New(server.URL)
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)
	db, err := gerarConexaoBanco()
	if err != nil {
		t.Errorf("Ocorreu um erro ao gerar conexão : %v", err)
	}
	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)

	consulta := &consultaNCM{
		consultaSefaz:          consultaSefaz,
		respotoryNCM:           repositoryNCM,
		repositoryNomenclatura: repositoryNomenclatura,
		modeloData:             "02/01/2006",
	}
	ncm := preencheNcmReceita()
	lista, err := consulta.listaNomenclatura(ncm.Nomenclaturas)
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

func Test_consultaNCM_listaNomenclatura(t *testing.T) {
	type fields struct {
		consultaSefaz          ConsultaNCMSefaz.IConsultaSefaz
		respotoryNCM           NCM.IRepositoryNCM
		repositoryNomenclatura NCM.IRepositoryNomenclatura
		modeloData             string
	}
	type args struct {
		listaNCM []ConsultaNCMSefaz.Nomenclatura
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []NCM.NomenclaturaBanco
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consulta := &consultaNCM{
				consultaSefaz:          tt.fields.consultaSefaz,
				respotoryNCM:           tt.fields.respotoryNCM,
				repositoryNomenclatura: tt.fields.repositoryNomenclatura,
				modeloData:             tt.fields.modeloData,
			}
			got, err := consulta.listaNomenclatura(tt.args.listaNCM)
			if (err != nil) != tt.wantErr {
				t.Errorf("listaNomenclatura() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listaNomenclatura() got = %v, want %v", got, tt.want)
			}
		})
	}
}
