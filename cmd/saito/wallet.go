package saito

import (
	"github.com/btcsuite/btcd/btcec"
)

type Wallet struct {
	balance     float64
	privatekey  *btcec.PrivateKey
	publickey   *btcec.PublicKey
	identifier  string
	inputs      []Slip
	outputs     []Slip
	spends      []int
	default_fee float64
	version     float32
	pending     []Transaction

	inputs_hmap                []byte
	inputs_hmap_counter        int64
	inputs_hmap_counter_limit  int64
	outputs_hmap               []byte
	outputs_hmap_counter       int64
	outputs_hmap_counter_limit int64

	is_testing bool
}

func NewWallet() *Wallet {
	wallet := Wallet{}
	// w                       = {};
	// this.wallet.balance               = "0.0";
	// this.wallet.privatekey            = "";
	// this.wallet.publickey             = "";
	// this.wallet.identifier            = "";
	// this.wallet.inputs                = [];
	// this.wallet.outputs               = [];
	// this.wallet.spends                = [];
	// this.wallet.default_fee           = 2;
	// this.wallet.version               = 2.06;
	// this.wallet.pending               = []; // sent but not seen

	// this.inputs_hmap                  = [];
	// this.inputs_hmap_counter 	    = 0;
	// this.inputs_hmap_counter_limit    = 100000;
	// this.outputs_hmap                 = [];
	// this.outputs_hmap_counter 	    = 0;
	// this.outputs_hmap_counter_limit   = 100000;

	// this.is_testing                   = false;
	wallet.balance = 0.0
	wallet.privatekey, _ = GenerateKeys()
	wallet.publickey = ReturnPublicKey(wallet.privatekey)
	wallet.identifier = ""
	wallet.default_fee = 2.0
	wallet.version = 2.06

	wallet.inputs_hmap_counter = 0
	wallet.inputs_hmap_counter_limit = 100000

	wallet.outputs_hmap_counter = 0
	wallet.outputs_hmap_counter_limit = 100000

	return &wallet
}

func (w *Wallet) ReturnPublicKey() *btcec.PublicKey {
	return w.publickey
}

func (w *Wallet) ReturnPrivateKey() *btcec.PrivateKey {
	return w.privatekey
}

func (w *Wallet) ReturnDefaultFee() float64 {
	return w.default_fee
}

func (w *Wallet) ReturnAvailableInputs() float64 {
	// var value   = Big(0.0);
	// this.purgeExpiredSlips();

	// lowest acceptable block_id for security (+1 because is next block, +1 for safeguard)
	// var lowest_block = this.app.blockchain.returnLatestBlockId() - this.app.blockchain.returnGenesisPeriod();
	//     lowest_block = lowest_block+2;

	// // calculate value
	for i := 0; i < len(w.inputs); i++ {
		if w.spends[i] == 0 && w.inputs[i].lc == 1 {
			//if w.inputs[i].lc == 1 && w.inputs[i].bid >= lowest_block {
			return w.inputs[i].amt
		}
	}

	// return value.toFixed(8);

	return 0.0
}

func (w *Wallet) ReturnAdequateInputs(amt float64) *[]Slip {
	// var utxiset = [];
	var utxiset = []Slip{}
	// var value   = Big(0.0);
	// var bigamt  = Big(amt);

	// More to cover with Blockchain store
	// var lowest_block = this.app.blockchain.returnLatestBlockId() - this.app.blockchain.returnGenesisPeriod();

	//
	// this adds a 1 block buffer so that inputs are valid in the future block included
	//
	// lowest_block = lowest_block+2;

	// this.purgeExpiredSlips();

	for i := 0; i < len(w.inputs)-1; i++ {
		if w.spends[i] == 0 || i >= len(w.spends) {
			slip := w.inputs[i]
			//if slip.lc == 1 && slip.bid >= lowest_block {
			if slip.lc == 1 {
				// hmap needs to be implemented
				//if this.app.mempool.transactions_inputs_hmap[slip.returnIndex()] != 1 {
				w.spends[i] = 1
				utxiset = append(utxiset, slip)
				if slip.amt >= amt {
					return &utxiset
				}
			}
		}
	}
	return &[]Slip{}
}

// Data stores are required to pass as a shared resource. This currently does not function well

// func (w *Wallet) purgeExpiredSlips() {
// 	gid := blockchain.returnGenesisBlockId()
// 	for i := len(w.inputs) - 1; i >= 0; i-- {
// 		if w.inputs[i].bid < gid {
// 			w.inputs.splice(i, 1)
// 			w.spends.splice(i, 1)
// 		}
// 	}
// 	for j := len(w.outputs) - 1; j >= 0; j-- {
// 		if w.outputs[j].bid < gid {
// 			w.outputs.splice(j, 1)
// 		}
// 	}
// }
