package main

import (
	"XianfengChain04/chain"
	"fmt"
	"github.com/bolt-master"
)

const BLOCKS = "xiangfengchain04.db"

func main() {

	//打开数据库文件
	db, err := bolt.Open(BLOCKS, 0600, nil)
	if err != nil {
		panic(err.Error())
	}
    defer db.Close()//xxx.db.lock

    blockChain := chain.CreateChain(db)
    //创世区块
    err = blockChain.CreateGensis([]byte("hello world"))
    if err != nil {
    	fmt.Println(err.Error())
		return
	}
	//新增一个区块
	err = blockChain.CreateNewBlock([]byte("hello"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//测试
	//lastBlock := blockChain.GetLastBlock()
	//fmt.Println("最新区块是：", lastBlock)
	//
	//blocks, err := blockChain.GetAllBlocks()
    //if err != nil {
    //	fmt.Println(err.Error())
	//	return
	//}
	//for index, block := range blocks{
	//	fmt.Printf("第%d个区块：\n", index)
	//	fmt.Println(block)
	//}

	//迭代器功能测试
	for blockChain.HasNext() {
		block := blockChain.Next()
		fmt.Printf("迭代到第%d个区块，区块高度：", block.Height)
		fmt.Printf("区块hash：%v", block.Hash)
		fmt.Printf("区块的信息：%s", string(block.Data))
		fmt.Println()
	}
}
