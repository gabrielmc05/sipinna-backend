package dboper

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func Connect() (dbsipinna *sql.DB, err error) {
	var dataSourceName string = "root:"

	if val, ok := os.LookupEnv("MARIADB_ROOT_PASSWORD"); ok {
		dataSourceName += val + "@("
	} else {
		log.Fatal("MARIADB_ROOT_PASSWORD environment variable not set.")
	}
	if val, ok := os.LookupEnv("MARIADB_ADDRESS"); ok {
		dataSourceName += val + ")/"
	} else {
		log.Fatal("MARIADB_ADDRESS environment variable not set.")
	}
	if val, ok := os.LookupEnv("MARIADB_DATABASE"); ok {
		dataSourceName += val
	} else {
		log.Fatal("MARIADB_DATABASE environment variable not set.")
	}

	return sql.Open("mysql", dataSourceName)
}
