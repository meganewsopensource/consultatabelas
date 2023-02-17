package main

import (
	"ConsultaTabelas/Banco"
	"ConsultaTabelas/Banco/NCM"
	"ConsultaTabelas/ConsultaHTTP"
	"ConsultaTabelas/ConsultaNCM"
	"ConsultaTabelas/ConsultaNCMSefaz"
	"ConsultaTabelas/LeituraVariaveis"
	"ConsultaTabelas/Web"
	_ "ConsultaTabelas/docs"
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Consulta Tabelas
// @version 1.0
// @description Testing Swagger APIs.
// @termsOfService http://swagger.io/terms/// @contact.name API Support

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io// @securityDefinitions.apiKey JWT

// @in header
// @name token// @license.name Apache 2.0

// @license.url http://www.apache.org/licenses/LICENSE-2.0.html// @host localhost:8081

// @BasePath /

// @schemes http
func main() {
	variaveis, err := LeituraVariaveis.NewLeVariavelAmbiente(".env")
	if err != nil {
		panic(err)
	}
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
	r := ConfigurarGin()
	controllerNcm := Web.NewControllerNCM(consulta)

	repositorySaude := NCM.NewRepositorySaude(db)
	verificaSaude := Banco.NewVerificaSaudeBanco(repositorySaude)
	controllerSaude := Web.NewControllerSaude(verificaSaude)

	public := r.Group("/")
	{
		public.POST("ncms/atualizar", controllerNcm.AtualizarNCM)
		public.GET("ncms", controllerNcm.ListarNCMS)
		public.GET("ncms/:data", controllerNcm.ListarNCMPorData)
		public.GET("atualizacoes/ultima", controllerNcm.DataUltimaAtualizacao)
		public.GET("metrics", prometheusHandler())
		public.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		public.GET("saude", controllerSaude.VerificarSaude)
	}

	runCronJobs(consulta.AtualizarNCM, variaveis.CronExpression())
	r.Run()
}

func ConfigurarGin() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "username", "login", "senha", "Access-Control-Allow-Credentials"},
	}))

	return r
}

func runCronJobs(consulta func() error, cronExpression string) {
	cron := ConsultaNCM.NewCronJob(cronExpression)
	cron.ConfigurarCron(consulta)
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
