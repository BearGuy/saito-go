package saito

type Saito struct {
	Mempool    Mempool
	Blockchain Blockchain
}

func InitSaito() *Saito {
	s := &Saito{}
	s.Mempool = NewMempool()
	s.Blockchain = NewBlockchain()
	return s
}
