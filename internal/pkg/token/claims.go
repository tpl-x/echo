package token

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Name string `json:"name,omitempty"`
	ID   any    `json:"id,omitempty"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	ID any `json:"id,omitempty"`
	jwt.RegisteredClaims
}
