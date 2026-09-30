package dboper

import (
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/keys"
	"log"
)

func EncryptUserPassword(password string) (encrypted string, err error) {
	var (
		passwordBytes []byte = []byte(password)
		gcm           cipher.AEAD
		nonce         []byte
	)

	if gcm, err = cipher.NewGCM(keys.GetDBCipherBlock()); err != nil {
		log.Printf("DB GCM could not be made: %s.", err)
		return "", err
	}

	nonce = make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		log.Printf("Nonce (random 32 bytes) could not be generated: %s.", err)
		return "", err
	}

	return string(gcm.Seal(nonce, nonce, passwordBytes, nil)), nil
}

func DecryptUserPassword(encrypted string) (password string, err error) {
	var (
		encryptedBytes []byte = []byte(encrypted)
		passwordBytes  []byte
		gcm            cipher.AEAD
	)

	if gcm, err = cipher.NewGCM(keys.GetDBCipherBlock()); err != nil {
		log.Printf("DB GCM could not be made: %s.", err)
		return "", err
	}

	if len(encryptedBytes) < gcm.NonceSize() {
		return "", fmt.Errorf("EncryptedBytes too short.")
	}

	passwordBytes, err = gcm.Open(nil, encryptedBytes[:gcm.NonceSize()], encryptedBytes[gcm.NonceSize():], nil)
	return string(passwordBytes), err
}
