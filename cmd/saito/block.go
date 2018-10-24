package saito

import (
	"time"
)

// Block struct, the building block *snerk* of Saito
type Block struct {
	id           int
	unixtime     int64
	prevhash     []byte
	merkle       []byte
	burnFee      float64
	transactions []Transaction
	isValid      bool
}

// NewBlock is a constructor for creaing new blocks
func NewBlock() Block {
	block := Block{}
	block.unixtime = time.Now().Unix()
	block.transactions = nil
	return block
}
