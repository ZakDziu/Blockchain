package service

import (
	"log"
	"time"

	"blockchain/block"
	"blockchain/db"
	"blockchain/user"
)

func StartMakeTransactions(mongo *db.Mongo) {
	time.Sleep(5 * time.Second)
	transaction := &block.Transaction{
		AddressSender:    user.SenderAddress,
		AddressRecipient: user.RecipientAddress,
		Sum:              1,
		CreatedAt:        time.Now().Unix(),
	}
	transaction.CalculateGas()
	transaction.AddTransactionHash()

	err := mongo.UpdatesWithCreateNewTransaction(transaction)
	if err != nil {
		log.Panic(err)
	}
	StartMakeTransactions(mongo)
}
