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
	spends      []Slip
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
	// this.wallet                       = {};
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
