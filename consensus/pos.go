package consensus

import (
	"XianfengChain04/chain"
	"fmt"
)

type PoS struct {
	chain.Block
}

func (pos PoS) FindNonce() int64 {
	fmt.Println("这是共识机制PoS算法的实现。。。")
	return 0
}
