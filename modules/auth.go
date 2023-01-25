package modules

import (
	"blockchain/auth"
	"blockchain/db"
	"blockchain/user"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

type TokenResponse struct {
	Token string `json:"token"`
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
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		})
		return
	}
	mongo := db.GetDB(ctx)
	b, err := mongo.CheckExistUser(body.Name)
	if b {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:  "User with this username exist",
			Status: http.StatusBadRequest,
		})
		return
	}
	body.CreatedAt = time.Now().Unix()
	body.GenerateUserAddress()
	u, err := mongo.CreateNewUser(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintln(w, ErrorResponse{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		})
		return
	}
	token, err := auth.GenerateJWT(u)
	response := TokenResponse{Token: token}
	json.NewEncoder(w).Encode(response)
}

func SignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var body user.User
	ctx := r.Context()
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		})
		return
	}
	mongo := db.GetDB(ctx)
	u, err := mongo.CheckUserCredentials(body.Name, body.Password)
	if u.ID == primitive.NilObjectID {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:  "User not exist",
			Status: http.StatusBadRequest,
		})
		return
	}
	token, err := auth.GenerateJWT(u)
	response := TokenResponse{Token: token}
	json.NewEncoder(w).Encode(response)

}
