package saito

import (
	"bytes"
	"fmt"
	"time"
)

// Index contains block metadata stored by the blockchain struct
type Index struct {
	bid          []int64
	hash         [][]byte
	prevhash     [][]byte
	blockID      []int64
	minTXID      []int64
	maxTXID      []int64
	timestamps   []int64
	longestChain []int
	burnfee      []float64
}

// Blockchain is the heart of saito
type Blockchain struct {
	heartbeat                 int64
	maxHeartbeat              int64
	genesisPeriod             int64
	genesisTS                 int64
	genesisBlockID            int64
	forkGuard                 int64
	forkID                    int64
	forIDMod                  int64
	oldLC                     int64
	index                     Index
	blocks                    []Block
	blockHashmap              map[string]int64
	lastHash                  []byte
	lastBid                   int64
	LCHashmap                 []int64
	longestChain              int
	LowestAcceptableTimestamp int64
	LowestAcceptableBlockId   int64
	LowestAcceptableHash      []byte
}

func NewBlockchain() Blockchain {
	bchain := Blockchain{}
	bchain.index = Index{}
	bchain.genesisPeriod = 12160
	bchain.blockHashmap = make(map[string]int64)
	return bchain
}

func (bchain *Blockchain) AddBlock(blk Block) {
	// verify that it's a valid block before appending
	bid := blk.id
	ts := blk.unixtime
	hash := blk.merkle
	prevhash := blk.prevhash

	// ignore pre-genesis blocks

	if ts < bchain.genesisTS || bid < bchain.genesisBlockID {
		// cannot add this block
		return
	}

	if bchain.IsHashIndexed(hash) {
		// block is already included in the chain
		return
	}

	// if the previous block hash was not indexed, we want to fetch it
	// if ts < bchain.LowestAcceptableTimestamp {
	// }

	// insert indexes
	pos := binaryInsert(bchain.index.timestamps, ts)
	if pos <= bchain.longestChain && len(bchain.index.longestChain) > 1 {
		bchain.longestChain++
	}

	index_len := len(bchain.index.hash)
	bchain.index.hash = append(bchain.index.hash[0:pos], hash)
	bchain.index.hash = append(bchain.index.hash, bchain.index.hash[pos+1:index_len-1]...)

	bchain.index.prevhash = append(bchain.index.prevhash[0:pos], prevhash)
	bchain.index.prevhash = append(bchain.index.prevhash, bchain.index.prevhash[pos+1:index_len-1]...)

	bchain.index.bid = append(bchain.index.bid[0:pos], bid)
	bchain.index.bid = append(bchain.index.bid, bchain.index.bid[pos+1:index_len-1]...)

	// bchain.index.maxTXID = append(bchain.index.maxTXID[0:pos], maxTXID)
	// bchain.index.maxTXID = append(bchain.index.maxTXID, bchain.index.maxTXID[pos+1:index_len-1]...)

	// bchain.index.minTXID = append(bchain.index.minTXID[0:pos], minTXID)
	// bchain.index.minTXID = append(bchain.index.minTXID, bchain.index.minTXID[pos+1:index_len-1]...)

	bchain.index.hash = append(bchain.index.hash[0:pos], hash)
	bchain.index.hash = append(bchain.index.hash, bchain.index.hash[pos+1:index_len-1]...)

	//   bchain.index.prevhash.splice(pos, 0, prevhash);
	//   bchain.index.bid.splice(pos, 0, bid);
	//   bchain.index.maxtid.splice(pos, 0, newblock.returnMaxTxId());
	//   bchain.index.mintid.splice(pos, 0, newblock.returnMinTxId());
	//   bchain.index.lc.splice(pos, 0, 0);
	//   bchain.index.bf.splice(pos, 0, newblock.returnBurnFeeValue());
	//   bchain.blocks.splice(pos, 0, newblock);

	hash_string := string(hash)
	bchain.blockHashmap[hash_string] = bid

	blk.isValid = true
	blk.prevhash = prevhash
	bchain.blocks = append(bchain.blocks, blk)

	// identify longest chain
	IAmTheLongestChain := false
	if bchain.longestChain == 0 && len(bchain.blocks) == 1 {
		if bchain.lastBid > 0 {
			if bytes.Equal(prevhash, bchain.lastHash) {
				IAmTheLongestChain = true
			}
		} else {
			bchain.index.longestChain[pos] = 1
			IAmTheLongestChain = true
		}

	}

	if bid >= bchain.index.bid[bchain.longestChain] {
		var search_pos int
		var search_bf float64
		var search_ts int64
		var search_hash []byte
		var search_prevhash []byte
		var sharedAncestorPos int

		if bytes.Equal(prevhash, bchain.index.hash[bchain.longestChain]) {
			IAmTheLongestChain = true
		} else {
			lchain_pos := bchain.longestChain
			nchain_pos := pos

			lchain_len := 0
			nchain_len := 0

			lchain_bf := bchain.index.burnfee[lchain_pos]
			nchain_bf := bchain.index.burnfee[nchain_pos]

			lchain_ts := bchain.index.timestamps[lchain_pos]
			nchain_ts := bchain.index.timestamps[nchain_pos]

			lchain_prevhash := bchain.index.prevhash[lchain_pos]
			nchain_prevhash := bchain.index.prevhash[nchain_pos]

			if nchain_ts >= lchain_ts {
				search_pos = nchain_pos - 1
			} else {
				search_pos = lchain_pos - 1
			}
			//
			// find the last shared ancestor
			//
			for search_pos >= 0 {
				search_ts = bchain.index.timestamps[search_pos]
				search_bf = bchain.index.burnfee[search_pos]
				search_hash = bchain.index.hash[search_pos]
				search_prevhash = bchain.index.prevhash[search_pos]

				if bytes.Equal(search_hash, lchain_prevhash) && bytes.Equal(search_hash, nchain_prevhash) {
					shared_ancestor_pos := search_pos
					search_pos = -1
				} else {
					if bytes.Equal(search_hash, lchain_prevhash) {
						lchain_len++
						lchain_prevhash := bchain.index.prevhash[search_pos]
						lchain_bf := lchain_bf + bchain.index.burnfee[search_pos]
					}

					if bytes.Equal(search_hash, nchain_prevhash) {
						nchain_prevhash = bchain.index.prevhash[search_pos]
						nchain_len++
						nchain_bf = nchain_bf + bchain.index.burnfee[search_pos]
					}

					sharedAncestorPos = search_pos
					search_pos--
				}
			}

			if nchain_len > lchain_len && nchain_bf >= lchain_bf {
				IAmTheLongestChain = true
			} else {
				// to prevent our system from being gamed, we
				// require the attacking chain to have equivalent
				// or greater aggreaget burn fee. Thsi ensures that
				// an attacker cannot lower difficulty, pump out a
				// ton of blocks, and then hike the difficulty only
				// at the last moment to claim the longest chain.

				// this is like the option above, except that we
				// have a choice of which blocl to support.
				if nchain_len == nchain_len && nchain_bf >= lchain_bf {
					//vote.prefersBlock(newBlock, this.returnLatestBlock())
					IAmTheLongestChain = true
				}
			}
		}
	} else {
		if bytes.Equal(blk.prevhash, bchain.lastHash) && bytes.Equal(blk.prevhash, nil) {
			// for { }
			// bchain.index.lc[] = true
			// storage.onChainReorganization
			// wallet.onChainReorganization
			// modules.onChainReorganization
			IAmTheLongestChain = true
			bchain.lastHash = hash
			bchain.longestChain = pos
			// modules.updateBalance
			bchain.LowestAcceptableTimestamp = ts
			bchain.LowestAcceptableBlockId = bid
			bchain.LowestAcceptableHash = hash

			// storage.saveOptions()
		}
	}

	if IAmTheLongestChain && len(bchain.index.hash) > 1 {
		bchain.longestChain = pos
		bchain.index.longestChain[pos] = 1
		bchain.blockHashmap[string(hash)] = blk.id

		// this.app.miner.stopMining()

		// edge case
		LCAtThisBid := false
		locOfLastBlock := -1

		if bchain.LowestAcceptableBlockId == (blk.id - 1) {
			if bytes.Equal(bchain.LowestAcceptableHash, blk.prevhash) {
				for h := pos - 1; h >= 0; h-- {
					if bytes.Equal(bchain.index.hash[h], blk.prevhash) {
						locOfLastBlock = h
					}
					if bchain.index.longestChain[h] == 1 && bchain.index.bid[h] == blk.id-1 {
						LCAtThisBid = false
					}
					if !LCAtThisBid && locOfLastBlock != -1 {
						if bchain.index.longestChain[locOfLastBlock] == 0 {
							bchain.index.longestChain[locOfLastBlock] = 1
							// storage.onChainReorganization()
							// wallet.onChainReorganization()
							// modules.onChainReorganization()
						}
					}
				}
			}
		}
		sharedAncestorHash := bchain.index.hash[sharedAncestorPos]
		newHashToHuntFor := blk.returnHash()
		var newBlockHashes [][]byte
		var newBlocksIdxs []int64
		var newBlockIds []int64
		var oldHashToHuntFor []byte
		var oldBlockHashes [][]byte
		var oldBlockIdxs []int64
		var oldBlockIds []int64
		var rewriteFromStart bool

		if sharedAncestorPos == -1 && bchain.index.longestChain[0] == 0 {
			// for g := pos - 1
			if rewriteFromStart && len(bchain.blocks) > 0 {
			}
		}

		if bytes.Equal(blk.prevhash, oldHashToHuntFor) {
			newBlockHashes.push(bchain.index.hash[pos])
		}

	} else {
		fmt.Println(" .... into wind:    ", time.Now().Unix())
	}
	fmt.Println("---- Added Block To Blockchain! ----")
	fmt.Println(blk)
}

func (bchain *Blockchain) ReturnLastBlock() Block {
	if len(bchain.blocks) != 0 {
		return bchain.blocks[len(bchain.blocks)-1]
	}
	return Block{}
}

func (bchain *Blockchain) IsHashIndexed(hash []byte) bool {
	for _, v := range bchain.index.hash {
		if bytes.Equal(hash, v) {
			return true
		}
	}

	return false
}

func binaryInsert(list []int64, item int64) int {
	start := 0
	end := len(list)

	for start < end {

		pos := (start + end) >> 1
		cmp := item - list[pos]

		if cmp == 0 {
			start = pos
			end = pos
			break
		} else if cmp < 0 {
			end = pos
		} else {
			start = pos + 1
		}
	}

	// if !search {
	// 	list.splice(start, 0, item)
	// }

	return start
}
