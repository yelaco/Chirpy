package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (cf *apiConfig) validatePassword(password string, email string) (LoginResponse, bool) {
	login, err := cf.db.GetUserByEmail(email)
	if err != nil {
		return LoginResponse{}, false
	}
	e := bcrypt.CompareHashAndPassword([]byte(login.Password), []byte(password))
	if e != nil {
		return LoginResponse{}, false
	}
	return LoginResponse{
		Id:    login.Id,
		Email: login.Email,
	}, true
}

func (cf *apiConfig) createAccessToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(1) * time.Hour)),
		Subject:   fmt.Sprintf("%d", userId),
	})
	accessToken, err := token.SignedString([]byte(cf.jwtSecret))
	if err != nil {
		return "", errors.New("createAccessToken: " + err.Error())
	}
	return accessToken, nil
}

func (cf *apiConfig) createRefreshToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-refresh",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(60*24) * time.Hour)),
		Subject:   fmt.Sprintf("%d", userId),
	})
	refreshToken, err := token.SignedString([]byte(cf.jwtSecret))
	if err != nil {
		return "", errors.New("createRefreshToken: " + err.Error())
	}
	return refreshToken, nil
}

func (cf *apiConfig) validateJWT(signedToken string) (string, error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cf.jwtSecret), nil
	})
	if err != nil {
		return "", errors.New("validateJWT 1: " + err.Error())
	}

	expirationTime, err := token.Claims.GetExpirationTime()
	if err != nil {
		return "", errors.New("validateJWT 2: " + err.Error())
	}

	if !token.Valid || time.Now().UTC().After(expirationTime.Time) {
		return "", errors.New("validateJWT 3: invalid or expired token")
	}

	userId, err := token.Claims.GetSubject()
	if err != nil {
		return "", errors.New("validateJWT 4: " + err.Error())
	}

	return userId, nil
}

// GetBearerToken -
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No Authorization header included")
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

func (cf *apiConfig) validateAccessToken(signedToken string) (bool, error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cf.jwtSecret), nil
	})
	if err != nil {
		return false, errors.New("validateAccessToken: " + err.Error())
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return false, errors.New("validateAccessToken: " + err.Error())
	}
	if issuer != "chirpy-access" {
		return false, errors.New("validateAccessToken: invalid issuer")
	}
	return true, nil
}

func (cf *apiConfig) validateRefreshToken(signedToken string) (bool, error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cf.jwtSecret), nil
	})
	if err != nil {
		return false, errors.New("validateRefreshToken: " + err.Error())
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return false, errors.New("validateRefreshToken: " + err.Error())
	}
	if issuer != "chirpy-refresh" {
		fmt.Println(issuer)
		return false, errors.New("validateRefreshToken: invalid issuer")
	}
	if revoked, _ := cf.revokedYet(signedToken); revoked {
		return false, errors.New("validateRefreshToken: revoked refresh token")
	}
	return true, nil
}

func (cf *apiConfig) userIdFromToken(signedToken string) (int, error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cf.jwtSecret), nil
	})
	if err != nil {
		return -1, errors.New("userIdFromToken: " + err.Error())
	}

	userId, err := token.Claims.GetSubject()
	if err != nil {
		return -1, errors.New("userIdFromToken: " + err.Error())
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		return -1, errors.New("userIdFromToken: " + err.Error())
	}
	return id, nil
}

func (cf *apiConfig) revokedYet(signedToken string) (bool, error) {
	_, err := cf.db.GetRefreshToken(signedToken)
	if err != nil {
		return false, errors.New("revokedYet: " + err.Error())
	}
	return true, nil
}
