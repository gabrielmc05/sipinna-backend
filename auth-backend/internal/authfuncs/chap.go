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
	expected      [32]byte
	challenge     *[]byte
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
func (chap *CHAP) Authenticate(answer [32]byte, time int64) (ok bool, err error) {
	if time < chap.unixTimeStart || time-chap.unixTimeStart > MAX_CHAP_TIME {
		return false, fmt.Errorf("Authentication time exceeded or invalid.")
	}
	if answer != chap.expected {
		return false, fmt.Errorf("Wrong answer!")
	}
	return true, nil
}
