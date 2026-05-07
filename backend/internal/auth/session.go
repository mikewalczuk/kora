package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const cookieName = "session"
const tokenTTL = 24 * time.Hour

func signToken(secret, username, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"uid": userID,
		"exp": time.Now().Add(tokenTTL).Unix(),
	})
	return token.SignedString([]byte(secret))
}

func verifyToken(secret, tokenStr string) (username, userID string, err error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}
	username, ok = claims["sub"].(string)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}
	userID, ok = claims["uid"].(string)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}
	return username, userID, nil
}
