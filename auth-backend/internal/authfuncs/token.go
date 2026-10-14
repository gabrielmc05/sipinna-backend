package authfuncs

import (
	"github.com/gabrielmc05/sipinna-backend/auth-backend/internal/dboper"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"strconv"
	"time"
)

// TOKEN_EXPIRATION_TIME es el tiempo máximo que debe durar un JWT
const TOKEN_EXPIRATION_TIME time.Duration = 24 * time.Hour // 24h0m0s

// GenerateToken crea un nuevo token con el ID del usuario
func GenerateToken(user_id uint32) (string, error) {
	var (
		err        error
		key        []byte
		token      *jwt.Token
		expiration *jwt.NumericDate = jwt.NewNumericDate(time.Now().Add(TOKEN_EXPIRATION_TIME))
	)

	if key, err = dboper.NewUserSignatureKey(user_id); err != nil {
		log.Printf("Error creating JWT for user with id %d: %s", user_id, err)
		return "", err
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: expiration,
		Subject:   strconv.FormatUint(uint64(user_id), 10),
	})

	return token.SignedString(key)
}

// VerifyToken verifica el JWT de un usuario con su firma.
// NOTA: no me di cuenta que se puede obtener el ID directo del JWT.
func VerifyToken(tokenString string, user_id uint32) (ok bool, err error) {
	var token *jwt.Token

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err = jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return dboper.GetUserSignatureKey(user_id)
	})
	if err != nil || !token.Valid {
		return false, err
	}

	return true, nil
}
