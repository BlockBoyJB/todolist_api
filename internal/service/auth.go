package service

import (
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	authServicePrefixLog = "/service/auth"
)

var defaultSignMethod = jwt.SigningMethodHS256

type TokenClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

type authService struct {
	signKey  []byte
	tokenTTL time.Duration
}

func newAuthService(signKey string, tokenTTL time.Duration) *authService {
	return &authService{
		signKey:  []byte(signKey),
		tokenTTL: tokenTTL,
	}
}

func (s *authService) CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(defaultSignMethod, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Username: username,
	})

	signedToken, err := token.SignedString(s.signKey)
	if err != nil {
		log.Errorf("%s/CreateToken error sign token: %s", authServicePrefixLog, err)
		return "", err
	}
	return signedToken, nil
}

func (s *authService) ParseToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrIncorrectSignMethod
		}
		return s.signKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, ErrCannotParseToken
	}
	return claims, nil
}
