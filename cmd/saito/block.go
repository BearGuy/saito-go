package saito

import (
	"time"
)

// Block struct, the building block *snerk* of Saito
type Block struct {
	id            int
	unixtime      int64
	prevhash      []byte
	merkle        []byte
	transactions  []Transaction
	burnFee       float64
	difficulty    float64
	paysplit      float64
	treasury      float64
	coinbase      float64
	reclaimed     float64
	paysplitVote  int
	segadd        []Transaction
	size          int
	atr           int
	atrLowerLimit int
	atrFeeCurve   int
	isValid       bool
	minTXid       int
	maxTXid       int
	filename      string
	hash          []byte
	// transactions      []Transaction
	confirmations     int //defaults to -1
	prevalidated      int //default 0 -- set to 1 to forceAdd to blockchain without running callbacks
	peerOrigin        string
	averageFee        string
	segAddMax         int
	segAddMap         []Transaction
	segAddEnabled     int
	segAddCompression int
	saveBlockID       int
	saveDatabaseID    int
}

// NewBlock is a constructor for creaing new blocks
func NewBlock() Block {
	block := Block{}
	block.unixtime = time.Now().Unix()
	return block
}

func (b *Block) bundBlock() {

}
