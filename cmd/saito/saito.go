package saito

type Saito struct {
	Miner      *Miner
	Mempool    *Mempool
	Blockchain *Blockchain
	Wallet     *Wallet
}

func InitSaito() *Saito {
	s := Saito{}
	s.Mempool = NewMempool()
	s.Blockchain = NewBlockchain()
	s.Wallet = NewWallet()
	s.Miner = NewMiner(s.Wallet)
	return &s
}
