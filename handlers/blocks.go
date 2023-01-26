package handlers

import (
	"blockchain/auth"
	"blockchain/block"
	"blockchain/db"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"time"
)

type NewTransactionResponse struct {
	Message string `json:"message"`
}

func NewTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	authToken := r.Header.Get("Authorization")
	var body block.Transaction
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
	body.AddressSender, err = auth.GetUserAddress(authToken)
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
	body.CreatedAt = time.Now().Unix()
	body.CalculateGas()
	err = mongo.UpdatesWithCreateNewTransaction(body)
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

	response := NewTransactionResponse{Message: fmt.Sprintf("%v send to the %v", body.Sum, body.AddressRecipient)}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Panic(err)
	}
}

func Blocks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	hashId := params.Get("hash")
	addressSender, err := strconv.Atoi(params.Get("addressSender"))
	if err != nil {
		err = json.NewEncoder(w).Encode(ErrorResponse{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		})
		if err != nil {
			log.Panic(err)
		}
	}
	addressRecipient, err := strconv.Atoi(params.Get("addressRecipient"))
	if err != nil {
		err = json.NewEncoder(w).Encode(ErrorResponse{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		})
		if err != nil {
			log.Panic(err)
		}
	}
	blockNumber, err := strconv.Atoi(params.Get("blockNumber"))
	if err != nil {
		err = json.NewEncoder(w).Encode(ErrorResponse{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		})
		if err != nil {
			log.Panic(err)
		}
	}
	transactionCreatedAt, err := strconv.Atoi(params.Get("transactionCreatedAt"))
	if err != nil {
		err = json.NewEncoder(w).Encode(ErrorResponse{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		})
		if err != nil {
			log.Panic(err)
		}
	}
	pageSize, err := strconv.Atoi(params.Get("pageSize"))
	if err != nil {
		err = json.NewEncoder(w).Encode(ErrorResponse{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		})
		if err != nil {
			log.Panic(err)
		}
	}
	page, err := strconv.Atoi(params.Get("page"))
	if err != nil {
		err = json.NewEncoder(w).Encode(ErrorResponse{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		})
		if err != nil {
			log.Panic(err)
		}
	}
	ctx := r.Context()
	mongo := db.GetDB(ctx)
	blocks := mongo.GetAllBlocks(hashId, uint32(addressSender), uint32(addressRecipient), blockNumber, int64(transactionCreatedAt), page, pageSize)
	err = json.NewEncoder(w).Encode(blocks)
	if err != nil {
		log.Panic(err)
	}
}
