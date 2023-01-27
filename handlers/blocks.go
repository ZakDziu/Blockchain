package handlers

import (
	"blockchain/auth"
	"blockchain/block"
	"blockchain/db"
	"blockchain/dto"
	"encoding/hex"
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

type TransactionRequest struct {
	AddressRecipient string  `json:"addressRecipient"`
	Sum              float64 `json:"sum"`
}

func NewTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	authToken := r.Header.Get("Authorization")
	var body TransactionRequest
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
	addressRecipient, _ := hex.DecodeString(body.AddressRecipient)
	t := block.Transaction{
		AddressSender:    nil,
		AddressRecipient: addressRecipient,
		Sum:              body.Sum,
		Gas:              0,
		CreatedAt:        0,
	}
	mongo := db.GetDB(ctx)
	t.AddressSender, err = auth.GetUserAddress(authToken)
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
	t.CreatedAt = time.Now().Unix()
	t.CalculateGas()
	t.AddTransactionHash()
	err = mongo.UpdatesWithCreateNewTransaction(t)
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
	var blockNumber int
	var transactionCreatedAt int
	var pageSize int
	var page int

	params := r.URL.Query()
	hashId := params.Get("hash")
	_, err := hex.DecodeString(hashId)
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
	addressSender := params.Get("addressSender")
	_, err = hex.DecodeString(addressSender)
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
	addressRecipient := params.Get("addressRecipient")
	_, err = hex.DecodeString(hashId)
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
	bN := params.Get("blockNumber")
	if bN != "" {
		blockNumber, err = strconv.Atoi(bN)
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
	}
	tCA := params.Get("transactionCreatedAt")
	if tCA != "" {
		transactionCreatedAt, err = strconv.Atoi(tCA)
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
	}
	pS := params.Get("pageSize")
	if pS != "" {
		pageSize, err = strconv.Atoi(pS)
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
	}
	p := params.Get("page")
	if p != "" {
		page, err = strconv.Atoi(p)
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
	}
	ctx := r.Context()
	mongo := db.GetDB(ctx)
	blocks := mongo.GetAllBlocks(hashId, addressSender, addressRecipient, blockNumber, int64(transactionCreatedAt), page, pageSize)
	transactions := dto.PrepareDataForTransactions(blocks, addressSender, addressRecipient)
	err = json.NewEncoder(w).Encode(transactions)
	if err != nil {
		log.Panic(err)
	}
}
