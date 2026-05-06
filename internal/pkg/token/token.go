package token

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/tpl-x/echo/internal/config"
	"github.com/tpl-x/echo/internal/pkg/usersource"
	"go.uber.org/fx"
)

type CreditGenerator struct {
	accessSecret       string
	accessTokenExpSec  int64
	refreshSecret      string
	refreshTokenExpSec int64
}

type AccessTokenProvider interface {
	CreateAccessToken(user usersource.UserProvider) (accessToken string, exp int64, err error)
	ParseAccessToken(accessToken any) (claims *JwtCustomClaims, err error)
	AccessTokenMiddleware() echo.MiddlewareFunc
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

func NewCreditGeneratorFromConfig(tokenConfig *config.JWTConfig) (*CreditGenerator, error) {
	if tokenConfig.AccessSecret == "" {
		return nil, fmt.Errorf("jwt access_secret cannot be empty")
	}
	if tokenConfig.RefreshSecret == "" {
		return nil, fmt.Errorf("jwt refresh_secret cannot be empty")
	}
	if tokenConfig.AccessTokenExpSec <= 0 {
		return nil, fmt.Errorf("jwt access_token_exp_sec must be greater than zero")
	}
	if tokenConfig.RefreshTokenExpSec <= 0 {
		return nil, fmt.Errorf("jwt refresh_token_exp_sec must be greater than zero")
	}
	return NewCreditGenerator(
		tokenConfig.AccessSecret,
		tokenConfig.AccessTokenExpSec,
		tokenConfig.RefreshSecret,
		tokenConfig.RefreshTokenExpSec,
	), nil
}

func (c *CreditGenerator) AccessTokenMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			authHeader := ctx.Request().Header.Get(echo.HeaderAuthorization)
			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or malformed jwt")
			}

			claims, err := c.ParseAccessToken(strings.TrimPrefix(authHeader, bearerPrefix))
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt").Wrap(err)
			}
			ctx.Set("user", &jwt.Token{
				Raw:    strings.TrimPrefix(authHeader, bearerPrefix),
				Method: jwt.SigningMethodHS256,
				Claims: claims,
				Valid:  true,
			})
			return next(ctx)
		}
	}
}

var Module = fx.Module("token",
	fx.Provide(
		func(appConfig *config.AppConfig) *config.JWTConfig {
			return &appConfig.JWT
		},
		NewCreditGeneratorFromConfig,
		func(generator *CreditGenerator) AccessTokenProvider {
			return generator
		},
		func(generator *CreditGenerator) RefreshTokenProvider {
			return generator
		},
	),
)
