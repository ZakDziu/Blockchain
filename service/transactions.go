package service

import (
	"blockchain/block"
	"blockchain/db"
	"blockchain/user"
	"fmt"
	"log"
	"time"
)

func StartMakeTransactions(mongo *db.Mongo) {
	time.Sleep(5 * time.Second)
	transaction := block.Transaction{
		AddressSender:    user.SenderAddress,
		AddressRecipient: user.RecipientAddress,
		Sum:              1,
		CreatedAt:        time.Now().Unix(),
	}
	transaction.CalculateGas()
	transaction.AddTransactionHash()
	fmt.Println(fmt.Sprintf("%x", transaction.ID))
	err := mongo.UpdatesWithCreateNewTransaction(transaction)
	if err != nil {
		log.Panic(err)
	}
	StartMakeTransactions(mongo)
}
