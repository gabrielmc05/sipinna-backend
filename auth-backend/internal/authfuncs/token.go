package authfuncs

import (
	"fmt"
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

func VerifyToken(tokenString string, user_id uint32) (bool, error) {
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claims := &jwt.RegisteredClaims{}

	parsed, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("algoritmo de firma no permitido")
			}

			return dboper.GetUserSignatureKey(user_id)
		},
	)
	if err != nil {
		return false, err
	}

	if !parsed.Valid || claims.ExpiresAt == nil {
		return false, fmt.Errorf("token inválido o sin expiración")
	}

	if claims.Subject != strconv.FormatUint(uint64(user_id), 10) {
		return false, fmt.Errorf("el token no corresponde al usuario")
	}

	return true, nil
}
