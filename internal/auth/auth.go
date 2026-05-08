package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer: "chirpy-access",
		IssuedAt: &jwt.NumericDate{time.Now()},
		ExpiresAt: &jwt.NumericDate{time.Now().Add(expiresIn)},
		Subject: userID.String(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedStr, err := jwtToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedStr, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	userID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.Parse(userID)
}

func GetBearerToken(headers http.Header) (string, error) {
	bearerToken := headers.Get("Authorization")
	if len(bearerToken) == 0 {
		return "", errors.New("No Auth token found")
	}

	
	return strings.TrimSpace(strings.TrimPrefix(bearerToken, "Bearer")), nil
}

func MakeRefreshToken() string {
	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	return hex.EncodeToString(randBytes)
}

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if len(apiKey) == 0 {
		return "", errors.New("No api key found")
	}

	
	return strings.TrimSpace(strings.TrimPrefix(apiKey, "ApiKey")), nil
}
