package saito

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type Miner struct {
	mining_active bool
	// mining_speed  int64
	wallet Wallet
}

// Need a channel for passing blocks through
// instigate at the top level

func NewMiner(wallet *Wallet) *Miner {
	m := Miner{mining_active: false, wallet: *wallet}
	return &m
}

func (m *Miner) StartMining(prevblk *Block) {
	// for {
	//
	// Check that we have a previous block to mine on.
	//
	if prevblk == nil {
		return
	}
	if !prevblk.isValid {
		return
	}

	if prevblk.id == 1 {
		return
	}
	publickey := m.wallet.ReturnPublicKey()

	gtSolutionChannel := make(chan GoldenTicketSolution)

	// for {
	go generateSolution(gtSolutionChannel, publickey.SerializeUncompressed())
	// }

	gt := NewGoldenTicket()

	proposedSolution := <-gtSolutionChannel
	fmt.Println(proposedSolution)
	if gt.isValidSolution(proposedSolution, prevblk) {
		gt.calculateSolution(proposedSolution, prevblk, publickey)
		fmt.Println("WE FOUND A WINNER :D :D :D")

		// gt.findWinner()
		winningNode := gt.FindWinner(prevblk)

		if winningNode == nil {
			winningNode = m.wallet.ReturnPublicKey()
		}

		// let total_fees_needed_for_prevblk    = await this.app.burnfee.returnBurnFeePaidForThisBlock(prevblk);
		// let total_fees_available_for_creator = prevblk.returnAvailableFees(prevblk.block.creator);
		// let total_fees_in_prevblk            = prevblk.returnFeesTotal();
		// let creator_surplus			 = Big(total_fees_available_for_creator).minus(Big(total_fees_needed_for_prevblk));
		// let total_fees_for_miners_and_nodes  = Big(total_fees_in_prevblk).minus(creator_surplus).plus(prevblk.returnCoinbase());

		// // miner and node shares
		// let miner_share = total_fees_for_miners_and_nodes.div(2).toFixed(8);
		// let node_share =  total_fees_for_miners_and_nodes.minus(Big(miner_share)).toFixed(8);

		// calculate the winnings

		// build the transaction
		tx := NewTransaction()
		tx.transaction.ts = time.Now().UnixNano()
		tx.transaction.typ = 1
		tx.transaction.msg, _ = json.Marshal(gt)

		//tx.transaction.msg = gt.solution

		availableFunds := m.wallet.ReturnAvailableInputs()
		defaultFee := m.wallet.ReturnDefaultFee()

		if availableFunds >= defaultFee && len(m.wallet.inputs) > 1 {
			//
			// set a low fee for testing
			//
			tx.transaction.from = *m.wallet.ReturnAdequateInputs(0.0001)
			if tx.transaction.from == nil {
				// can't attempt a solution with no slips
				return
			}
		} else {
			tx.transaction.from = append(tx.transaction.from, *NewSlip(*m.wallet.ReturnPublicKey(), 0.0, 1))
			// new saito.slip(this.app.wallet.returnPublicKey(), 0.0, 1));
		}

		var total_submitted float64
		for r := 0; r < len(tx.transaction.from); r++ {
			total_submitted += tx.transaction.from[r].amt
		}

		// to slips
		tx.transaction.to = append(tx.transaction.to, *NewSlip(*m.wallet.ReturnPublicKey(), 1000, 1))
		tx.transaction.to = append(tx.transaction.to, *NewSlip(*winningNode, 1000, 1))

		// set slip ids
		for i := 0; i < len(tx.transaction.to); i++ {
			tx.transaction.to[i].sid = i
		}

		// change address
		change_amount := total_submitted - m.wallet.ReturnDefaultFee()

		//
		// set a low fee for testing
		// TODO: remove one module debugged
		//
		change_amount = total_submitted - 0.0001

		if change_amount > 0 {
			tx.transaction.to = append(tx.transaction.to, *NewSlip(*m.wallet.ReturnPublicKey(), change_amount, 0))
			tx.transaction.to[len(tx.transaction.to)-1].gt = 0
		}

		// tx = this.app.wallet.signTransaction(tx);

		// this.app.network.propagateTransaction(tx);

		// propagate the transaction
	}
	// }
}

func generateSolution(c chan GoldenTicketSolution, publickey []byte) {
	for i := 0; i < 10; i++ {
		go func() {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			solutionSeed := [][]byte{publickey, float64ToByte(r.Float64())}
			c <- GoldenTicketSolution{hash: DoubleHashB(bytes.Join(solutionSeed, []byte(", "))), random: r.Float64()}
		}()
	}
}
