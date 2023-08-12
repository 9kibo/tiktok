package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"tiktok/config"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
type JWT struct {
	JwtKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(config.JwtKey),
	}
}

type MyClaims struct {
	Userid             uint `json:"userid"`
	jwt.StandardClaims      //jwt默认字段
	// 可不填
	//Issuer    string          签发者
	//Subject   string			签发对象
	//Audience  ClaimStrings	签发受众
	//ExpiresAt *NumericDate	过期时间
	//NotBefore *NumericDate	最早使用时间
	//IssuedAt  *NumericDate	签发时间
}

// CreateToken 生成token
func (j *JWT) CreateToken(claims MyClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JwtKey)
}

// ParserToken 解析token
func (j *JWT) ParserToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.JwtKey, nil
	})
	if err != nil {
		return nil, errors.Errorf("无效的Token,err:%s", err)
	}
	if token != nil {
		if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, errors.Errorf("无效的Token")
	}

	return nil, errors.Errorf("无效的Token,err")
}

// JwtToken jwt中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		// 视频投稿接口 token 参数存放在body中
		if len(token) == 0 {
			token = c.Request.PostFormValue("token")
		}
		if len(token) == 0 {
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  "need token",
			})
			c.Abort()
			return
		}
		j := NewJWT()
		claims, err := j.ParserToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				StatusMsg:  err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("userId", claims.Userid)
		c.Next()
	}
}

// JwtWithOutLogin 获取视频流接口需要特殊处理,当未携带token时将userId设为0，方便处理
func JwtWithOutLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		var userId uint
		if len(token) == 0 {
			userId = 0
		} else {
			j := NewJWT()
			claims, err := j.ParserToken(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, Response{
					StatusCode: -1,
					StatusMsg:  err.Error(),
				})
				c.Abort()
				return
			}
			userId = claims.Userid

		}
		c.Set("userId", userId)
		c.Next()
	}
}
