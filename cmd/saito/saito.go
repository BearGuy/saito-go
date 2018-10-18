package saito

type Saito struct {
	Mempool Mempool
}

func InitSaito() *Saito {
	s := &Saito{}
	s.Mempool = NewMempool()
	return s
}
