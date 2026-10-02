package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// DB es el pool de conexiones compartido para que el resto del paquete pueda
// ejecutar consultas. *sql.DB es un tipo puntero que administra conexiones,
// no representa necesariamente una única conexión física. 

var DB *sql.DB

// init se ejecuta automáticamente al cargar este paquete, antes de main.
// Carga variables definidas en .env al entorno del proceso.

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s.", err)
	}
}

// Connect construye la dirección de conexión usando variables de entorno,
// inicializa el pool MySQL y comprueba que la base responda.
func Connect() error {
	// Estos nombres deben existir en el entorno o en el archivo .env.
	user := "root"
	pass := os.Getenv("MARIADB_ROOT_PASSWORD")
	host := os.Getenv("MARIADB_HOST")
	dbname := os.Getenv("MARIADB_DATABASE")

	// DSN significa Data Source Name: contiene usuario, contraseña, host y base.
	// parseTime=true permite convertir columnas de fecha/hora a tipos time.Time.
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, pass, host, dbname)

	var err error
	// El import del driver MySQL se hace con _ porque solo necesitamos su efecto:
	// registrar "mysql" para que database/sql pueda usarlo. sql.Open prepara el pool;
	// la conexión efectiva se comprueba en el Ping de abajo.
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// Ping intenta comunicarse con la base de datos. Si devuelve error, Connect
	// informa el fallo a main, que decide no iniciar el servidor.
	return DB.Ping()
}
