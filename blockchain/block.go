package blockchain

import (
	"github.com/poslegm/blockchain-chat/network"
	"time"
)

type Block struct {
	//parent of this block
	parent     *Block

	//children of this block, should only be one, but if chain splits we
	//need to know which we should select
	children   []*Block

	//double sha256 hash of this block
	hash       network.ShaHash

	//double sha256 hash of parent of this block
	parentHash network.ShaHash

	//block position in chain
	height     int32

	//creation time
	timestamp  time.Time
}
