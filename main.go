package main

import (
	"ConsultaTabelas/Banco/NCM"
	"ConsultaTabelas/ConsultaHTTP"
	"ConsultaTabelas/ConsultaNCM"
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
)

func main() {
	sqlDB, err := sql.Open("pgx", "postgres://admin:admin@localhost:5432/tabelas")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	migrations := NCM.NewMigration(db)
	err = migrations.Executar()
	if err != nil {
		panic(err)
	}
	
	consultaHttp := ConsultaHTTP.New("https://portalunico.siscomex.gov.br/classif/api/publico/nomenclatura/download/json")
	consulta := ConsultaNCM.New(consultaHttp)
	ncm, _ := consulta.ConsultarNCM()
	fmt.Println("Data do Arquivo : " + ncm.DataUltimaAtualizacaoNcm)
	fmt.Println("Quantidade de nomenclaturas : " + strconv.Itoa(len(ncm.Nomenclaturas)))
}
