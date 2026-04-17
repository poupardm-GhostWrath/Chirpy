package auth

import (
	"errors"
	"strings"
	"net/http"

	"github.com/alexedwards/argon2id"
)

const tokenIssuer = "chirpy-access"

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	if header == "" {
		return "", errors.New("No Auth header included in request")
	}
	splitHeader := strings.Split(header, " ")
	if len(splitHeader) < 2 || splitHeader[0] != "Bearer" {
		return "", errors.New("Malformed authorization header")
	}
	return splitHeader[1], nil
}