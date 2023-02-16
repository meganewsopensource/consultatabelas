package LeituraVariaveis

import (
	"github.com/joho/godotenv"
	"os"
)

type leVariavelAmbiente struct {
}

type IVariavelAmbiente interface {
	ConnectionString() string
	ConnectionHttp() string
	CronExpression() string
}

func NewLeVariavelAmbiente(nomeArquivo string) (IVariavelAmbiente, error) {
	var err error
	if _, err = os.Stat(nomeArquivo); err == nil {
		err = godotenv.Load(nomeArquivo)
	}

	return &leVariavelAmbiente{}, err
}

func (variavel *leVariavelAmbiente) ConnectionString() string {
	return os.Getenv("CONNSTRING")
}

func (variavel *leVariavelAmbiente) ConnectionHttp() string {
	return os.Getenv("CONNHTTP")
}

func (variavel *leVariavelAmbiente) CronExpression() string {
	return os.Getenv("CRONEXPRESSION")
}
