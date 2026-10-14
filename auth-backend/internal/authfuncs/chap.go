// Es una implementación chafa del protocolo CHAP, para la
// autenticación. Se utiliza dentro del paquete, aunque los objetos
// estén públicos.
package authfuncs

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"hash"
	"log"
	"time"
)

// MAX_CHAP_TIME define el tiempo máximo para completar la
// autenticación antes de expirar.
const MAX_CHAP_TIME time.Duration = 5 * time.Minute // 5m0s

// CHAP es para iniciar un nuevo intento de autenticación
type CHAP struct {
	expirationTime time.Time
	challenge      []byte
	expected       [sha256.Size]byte
}

// NewCHAP crea un nuevo intento de autenticación con un tiempo de
// expiración con base en el tiempo actual generarse cada que se
// quiera un nuevo intento de autenticación.
func NewCHAP(passwordHash []byte) (CHAP, error) {
	var (
		mac       hash.Hash = hmac.New(sha256.New, passwordHash)
		challenge []byte    = make([]byte, 64)
		expected  [sha256.Size]byte
	)

	if _, err := rand.Read(challenge); err != nil {
		err = fmt.Errorf("Error while generating challenge: %s", err)
		log.Print(err)
		return CHAP{}, err
	}

	mac.Write(challenge)
	copy(expected[:], mac.Sum(nil))

	return CHAP{
		expirationTime: time.Now().Add(MAX_CHAP_TIME),
		challenge:      challenge,
		expected:       expected,
	}, nil
}

// Challenge es un método getter que regresa el challenge
func (chap *CHAP) Challenge() []byte {
	return chap.challenge
}

// Authenticate es el método que compara la respuesta del intento de
// autenticación con la respuesta esperada.
func (chap *CHAP) Authenticate(response [32]byte, authTime time.Time) (ok bool, err error) {
	if authTime.After(chap.expirationTime) {
		return false, fmt.Errorf(
			"Authentication time (expiration at: %s) exceeded or invalid: %s.",
			chap.expirationTime.UTC(),
			authTime.UTC())
	}
	if subtle.ConstantTimeCompare(response[:], chap.expected[:]) != 1 {
		return false, fmt.Errorf("Wrong response!")
	}
	return true, nil
}
