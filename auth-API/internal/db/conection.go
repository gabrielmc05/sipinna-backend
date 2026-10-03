package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)


var DB *sql.DB


// cargamos variables del .env al entorno del proceso
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s.", err)
	}
}

// Connect construye la dirección de conexión usando variables de entorno,
// inicializa el pool MySQL y comprueba que la base responda.
func Connect() error {
	
	user := "root"
	pass := os.Getenv("MARIADB_ROOT_PASSWORD")
	host := os.Getenv("MARIADB_HOST")
	dbname := os.Getenv("MARIADB_DATABASE")

	// parseTime=true convierte columnas de fecha/hora a tipos time.Time.
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, pass, host, dbname)

	var err error


	// registrar "mysql" para que database/sql pueda usarlo. sql.Open para preparar el pool;

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// Probar conexión
	return DB.Ping()
}
