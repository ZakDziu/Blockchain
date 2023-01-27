package service

import (
	"blockchain/db"
	"blockchain/user"
	"log"
	"time"
)

func CreateUsers(mongo *db.Mongo) {
	sender := user.User{
		Name:      "sender",
		Password:  "myPass",
		Balance:   1000000,
		CreatedAt: time.Now().Unix(),
	}
	sender.AddAddress()
	recipient := user.User{
		Name:      "recipient",
		Password:  "myPass",
		Balance:   0,
		CreatedAt: time.Now().Unix(),
	}
	recipient.AddAddress()
	admin := user.User{
		Name:      "admin",
		Password:  "myPass",
		Balance:   0,
		CreatedAt: time.Now().Unix(),
	}
	admin.AddAddress()

	if !mongo.GetUserByName(sender.Name) {
		_, err := mongo.CreateNewUser(sender)
		if err != nil {
			log.Panic(err)
		}
	}
	if !mongo.GetUserByName(recipient.Name) {
		_, err := mongo.CreateNewUser(recipient)
		if err != nil {
			log.Panic(err)
		}
	}
	if !mongo.GetUserByName(admin.Name) {
		_, err := mongo.CreateNewUser(admin)
		if err != nil {
			log.Panic(err)
		}
	}

}
