package main

import (
	"XianfengChain04/chain"
	"fmt"
)

func main() {
	fmt.Println("hello world")

	block0 := chain.CreateGenesis([]byte("hello world"))
	block1 := chain.NewBlock(block0.Height, block0.Haah, []byte("hello world"))

	fmt.Println(block0)
	fmt.Println(block1)
}
