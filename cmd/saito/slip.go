package saito

import "github.com/btcsuite/btcd/btcec"

// Slip struct contains outline for UTXO used in Saito
type Slip struct {
	add   btcec.PublicKey
	amt   float64
	gt    int
	bid   uint64
	tid   uint64
	sid   int
	bhash uint64
	lc    uint64
	ft    uint64
}

func NewSlip(publickey btcec.PublicKey, amt float64, gt int) *Slip {
	slip := Slip{}
	slip.add = publickey
	slip.amt = amt
	slip.gt = gt
	return &slip
}
