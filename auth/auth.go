package auth

import (
	"errors"
	"time"

	"blockchain/user"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var sampleSecretKey = "this is my secret key"

type Claims struct {
	Username string             `json:"username"`
	Exp      time.Time          `json:"exp"`
	Id       primitive.ObjectID `json:"id"`
	Address  []byte             `json:"address"`
	jwt.StandardClaims
}

func GenerateJWT(u *user.User) (string, error) {
	claims := Claims{
		Username: u.Name,
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

func GetUserAddress(c *gin.Context) ([]byte, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(c.Request.Header["Authorization"][0], claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(sampleSecretKey), nil
	})
	if err != nil {
		return []byte{}, err
	}
	if claims.Exp.Unix() < time.Now().Unix() {
		return []byte{}, errors.New("token is expired")
	}
	if !token.Valid {
		return []byte{}, errors.New("token is not valid")
	}
	return claims.Address, nil
}
