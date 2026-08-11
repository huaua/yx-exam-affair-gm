package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/cst"
	"github.com/sealsee/web-base/public/setting"
	"github.com/sealsee/web-base/public/utils/redis"
	"time"
)

type JWT struct {
	TokenId string
	jwt.StandardClaims
}

var mySecret []byte
var expireTime int64
var issuer string

func Init() {
	mySecret = []byte(setting.Conf.TokenConfig.Secret)
	expireTime = setting.Conf.TokenConfig.ExpireTime
	issuer = setting.Conf.TokenConfig.Issuer
}

// GenToken 生成JWT
func GenToken(tokenId string) string {
	// 创建一个我们自己的声明的数据
	c := JWT{
		TokenId: tokenId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(setting.Conf.JWTMaxLifetimeMin) * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signedString, err := token.SignedString(mySecret)
	if err != nil {
		panic(err)
	}
	return signedString
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (sessionUser *context.SessionUser, err error) {
	var mc = new(JWT)
	_, err = jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected jwt signing method")
		}
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !mc.VerifyIssuer(issuer, true) {
		return nil, errors.New("invalid jwt issuer")
	}
	sessionUser, err = redis.GetStruct(cst.LoginTokenKey+mc.TokenId, sessionUser)

	return
}
