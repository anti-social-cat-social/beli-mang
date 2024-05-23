package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Custom claim used to generate token
type CustomClaim struct {
	Uuid string
	Role string
	jwt.RegisteredClaims
}

var (
	key            string
	expirationTime int64
)

func getKey() []byte {
	key = os.Getenv("JWT_SECRET")

	return []byte(key)
}

type TokenData struct {
	ID   string
	Name string
	Role string
}

func GenerateToken(data TokenData) (string, error) {
	// Set expiration time
	expirationTime = 8

	claims := CustomClaim{
		data.ID,
		data.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationTime) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "halosuster",
			Subject:   data.Name,
			ID:        data.ID,
		},
	}

	// Create token based on the claim above
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(getKey())
	if err != nil {
		return "", err
	}

	return s, nil
}

// Validate token
func ValidateToken(tokenChecked string) (*CustomClaim, error) {
	// Get user token from environment variable
	tokenString := getKey()

	// Parse token
	token, err := jwt.ParseWithClaims(tokenChecked, &CustomClaim{}, func(t *jwt.Token) (interface{}, error) {
		// Validate the token alg
		// If alg is not valid, return error
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signature")
		}

		return []byte(tokenString), nil
	})

	// Handle if there is any error from parsed token
	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		return nil, errors.New("invalid token form")

	case errors.Is(err, jwt.ErrSignatureInvalid):
		return nil, errors.New("invalid token signature")

	case errors.Is(err, jwt.ErrTokenExpired):
		return nil, errors.New("your time has come")

	case err != nil:
		return nil, err
	}

	// If claims is valid, return it
	if claims, ok := token.Claims.(*CustomClaim); ok {
		return claims, nil
	} else {
		// Otherwise, return error message
		return nil, errors.New("cannot handle token")
	}
}
