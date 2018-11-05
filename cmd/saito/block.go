package saito

import (
	"time"
)

// Block struct, the building block *snerk* of Saito
type Block struct {
	id           int64
	unixtime     int64
	prevhash     []byte
	merkle       []byte
	burnFee      float64
	transactions []Transaction
	isValid      bool
	lc           int64
}

// NewBlock is a constructor for creaing new blocks
func NewBlock() Block {
	block := Block{}
	block.id = 1
	block.unixtime = time.Now().Unix()
	block.transactions = nil
	block.lc = 0
	return block
}
