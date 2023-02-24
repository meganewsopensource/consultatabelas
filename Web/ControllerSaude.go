package Web

import (
	"ConsultaTabelas/Banco"
	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
	"net/http"
)

type controllerSaude struct {
	verificaSaude Banco.IVerificaSaudeBanco
}

func NewControllerSaude(verificaSaude Banco.IVerificaSaudeBanco) IControllerSaude {

	controller := controllerSaude{verificaSaude}
	return &controller
}

type IControllerSaude interface {
	VerificarSaude(context *gin.Context)
}

// @BasePath /

// VerificarSaude
// @Summary Verificar a saude do banco de dados
// @Schemes
// @Description Verificar a saude do banco de dados
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} NCM
// @Router /saude [get]
func (controller *controllerSaude) VerificarSaude(context *gin.Context) {
	saude, saudavel := controller.verificaSaude.VerificarSaude()
	if saudavel {
		context.JSON(http.StatusOK, saude)
	} else {
		context.JSON(512, saude)
	}
}
