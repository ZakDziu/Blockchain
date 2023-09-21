package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"time"

	"blockchain/auth"
	"blockchain/model"
	"blockchain/utils"

	"github.com/gin-gonic/gin"
)

type Transaction struct {
	ID               []byte
	AddressSender    []byte
	AddressRecipient []byte
	Sum              float64
	Gas              float64
	CreatedAt        int64
	Nonce            int
}

func (t *Transaction) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			t.AddressRecipient,
			t.AddressSender,
			utils.IntToHex(t.CreatedAt),
			utils.IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (t *Transaction) AddTransactionHash() {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	for nonce < maxNonce {
		data := t.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}
	}
	h := hash[:]
	t.ID = h[:]
}

func CreateNewTransaction(c *gin.Context, req model.TransactionRequest) (*Transaction, error) {
	var err error

	addressRecipient, _ := hex.DecodeString(req.AddressRecipient)
	transaction := &Transaction{
		AddressSender:    nil,
		AddressRecipient: addressRecipient,
		Sum:              req.Sum,
		Gas:              0,
		CreatedAt:        0,
	}

	transaction.AddressSender, err = auth.GetUserAddress(c)
	if err != nil {
		return nil, err
	}

	transaction.CreatedAt = time.Now().Unix()
	transaction.CalculateGas()
	transaction.AddTransactionHash()

	return transaction, nil
}
