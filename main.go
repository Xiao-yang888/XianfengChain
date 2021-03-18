package main

import (
	"XianfengChain04/chain"
	"XianfengChain04/client"
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
    cmdClient := client.CmdClient{blockChain}


    //cmdClient.Help()
    cmdClient.Run()
}


