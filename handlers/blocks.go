package handlers

import (
	"blockchain/block"
	"blockchain/dto"
	"blockchain/internal/types"
	"blockchain/model"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BlocksHandlerInterface interface {
	types.Controller
	NewTransaction(c *gin.Context)
	Blocks(c *gin.Context)
}

type BlocksHandler struct {
	api *Api
}

func NewBlocksHandler(a *Api) BlocksHandlerInterface {
	return &BlocksHandler{
		api: a,
	}
}

func (ctr *BlocksHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/transaction", ctr.NewTransaction)
	r.POST("/blocks", ctr.Blocks)

}

func (ctr *BlocksHandler) NewTransaction(c *gin.Context) {
	var body model.TransactionRequest

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	transaction, err := block.CreateNewTransaction(c, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	err = ctr.api.MongoDB.UpdatesWithCreateNewTransaction(transaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	response := model.NewTransactionResponse{Message: fmt.Sprintf("%v send to the %v", body.Sum, body.AddressRecipient)}

	c.JSON(http.StatusOK, response)
}

func (ctr *BlocksHandler) Blocks(c *gin.Context) {
	var body model.BlockRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	blocks := ctr.api.MongoDB.GetAllBlocks(body)
	transactions := dto.PrepareDataForTransactions(blocks, body.AddressSender, body.AddressRecipient)
	c.JSON(http.StatusOK, transactions)
}
