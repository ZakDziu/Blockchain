package handlers

import (
	"fmt"
	"net/http"
	"time"

	"blockchain/auth"
	"blockchain/internal/types"
	"blockchain/model"
	"blockchain/user"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandlerInterface interface {
	types.Controller
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
}

type AuthHandler struct {
	api *Api
}

func NewAuthHandler(a *Api) AuthHandlerInterface {
	return &AuthHandler{
		api: a,
	}
}

func (ctr *AuthHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/sign-up", ctr.SignUp)
	r.POST("/sign-up", ctr.SignIn)

}

func (ctr *AuthHandler) SignUp(c *gin.Context) {
	var body user.User

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrUnhealthy)
		return
	}

	b, err := ctr.api.MongoDB.CheckExistUser(body.Name)
	if b {
		c.JSON(http.StatusBadRequest, model.ErrUserWithThisNameExist)

		return
	}

	body.CreatedAt = time.Now().Unix()
	body.AddAddress()

	u, err := ctr.api.MongoDB.CreateNewUser(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewError(http.StatusInternalServerError, err.Error()))

		return
	}

	token, err := auth.GenerateJWT(u)
	response := model.TokenResponse{Token: token, Address: fmt.Sprintf("%x", body.Address)}

	c.JSON(http.StatusOK, response)
}

func (ctr *AuthHandler) SignIn(c *gin.Context) {
	var body user.User
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrUnhealthy)
		return
	}

	u, err := ctr.api.MongoDB.CheckUserCredentials(body.Name, body.Password)
	if u.ID == primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, model.ErrUserNotExist)
		return
	}

	token, err := auth.GenerateJWT(u)
	response := model.TokenResponse{Token: token, Address: fmt.Sprintf("%x", u.Address)}

	c.JSON(http.StatusOK, response)
}
