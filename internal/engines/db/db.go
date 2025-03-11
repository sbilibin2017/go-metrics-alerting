package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	*sql.DB
}

func NewDB() *DB {
	return &DB{}
}

func (e *DB) Open(dsn string) {
	if e.DB != nil {
		return
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	e.DB = db
}

func (e *DB) Close() {
	if e.DB == nil {
		return
	}
	err := e.DB.Close()
	if err != nil {
		panic(err)
	}
}
