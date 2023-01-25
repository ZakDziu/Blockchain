package block

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Block struct {
	Timestamp             int64
	Data                  []Transaction
	PrevBlockHash         []byte
	Hash                  []byte
	NumberOfConfirmations int
	BlockNumber           int
	Nonce                 int
}

type Transaction struct {
	AddressSender    uint32  `json:"-"`
	AddressRecipient uint32  `json:"addressRecipient"`
	Sum              float64 `json:"sum"`
	Gas              float64 `json:"-"`
	CreatedAt        int64   `json:"createdAt"`
}

func NewBlock(ctx context.Context, db *mongo.Collection, prevBlockHash []byte, prevBlockNumber int) *Block {
	block := &Block{
		time.Now().Unix(),
		[]Transaction{},
		prevBlockHash,
		[]byte{},
		0,
		prevBlockNumber + 1,
		0}
	pow := block.NewProofOfWork()
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	_, err := db.InsertOne(ctx, block)
	if err != nil {
		log.Fatal(err)
	}

	return block
}

func (t *Transaction) CalculateGas() {
	t.Gas = t.Sum / 100 * 12.5
}
