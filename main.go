package main

import (
	"blockchain/block"
	"blockchain/db"
	"blockchain/handlers"
	"blockchain/service"
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	ctx := context.TODO()
	mongo := db.GetDB(ctx)
	defer mongo.DeleteDataAndCloseConnection(ctx)

	bc := block.NewBlockchain(ctx, mongo.DB.Block)
	service.CreateUsers(mongo)
	go service.StartAddBlockService(bc)
	go service.StartMakeTransactions(mongo)

	fmt.Println("Server listening!")
	log.Fatal(http.ListenAndServe(":8989", router()))
}

func router() *httprouter.Router {
	r := httprouter.New()
	r.POST("/sign-up", handlers.SignUp)
	r.POST("/sign-in", handlers.SignIn)
	r.POST("/transaction", handlers.NewTransaction)
	r.GET("/blocks", handlers.Blocks)

	return r
}
