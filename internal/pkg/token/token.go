package token

import "github.com/tpl-x/echo/internal/pkg/usersource"

type CreditGenerator struct {
	accessSecret       string
	accessTokenExpSec  int64
	refreshSecret      string
	refreshTokenExpSec int64
}

type AccessTokenProvider interface {
	CreateAccessToken(user usersource.UserProvider) (accessToken string, exp int64, err error)
	ParseAccessToken(accessToken any) (claims *JwtCustomClaims, err error)
}

type RefreshTokenProvider interface {
	CreateRefreshToken(user usersource.UserProvider) (freshToken string, exp int64, err error)
	ParseRefreshToken(refreshToken any) (claims *JwtCustomRefreshClaims, err error)
}

func NewCreditGenerator(
	accessSecret string,
	accessTokenExpSec int64,
	refreshSecret string,
	refreshTokenExpSec int64) *CreditGenerator {
	return &CreditGenerator{
		accessSecret:       accessSecret,
		accessTokenExpSec:  accessTokenExpSec,
		refreshSecret:      refreshSecret,
		refreshTokenExpSec: refreshTokenExpSec,
	}
}
