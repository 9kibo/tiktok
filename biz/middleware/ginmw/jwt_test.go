package ginmw

import (
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
)

func initJwt() *Jwt {
	return &Jwt{
		Alg:       jwt.SigningMethodHS256,
		SecretKey: []byte("123456789"),
		TokenKey:  "token",
		Issuer:    "Issuer",
		Audience:  "Audience",
		ExpireDay: 2 * time.Minute,
	}
}
func TestParseToken(t *testing.T) {
	j := initJwt()
	token, err := j.CreateToken(&PublicClaims{
		UserId: 123,
	})
	if err != nil {
		panic(err)
	}
	t.Log("token=", token)
	claims, err := j.ParseToken(token)
	if err != nil {
		panic(err)
	}
	t.Logf("claims=%#v", claims)
}
