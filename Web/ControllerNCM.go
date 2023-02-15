package Web

import (
	"ConsultaTabelas/ConsultaNCM"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type controllerNCM struct {
	consulta ConsultaNCM.IConsultaNCM
}

func NewControllerNCM(consulta ConsultaNCM.IConsultaNCM) IControllerNCM {
	controller := controllerNCM{
		consulta: consulta,
	}
	return &controller
}

type IControllerNCM interface {
	AtualizarNCM(context *gin.Context)
	ListarNCMS(context *gin.Context)
	DataUltimaAtualizacao(context *gin.Context)
	ListarNCMPorData(context *gin.Context)
}

func (controller *controllerNCM) AtualizarNCM(context *gin.Context) {
	err := controller.consulta.AtualizarNCM()
	if err != nil {
		context.JSON(http.StatusInternalServerError, fmt.Sprintf("Não foi possivel atualiza os dados de NCM! Erro : %v", err))
	} else {
		context.JSON(http.StatusOK, "Os dados de NCM foram atualizados!")
	}
}

func (controller *controllerNCM) ListarNCMS(context *gin.Context) {
	lista, err := controller.consulta.ListarNCMs()
	if err != nil {
		context.JSON(http.StatusInternalServerError, "Ocorreu um erro ao buscar os registros no banco de dados!")
	} else {
		context.JSON(http.StatusOK, lista)
	}
}

func (controller *controllerNCM) DataUltimaAtualizacao(context *gin.Context) {
	data, err := controller.consulta.UltimaAtualizacao()
	if err != nil {
		context.JSON(http.StatusInternalServerError, "Ocorreu um erro ao buscar a última data de atualização!")
	} else {
		context.JSON(http.StatusOK, data)
	}
}

func (controller *controllerNCM) ListarNCMPorData(context *gin.Context) {
	data := context.Param("data")
	lista, err := controller.consulta.ListarNCMPorData(data)
	if err != nil {
		context.JSON(http.StatusInternalServerError, "Ocorreu um erro ao consultar NCMs por data!")
	} else {
		context.JSON(http.StatusOK, lista)
	}
}
