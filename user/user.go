package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hash/fnv"
)

const AdminName = "zakhar"

type MyObjectID string

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"-"`
	Name      string             `bson:"name" json:"username"`
	Password  string             `json:"password"`
	CreatedAt int64
	Address   uint32
	Balance   float64 `bson:"balance"`
}

func (u *User) GenerateUserAddress() {
	h := fnv.New32a()
	h.Write([]byte(u.Name + u.Password))
	u.Address = h.Sum32()
}
