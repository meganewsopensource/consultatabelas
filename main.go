package main

import (
	"ConsultaTabelas/Banco/NCM"
	"ConsultaTabelas/ConsultaHTTP"
	"ConsultaTabelas/ConsultaNCM"
	"ConsultaTabelas/ConsultaNCMSefaz"
	"ConsultaTabelas/LeituraVariaveis"
	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	variaveis := LeituraVariaveis.NewLeVariavelAmbiente()
	sqlDB, err := sql.Open("pgx", variaveis.ConnectionString())
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

	consultaHttp := ConsultaHTTP.New(variaveis.ConnectionHttp())
	consultaSefaz := ConsultaNCMSefaz.New(consultaHttp)
	repositoryNCM := NCM.NewRepositoryNCM(db)
	repositoryNomenclatura := NCM.NewRepositoryNomenclatura(db)
	consulta := ConsultaNCM.NewConsultaNCM(consultaSefaz, repositoryNCM, repositoryNomenclatura)
	err = consulta.AtualizarNCM()
	if err != nil {
		log.Fatal(err)
	}
}
