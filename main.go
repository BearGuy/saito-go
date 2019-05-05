package main

import (
	"time"

	"github.com/bearguy/saito-go/cmd/saito"
)

func main() {
	saito := saito.InitSaito()
	for {
		newblock := saito.Mempool.Bundle(saito.Blockchain.ReturnLastBlock())
		newblock.AddCreator(saito.Wallet.ReturnPublicKey())
		saito.Blockchain.AddBlock(newblock)

		saito.Mempool.BurnFee = 2
		saito.Mempool.BundlingFeesNeeded = 2
		saito.Mempool.Starttime = time.Now().Unix()
	}
}
