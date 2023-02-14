package LeituraVariaveis

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type leVariavelAmbiente struct {
}

type IVariavelAmbiente interface {
	ConnectionString() string
	ConnectionHttp() string
}

func NewLeVariavelAmbiente() IVariavelAmbiente {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}

	return &leVariavelAmbiente{}
}

func (variavel *leVariavelAmbiente) ConnectionString() string {
	return os.Getenv("CONNSTRING")
}

func (variavel *leVariavelAmbiente) ConnectionHttp() string {
	return os.Getenv("CONNHTTP")
}
