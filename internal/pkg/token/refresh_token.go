package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tpl-x/echo/internal/pkg/usersource"
	"time"
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
	refreshToken, err := t.SignedString([]byte(c.accessSecret))
	if err != nil {
		return "", 0, err
	}
	return refreshToken, refreshTokenExp.Unix(), nil
}

func (c *CreditGenerator) ParseRefreshToken(refreshToken any) (claims *JwtCustomRefreshClaims, err error) {
	if token, success := refreshToken.(*jwt.Token); success {
		if claims, ok := token.Claims.(*JwtCustomRefreshClaims); ok {
			return claims, nil
		}
	}
	return nil, errors.New("invalid refreshToken")
}
