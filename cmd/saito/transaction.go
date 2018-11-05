package saito

import (
	"crypto/sha256"

	"github.com/cbergoon/merkletree"
)

// Transaction struct
type Transaction struct {
	transaction TxData
	size        uint64
	dmsg        string
	cfee        string
	ufee        string
	fee         string
	isValid     bool
	atr         uint64
	//trapdoor
}

// TData struct for specific transaction data
type TxData struct {
	id   uint64
	from []Slip
	to   []Slip
	ts   int64
	sig  string
	path []TxData
	gt   int
	ft   int
	msg  string
	msig string
	ps   int
	rb   int
}

// AddFrom adds Slips to the current transaction
func (tx *Transaction) AddFrom(fromAddress []byte, fromAmount float64) {
	slip := Slip{}
	slip.add = fromAddress
	slip.amt = fromAmount
	tx.transaction.from = append(tx.transaction.from, slip)
}

// AddTo adds Slips to the current transaction
func (tx *Transaction) AddTo(toAddress []byte, toAmount float64) {
	slip := Slip{}
	slip.add = toAddress
	slip.amt = toAmount
	tx.transaction.to = append(tx.transaction.to, slip)
}

// CalculateHash is needed for MerkleTree calculations
func (tx Transaction) CalculateHash() ([]byte, error) {
	h := sha256.New()
	serializedTX := []byte(tx.transaction.sig)
	if _, err := h.Write(serializedTX); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// Equals checks equality between merkleRoot.content
func (tx Transaction) Equals(other merkletree.Content) (bool, error) {
	return tx.transaction.sig == other.(Transaction).transaction.sig, nil
}
