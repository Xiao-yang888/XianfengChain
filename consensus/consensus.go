package consensus

import (
	"math/big"
)

type Consensus interface {
	FindNonce() ([32]byte, int64)
}

/**
 *定义区块结构体的接口标准
 */
type BlockInterface interface {
	GetHeight()    int64
	GetVersion()   int64
	GetTimeStamp() int64
	GetPrevHash()  [32]byte
	GetData()      []byte
}
func NewPoW(block BlockInterface) Consensus {
	init := big.NewInt(1)
	init.Lsh(init, 255 - DIFFICULTY)
	return PoW{block, init}
}

//func NewPoS(block BlockInterface) Consensus {
//	return PoS{Block: block}
//}

