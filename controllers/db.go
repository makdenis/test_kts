package controllers

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Handle struct {
	Db *sql.DB
}

func (handle *Handle) InitDB(dialect, connStr string) {
	db, err := sql.Open(dialect, connStr)
	if err != nil {
		log.Println(err)
	}
	handle.Db = db
}
