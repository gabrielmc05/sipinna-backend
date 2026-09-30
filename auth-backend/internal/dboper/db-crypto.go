package dboper

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/keys"
	"log"
)

// EncryptUserPassword se utiliza para encriptar las contraseñas con la
// llave simétrica de la base de datos.
func EncryptUserPassword(password string) (encrypted string, err error) {
	var (
		passwordBytes []byte = []byte(password)
		nonce         []byte
	)

	nonce = make([]byte, keys.GetDBCipherGCM().NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		log.Printf("Nonce (random 32 bytes) could not be generated: %s.", err)
		return "", err
	}

	return hex.EncodeToString(keys.GetDBCipherGCM().Seal(nonce, nonce, passwordBytes, nil)), err
}

// DecryptUserPassword se utiliza para desencriptar las contraseñas con la
// llave simétrica de la base de datos.
func DecryptUserPassword(encrypted []byte) (password string, err error) {
	var passwordBytes []byte

	if len(encrypted) < keys.GetDBCipherGCM().NonceSize() {
		return "", fmt.Errorf("EncryptedBytes too short.")
	}

	passwordBytes, err = keys.GetDBCipherGCM().Open(
		nil,
		encrypted[:keys.GetDBCipherGCM().NonceSize()],
		encrypted[keys.GetDBCipherGCM().NonceSize():],
		nil,
	)
	return string(passwordBytes), err
}
