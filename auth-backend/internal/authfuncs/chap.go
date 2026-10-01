package authfuncs

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"time"
)

const MAX_CHAP_TIME int64 = 300

// Para iniciar un nuevo intento de autenticación
type CHAP struct {
	unixTimeStart int64
	challenge     *[]byte
	expected      [32]byte
}

// Debe generarse cada que se quiera un nuevo intento
// de autenticación
func NewCHAP(password string) CHAP {
	var challengeBuffer []byte = make([]byte, 64)
	if _, err := rand.Read(challengeBuffer); err != nil {
		log.Fatalf("Error while generating challenge: %s", err)
	}

	return CHAP{
		unixTimeStart: time.Now().Unix(),
		challenge:     &challengeBuffer,
		expected:      sha256.Sum256(append([]byte(password), challengeBuffer...)),
	}
}

// Función getter que regresa el challenge
func (chap *CHAP) Challenge() []byte {
	return *chap.challenge
}

// Función que compara la respuesta con el reto
func (chap *CHAP) Authenticate(response [32]byte, authTime int64) (ok bool, err error) {
	if authTime < chap.unixTimeStart || authTime-chap.unixTimeStart > MAX_CHAP_TIME {
		return false, fmt.Errorf(
			"Authentication time (started: %s) exceeded or invalid: %s.",
			time.Unix(chap.unixTimeStart, 0).UTC(),
			time.Unix(authTime, 0).UTC())
	}
	if response != chap.expected {
		return false, fmt.Errorf("Wrong response!")
	}
	return true, nil
}
