package main

import (
	"fmt"

	"blockchain/block"
	"blockchain/handlers"
	"blockchain/routes"
	"blockchain/service"

	"github.com/gin-gonic/gin"
)

func main() {
	api := handlers.NewApi()
	defer api.MongoDB.DeleteDataAndCloseConnection(api.Ctx)

	bc := block.NewBlockchain(api.Ctx, api.MongoDB.DB.Block)
	service.CreateUsers(api.MongoDB)
	go service.StartAddBlockService(bc)
	go service.StartMakeTransactions(api.MongoDB)

	fmt.Println("Server listening!")
	authHandler := handlers.NewAuthHandler(api)
	blockHandler := handlers.NewBlocksHandler(api)

	r := routes.NewRouter(gin.Default(),
		authHandler,
		blockHandler,
	)

	r.SetupRouter()

	r.Server.Run(":8000")
}
