package saito

type Saito struct {
	Mempool    *Mempool
	Blockchain *Blockchain
	Wallet     *Wallet
}

func InitSaito() *Saito {
	s := Saito{}
	s.Mempool = NewMempool()
	s.Blockchain = NewBlockchain()
	s.Wallet = NewWallet()
	return &s
}
