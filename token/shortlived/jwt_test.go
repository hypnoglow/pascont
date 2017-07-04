package shortlived

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

func TestNewJWTString(t *testing.T) {
	type jwtData struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}

	issuer := "theIssuer"
	subject := "theSubject"
	key := []byte("key")

	cases := []struct {
		payload       interface{}
		expectedError bool
	}{
		// Successful.
		{
			payload:       jwtData{ID: 123, Name: "SomeName"},
			expectedError: false,
		},
		// Illegal type for jwt claims.
		{
			payload:       make(chan int),
			expectedError: true,
		},
	}

	for i, c := range cases {
		_, err := NewJWTString(issuer, subject, c.payload, key)

		// TODO: validate token.

		if c.expectedError && err == nil {
			t.Errorf(
				"testcase %d: Expected error but got %v",
				i,
				err,
			)
		}

		if !c.expectedError && err != nil {
			t.Errorf(
				"testcase %d: Expected no error but got %v",
				i,
				err,
			)
		}
	}
}

func TestParseJWTString(t *testing.T) {
	issuer := "theIssuer"
	subject := "theSubject"
	key := []byte("key")
	exampleClaims := map[string]interface{}{"a": "value"}
	exampleTokenHMAC, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   issuer,
		"sub":   subject,
		subject: exampleClaims,
	}).SignedString(key)

	cases := []struct {
		token          string
		issuer         string
		key            []byte
		expectedClaims jwt.MapClaims
		expectedError  error
	}{
		// Malformed claims
		{
			token:         "abc",
			expectedError: jwt.NewValidationError("token contains an invalid number of segments", jwt.ValidationErrorMalformed),
		},
		// Invalid key
		{
			token:         exampleTokenHMAC,
			issuer:        issuer,
			key:           []byte("invalid key"),
			expectedError: &jwt.ValidationError{Inner: fmt.Errorf("signature is invalid"), Errors: jwt.ValidationErrorSignatureInvalid},
		},
		// Invalid issuer
		{
			token:         exampleTokenHMAC,
			issuer:        "invalid issuer",
			key:           key,
			expectedError: fmt.Errorf("Invalid token %s", exampleTokenHMAC),
		},
		// TODO: test not a HMAC sign.
		// Successful
		{
			token:  exampleTokenHMAC,
			issuer: issuer,
			key:    key,
			expectedClaims: jwt.MapClaims{
				"iss":   issuer,
				"sub":   subject,
				subject: exampleClaims,
			},
		},
	}

	for i, c := range cases {
		claims, err := ParseJWTString(c.token, c.issuer, c.key)

		// Do not verify time-related fields.
		delete(claims, "iat")
		delete(claims, "exp")

		if !reflect.DeepEqual(claims, c.expectedClaims) {
			t.Errorf(
				"testcase %d: Expected claims to be %#v but got %#v",
				i,
				c.expectedError,
				err,
			)
		}

		if !reflect.DeepEqual(err, c.expectedError) {
			t.Errorf(
				"testcase %d: Expected error to be %#v but got %#v",
				i,
				c.expectedError,
				err,
			)
		}
	}
}
