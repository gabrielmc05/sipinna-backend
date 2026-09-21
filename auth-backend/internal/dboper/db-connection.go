package dboper

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var db *sql.DB

func init() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if db, err = connect(); err != nil {
		log.Fatal(err)
	}
}

func connect() (dbsipinna *sql.DB, err error) {
	var dataSourceName string = "root:"

	if val, ok := os.LookupEnv("MARIADB_ROOT_PASSWORD"); ok {
		dataSourceName += val + "@/"
	} else {
		log.Fatal("MARIADB_ROOT_PASSWORD environment variable not set.")
	}
	if val, ok := os.LookupEnv("MARIADB_DATABASE"); ok {
		dataSourceName += val
	} else {
		log.Fatal("MARIADB_DATABASE environment variable not set.")
	}

	return sql.Open("mysql", dataSourceName)
}
