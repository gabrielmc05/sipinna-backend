package dboper

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

func init() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

// Se usa esta función para hacer la conexión con la respectiva
// base de datos cuya fuente esté en las variables de entorno,
// meramente una función wrapper para realizarlo
func Connect() (dbsipinna *sql.DB, err error) {
	var dataSourceName strings.Builder

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

	//log.Println(dataSourceName.String())
	return sql.Open("mysql", dataSourceName.String())
}
