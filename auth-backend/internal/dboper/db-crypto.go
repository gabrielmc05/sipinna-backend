package dboper

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/keys"
	"log"
)

// EncryptHashUserPassword se utiliza para aplicarles el algoritmo Argon2id
// y luego, encriptar las contraseñas con la llave simétrica de la base de
// datos, se usa en el registro.
func EncryptHashUserPassword(password string) (encryptedHash string, err error) {
	var (
		passwordBytes [32]byte = sha256.Sum256([]byte(password))
		nonce         []byte
	)

	nonce = make([]byte, keys.GetDBCipherGCM().NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		log.Printf("Nonce (random 32 bytes) could not be generated: %s.", err)
		return "", err
	}

	return hex.EncodeToString(keys.GetDBCipherGCM().Seal(nonce, nonce, passwordBytes[:], nil)), err
}

// DecryptUserPasswordHash se utiliza para desencriptar los Hash (Argon2id)
// de contraseñas con la llave simétrica de la base de datos.
func DecryptUserPasswordHash(encrypted []byte) (password string, err error) {
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
