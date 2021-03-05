package main

import (
	"XianfengChain04/chain"
	"fmt"
)

func main() {
	blockChain := chain.CreateChainWithGensis([]byte("hello world"))
	blockChain.CreateNewBlock([]byte("hello"))
	fmt.Println("区块链中的区块个数:", len(blockChain.Blocks))

	fmt.Println("区块0：", blockChain.Blocks[0])
	//fmt.Println("区块1：", blockChain.Blocks[1])

	firstBlock := blockChain.Blocks[0]
	firstBytes, err := firstBlock.Serialize()
	if err != nil {
		panic(err.Error())
	}
	//反序列化，验证逆过程
	deFirstBlock, err := chain.Deserialize(firstBytes)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(deFirstBlock.Data))
}
