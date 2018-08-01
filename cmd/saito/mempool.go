package saito

import (
	"fmt"
	"time"
	"github.com/bearguy/saito-go/cmd/crypt"
)

// Mempool is our queueing solution for transactions and producing blocks
type Mempool struct {
	//this.data_directory             = path.join(__dirname, '../data/');

	//ioutil.ReadFile("/tmp/dat")

	//this.transactions               = []; // array
	transactions []Transaction
	downloads    []Block
	blocks       []Block
	recovered    []Transaction

	// mempool safety caps
	//
	// temporary safety limits designed to help avoid
	// spam attacks taking down everyone in the network
	// by preventing propoagation past a certain point
	//

	//this.transaction_size_cap       = 1024000000;// bytes hardcap 1GB
	transactionSizeCap uint64

	//this.transaction_size_current   = 0.0;
	transactionSizeCurrent uint64

	//this.block_size_cap             = 1024000000; // bytes hardcap 1GB
	blockSizeCap uint64

	//this.block_size_current         = 0.0;
	blockSizeCurrent float64

	//this.download_size_cap          = 512000000; // bytes hardcap 1GB
	downloadSizeCap uint64

	//this.download_size_current      = 0.0;
	downloadSizeCurrent float64

	//this.currently_processing       = 0;
	currentlyProcessing bool

	// this.processing_timer           = null;
	ProcessingTimer int

	// this.processing_speed           = 50; // 0.1 seconds (add blocks)
	ProcessingSpeed float64

	// this.currently_downloading      = 0;
	CurrentlyDownloading bool

	// this.downloading_timer          = null;
	DownloadingTimer int64

	// this.downloading_speed          = 300; // 0.3 seconds
	DownloadingSpeed int64

	// this.currently_bundling         = 0;
	CurrentlyBundling bool

	// this.bundling_timer             = null;
	BundlingTimer int64

	// this.bundling_speed             = 400; // 0.4 seconds
	BundlingSpeed int64

	// this.currently_clearing         = 0;   // 1 when clearing blks + txs
	CurrentlyClearing bool

	//this.currently_creating         = 0;   // 1 when creating a block from our txs
	CurrentlyCreating bool

	//this.bundling_fees_needed       = "-1";
	BundlingFeesNeeded float64

	// this.load_time                  = new Date().getTime();
	LoadTime int64

	//this.load_delay                 = 4000; // delay on startup so we have time
	LoadDelay int64
}

func (m *Mempool) init() {

}

func (m *Mempool) run(bchain *Blockchain) {
	for {
		if m.BundlingFeesNeeded == 0 {
			// newblock := m.produceBlock(lastBlock)
			// var newBlock Block
			if len(bchain.blocks) != 0 {
				lastBlock := bchain.returnLastBlock()
				newBlock := m.produceBlock(lastBlock)
				bchain.addBlock(newBlock)
			} else {
				newBlock := NewBlock()
				newBlock.id = 0
				bchain.addBlock(newBlock)
			}
		} else {
			fmt.Println("%d: %f", time.Now().Unix(), m.BundlingFeesNeeded)
			m.BundlingFeesNeeded = m.BundlingFeesNeeded - bchain.feestep
		}
	}
}

// func (m *Mempool) addBlock(blk Block) bool {
// 	if len(blk.hash) != 0 {
// 		m.blocks = append(m.blocks, blk)
// 		return true
// 	}
// 	return false
// }

// func (m *Mempool) processBlock(blk *Block) {

// }

func (m *Mempool) produceBlock(lastBlock Block) Block {
	blk := NewBlock()
	blk.prevhash = lastBlock.hash
	blk.hash := crypt.ReturnMerkleRoot(blk.transactions)
	return blk
}
