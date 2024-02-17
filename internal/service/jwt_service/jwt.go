package jwtservice

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	salt []byte
}

func NewJwtService(salt string) *JwtService {
	return &JwtService{salt: []byte(salt)}
}

func (j *JwtService) GetJwt(id int) (string, error) {
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	})
	token, err := jwt.SignedString(j.salt)

	return token, err
}

func (j *JwtService) ParseJwt(token string) (int, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return j.salt, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("map claims error")
	}

	idInterface := claims["id"]
	value, ok := idInterface.(float64)
	if !ok {
		return 0, errors.New("id is not int")
	}

	return int(value), nil
}
