package handlers

import (
	"blockchain/auth"
	"blockchain/db"
	"blockchain/user"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

type TokenResponse struct {
	Token   string `json:"token"`
	Address string `json:"address"`
}

type SignRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var body user.User
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(ErrorResponse{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		})
		if err != nil {
			log.Panic(err)
		}
		return
	}
	mongo := db.GetDB(ctx)
	b, err := mongo.CheckExistUser(body.Name)
	if b {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(ErrorResponse{
			Error:  "User with this username exist",
			Status: http.StatusBadRequest,
		})
		if err != nil {
			log.Panic(err)
		}
		return
	}
	body.CreatedAt = time.Now().Unix()
	body.AddAddress()
	u, err := mongo.CreateNewUser(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(ErrorResponse{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		})
		if err != nil {
			log.Panic(err)
		}
		return
	}
	token, err := auth.GenerateJWT(u)
	response := TokenResponse{Token: token, Address: fmt.Sprintf("%x", body.Address)}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Panic(err)
	}
	return
}

func SignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var body user.User
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(ErrorResponse{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		})
		if err != nil {
			log.Panic(err)
		}
		return
	}
	mongo := db.GetDB(ctx)
	u, err := mongo.CheckUserCredentials(body.Name, body.Password)
	if u.ID == primitive.NilObjectID {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(ErrorResponse{
			Error:  "User not exist",
			Status: http.StatusBadRequest,
		})
		if err != nil {
			log.Panic(err)
		}
		return
	}
	token, err := auth.GenerateJWT(u)
	response := TokenResponse{Token: token, Address: fmt.Sprintf("%x", u.Address)}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Panic(err)
	}
	return
}
