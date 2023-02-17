package MockTestes

import (
	"ConsultaTabelas/Banco/NCM"
	"ConsultaTabelas/ConsultaNCMSefaz"
	"encoding/json"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"os"
)

func GerarConexaoBanco() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	migracao := NCM.NewMigration(db)
	_ = migracao.Executar()
	return db, err
}

func CriarServidor() *httptest.Server {
	ncm := PreencheNcmReceita()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonNCM, _ := json.Marshal(ncm)
		fmt.Fprintf(w, string(jsonNCM[:]))
	}))
	return server
}

func DeletarBanco(db *gorm.DB) {
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

func PreencheNcmReceita() ConsultaNCMSefaz.NcmReceita {
	return ConsultaNCMSefaz.NcmReceita{
		DataUltimaAtualizacaoNcm: "01/01/2023",
		Nomenclaturas:            PreencheNomenclaturas(),
	}
}

func PreencheNomenclaturas() []ConsultaNCMSefaz.Nomenclatura {
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
