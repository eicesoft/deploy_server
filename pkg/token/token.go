package token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var _ Token = (*token)(nil)

type Token interface {
	d()

	// JwtSign JWT 签名方式
	JwtSign(userId int32, userName string, role string, expireDuration time.Duration) (tokenString string, err error)
	// JwtParse JWT解析
	JwtParse(tokenString string) (*claims, error)
}

type token struct {
	secret string
}

type claims struct {
	UserID   int32
	UserName string
	Role     string
	jwt.StandardClaims
}

func New(secret string) Token {
	return &token{
		secret: secret,
	}
}

func (t *token) d() {}
