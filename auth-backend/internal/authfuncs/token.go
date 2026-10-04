package authfuncs

import (
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/dboper"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"strconv"
	"time"
)

const TOKEN_EXPIRATION_TIME time.Duration = 24 * time.Hour // 24h0m0s

func GenerateToken(user_id uint32) (string, error) {
	var (
		err        error
		key        []byte
		token      *jwt.Token
		expiration *jwt.NumericDate = jwt.NewNumericDate(time.Now().Add(TOKEN_EXPIRATION_TIME))
	)

	if key, err = dboper.NewUserSignatureKey(user_id); err != nil {
		log.Printf("Error creating JWT for user with id %d: %s", user_id, err)
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: expiration,
		Subject:   strconv.FormatUint(uint64(user_id), 10),
	})

	return token.SignedString(key)
}

/* TESTING PURPOSES ONLY */
//func GenerateToken(user_id uint32) (key []byte, tokenStr string, err error) {
//	var (
//		token      *jwt.Token
//		expiration *jwt.NumericDate = jwt.NewNumericDate(time.Now().Add(TOKEN_EXPIRATION_TIME))
//	)
//
//	if key, err = dboper.NewUserSignatureKey(user_id); err != nil {
//		log.Printf("Error creating JWT for user with id %d: %s", user_id, err)
//	}
//
//	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
//		ExpiresAt: expiration,
//		Subject:   strconv.FormatUint(uint64(user_id), 10),
//	})
//
//	tokenStr, err = token.SignedString(key)
//	log.Println(base64.StdEncoding.EncodeToString(key))
//
//	return key, tokenStr, err
//}
