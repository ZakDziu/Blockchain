package service

import (
	"blockchain/block"
	"time"
)

func StartAddBlockService(bc *block.Blockchain) {
	time.Sleep(40 * time.Second)
	bc.AddBlock()
	StartAddBlockService(bc)
}
