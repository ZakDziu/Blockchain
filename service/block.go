package service

import (
	"blockchain/block"
	"blockchain/db"
	"time"
)

func StartAddBlockService(bc *block.Blockchain) {
	bc.AddBlock()
	time.Sleep(40 * time.Second)
	StartAddBlockService(bc)
}

func StartMakeTransactions(mongo *db.Mongo) {
	time.Sleep(5 * time.Second)
	transaction := block.Transaction{
		AddressSender:    3675513191,
		AddressRecipient: 441489459,
		Sum:              1,
		CreatedAt:        time.Now().Unix(),
	}
	transaction.CalculateGas()
	mongo.UpdatesWithCreateNewTransaction(transaction)
	StartMakeTransactions(mongo)
}
