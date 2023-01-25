package service

import (
	"blockchain/block"
	"blockchain/db"
	"blockchain/user"
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

func CreateUsers(mongo *db.Mongo) {
	sender := user.User{
		Name:      "sender",
		Password:  "myPass",
		Address:   3675513191,
		Balance:   1000000,
		CreatedAt: time.Now().Unix(),
	}
	recipient := user.User{
		Name:      "recipient",
		Password:  "myPass",
		Address:   441489459,
		Balance:   0,
		CreatedAt: time.Now().Unix(),
	}
	admin := user.User{
		Name:      user.AdminName,
		Password:  "myPass",
		Address:   2497565411,
		Balance:   0,
		CreatedAt: time.Now().Unix(),
	}
	mongo.CreateNewUser(sender)
	mongo.CreateNewUser(recipient)
	mongo.CreateNewUser(admin)

}
