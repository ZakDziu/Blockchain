package block

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Block struct {
	Timestamp             int64
	Data                  []*Transaction
	PrevBlockHash         []byte
	Hash                  []byte
	NumberOfConfirmations int
	BlockNumber           int
	Nonce                 int
}

func NewBlock(ctx context.Context, db *mongo.Collection, prevBlockHash []byte, prevBlockNumber int) *Block {
	block := &Block{
		time.Now().Unix(),
		[]*Transaction{},
		prevBlockHash,
		[]byte{},
		0,
		prevBlockNumber + 1,
		0}
	pow := block.NewProofOfWork()
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	_, err := db.UpdateMany(ctx, bson.M{}, bson.M{"$inc": bson.M{"numberofconfirmations": 1}})
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.InsertOne(ctx, block)
	if err != nil {
		log.Fatal(err)
	}

	return block
}

func (t *Transaction) CalculateGas() {
	t.Gas = t.Sum / 100 * 12.5
}
