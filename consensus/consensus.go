package consensus

import (
	"XianfengChain04/chain"
	"math/big"
)

type Consensus interface {
	FindNonce() int64
}

func NewPoW(block chain.Block) Consensus {
	init := big.NewInt(1)
	init.Lsh(init, 255 - DIFFICULTY)
	return PoW{block, init}
}

func NewPoS(block chain.Block) Consensus {
	return PoS{Block: block}
}

