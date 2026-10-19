package token

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gabrielmc05/sipinna-backend/comms-api/internal/dboper"
	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid bearer token")

func Authenticate(authorization string) (uint32, error) {
	if !strings.HasPrefix(authorization, "Bearer ") {
		return 0, ErrInvalidToken
	}
	tokenString := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
	if tokenString == "" {
		return 0, ErrInvalidToken
	}

	claims := &jwt.RegisteredClaims{}
	var keyLookupErr error
	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(parsed *jwt.Token) (any, error) {
		if parsed.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		userID, err := strconv.ParseUint(claims.Subject, 10, 32)
		if err != nil || userID == 0 {
			return nil, ErrInvalidToken
		}
		key, err := dboper.GetUserSignatureKey(uint32(userID))
		if err != nil {
			keyLookupErr = err
			return nil, err
		}
		return key, nil
	})
	if err != nil {
		if keyLookupErr != nil {
			return 0, keyLookupErr
		}
		return 0, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	if !parsed.Valid || claims.ExpiresAt == nil {
		return 0, ErrInvalidToken
	}
	userID, err := strconv.ParseUint(claims.Subject, 10, 32)
	if err != nil || userID == 0 {
		return 0, ErrInvalidToken
	}
	return uint32(userID), nil
}

func VerifyToken(tokenString string, user_id uint32) (bool, error) {
	// Verifica la firma, la expiración y que la cuenta esté activa.
	tokenUserID, err := Authenticate(tokenString)
	if err != nil {
		return false, err
	}

	// El ID del JWT debe coincidir con el usuario solicitado.
	if tokenUserID != user_id {
		return false, ErrInvalidToken
	}

	return true, nil
}
