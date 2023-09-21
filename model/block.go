package model

import (
	"encoding/hex"
)

type BlockRequest struct {
	HashID               string `json:"hashID"`
	AddressSender        string `json:"addressSender"`
	AddressRecipient     string `json:"addressRecipient"`
	BlockNumber          int    `json:"blockNumber"`
	TransactionCreatedAt int64  `bson:"transactionCreatedAt"`
	PageSize             int    `json:"pageSize"`
	Page                 int    `json:"page"`
}

func (b *BlockRequest) Validate() error {
	_, err := hex.DecodeString(b.HashID)
	if err != nil {
		return err
	}

	_, err = hex.DecodeString(b.AddressSender)
	if err != nil {
		return err
	}

	_, err = hex.DecodeString(b.AddressRecipient)
	if err != nil {
		return err
	}

	return nil
}
