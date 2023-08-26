package ginmw

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"tiktok/biz/model"
	"tiktok/pkg/constant"
	"tiktok/pkg/errno"
	"tiktok/pkg/utils"
	"time"
)

var JWT *Jwt

// WithJwtAuth jwt中间件
func WithJwtAuth(jwt *Jwt) gin.HandlerFunc {
	JWT = jwt
	return func(c *gin.Context) {
		token := c.Query(JWT.TokenKey)
		// 视频投稿接口 token 参数存放在body中
		if len(token) == 0 {
			token = c.Request.PostFormValue(JWT.TokenKey)
		}
		if len(token) == 0 {
			c.JSON(http.StatusUnauthorized, model.BuildBaseResp(errno.AuthorizationFailed))
			c.Abort()
			return
		}
		claims, err := JWT.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.BuildBaseResp(errno.AuthorizationFailed.AppendMsg(err.Error())))
			c.Abort()
			return
		}
		c.Set(constant.UserId, claims.Public.UserId)
		c.Next()
	}
}

type Jwt struct {
	Alg       jwt.SigningMethod
	SecretKey []byte
	TokenKey  string
	Issuer    string
	Audience  string
	ExpireDay time.Duration
}
type PublicClaims struct {
	UserId int64
}
type MClaims struct {
	jwt.RegisteredClaims
	Public *PublicClaims
}

func (c MClaims) Validate() error {
	if c.Public.UserId == 0 {
		return errors.New("invalid userId in token")
	}
	return nil
}

// CreateToken 生成token
func (t Jwt) CreateToken(public *PublicClaims) (string, error) {
	token := jwt.NewWithClaims(t.Alg, &MClaims{
		Public: public,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        utils.UUID4(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(t.ExpireDay) * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    t.Issuer,
			Audience:  append(jwt.ClaimStrings{}, t.Audience),
			Subject:   strconv.FormatInt(public.UserId, 10),
		},
	})
	return token.SignedString(t.SecretKey)
}

// ParseToken 解析token
// 需要强secret，签名算法校验, 过期校验, 签名校验，过期校验，private claims校验
// jwt.ParseWithClaims 默认
//
//	在p.ParseUnverified要求有alg
//	token.Method.Verify
//	p.skipClaimsValidation=false
//		但是默认只有时间检验， 其他要配置, 自定义校验需要实现ClaimsValidator即实现Claims和Validate() error
func (t Jwt) ParseToken(tokenString string) (*MClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&MClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return t.SecretKey, nil
		},
		jwt.WithValidMethods([]string{t.Alg.Alg()}),
		jwt.WithIssuer(t.Issuer),
		jwt.WithAudience(t.Audience),
	)
	if err != nil || !token.Valid {
		return nil, errors.Errorf("invalid Token, err:%s", err)
	}
	if claims, ok := token.Claims.(*MClaims); ok {
		return claims, nil
	} else {
		return nil, errors.Errorf("invalid token, type is not *MClaims")
	}
}
