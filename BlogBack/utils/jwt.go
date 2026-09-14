package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	SessionCookieName = "cspona_token"
	SessionMaxAge     = 7 * 24 * time.Hour
)

type SessionClaims struct {
	UserID   uint   `json:"uid"`
	Username string `json:"uname"`
	RoleQx   string `json:"qx"`
	jwt.RegisteredClaims
}

func jwtSecret() ([]byte, error) {
	if len(aesKey) == 0 {
		return nil, errors.New("AES/JWT 密钥未初始化")
	}
	return aesKey, nil
}

func SignSessionToken(userID uint, username, roleQx string) (string, error) {
	secret, err := jwtSecret()
	if err != nil {
		return "", err
	}
	now := time.Now()
	claims := SessionClaims{
		UserID:   userID,
		Username: username,
		RoleQx:   roleQx,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(SessionMaxAge)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "cspona",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseSessionToken(tokenStr string) (*SessionClaims, error) {
	secret, err := jwtSecret()
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(tokenStr, &SessionClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*SessionClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
