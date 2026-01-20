package auth

import (
	"cluster-agent/internal/auth/permissions"
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserId      string                   `json:"sub"`
	Permissions []permissions.Permission `json:"permissions"`
}

func ParseToken(token string, claims *UserClaims, pubKey *rsa.PublicKey) (*jwt.Token, error) {
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	return jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return pubKey, nil
	})
}
