package Banco

import (
	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IConexaoBanco interface {
	Conectar() error
	Conexao() *gorm.DB
}

type conexaoBanco struct {
	db               *gorm.DB
	connectionString string
}

func NewBanco(stringConnection string) IConexaoBanco {
	return &conexaoBanco{connectionString: stringConnection}
}

func (banco *conexaoBanco) Conectar() error {
	if banco.db != nil {
		return nil
	}

	sqlDB, err := sql.Open("pgx", banco.connectionString)
	if err != nil {
		return err
	}
	banco.db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

func (banco *conexaoBanco) Conexao() *gorm.DB {
	return banco.db
}
