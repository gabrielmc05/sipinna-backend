package dboper

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/dboper"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
)

var db *sql.DB

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s.", err)
	}

	if initialPassword, ok := os.LookupEnv("ROOT_ADMIN_INITIAL_PASSWORD"); ok {
		dboper.ChangeUserPassword(3, initialPassword)
	} else {
		dboper.ChangeUserPassword(3, "administrador12345")
	}
}

// Setup se usa para hacer la conexión con la respectiva base de datos
// cuya fuente esté en las variables de entorno, meramente una función
// wrapper para realizarlo.
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
	return err
}

// CloseDB se usa en caso de tener que cerrar la conexión.
func CloseDB() error {
	return db.Close()
}
