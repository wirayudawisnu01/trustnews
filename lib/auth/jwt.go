package auth

import (
	"trustnews/config"
	"trustnews/internal/core/domain/entity"
	"time"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type  Jwt interface {
	GenerateToken(data *entity.JwtData)(string, int64, error)
	VerifyAccessToken(token string)(*entity.JwtData, error)
}

type Options struct {
	signingKey string
	issuer string
}

func (o *Options) GenerateToken(data *entity.JwtData) (string, int64, error) {
	now := time.Now().Local()
	expiresAt := now.Add(time.Hour * 24)
	data.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	data.RegisteredClaims.Issuer = o.issuer
	data.RegisteredClaims.NotBefore = jwt.NewNumericDate(now)
	acToken := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	accessToken, err := acToken.SignedString([]byte(o.signingKey))
	if err != nil {
		return "", 0, err
	} 

	return accessToken, expiresAt.Unix(), nil
}

func (o *Options) VerifyAccessToken(token string) (*entity.JwtData, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token)(interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Signing method invalid")
		}

		return []byte(o.signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	if parsedToken.Valid {
		claim, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok || !parsedToken.Valid {
			return nil, err
		}

		jwtData := &entity.JwtData {
			UserID: claim["user_id"].(float64),
		}

		return jwtData, nil
	}

	return nil, fmt.Errorf("Token is not valid")
}

func NewJwt(cfg *config.Config) Jwt {
	opt := new(Options)
	opt.signingKey = cfg.App.JwtSecretKey
	opt.issuer = cfg.App.JwtIssuer

	return opt
}