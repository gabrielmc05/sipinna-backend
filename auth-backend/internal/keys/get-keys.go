/* TODO */
package keys

import (
	"encoding/hex"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(".env.keys"); err != nil {
		log.Printf("Error loading .env.keys file: %s.\n", err)
	}
}

func GetSymmetricDBKey() (dbkey []byte) {
	if val, ok := os.LookupEnv("MARIADB_DATABASE_SYMMETRIC_KEY"); ok {
		dbkey, _ = hex.DecodeString(val)
		log.Printf("Key: %x, key size: %d", dbkey, len(dbkey))
	} else {
		log.Fatal("MARIADB_DATABASE_SYMMETRIC_KEY environment variable not set.")
	}

	return dbkey
}
