package dto

import (
	"blockchain/block"
	"fmt"
	"time"
)

type TransactionDto struct {
	ID                    string    `json:"transactionId"`
	AddressSender         string    `json:"addressSender"`
	AddressRecipient      string    `json:"addressRecipient"`
	BlockNumber           int       `json:"blockNumber"`
	NumberOfConfirmations int       `json:"numberOfConfirmations"`
	CreatedAt             time.Time `json:"createdAt"`
	Sum                   float64   `json:"sum"`
	Gas                   float64   `json:"gas"`
}

func PrepareDataForTransactions(blocks []*block.Block, addressSender, addressRecipient string) []TransactionDto {
	var transactions []TransactionDto
	for _, b := range blocks {
		if len(b.Data) != 0 {
			for _, t := range b.Data {
				transaction := TransactionDto{
					NumberOfConfirmations: b.NumberOfConfirmations,
					BlockNumber:           b.BlockNumber,
					ID:                    fmt.Sprintf("%x", t.ID),
					AddressSender:         fmt.Sprintf("%x", t.AddressSender),
					AddressRecipient:      fmt.Sprintf("%x", t.AddressRecipient),
					Sum:                   t.Sum,
					Gas:                   t.Gas,
					CreatedAt:             time.UnixMicro(t.CreatedAt),
				}
				if addressSender != "" || addressRecipient != "" {
					if addressSender != "" && addressRecipient != "" {
						if transaction.AddressSender == addressSender && transaction.AddressRecipient == addressRecipient {
							transactions = append(transactions, transaction)
						} else {
							continue
						}
					}
					if transaction.AddressSender == addressSender || transaction.AddressRecipient == addressRecipient {
						transactions = append(transactions, transaction)
					} else {
						continue
					}
				}
				transactions = append(transactions, transaction)
			}
		}
	}
	return transactions
}
