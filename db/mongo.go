package db

import (
	"blockchain/block"
	"blockchain/user"
	"context"
	"encoding/hex"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
)

type Mongo struct {
	ctx    context.Context
	DB     *Collections
	Client *mongo.Client
}

type Collections struct {
	User  *mongo.Collection
	Block *mongo.Collection
}

func GetDB(ctx context.Context) *Mongo {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
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
	collBlocks := client.Database("db").Collection("blocks")
	collUsers := client.Database("db").Collection("users")
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

func (m *Mongo) GetLastBlock() (block.Block, error) {
	opts := options.FindOne().SetSort(bson.M{"$natural": -1})
	var lastBlock block.Block
	err := m.DB.Block.FindOne(m.ctx, bson.M{}, opts).Decode(&lastBlock)
	return lastBlock, err
}

func (m *Mongo) GetAllBlocks(hashId string, addressSender, addressRecipient string, blockNumer int, transactionCreatedAt int64, page, pageSize int) []*block.Block {
	var blocks []*block.Block
	filter := bson.M{}
	if addressSender != "" {
		// 1,0,0,0
		addressS, _ := hex.DecodeString(addressSender)

		filter["data"] = bson.M{"$elemMatch": bson.M{"addresssender": addressS}}
		if addressRecipient != "" {
			//1,1,0,0
			addressR, _ := hex.DecodeString(addressRecipient)

			filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addresssender": addressS}, {"addressrecipient": addressR}}}}
			if transactionCreatedAt != 0 {
				//1,1,1,0
				filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addresssender": addressS}, {"addressrecipient": addressR}, {"createdat": transactionCreatedAt}}}}
				if hashId != "" {
					//1,1,1,1
					h, _ := hex.DecodeString(hashId)

					filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addresssender": addressS}, {"addressrecipient": addressR}, {"createdat": transactionCreatedAt}, {"id": h}}}}
				}
			}
			if hashId != "" {
				//1,1,0,1
				h, _ := hex.DecodeString(hashId)

				filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addresssender": addressS}, {"addressrecipient": addressR}, {"id": h}}}}
			}
		}
		if transactionCreatedAt != 0 {
			//1,0,1,0
			filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addresssender": addressS}, {"createdat": transactionCreatedAt}}}}
			if hashId != "" {
				//1,0,1,1
				h, _ := hex.DecodeString(hashId)

				filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addresssender": addressS}, {"createdat": transactionCreatedAt}, {"id": h}}}}
			}
		}
		if hashId != "" {
			//1,0,0,1
			h, _ := hex.DecodeString(hashId)

			filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addresssender": addressS}, {"id": h}}}}
		}
	} else {
		if addressRecipient != "" {
			//0,1,0,0
			addressR, _ := hex.DecodeString(addressRecipient)

			filter["data"] = bson.M{"$elemMatch": bson.M{"addressrecipient": addressR}}
			if transactionCreatedAt != 0 {
				//0,1,1,0
				filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addressrecipient": addressR}, {"createdat": transactionCreatedAt}}}}
				if hashId != "" {
					//0,1,1,1
					h, _ := hex.DecodeString(hashId)

					filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addressrecipient": addressR}, {"createdat": transactionCreatedAt}, {"id": h}}}}
				}
			}
			if hashId != "" {
				//0,1,0,1
				h, _ := hex.DecodeString(hashId)

				filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"addressrecipient": addressR}, {"id": h}}}}
			}
		} else {
			if transactionCreatedAt != 0 {
				//0,0,1,0
				filter["data"] = bson.M{"$elemMatch": bson.M{"createdat": transactionCreatedAt}}
				if hashId != "" {
					//0,0,1,1
					h, _ := hex.DecodeString(hashId)

					filter["data"] = bson.M{"$elemMatch": bson.M{"$and": []bson.M{{"createdat": transactionCreatedAt}, {"id": h}}}}
				}
			}
		}
	}
	if hashId != "" {
		//0,0,0,1
		h, _ := hex.DecodeString(hashId)
		filter["data"] = bson.M{"$elemMatch": bson.M{"id": h}}
	}
	if blockNumer != 0 {
		filter["blocknumber"] = blockNumer
	}
	opt := options.Find()
	if page != 0 && pageSize != 0 {
		if page == 1 {
			opt.SetSkip(0)
			opt.SetLimit(int64(pageSize))
		}
		opt.SetSkip(int64((page - 1) * pageSize))
		opt.SetLimit(int64(pageSize))
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

func (m *Mongo) CreateNewUser(newUser user.User) (user.User, error) {
	newUser.ID = primitive.NewObjectID()
	_, err := m.DB.User.InsertOne(m.ctx, newUser)
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

func (m *Mongo) CheckUserCredentials(name string, password string) (user.User, error) {
	var registeredUser user.User
	err := m.DB.User.FindOne(m.ctx, bson.M{"name": name, "password": password}).Decode(&registeredUser)
	return registeredUser, err
}

func (m *Mongo) UpdatesWithCreateNewTransaction(t block.Transaction) error {
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
