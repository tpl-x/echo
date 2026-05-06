package token

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tpl-x/echo/internal/pkg/usersource"
)

var _ RefreshTokenProvider = (*CreditGenerator)(nil)

func (c *CreditGenerator) CreateRefreshToken(user usersource.UserProvider) (freshToken string, exp int64, err error) {
	refreshTokenExp := time.Now().Add(time.Duration(c.refreshTokenExpSec) * time.Second)
	claims := &JwtCustomRefreshClaims{
		ID: user.UserId(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExp),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := t.SignedString([]byte(c.refreshSecret))
	if err != nil {
		return "", 0, err
	}
	return refreshToken, refreshTokenExp.Unix(), nil
}

func (c *CreditGenerator) ParseRefreshToken(refreshToken any) (claims *JwtCustomRefreshClaims, err error) {
	switch token := refreshToken.(type) {
	case string:
		return c.parseRefreshTokenString(token)
	case *jwt.Token:
		return refreshClaimsFromToken(token)
	default:
		return nil, errors.New("invalid refreshToken")
	}
}

func (c *CreditGenerator) parseRefreshTokenString(refreshToken string) (*JwtCustomRefreshClaims, error) {
	token, err := jwt.ParseWithClaims(
		strings.TrimSpace(refreshToken),
		&JwtCustomRefreshClaims{},
		func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected jwt signing method: %s", token.Method.Alg())
			}
			return []byte(c.refreshSecret), nil
		},
	)
	if err != nil {
		return nil, err
	}
	return refreshClaimsFromToken(token)
}

func refreshClaimsFromToken(token *jwt.Token) (*JwtCustomRefreshClaims, error) {
	if token == nil || !token.Valid {
		return nil, errors.New("invalid refreshToken")
	}
	claims, ok := token.Claims.(*JwtCustomRefreshClaims)
	if !ok {
		return nil, errors.New("invalid refreshToken claims")
	}
	return claims, nil
}
