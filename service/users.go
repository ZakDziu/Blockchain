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
