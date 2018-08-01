package saito

// Index contains block metadata stored by the blockchain struct
type Index struct {
	hash []byte
	// 	prevhash:    [],                 // hash of previous block
	prevHash []byte
	// 	block_id:    [],                 // block id
	blockID int64
	// 	mintid:      [],                 // min tid
	minTXID int64
	// 	maxtid:      [],                 // max tid
	maxTXID int64
	// 	ts:          [],                 // timestamps
	timestamps int64
	// 	lc:          [],                 // is longest chain (0 = no, 1 = yes)
	longestChain bool
}

// Blockchain is the heart of saito
type Blockchain struct {
	//this.heartbeat               = 30;        // expect new block every 30 seconds
	heartbeat int64

	//this.max_heartbeat           = 120;       // burn fee hits zero every 120 seconds
	maxHeartbeat int64

	//this.genesis_period          = 12160;     // number of blocks before money disappears
	genesisPeriod int64
	// 90,000 is roughly a 30 day transient blockchain.
	//this.genesis_ts              = 0;         // unixtime of earliest block
	genesisTS int64

	//this.genesis_block_id        = 0;         // earliest block_id we care about
	genesisBlockID int64

	//this.fork_guard            = 120;       // discard forks that fall N blocks behind, this can
	forkGuard int64
	// result in a chain fork, so this needs to be long
	// enough that we reasonably decide that nodes that
	// cannot keep up-to-date with the network must resync
	//
	// the fork guard is used primarily when identifying
	// what blocks we can delete, since we must have the
	// full genesis period, plus whatever fork guard limit
	// suggests that someone can re-write the genesis chain
	//
	//this.fork_id                 = "";        // a string we use to identify our longest-chain
	forkID int64
	// generated deterministically from the block hashes
	// and thus unique for every fork
	//
	//this.fork_id_mod             = 10;	    // update fork id every 10 blocks
	forIDMod int64

	//this.old_lc                  = -1;	    // old longest-chain when processing new block
	oldLC int64
	// this will be set to the position of the current
	// head of the longest chain in our indexes before
	// we try to validate the newest block, so that we
	// can gracefully reset to the known-good block if
	// there are problems

	/////////////
	// Indexes //
	/////////////
	//
	// These hold the most important data needed to interact
	// with the blockchain objects, and must be kept for the
	// entire period the block is part of the transient
	// blockchain.
	//
	// If we add or delete these items, we must make changes
	// to the following functions
	//
	//    addBlockToBlockchain (add)
	//    addBlockToBlockchainPartTwo (lc_hashmap only)
	//    purgeArchivedData (delete)
	//
	// this.index = {
	// 	hash:        [],                 // hashes
	// 	prevhash:    [],                 // hash of previous block
	// 	block_id:    [],                 // block id
	// 	mintid:      [],                 // min tid
	// 	maxtid:      [],                 // max tid
	// 	ts:          [],                 // timestamps
	// 	lc:          [],                 // is longest chain (0 = no, 1 = yes)
	// 	burnfee:     [],                 // burnfee per block
	// 	feestep:     []                  // feestep per block
	// };
	//this.blocks         = [];
	index        []Index
	blocks       []Block
	LCHashmap    []int64
	longestChain int64

	// 	burnfee:     [],                 // burnfee per block
	burnfee float64
	// 	feestep:     []                  // feestep per block
	feestep float64

	//this.block_hashmap  = [];
	// this.lc_hashmap     = []; 	     // hashmap index is the  block hash and contains
	// 					// 1 or 0 depending on if they are the longest
	// 					// chain or not.
	// this.longestChain   = -1;          // position of longest chain in indices

	///////////////////
	// monitor state //
	///////////////////
	// this.currently_reclaiming = 0;
	currentlyReclaiming bool

	//this.currently_indexing = 0;
	currentlyIndexing bool

	//this.block_saving_timer = null;
	blockSavingTimer int64

	//this.block_saving_timer_speed = 10;
	blockSavingTimerSpeed int64

	//
	// this are set to the earliest block that we process
	// to ensure that we don't load missing blocks endlessly
	// into the past.
	//
	// the blk_limit is checked in the storage class when
	// validating slips as part of its sanity check so that
	// it does not cry foul if it lacks a full genesis period
	// worth of blocks but cannot validate slips.
	//
	//this.ts_limit = -1;
	tsLimit int64

	//this.blk_limit = -1;
	blockLimit int64

	///////////////
	// Callbacks //
	///////////////
	//this.callback_limit   = 100;        // only run callbacks on the last X blocks
	// this should be at least 10 to be safe, as
	// we delete data from blks past this limit
	// and that can complicate syncing to lite-clients

	// if we send full blocks right away
	//this.run_callbacks 	  = 1;	     // 0 for central nodes focused on scaling

	//
	// these are used to tell the blockchain class from
	// what block it should start syncing. we load these
	// from our options file on initialize.
	//
	//this.previous_block_id = -1;
	//this.previous_ts_limit = -1;
	//this.previous_block_hash = "";
}

func (bchain *Blockchain) addBlock(blk Block) {
	bchain.blocks = append(bchain.blocks, blk)
}

func (bchain *Blockchain) returnLastBlock() Block {
	if len(bchain.blocks) != 0 {
		return bchain.blocks[len(bchain.blocks)-1]
	}
	return Block{}
}
