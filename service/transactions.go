package service

import (
	"blockchain/block"
	"blockchain/db"
	"log"
	"time"
)

func StartMakeTransactions(mongo *db.Mongo) {
	time.Sleep(5 * time.Second)
	transaction := block.Transaction{
		AddressSender:    3675513191,
		AddressRecipient: 441489459,
		Sum:              1,
		CreatedAt:        time.Now().Unix(),
	}
	transaction.CalculateGas()
	err := mongo.UpdatesWithCreateNewTransaction(transaction)
	if err != nil {
		log.Panic(err)
	}
	StartMakeTransactions(mongo)
}
