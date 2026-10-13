package dboper

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

var db *sql.DB

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s.", err)
	}
}

// Se usa esta función para hacer la conexión con la respectiva
// base de datos cuya fuente esté en las variables de entorno,
// meramente una función wrapper para realizarlo
func Setup() error {
	var (
		dataSourceName strings.Builder
		err            error
	)

	dataSourceName.WriteString("root:")
	if val, ok := os.LookupEnv("MARIADB_ROOT_PASSWORD"); ok {
		dataSourceName.WriteString(val)
		dataSourceName.WriteString("@(")
	} else {
		log.Fatal("MARIADB_ROOT_PASSWORD environment variable not set.")
	}
	if val, ok := os.LookupEnv("MARIADB_ADDRESS"); ok {
		dataSourceName.WriteString(val)
		dataSourceName.WriteString(")/")
	} else {
		log.Fatal("MARIADB_ADDRESS environment variable not set.")
	}
	if val, ok := os.LookupEnv("MARIADB_DATABASE"); ok {
		dataSourceName.WriteString(val)
	} else {
		log.Fatal("MARIADB_DATABASE environment variable not set.")
	}

	db, err = sql.Open("mysql", dataSourceName.String())
	return err
}

// CloseDB se usa en caso de tener que cerrar la conexión.
func CloseDB() error {
	return db.Close()
}
