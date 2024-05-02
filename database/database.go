package database

import (
	"database/sql"
	"net/url"

	"github.com/gioCuesta25/employees-manager-backend/config"
	_ "github.com/lib/pq"
)

func NewDbConnection(env config.Environment) (*sql.DB, error) {
	// data source name
	dsn := url.URL{
		Scheme: "postgres",
		Host:   env.DbHost,
		User:   url.UserPassword(env.DbUser, env.DbPassword),
		Path:   env.DbName,
	}

	q := dsn.Query()
	q.Add("sslmode", "disable")

	dsn.RawQuery = q.Encode()

	db, err := sql.Open("postgres", dsn.String())

	if err != nil {
		return nil, err
	}

	return db, nil
}
