package block

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Blockchain struct {
	ctx    context.Context
	Blocks []*Block
	db     *mongo.Collection
}

func (bc *Blockchain) AddBlock() {
	opts := options.FindOne().SetSort(bson.M{"$natural": -1})
	var lastBlock Block
	err := bc.db.FindOne(bc.ctx, bson.M{}, opts).Decode(&lastBlock)
	if err != nil {
		log.Fatal(err)
	}
	NewBlock(bc.ctx, bc.db, lastBlock.Hash, lastBlock.BlockNumber)
}

func NewGenesisBlock(ctx context.Context, db *mongo.Collection) *Block {
	opts := options.FindOne().SetSort(bson.M{"$natural": -1})
	var lastBlock Block
	err := db.FindOne(ctx, bson.M{}, opts).Decode(&lastBlock)
	if err != nil {
		return NewBlock(ctx, db, []byte{}, -1)
	}
	if lastBlock.Timestamp != 0 {
		return NewBlock(ctx, db, lastBlock.Hash, lastBlock.BlockNumber)
	}
	return NewBlock(ctx, db, []byte{}, -1)

}

func NewBlockchain(ctx context.Context, db *mongo.Collection) *Blockchain {
	return &Blockchain{ctx, []*Block{NewGenesisBlock(ctx, db)}, db}
}
