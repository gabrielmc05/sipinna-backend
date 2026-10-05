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

const MAX_CHAP_TIME time.Duration = 300 * time.Second // 5m0s

// Para iniciar un nuevo intento de autenticación
type CHAP struct {
	expirationTime time.Time
	challenge      []byte
	expected       [sha256.Size]byte
}

// Debe generarse cada que se quiera un nuevo intento
// de autenticación
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

// Función getter que regresa el challenge
func (chap *CHAP) Challenge() []byte {
	return chap.challenge
}

// Función que compara la respuesta con el reto
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
