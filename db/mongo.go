package db

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"blockchain/block"
	"blockchain/config"
	"blockchain/model"
	"blockchain/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const mongoCollectionBlocks = "blocks"
const mongoCollectionUsers = "users"

type Mongo struct {
	ctx    context.Context
	DB     *Collections
	Client *mongo.Client
}

type Collections struct {
	User  *mongo.Collection
	Block *mongo.Collection
}

func GetDB(ctx context.Context, config *config.DBConfig) *Mongo {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	collBlocks := client.Database(config.MongoDB).Collection(mongoCollectionBlocks)
	collUsers := client.Database(config.MongoDB).Collection(mongoCollectionUsers)
	return &Mongo{ctx: ctx, Client: client, DB: &Collections{
		User:  collUsers,
		Block: collBlocks,
	}}
}

func (m *Mongo) DeleteDataAndCloseConnection(ctx context.Context) {
	filter := bson.M{}
	_, err := m.DB.Block.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	_, err = m.DB.User.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Client.Disconnect(ctx)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func (m *Mongo) GetLastBlock() (*block.Block, error) {
	opts := options.FindOne().SetSort(bson.M{"$natural": -1})
	var lastBlock *block.Block
	err := m.DB.Block.FindOne(m.ctx, bson.M{}, opts).Decode(&lastBlock)
	return lastBlock, err
}

func (m *Mongo) GetAllBlocks(req model.BlockRequest) []*block.Block {
	var blocks []*block.Block
	filter := bson.M{}
	if req.AddressSender != "" {
		// 1,0,0,0
		addressS, _ := hex.DecodeString(req.AddressSender)

		filter["data"] = bson.M{"$elemMatch": bson.M{"addresssender": addressS}}
		if req.AddressRecipient != "" {
			//1,1,0,0
			addressR, _ := hex.DecodeString(req.AddressRecipient)

			filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addresssender": addressS}, {"addressrecipient": addressR}}}}
			if req.TransactionCreatedAt != 0 {
				//1,1,1,0
				filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addresssender": addressS}, {"addressrecipient": addressR}, {"createdat": req.TransactionCreatedAt}}}}
				if req.HashID != "" {
					//1,1,1,1
					h, _ := hex.DecodeString(req.HashID)

					filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addresssender": addressS}, {"addressrecipient": addressR}, {"createdat": req.TransactionCreatedAt}, {"id": h}}}}
				}
			}
			if req.HashID != "" {
				//1,1,0,1
				h, _ := hex.DecodeString(req.HashID)

				filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addresssender": addressS}, {"addressrecipient": addressR}, {"id": h}}}}
			}
		}
		if req.TransactionCreatedAt != 0 {
			//1,0,1,0
			filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addresssender": addressS}, {"createdat": req.TransactionCreatedAt}}}}
			if req.HashID != "" {
				//1,0,1,1
				h, _ := hex.DecodeString(req.HashID)

				filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addresssender": addressS}, {"createdat": req.TransactionCreatedAt}, {"id": h}}}}
			}
		}
		if req.HashID != "" {
			//1,0,0,1
			h, _ := hex.DecodeString(req.HashID)

			filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addresssender": addressS}, {"id": h}}}}
		}
	} else {
		if req.AddressRecipient != "" {
			//0,1,0,0
			addressR, _ := hex.DecodeString(req.AddressRecipient)

			filter["data"] = bson.M{"$elemMatch": bson.M{"addressrecipient": addressR}}
			if req.TransactionCreatedAt != 0 {
				//0,1,1,0
				filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addressrecipient": addressR}, {"createdat": req.TransactionCreatedAt}}}}
				if req.HashID != "" {
					//0,1,1,1
					h, _ := hex.DecodeString(req.HashID)

					filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addressrecipient": addressR}, {"createdat": req.TransactionCreatedAt}, {"id": h}}}}
				}
			}
			if req.HashID != "" {
				//0,1,0,1
				h, _ := hex.DecodeString(req.HashID)

				filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addressrecipient": addressR}, {"id": h}}}}
			}
		} else {
			if req.TransactionCreatedAt != 0 {
				//0,0,1,0
				filter["data"] = bson.M{"$elemMatch": bson.M{"createdat": req.TransactionCreatedAt}}
				if req.HashID != "" {
					//0,0,1,1
					h, _ := hex.DecodeString(req.HashID)

					filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"createdat": req.TransactionCreatedAt}, {"id": h}}}}
				}
			}
		}
	}
	if req.HashID != "" {
		//0,0,0,1
		h, _ := hex.DecodeString(req.HashID)
		filter["data"] = bson.M{"$elemMatch": bson.M{"id": h}}
	}
	if req.BlockNumber != 0 {
		filter["blocknumber"] = req.BlockNumber
	}
	opt := options.Find()
	if req.Page != 0 && req.PageSize != 0 {
		if req.Page == 1 {
			opt.SetSkip(0)
			opt.SetLimit(int64(req.PageSize))
		}
		opt.SetSkip(int64((req.Page - 1) * req.PageSize))
		opt.SetLimit(int64(req.PageSize))
	} else {
		opt.SetSkip(0)
		opt.SetLimit(1000)
	}
	cur, err := m.DB.Block.Find(m.ctx, filter, opt)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var elem block.Block
		err = cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		blocks = append(blocks, &elem)
	}
	return blocks
}

