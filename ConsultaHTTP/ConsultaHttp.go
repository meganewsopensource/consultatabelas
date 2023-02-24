package ConsultaHTTP

import (
	"io"
	"net/http"
)

type IConsultaHttp interface {
	Consultar() ([]byte, error)
}

type ConsultaHttp struct {
	url string
}

func New(url string) IConsultaHttp {
	return &ConsultaHttp{url}
}

func (consultaHttp ConsultaHttp) Consultar() ([]byte, error) {
	resp, err := http.Get(consultaHttp.url)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}
