package dboper

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
	"time"
)

var db *sql.DB

func init() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
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

	password, ok := os.LookupEnv("MARIADB_ROOT_PASSWORD")
	if !ok {
		return fmt.Errorf("MARIADB_ROOT_PASSWORD environment variable not set")
	}
	address, ok := os.LookupEnv("MARIADB_ADDRESS")
	if !ok {
		return fmt.Errorf("MARIADB_ADDRESS environment variable not set")
	}
	database, ok := os.LookupEnv("MARIADB_DATABASE")
	if !ok {
		return fmt.Errorf("MARIADB_DATABASE environment variable not set")
	}

	dataSourceName.WriteString("root:")
	dataSourceName.WriteString(password)
	dataSourceName.WriteString("@(")
	dataSourceName.WriteString(address)
	dataSourceName.WriteString(")/")
	dataSourceName.WriteString(database)
	dataSourceName.WriteString("?parseTime=true")

	db, err = sql.Open("mysql", dataSourceName.String())
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		db = nil
		return fmt.Errorf("connect to database: %w", err)
	}
	return nil
}

func CloseDB() error {
	if db == nil {
		return nil
	}
	return db.Close()
}
