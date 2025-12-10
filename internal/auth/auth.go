package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"runtime"
	"strings"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	params := &argon2id.Params{
	Memory:      128 * 1024,
	Iterations:  4,
	Parallelism: uint8(runtime.NumCPU()),
	SaltLength:  16,
	KeyLength:   32,
}
	hash, err := argon2id.CreateHash(password, params)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	result, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return result, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	errMsg := errors.New("bad authorization header")
	bearerToken := headers.Get("Authorization")
	if bearerToken == "" {
		return "", errMsg
	}
	bearerToken = strings.TrimSpace(bearerToken)
	if !strings.HasPrefix(bearerToken, "Bearer ") {
		return "", errMsg
	}
	token := strings.TrimSpace(strings.TrimPrefix(bearerToken, "Bearer"))
	if token == "" {
		return "",errMsg
	}
	return token, nil
}

func MakeRefreshToken() (string, error) {
	random := make([]byte, 32)
	if bytes, err := rand.Read(random); bytes != len(random) || err != nil {
		return "", err
	}
	return hex.EncodeToString(random), nil
}