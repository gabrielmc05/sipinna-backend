/* TODO */
package keys

import (
	"crypto/aes"
	"encoding/hex"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	var err error
	if err = godotenv.Load(".env.keys"); err != nil {
		log.Fatalf("Error loading .env.keys file: %s.\n", err)
	}

	if dbBlock, err = aes.NewCipher(getSymmetricDBKey()); err != nil {
		log.Fatalf("DB cipher block could not be made: %s.\n", err)
	}
}

func getSymmetricDBKey() (dbkey []byte) {
	if val, ok := os.LookupEnv("MARIADB_DATABASE_SYMMETRIC_KEY"); ok {
		dbkey, _ = hex.DecodeString(val)
		log.Printf("Key: %x, key size: %d", dbkey, len(dbkey))
	} else {
		log.Fatal("MARIADB_DATABASE_SYMMETRIC_KEY environment variable not set.")
	}

	return dbkey
}
