package auth

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"myproject/user"
	"time"
)

var sampleSecretKey = "this is my secret key"

type Claims struct {
	Username string             `json:"username"`
	Password string             `json:"password"`
	Exp      time.Time          `json:"exp"`
	Id       primitive.ObjectID `json:"id"`
	Address  uint32             `json:"address"`
	jwt.StandardClaims
}

func GenerateJWT(u user.User) (string, error) {
	claims := Claims{
		Username: u.Name,
		Password: u.Password,
		Exp:      time.Now().Add(10 * time.Minute),
		Id:       u.ID,
		Address:  u.Address,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(sampleSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserAddress(tokenString string) (uint32, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(sampleSecretKey), nil
	})
	if err != nil {
		return 0, err
	}
	if claims.Exp.Unix() < time.Now().Unix() {
		return 0, errors.New("token is expired")
	}
	if !token.Valid {
		return 0, errors.New("token is not valid")
	}
	return claims.Address, nil
}
