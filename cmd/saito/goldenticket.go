package saito

import (
	"encoding/hex"
	"math"
	"strconv"

	"github.com/btcsuite/btcd/btcec"
)

type GoldenTicket struct {
	name      string
	vote      bool
	target    []byte
	random    float64
	publickey *btcec.PublicKey
}

type GoldenTicketSolution struct {
	hash   []byte
	random float64
}

func NewGoldenTicket() *GoldenTicket {
	return &GoldenTicket{}
}

func (gt *GoldenTicket) isValidSolution(soln GoldenTicketSolution, blk *Block) bool {
	// Set the number of digits to match - difficultyOrder
	// As our difficulty is granular to several decimal places
	// we check that the last digit in the test is within the
	// difficultyGrain - the decimal part of the difficulty
	//
	difficulty := blk.ReturnDifficulty()
	difficultyOrder, difficultyGrain := math.Modf(difficulty)
	// := math.Floor(difficulty)

	// We are testing our generated hash against the hash of the previous block.
	// th is the test hash.
	// ph is the previous hash.
	//
	solnHex := hex.EncodeToString(soln.hash)
	blockHashHex := hex.EncodeToString(blk.ReturnHash())

	th, _ := strconv.ParseInt(solnHex[0:int(difficultyOrder)+1], 0, 64)
	ph, _ := strconv.ParseInt(blockHashHex[0:int(difficultyOrder)+1], 0, 64)

	if th >= ph && float64((th-ph)/16) <= difficultyGrain {
		return true
	}
	return false
}

func (gt *GoldenTicket) calculateSolution(soln GoldenTicketSolution, blk *Block, publickey *btcec.PublicKey) {
	gt.name = "golden ticket"
	gt.target = blk.ReturnHash()
	gt.vote = false //this.app.voter.returnDifficultyVote(blk.block.difficulty)
	gt.random = soln.random
	gt.publickey = publickey
}

func (gt *GoldenTicket) FindWinner(blk *Block) *btcec.PublicKey {
	//
	// sanity check
	//
	if blk == nil {
		return nil
	}
	if blk.transactions == nil {
		return nil
	}
	if len(blk.transactions) == 0 {
		return blk.creator
	}

	// winnerHash := DoubleHashB(float64ToByte(gt.random))
	// winnerHex := hex.EncodeToString(winnerHash)

	// maxNum, _ := strconv.ParseInt("0xffffffffffffffff", 0, 64)
	// winnerNum, _ := strconv.ParseInt(winnerHex[0:16], 0, 64)

	// winnerDecimal := winnerNum / maxNum

	// Add winning usable fee number
	// fmt.Println(winnerDecimal)
	return nil
}
