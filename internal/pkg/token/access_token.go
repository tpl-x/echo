package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tpl-x/echo/internal/pkg/usersource"
	"time"
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
	if token, success := accessToken.(*jwt.Token); success {
		if claims, ok := token.Claims.(*JwtCustomClaims); ok {
			return claims, nil
		}
	}
	return nil, errors.New("invalid accessToken")
}
