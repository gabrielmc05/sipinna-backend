package token

import (
	"github.com/gabrielmc05/sipinna-backend/comms-api/internal/dboper"
	"github.com/golang-jwt/jwt/v5"
)

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
