package saito

import (
	"fmt"
	"os"
	"time"
)

// Mempool is our queueing solution for transactions and producing blocks
type Mempool struct {
	transactions       []Transaction
	BundlingFeesNeeded float64
	FeeStep            float64
}

func (m *Mempool) Bundle() {
	for {
		if m.BundlingFeesNeeded <= 0 {
			fmt.Println("YOU CAN PRODUCE A BLOCK! :D")
			os.Exit(3)
		} else {
			fmt.Println("%d: %f", time.Now().Unix(), m.BundlingFeesNeeded)
			m.BundlingFeesNeeded = m.BundlingFeesNeeded - m.FeeStep
		}
	}
}

func NewMempool() Mempool {
	var mempool = Mempool{}
	mempool.BundlingFeesNeeded = 2
	mempool.FeeStep = 0.001
	return mempool
}

// func (m *Mempool) produceBlock(lastBlock Block) Block {
// 	blk := NewBlock()
// 	blk.prevhash = lastBlock.hash
// 	blk.hash := crypt.ReturnMerkleRoot(blk.transactions)
// 	return blk
// }
