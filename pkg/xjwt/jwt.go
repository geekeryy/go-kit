package xjwt

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

const DefaultExpireDuration = time.Hour * 24 * 30

var (
	ErrTokenInvalid = errors.New("couldn't handle this token")
	signKey         = []byte("PCvWzZnAwIvvtjNI")
)

type Business struct {
	UUID   string `json:"uuid"`
	Role   uint64 `json:"role"`
	Extend string `json:"extend"`
}

type CustomClaims struct {
	Business string
	jwt.StandardClaims
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
	var expiresAt time.Time
	if expires > time.Second {
		expiresAt = time.Now().Add(expires)
	} else {
		expiresAt = time.Now().Add(DefaultExpireDuration)
	}
	claims := &CustomClaims{
		Business: bus,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: &jwt.Time{Time: expiresAt},
			IssuedAt:  &jwt.Time{Time: time.Now()},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// ParseToken 解析Token
func ParseToken(tokenString string, bus interface{}) error {
	customClaims := CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &customClaims, func(token *jwt.Token) (interface{}, error) {
		return signKey, nil
	})
	if err != nil {
		return err
	}
	if token == nil || !token.Valid {
		return ErrTokenInvalid
	}
	if err := token.Claims.Valid(nil); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(customClaims.Business), bus); err != nil {
		return err
	}
	return nil

}
