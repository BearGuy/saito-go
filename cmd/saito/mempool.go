package saito

import (
	"fmt"
	"os"
	"time" // "github.com/bearguy/saito-go/cmd/crypto"
)

// Mempool is our queueing solution for transactions and producing blocks
type Mempool struct {
	transactions       []Transaction
	BurnFee            float64
	BundlingFeesNeeded float64
	FeesAcquired       float64
	Heartbeat          float64
	Starttime          int64
}

func NewMempool() Mempool {
	var mempool = Mempool{}
	mempool.BundlingFeesNeeded = 2
	mempool.BurnFee = 2
	mempool.Heartbeat = 0.06666666666
	mempool.FeesAcquired = 0
	mempool.Starttime = time.Now().Unix()
	return mempool
}

func (m *Mempool) Bundle() {
	for {
		if m.BundlingFeesNeeded <= 0 {
			fmt.Println("YOU CAN PRODUCE A BLOCK!")
			m.produceBlock(&Block{})
			os.Exit(3)
		} else {
			m.BundlingFeesNeeded = m.BurnFee - m.Heartbeat*float64(time.Now().Unix()-m.Starttime)
			DisplayBurnFeeCountdown(m.BundlingFeesNeeded)
		}
	}
}

func DisplayBurnFeeCountdown(bundling_fees float64) {
	t := time.Now()
	time.Sleep(1 * time.Second)

	// fmt.Println(time.Now().Format(time.RFC850), "----", bundling_fees)
	fmt.Printf("%d-%02d-%02dT%02d:%02d:%02d-00:00 ---- %0.10f\n",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), bundling_fees)
}

func (m *Mempool) produceBlock(lastBlock *Block) Block {
	blk := NewBlock()
	if lastBlock != nil {
		blk.prevhash = lastBlock.merkle
	}

	blk.merkle = ReturnMerkleTreeRoot(blk.transactions)
	// if err {
	// 	fmt.Println(err)
	// }

	return blk
}
