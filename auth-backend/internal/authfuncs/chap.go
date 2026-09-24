package authfuncs

import (
	"crypto/rand"
	"crypto/sha256"
	"log"
)

// Para iniciar un nuevo intento de autenticación
type CHAP struct {
	correctAnswer [32]byte
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
		challenge:     &challengeBuffer,
		correctAnswer: sha256.Sum256(append([]byte(password), challengeBuffer...)),
	}
}

// Función getter que regresa el challenge
func (chap *CHAP) Challenge() []byte {
	return *chap.challenge
}

// Función que compara la respuesta con el reto
func (chap *CHAP) Authenticate(answer [32]byte) bool {
	return answer == chap.correctAnswer
}
