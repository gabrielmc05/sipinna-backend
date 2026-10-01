package authfuncs

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"time"
)

const MAX_CHAP_TIME time.Duration = time.Duration(300 * time.Second) // 5m0s

// Para iniciar un nuevo intento de autenticación
type CHAP struct {
	expirationTime time.Time
	challenge      *[]byte
	expected       [32]byte
}

// Debe generarse cada que se quiera un nuevo intento
// de autenticación
func NewCHAP(passwordHash string) CHAP {
	var challengeBuffer []byte = make([]byte, 64)
	if _, err := rand.Read(challengeBuffer); err != nil {
		log.Fatalf("Error while generating challenge: %s", err)
	}

	return CHAP{
		expirationTime: time.Now().Add(MAX_CHAP_TIME),
		challenge:      &challengeBuffer,
		expected:       sha256.Sum256(append([]byte(passwordHash), challengeBuffer...)),
	}
}

// Función getter que regresa el challenge
func (chap *CHAP) Challenge() []byte {
	return *chap.challenge
}

// Función que compara la respuesta con el reto
func (chap *CHAP) Authenticate(response [32]byte, authTime time.Time) (ok bool, err error) {
	if chap.expirationTime.Compare(authTime) == -1 {
		return false, fmt.Errorf(
			"Authentication time (expiration at: %s) exceeded or invalid: %s.",
			chap.expirationTime.UTC(),
			authTime.UTC())
	}
	if response != chap.expected {
		return false, fmt.Errorf("Wrong response!")
	}
	return true, nil
}
