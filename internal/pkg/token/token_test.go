package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type testUser struct {
	id   string
	name string
}

func (u testUser) UserId() any {
	return u.id
}

func (u testUser) UserName() string {
	return u.name
}

func TestCreditGeneratorParseAccessTokenString(t *testing.T) {
	generator := NewCreditGenerator("access-secret", 60, "refresh-secret", 120)

	signedToken, exp, err := generator.CreateAccessToken(testUser{id: "user-1", name: "alice"})
	if err != nil {
		t.Fatalf("CreateAccessToken() error = %v", err)
	}
	if exp <= time.Now().Unix() {
		t.Fatalf("CreateAccessToken() exp = %d, want future timestamp", exp)
	}

	claims, err := generator.ParseAccessToken(signedToken)
	if err != nil {
		t.Fatalf("ParseAccessToken() error = %v", err)
	}
	if claims.Name != "alice" || claims.ID != "user-1" {
		t.Fatalf("ParseAccessToken() claims = %#v, want name alice and id user-1", claims)
	}
}

func TestCreditGeneratorRefreshTokenUsesRefreshSecret(t *testing.T) {
	generator := NewCreditGenerator("access-secret", 60, "refresh-secret", 120)

	signedToken, _, err := generator.CreateRefreshToken(testUser{id: "user-1", name: "alice"})
	if err != nil {
		t.Fatalf("CreateRefreshToken() error = %v", err)
	}

	if _, err := generator.ParseRefreshToken(signedToken); err != nil {
		t.Fatalf("ParseRefreshToken() error = %v", err)
	}

	_, err = jwt.ParseWithClaims(signedToken, &JwtCustomRefreshClaims{}, func(token *jwt.Token) (any, error) {
		return []byte("access-secret"), nil
	})
	if err == nil {
		t.Fatal("refresh token validated with access secret")
	}
}

func TestCreditGeneratorRejectsExpiredAccessToken(t *testing.T) {
	generator := NewCreditGenerator("access-secret", -1, "refresh-secret", 120)

	signedToken, _, err := generator.CreateAccessToken(testUser{id: "user-1", name: "alice"})
	if err != nil {
		t.Fatalf("CreateAccessToken() error = %v", err)
	}

	if _, err := generator.ParseAccessToken(signedToken); err == nil {
		t.Fatal("ParseAccessToken() accepted expired token")
	}
}

func TestCreditGeneratorRejectsUnexpectedSigningMethod(t *testing.T) {
	generator := NewCreditGenerator("access-secret", 60, "refresh-secret", 120)
	claims := &JwtCustomClaims{
		Name: "alice",
		ID:   "user-1",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	signedToken, err := token.SignedString([]byte("access-secret"))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}

	if _, err := generator.ParseAccessToken(signedToken); err == nil {
		t.Fatal("ParseAccessToken() accepted unexpected signing method")
	}
}
