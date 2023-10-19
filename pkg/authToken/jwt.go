package authToken

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserTokenInfo struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

type JWTUserInfoClaims struct {
	jwt.RegisteredClaims
	User *UserTokenInfo `json:"user,omitempty"`
}

func NewToken(secret string, expirationTime int, userInfo *UserTokenInfo) (string, error) {
	claims := JWTUserInfoClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Duration(expirationTime) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Local()),
		},
		userInfo,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signetToken, err := token.SignedString([]byte(secret))
	return signetToken, err
}

func VerifyToken(secret string, token string) (*UserTokenInfo, bool) {
	t, err := jwt.ParseWithClaims(token, &JWTUserInfoClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, false
	}

	expTime, err := t.Claims.GetExpirationTime()
	if err != nil {
		return nil, false
	}

	if !t.Valid || expTime.Before(time.Now().Local()) {
		return nil, false
	}

	if userInfo, ok := t.Claims.(*JWTUserInfoClaims); ok {
		return userInfo.User, true
	}

	return nil, false
}
