package main

import (
	"database/sql"
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/dboper"
	"log"
)

var db *sql.DB

func main() {
	var err error
	db, err = dboper.Connect()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dboper.Probando(db)
}
