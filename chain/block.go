package chain

import (
	"XianfengChain04/consensus"
	"XianfengChain04/transaction"
	"XianfengChain04/utils"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

const VERSION = 0x00

/**
 *区块的结构体定义
 */
type Block struct {
	Height    int64 //高度
	Version   int64
	PrevHash  [32]byte
	Hash      [32]byte
	//默克尔根
	TimeStamp int64
	//Difficulty int64
	Nonce     int64
	//区块体
	Transactions []transaction.Transaction
}

func (block Block) GetHeight() int64 {
	return block.Height
}

func (block Block) GetVersion() int64 {
	return block.Version
}

func (block Block) GetTimeStamp() int64 {
	return block.TimeStamp
}

func (block Block) GetPrevHash() [32]byte {
	return block.PrevHash
}

func (block Block) GetTransactions() []transaction.Transaction {
	return block.Transactions
}
/**
 *计算区块的hash值并进行赋值
 */
func (block *Block) CalculateBlockHash()  {
	heightByte, _ := utils.Int2Byte(block.Height)
	versionByte, _ := utils.Int2Byte(block.Version)
	timeByte, _ := utils.Int2Byte(block.TimeStamp)
	nonceByte, _ := utils.Int2Byte(block.Nonce)

	blockByte := bytes.Join([] []byte{heightByte, versionByte, block.PrevHash[:], timeByte, nonceByte, block.Transactions}, []byte{})
    //为区块的hash字段赋值
	block.Hash = sha256.Sum256(blockByte)
}

/**
 *区块的序列化方法
 */
func (block *Block) Serialize() ([]byte, error) {
	//缓冲区
	buff := new(bytes.Buffer)
	encoder := gob.NewEncoder(buff)
	err := encoder.Encode(&block)
	return buff.Bytes(), err
}

/**
 *反序列化函数
 */
func Deserialize(data []byte) (Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	return block, err
}

/**
 *生成创世区块的函数
 */
func CreateGenesis(data []byte) Block {
	fmt.Println("创建创世区块数据并未存储到交易中")
	tx := transaction.Transaction{}
	genesis := Block{
		Height:    0,
		Version:   VERSION,
		PrevHash:  [32]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
		TimeStamp: time.Now().Unix(),
		Transactions:      []transaction.Transaction{tx},
	}

	//调用pow，实现hash计算和寻找nonce
	proof := consensus.NewPoW(gensis)
	hash, nonce := proof.FindNonce()
    genesis.Hash = hash
    genesis.Nonce = nonce

	return genesis
}

/**
 *生成新区块的功能函数
 */
func NewBlock(height int64, prev [32]byte, data []byte) Block {
	tx := transaction.Transaction{}
	newBlock := Block{
		Height:    height + 1,
		Version:   VERSION,
		PrevHash:  prev,
		TimeStamp: time.Now().Unix(),
		Transactions:      []transaction.Transaction{tx},
	}

	proof := consensus.NewPoW(newBlock)
	hash, nonce := proof.FindNonce()
	newBlock.Hash = hash
	newBlock.Nonce = nonce
	return newBlock
}