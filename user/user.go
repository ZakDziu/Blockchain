package user

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"

	"blockchain/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var SenderAddress = []byte{0, 1, 225, 194, 49, 70, 157, 230, 199, 233, 13, 221, 133, 106, 143, 114, 254, 26, 21, 17, 13, 103, 50, 154, 65, 118, 178, 144, 155, 48, 240, 231}
var RecipientAddress = []byte{0, 0, 6, 158, 8, 224, 247, 0, 253, 241, 22, 247, 242, 255, 79, 236, 61, 56, 165, 167, 164, 49, 210, 103, 239, 92, 123, 183, 0, 228, 243, 255}
var AdminAddress = []byte{0, 1, 157, 252, 174, 99, 4, 41, 240, 19, 23, 168, 113, 127, 177, 147, 221, 130, 83, 9, 99, 92, 211, 65, 87, 125, 91, 233, 171, 157, 175, 228}

const targetBits = 15

var (
	maxNonce = math.MaxInt64
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"-"`
	Name      string             `bson:"name" json:"username"`
	Password  string             `json:"password"`
	CreatedAt int64
	Address   []byte
	Balance   float64 `bson:"balance"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (u *User) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			[]byte(u.ID.Hex()),
			[]byte(u.Name),
			[]byte(u.Password),
			utils.IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (u *User) AddAddress() {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	for nonce < maxNonce {
		data := u.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}
	}
	h := hash[:]
	u.Address = h[:]
}
