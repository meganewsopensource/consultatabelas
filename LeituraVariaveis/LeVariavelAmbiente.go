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
	ConnectionPort() string
}

func NewLeVariavelAmbiente(nomeArquivo string) (IVariavelAmbiente, error) {
	if _, err := os.Stat(nomeArquivo); err == nil {
		err := godotenv.Load(nomeArquivo)
		if err != nil {
			return &leVariavelAmbiente{}, err
		}
	}

	return &leVariavelAmbiente{}, nil
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

func (variavel *leVariavelAmbiente) ConnectionPort() string {
	return os.Getenv("CONNECTIONPORT")
}
