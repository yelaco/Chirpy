package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yelaco/Chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

var ErrNoAuthHeaderIncluded = errors.New("not auth header included in request")

func CheckPasswordHash(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func CreateRefreshToken() (string, error) {
	r := make([]byte, 32)
	n, err := rand.Read(r)
	if err != nil || n != 32 {
		return "", err
	}

	return hex.EncodeToString(r), nil
}

func ValidateRefreshToken(info database.TokenInfo) (int, error) {
	if info.ExpiresAt.Before(time.Now()) {
		return 0, errors.New("refresh token expired")
	}
	return info.UserID, nil
}

func CreateAccessToken(userID, expiresInSecond int, tokenSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresInSecond) * time.Second)),
		Subject:   fmt.Sprint(userID),
	})

	return token.SignedString([]byte(tokenSecret))
}

func ValidateAccessToken(tokenString, tokenSecret string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return "", err
	}

	expiresAt, err := token.Claims.GetExpirationTime()
	if err != nil {
		return "", err
	}

	if expiresAt.Before(time.Now()) {
		return "", errors.New("token expired")
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil || issuer != "chirpy" {
		return "", errors.New("invalid token issuer")
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", errors.New("invalid token subject")
	}

	return userIDString, nil
}

func GetBearerToken(header http.Header) (string, error) {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

func GetApiKey(header http.Header) (string, error) {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}
