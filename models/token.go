package models

import (
	"encoding/json"
	"errors"
	"shortener/configuration"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Token struct represents JWT token
type Token struct {
	Value string `json:"token"`
}

// Claims represents JWT token claims which are used in the app
type Claims struct {
	User
}

// GenerateAuthToken returns auth token for the user
func GenerateAuthToken(user User) (Token, error) {
	ttl := configuration.GetConfiguration().TokenTTL
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": struct {
			Login string `json:"login"`
			ID    string `json:"id"`
		}{
			Login: user.Login,
			ID:    user.ID,
		},
		"exp": time.Now().Add(time.Second * time.Duration(ttl)).Unix(),
	})
	secret := configuration.GetConfiguration().TokenSecret
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return Token{}, err
	}

	return Token{Value: tokenString}, nil
}

// GetClaims return major claims from token
func (t *Token) GetClaims() (*Claims, error) {
	token, err := jwt.Parse(t.Value, func(parsedToken *jwt.Token) (interface{}, error) {
		if _, ok := parsedToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Incorrect method in the token")
		}

		return []byte(configuration.GetConfiguration().TokenSecret), nil
	})

	if !token.Valid || err != nil {
		return nil, errors.New("Incorrect authentication token")
	}

	if err = token.Claims.Valid(); err != nil {
		return nil, err
	}

	rawClaims, ok := token.Claims.(jwt.MapClaims)

	if !rawClaims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, errors.New("Invalid expires at claim in the token")
	}

	if !ok {
		return nil, errors.New("Unknown type of claims")
	}

	rawUser, ok := rawClaims["user"]

	if !ok {
		return nil, errors.New("User claim is missing")
	}

	var user User

	stringifiedUser, err := json.Marshal(rawUser)

	if err != nil {
		return nil, errors.New("User claim has unknown type")
	}

	err = json.Unmarshal([]byte(stringifiedUser), &user)

	if err != nil {
		return nil, errors.New("User claim has unknown type")
	}

	return &Claims{User: user}, nil
}
