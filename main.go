package main

import (
	"ConsultaTabelas/Banco/NCM"
	"ConsultaTabelas/ConsultaHTTP"
	"ConsultaTabelas/ConsultaNCM"
	"ConsultaTabelas/ConsultaNCMSefaz"
	"ConsultaTabelas/LeituraVariaveis"
	"ConsultaTabelas/Web"
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
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
	r := ConfigurarGin()
	controllerNcm := Web.NewControllerNCM(consulta)

	public := r.Group("/")
	{
		public.POST("AtualizarNCM", controllerNcm.AtualizarNCM)
		public.GET("ncms", controllerNcm.ListarNCMS)
		public.GET("ncms/:data", controllerNcm.ListarNCMPorData)
		public.GET("atualizacoes/ultima", controllerNcm.DataUltimaAtualizacao)
		public.GET("saude", func(c *gin.Context) {
			sqlDB, err := db.DB()
			if err != nil {
				c.JSON(512, "Unhealthy")
				return
			}
			err = sqlDB.Ping()
			if err != nil {
				c.JSON(512, "Unhealthy")
				return
			}
			c.JSON(200, "Healthy")
		})
	}

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		err := http.ListenAndServe(":9090", nil)
		if err != nil {
			os.Exit(1)
		}
	}()

	runCronJobs(consulta.AtualizarNCM, variaveis.CronExpression())
	r.Run()
}

func ConfigurarGin() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
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