func (m *Mongo) CreateNewUser(newUser *user.User) (*user.User, error) {
	newUser.ID = primitive.NewObjectID()
	err := newUser.HashPassword()
	if err != nil {
		return nil, err
	}

	_, err = m.DB.User.InsertOne(m.ctx, newUser)
	return newUser, err
}

func (m *Mongo) GetUserByName(name string) bool {
	var u user.User
	_ = m.DB.User.FindOne(m.ctx, bson.M{"name": name}).Decode(&u)
	if u.ID != primitive.NilObjectID {
		return true
	}
	return false
}

func (m *Mongo) CheckExistUser(name string) (bool, error) {
	var registeredUser user.User
	err := m.DB.User.FindOne(m.ctx, bson.M{"name": name}).Decode(&registeredUser)
	if registeredUser.ID != primitive.NilObjectID {
		return true, err
	}
	return false, err
}

func (m *Mongo) CheckUserCredentials(req user.User) (*user.User, error) {
	var registeredUser *user.User
	err := req.HashPassword()
	if err != nil {
		return nil, err
	}
	err = m.DB.User.FindOne(m.ctx, bson.M{"name": req.Name, "password": req.Password}).Decode(&registeredUser)
	return registeredUser, err
}

func (m *Mongo) UpdatesWithCreateNewTransaction(t *block.Transaction) error {
	var sender user.User
	var recipient user.User
	var admin user.User
	lastBlock, err := m.GetLastBlock()
	if err != nil {
		return err
	}
	lastBlock.Data = append(lastBlock.Data, t)
	filterBlock := bson.D{{"blocknumber", lastBlock.BlockNumber}}

	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := m.Client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(m.ctx)

	filterSender := bson.D{{"address", t.AddressSender}}
	err = m.DB.User.FindOne(m.ctx, filterSender).Decode(&sender)

	filterRecipient := bson.D{{"address", t.AddressRecipient}}
	err = m.DB.User.FindOne(m.ctx, filterRecipient).Decode(&recipient)

	filterAdmin := bson.D{{"address", user.AdminAddress}}
	err = m.DB.User.FindOne(m.ctx, filterAdmin).Decode(&admin)

	updateSender := bson.M{"$set": bson.M{"balance": sender.Balance - t.Sum - t.Gas}}
	updateRecipient := bson.M{"$set": bson.M{"balance": recipient.Balance + t.Sum - t.Gas}}
	updateAdmin := bson.M{"$set": bson.M{"balance": admin.Balance + t.Gas}}
	updateBlock := bson.M{"$set": bson.M{"data": lastBlock.Data}}

	_, err = session.WithTransaction(m.ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		_, err = m.DB.Block.UpdateOne(m.ctx, filterBlock, updateBlock)
		if err != nil {
			err = session.AbortTransaction(ctx)
			if err != nil {
				log.Panic(err)
			}
		}
		_, err = m.DB.User.UpdateOne(m.ctx, filterRecipient, updateRecipient)
		if err != nil {
			err = session.AbortTransaction(ctx)
			if err != nil {
				log.Panic(err)
			}
		}
		_, err = m.DB.User.UpdateOne(m.ctx, filterAdmin, updateAdmin)
		if err != nil {
			err = session.AbortTransaction(ctx)
			if err != nil {
				log.Panic(err)
			}
		}
		result, err := m.DB.User.UpdateOne(m.ctx, filterSender, updateSender)
		if err != nil {
			err = session.AbortTransaction(ctx)
			if err != nil {
				log.Panic(err)
			}
		}

		return result, err
	}, txnOptions)

	return err
}
