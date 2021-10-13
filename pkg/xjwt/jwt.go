package xjwt

import (
	"encoding/json"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go/v4"
	"github.com/pkg/errors"
)

const DefaultExpireDuration = time.Hour * 24 * 30

var (
	ErrTokenInvalid     = errors.New("Couldn't handle this token")
	signKey             = []byte("PCvWzZnAwIvvtjNI")
)

type Business struct {
	UUID string `json:"uuid"`
	Role uint64 `json:"role"`
}

type CustomClaims struct {
	Business string
	jwtgo.StandardClaims
}

type TokenResp struct {
	Token     string `json:"token"`
	ExpiredAt uint64 `json:"expired_at"`
}

func Init(key string) {
	signKey = []byte(key)
}

// CreateToken 创建Token
func CreateToken(bus string, expires time.Duration) (string, error) {
	expiresAt := time.Now().Add(DefaultExpireDuration)
	if expires != 0 {
		expiresAt = time.Now().Add(expires)
	}
	claims := &CustomClaims{
		Business: bus,
		StandardClaims: jwtgo.StandardClaims{
			ExpiresAt: &jwtgo.Time{Time: expiresAt},
		},
	}
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// ParseToken 解析Token
func ParseToken(tokenString string, v interface{}) error {
	customClaims := CustomClaims{}
	token, err := jwtgo.ParseWithClaims(tokenString, &customClaims, func(token *jwtgo.Token) (interface{}, error) {
		return signKey, nil
	})
	if err != nil {
		return err
	}
	if token == nil || !token.Valid {
		return ErrTokenInvalid
	}

	if err := json.Unmarshal([]byte(customClaims.Business), v); err != nil {
		return err
	}
	return nil

}
