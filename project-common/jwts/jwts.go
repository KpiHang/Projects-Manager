package jwts

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JwtToken struct {
	AccessToken  string // 访问令牌
	RefreshToken string // 刷新令牌
	AccessExp    int64  // 访问令牌的过期时间（Unix 时间戳）
	RefreshExp   int64  // 刷新令牌的过期时间（Unix 时间戳）
}

func CreateToken(val string, exp time.Duration, secret string, refreshExp time.Duration, refreshSecret string) *JwtToken {
	// 计算访问令牌的过期时间
	aExp := time.Now().Add(exp).Unix()
	// 创建访问令牌
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": val,
		"exp":   aExp,
	})
	// 签名访问令牌
	aToken, _ := accessToken.SignedString([]byte(secret))

	// 计算刷新令牌的过期时间
	rExp := time.Now().Add(refreshExp).Unix()
	// 创建刷新令牌
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": val,
		"exp":   rExp,
	})
	// 签名刷新令牌
	rToken, _ := refreshToken.SignedString([]byte(refreshSecret))

	// 返回包含访问令牌和刷新令牌的 JwtToken 实例
	return &JwtToken{
		AccessToken:  aToken,
		AccessExp:    aExp,
		RefreshToken: rToken,
		RefreshExp:   rExp,
	}
}

func ParseToken(tokenString string, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		val := claims["token"].(string)
		exp := int64(claims["exp"].(float64))
		if time.Now().Unix() > exp {
			return "", errors.New("token expired")
		}
		return val, nil
	} else {
		return "", err
	}
}
