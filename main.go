package main

import (
	"ConsultaTabelas/Banco"
	"ConsultaTabelas/ConsultaHTTP"
	"ConsultaTabelas/ConsultaNCM"
	"fmt"
	"strconv"
)

func main() {
	conexao := Banco.NewBanco("postgres://admin:admin@localhost:5432/tabelas")
	err := conexao.Conectar()
	if err != nil {
		panic(err)
	}

	consultaHttp := ConsultaHTTP.New("https://portalunico.siscomex.gov.br/classif/api/publico/nomenclatura/download/json")
	consulta := ConsultaNCM.New(consultaHttp)
	ncm, _ := consulta.ConsultarNCM()
	fmt.Println("Data do Arquivo : " + ncm.DataUltimaAtualizacaoNcm)
	fmt.Println("Quantidade de nomenclaturas : " + strconv.Itoa(len(ncm.Nomenclaturas)))
}
