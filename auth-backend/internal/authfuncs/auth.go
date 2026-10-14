// Package authfuncs implementa los objetos y librerías necesarias
// para los protocolos de autenticación (el protocolo CHAP y el uso de
// JSON Web Tokens).
package authfuncs

import (
	//"encoding/base64"
	"fmt"
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/dboper"
	"log"
	"sync"
	"time"
)

// CHAPS_CLEANUP_TIME define el tiempo máximo para intentar completar
// la autenticación antes de ser limpiado de la memoria.
const CHAPS_CLEANUP_TIME time.Duration = 10 * time.Minute // 10m0s

// CHAPS es un mapa (diccionario) que relaciona el ID de un usuario
// con su intento de autenticación. Si el mapa no encuentra el ID
// entonces el intento es inválido directamente, si lo encuentra hace
// la comprobación con el intento de autenticación.
// Se borran los intentos de CHAPS después de 10 minutos de la vida
// del intento, aunque el máximo sean 5 minutos.
var (
	CHAPS      map[uint32]CHAP = make(map[uint32]CHAP)
	CHAPSMutex sync.RWMutex
)

// Inicia el limpiador de la memoria de intentos de autenticación
func init() {
	var ticker *time.Ticker = time.NewTicker(CHAPS_CLEANUP_TIME)
	go cleanCHAPS(ticker)
}

// NewAuthTry permite iniciar un intento de autenticación para un
// usuario con su ID.
func NewAuthTry(user_id uint32) (challenge []byte, err error) {
	var (
		chap         CHAP
		passwordHash []byte
	)

	if passwordHash, err = dboper.GetUserPasswordHash(user_id); err != nil {
		log.Printf("The new authentication attempt, by user with id: %d, has failed: %s", user_id, err)
		return nil, err
	}

	if chap, err = NewCHAP(passwordHash); err != nil {
		return nil, err
	}

	CHAPSMutex.Lock()
	CHAPS[user_id] = chap
	CHAPSMutex.Unlock()

	return chap.Challenge(), nil
}

// Authenticate busca si el usuario con ID realizó un intento de
// autenticación, y si sí, borra el intento de autenticación de la
// memoria y hace la comprobación del intento.
func Authenticate(user_id uint32, response [32]byte) (ok bool, err error) {
	var (
		authAttemptTime time.Time = time.Now()
		chap            CHAP
	)

	if chap, ok = getCHAP(user_id); !ok {
		return false, fmt.Errorf(
			"The authenticacion attempt, by user with id: %d, has already expired, or it's incorrect.",
			user_id,
		)
	}
	deleteCHAP(user_id)

	return chap.Authenticate(response, authAttemptTime)
}

// getCHAP es una función helper que permite obtener el intento de
// autenticación con el ID del usuario.
func getCHAP(user_id uint32) (chap CHAP, ok bool) {
	CHAPSMutex.RLock()
	chap, ok = CHAPS[user_id]
	CHAPSMutex.RUnlock()

	return chap, ok
}

// deleteCHAP es una función helper que permite borrar el intento de
// autenticación con el ID del usuario.
func deleteCHAP(user_id uint32) {
	CHAPSMutex.Lock()
	delete(CHAPS, user_id)
	CHAPSMutex.Unlock()
}

// cleanCHAPS es la función que periódicamente itera por el mapa de
// intentos de autenticación, y elimina aquellos cuya expiración haya
// sido superada pero cuyos intentos no hayan sido autoeliminados.
func cleanCHAPS(ticker *time.Ticker) {
	for {
		<-ticker.C
		var currentTime time.Time = time.Now()
		CHAPSMutex.Lock()
		for i, c := range CHAPS {
			// Compara si el tiempo actual supera al tiempo de
			// expiración del intento de autenticación
			if c.expirationTime.Compare(currentTime) > -1 {
				delete(CHAPS, i)
			}
		}
		CHAPSMutex.Unlock()
		log.Print("Los intentos de autenticación expirados, han sido limpiados.")
	}
}
