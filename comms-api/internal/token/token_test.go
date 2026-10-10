package token

import (
	"errors"
	"testing"
)

func TestAuthenticateRejectsMissingOrMalformedBearerToken(t *testing.T) {
	for _, authorization := range []string{"", "token", "Bearer ", "Bearer not-a-jwt"} {
		if _, err := Authenticate(authorization); !errors.Is(err, ErrInvalidToken) {
			t.Errorf("Authenticate(%q) error = %v, want ErrInvalidToken", authorization, err)
		}
	}
}
