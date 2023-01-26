package service

import (
	"blockchain/db"
	"blockchain/user"
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
		_, _ = mongo.CreateNewUser(sender)
	}
	if !mongo.GetUserByName(recipient.Name) {
		_, _ = mongo.CreateNewUser(recipient)
	}
	if !mongo.GetUserByName(admin.Name) {
		_, _ = mongo.CreateNewUser(admin)
	}

}
