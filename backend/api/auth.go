package api

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

// Auth contains the secret key for JWT token generation and validation.
type Auth struct {
	secretKey []byte
}

// Constructs a new Auth instance with the provided secret key.
func NewAuth(secretKey string) *Auth {
	return &Auth{secretKey: []byte(secretKey)}
}

// Creates a JWT access token using an authenticated user id, returns the encoded access token.
func (a *Auth) GenerateAccessToken(userID int) (string, error) {

	claims := jwt.MapClaims{
		"user_id": strconv.Itoa(userID),
		"exp":     time.Now().Add(time.Minute * 45).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the access token with the secret key.
	tokenString, err := token.SignedString(a.secretKey)
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func (a *Auth) ValidateAccessToken(tokenString string) (int, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if token.Method != jwt.SigningMethodHS256 {
			fmt.Println("error is HERE at signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.secretKey, nil
	})

	if err != nil {
		fmt.Println("error parsing token")
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// Extract the user id claim from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}
	userID, ok := claims["user_id"].(string)

	if !ok {
		return 0, fmt.Errorf("invalid user id claim")
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return 0, fmt.Errorf("invalid token claims")
	}

	return id, nil
}
