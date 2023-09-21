package handlers

import (
	"context"
	"fmt"

	"blockchain/config"
	"blockchain/db"

	"github.com/joho/godotenv"
)

type Api struct {
	Ctx     context.Context
	MongoDB *db.Mongo
}

func NewApi() *Api {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file", err)
	}

	ctx := context.Background()
	configDB := config.BuildDBConfig()

	return &Api{
		Ctx:     ctx,
		MongoDB: db.GetDB(ctx, configDB),
	}
}
