package modules

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"myproject/auth"
	"myproject/block"
	"myproject/db"
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
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		})
		return
	}
	mongo := db.GetDB(ctx)
	body.AddressSender, err = auth.GetUserAddress(authToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		})
		return
	}
	body.CreatedAt = time.Now().Unix()
	body.CalculateGas()
	err = mongo.UpdatesWithCreateNewTransaction(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		})
		return
	}

	response := NewTransactionResponse{Message: fmt.Sprintf("%v send to the %v", body.Sum, body.AddressRecipient)}
	json.NewEncoder(w).Encode(response)
}

func Blocks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	hashId := params.Get("hash")
	addressSender, _ := strconv.Atoi(params.Get("addressSender"))
	addressRecipient, _ := strconv.Atoi(params.Get("addressRecipient"))
	blockNumber, _ := strconv.Atoi(params.Get("blockNumber"))
	transactionCreatedAt, _ := strconv.Atoi(params.Get("transactionCreatedAt"))
	pageSize, _ := strconv.Atoi(params.Get("pageSize"))
	page, _ := strconv.Atoi(params.Get("page"))
	ctx := r.Context()
	mongo := db.GetDB(ctx)
	blocks := mongo.GetAllBlocks(hashId, uint32(addressSender), uint32(addressRecipient), blockNumber, int64(transactionCreatedAt), page, pageSize)

	json.NewEncoder(w).Encode(blocks)
}
