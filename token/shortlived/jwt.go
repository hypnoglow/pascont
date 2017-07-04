// THIS IS WIP.

package shortlived

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	jwtVersion         = "0.1.0"
	jwtSessionDuration = time.Minute * 5
)

// NewJWTString returns new encoded jwt.
func NewJWTString(issuer, subject string, payload interface{}, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// Reserved claims
		//"jti": "TODO TOKEN ID ??",
		"iss": issuer,
		"sub": subject,
		"exp": time.Now().Add(jwtSessionDuration).Unix(),
		"iat": time.Now().Unix(),
		// Private claims
		"v":     jwtVersion,
		subject: payload,
	})

	return token.SignedString(key)
}

// ParseJWTString parses and validates the tokenString and returns claims if valid.
func ParseJWTString(tokenString, issuer string, key []byte) (jwt.MapClaims, error) {
	// Parse the tokenString.
	// WARNING! This also validates claims (for example "exp" claim).
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid || !claims.VerifyIssuer(issuer, true) {
		return nil, fmt.Errorf("Invalid token %s", tokenString)
	}

	return claims, nil
}
