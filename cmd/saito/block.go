package saito

import (
	"bytes"
	"encoding/binary"
	"math"
	"time"
)

// Block struct, the building block *snerk* of Saito
type Block struct {
	id           int64
	unixtime     int64
	hash         []byte
	prevhash     []byte
	prehash      []byte
	merkle       []byte
	burnFee      float64
	transactions []Transaction
	isValid      bool
	lc           int64
}

// NewBlock is a constructor for creaing new blocks
func NewBlock() Block {
	block := Block{}
	block.id = 1
	block.unixtime = time.Now().Unix()
	block.transactions = nil
	block.lc = 0
	return block
}

func (b *Block) ReturnHash() []byte {
	if len(b.hash) > 1 {
		return b.hash
	}
	b.prehash = DoubleHashB(b.ReturnSignatureSource())
	b.hash = DoubleHashB(
		bytes.Join([][]byte{b.prehash, b.prevhash}, []byte(", ")),
	)
	return b.hash
}

func (b *Block) ReturnSignatureSource() []byte {
	s := [][]byte{
		int64ToByte(b.unixtime),
		b.merkle,
		int64ToByte(b.id),
		float64ToByte(b.burnFee),
		b.hash,
	}
	return bytes.Join(s, []byte(", "))

	// return this.block.ts
	//   + this.block.creator
	//   + this.block.merkle
	//   + this.block.id
	//   + JSON.stringify(this.block.bf)
	//   + this.block.difficulty
	//   + this.block.paysplit
	//   + this.block.treasury
	//   + this.block.coinbase
	//   + this.block.vote
	//   + this.block.reclaimed;
}

func float64ToByte(f float64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

func int64ToByte(i int64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(i))
	return buf[:]
}
