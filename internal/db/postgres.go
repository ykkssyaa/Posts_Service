package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type PostgresOptions struct {
	Host, Port, User, Password, Name string
}

func NewPostgresDB(options PostgresOptions) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			options.User,
			options.Password,
			options.Host,
			options.Port,
			options.Name))

	if err != nil {
		return nil, err
	}

	return db, nil
}
