package authfuncs

import (
	//"encoding/base64"
	"fmt"
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/dboper"
	"log"
	"sync"
	"time"
)

const CHAPS_CLEANUP_TIME time.Duration = 600 * time.Second // 10m0s

var (
	CHAPS      map[uint32]CHAP = make(map[uint32]CHAP)
	CHAPSMutex sync.RWMutex
)

func init() {
	var ticker *time.Ticker = time.NewTicker(CHAPS_CLEANUP_TIME)
	go cleanCHAPS(ticker)
}

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

func getCHAP(user_id uint32) (chap CHAP, ok bool) {
	CHAPSMutex.RLock()
	chap, ok = CHAPS[user_id]
	CHAPSMutex.RUnlock()

	return chap, ok
}

func deleteCHAP(user_id uint32) {
	CHAPSMutex.Lock()
	delete(CHAPS, user_id)
	CHAPSMutex.Unlock()
}

func cleanCHAPS(ticker *time.Ticker) {
	for {
		<-ticker.C
		var currentTime time.Time = time.Now()
		CHAPSMutex.Lock()
		for i, c := range CHAPS {
			if c.expirationTime.Compare(currentTime) > -1 {
				delete(CHAPS, i)
			}
		}
		CHAPSMutex.Unlock()
		log.Print("La lista de CHAP ha sido limpiada :3")
	}
}
