package Web

import (
	"ConsultaTabelas/ConsultaNCM"
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
}

func (controller *controllerNCM) AtualizarNCM(context *gin.Context) {
	err := controller.consulta.AtualizarNCM()
	if err != nil {
		context.JSON(http.StatusInternalServerError, "Não foi possivel atualiza os dados de NCM!")
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
