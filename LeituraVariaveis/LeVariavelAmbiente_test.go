package LeituraVariaveis

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func CriarArquivoEnv(nome string, texto string) error {
	file, err := os.Create(nome)
	if err != nil {
		return err
	}
	defer file.Close()
	data := []byte(texto)
	_, err = file.Write(data)
	return err
}

func ApagarArquivo(nome string) error {
	return os.Remove(nome)
}

func CriarArquivoZip(nome string) error {
	fileUncompressed, _ := os.Open(nome)
	read := bufio.NewReader(fileUncompressed)
	data, _ := io.ReadAll(read)
	name_of_file := ConverterNomeParaZip(nome)
	fileCompressed, _ := os.Create(name_of_file)
	w := gzip.NewWriter(fileCompressed)
	_, err := w.Write(data)
	err = w.Close()
	err = fileCompressed.Close()
	err = fileUncompressed.Close()
	return err
}

func ApagarArquivoZip(nome string) error {
	return ApagarArquivo(ConverterNomeParaZip(nome))
}

func ConverterNomeParaZip(nome string) string {
	return strings.Replace(nome, ".txt", ".gz", -1)
}

func TestNewLeVariavelAmbiente(t *testing.T) {
	nome := ".teste"
	texto := "CONNSTRING : \"teste1\"\nCONNHTTP : \"teste2\"\nCRONEXPRESSION: \"* * * * *\""
	err := CriarArquivoEnv(nome, texto)
	defer ApagarArquivo(nome)
	if err != nil {
		t.Errorf("Ocorreu um erro ao criar o arquivo %snome. %v", nome, err)
	}
	leitor, err := NewLeVariavelAmbiente(nome)
	if err != nil {
		t.Errorf("Ocorreu um erro ao ler o arquivo %snome. %v", nome, err)
	}
	if got, _ := NewLeVariavelAmbiente(nome); !reflect.DeepEqual(got, leitor) {
		t.Errorf("NewLeVariavelAmbiente() = %v diferente de %v", got, nome)
	}
}

func TestNewLeVariavelAmbiente_Fail(t *testing.T) {
	nome := "teste.txt"
	texto := ""
	err := CriarArquivoEnv(nome, texto)
	CriarArquivoZip(nome)
	defer func() {
		err = ApagarArquivo(nome)
		err = ApagarArquivoZip(ConverterNomeParaZip(nome))
	}()
	if err != nil {
		t.Errorf("Ocorreu um erro ao criar o arquivo %snome. %v", nome, err)
	}
	_, err = NewLeVariavelAmbiente(ConverterNomeParaZip(nome))
	if err == nil {
		t.Errorf("Não ocorreu um erro ao ler o arquivo %s %v", nome, err)
	}
}

func Test_leVariavelAmbiente_ConnectionHttp(t *testing.T) {
	nome := "teste.txt"
	texto := "CONNHTTP :   \"https://teste123.com\""
	err := CriarArquivoEnv(nome, texto)
	if err != nil {
		t.Errorf("Ocorreu um erro ao criar o arquivo de ambiente! %v", err)
	}
	defer ApagarArquivo(nome)
	variavel, err := NewLeVariavelAmbiente(nome)
	if err != nil {
		t.Errorf("Ocorreu um erro ao carregar os valores de ambiente! %v", err)
	}

	stringhttp := variavel.ConnectionHttp()
	if strings.Contains(stringhttp, texto) {
		t.Errorf("Ocorreu um erro, o valor de conexão carregado foi %v, "+
			" quando deveria ser %s", stringhttp, texto)
	}
}

func Test_leVariavelAmbiente_ConnectionString(t *testing.T) {
	nome := "teste.txt"
	texto := "CONNSTRING : \"banco://teste:teste@localhost:5432/teste\""
	err := CriarArquivoEnv(nome, texto)
	if err != nil {
		t.Errorf("Ocorreu um erro ao criar o arquivo de ambiente! %v", err)
	}
	defer ApagarArquivo(nome)
	variavel, err := NewLeVariavelAmbiente(nome)
	if err != nil {
		t.Errorf("Ocorreu um erro ao carregar os valores de ambiente! %v", err)
	}

	stringConnection := variavel.ConnectionString()
	if strings.Contains(stringConnection, texto) {
		t.Errorf("Ocorreu um erro, o valor de conexão carregado foi %v, "+
			" quando deveria ser %s", stringConnection, texto)
	}
}

func Test_leVariavelAmbiente_CronExpression(t *testing.T) {
	nome := "teste.txt"
	texto := "CRONEXPRESSION: \"* * * * *\""
	err := CriarArquivoEnv(nome, texto)
	if err != nil {
		t.Errorf("Ocorreu um erro ao criar o arquivo de ambiente! %v", err)
	}
	defer ApagarArquivo(nome)
	variavel, err := NewLeVariavelAmbiente(nome)
	if err != nil {
		t.Errorf("Ocorreu um erro ao carregar os valores de ambiente! %v", err)
	}

	stringCron := variavel.CronExpression()
	if strings.Contains(stringCron, texto) {
		t.Errorf("Ocorreu um erro, o valor de conexão carregado foi %v, "+
			" quando deveria ser %s", stringCron, texto)
	}
}

func Test_leVariavelAmbiente_ConnectionPort(t *testing.T) {
	nome := "teste.txt"
	texto := "CONNECTIONPORT: \"80\""
	err := CriarArquivoEnv(nome, texto)
	if err != nil {
		t.Errorf("Ocorreu um erro ao criar o arquivo de ambiente! %v", err)
	}
	defer ApagarArquivo(nome)
	variavel, err := NewLeVariavelAmbiente(nome)
	if err != nil {
		t.Errorf("Ocorreu um erro ao carregar os valores de ambiente! %v", err)
	}

	stringCron := variavel.ConnectionPort()
	if strings.Contains(stringCron, texto) {
		t.Errorf("Ocorreu um erro, o valor de conexão carregado foi %v, "+
			" quando deveria ser %s", stringCron, texto)
	}
}
