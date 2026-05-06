package token

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tpl-x/echo/internal/pkg/usersource"
)

var _ AccessTokenProvider = (*CreditGenerator)(nil)

func (c *CreditGenerator) CreateAccessToken(user usersource.UserProvider) (accessToken string, exp int64, err error) {
	accessTokenExp := time.Now().Add(time.Duration(c.accessTokenExpSec) * time.Second)
	claims := &JwtCustomClaims{
		user.UserName(),
		user.UserId(),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExp),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err = t.SignedString([]byte(c.accessSecret))
	if err != nil {
		return "", 0, err
	}
	return accessToken, accessTokenExp.Unix(), nil
}

func (c *CreditGenerator) ParseAccessToken(accessToken any) (*JwtCustomClaims, error) {
	switch token := accessToken.(type) {
	case string:
		return c.parseAccessTokenString(token)
	case *jwt.Token:
		return accessClaimsFromToken(token)
	default:
		return nil, errors.New("invalid accessToken")
	}
}

func (c *CreditGenerator) parseAccessTokenString(accessToken string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		strings.TrimSpace(accessToken),
		&JwtCustomClaims{},
		func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected jwt signing method: %s", token.Method.Alg())
			}
			return []byte(c.accessSecret), nil
		},
	)
	if err != nil {
		return nil, err
	}
	return accessClaimsFromToken(token)
}

func accessClaimsFromToken(token *jwt.Token) (*JwtCustomClaims, error) {
	if token == nil || !token.Valid {
		return nil, errors.New("invalid accessToken")
	}
	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok {
		return nil, errors.New("invalid accessToken claims")
	}
	return claims, nil
}
