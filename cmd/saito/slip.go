package saito

// Slip struct contains outline for UTXO used in Saito
type Slip struct {
	add   []byte
	amt   float64
	gt    uint64
	bid   uint64
	tid   uint64
	sid   uint64
	bhash uint64
	lc    uint64
	ft    uint64
}
