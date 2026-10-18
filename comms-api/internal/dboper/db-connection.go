package dboper

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
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
		rootCertPool   *x509.CertPool = x509.NewCertPool()
		pem            []byte
	)

	if pem, err = os.ReadFile("./global-bundle.pem"); err != nil {
		log.Fatalf("Error, failed to read cert: %s.", err)
	}

	rootCertPool.AppendCertsFromPEM(pem)
	mysql.RegisterTLSConfig("custom", &tls.Config{RootCAs: rootCertPool})

	dataSourceName.WriteString("sipinna_admin:")
	if val, ok := os.LookupEnv("MARIADB_ROOT_PASSWORD"); ok {
		dataSourceName.WriteString(val)
		dataSourceName.WriteString("@tcp(")
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
		dataSourceName.WriteString("?tls=custom&parseTime=true")
	} else {
		log.Fatal("MARIADB_DATABASE environment variable not set.")
	}

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
