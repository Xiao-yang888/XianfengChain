package consensus

import (
	"XianfengChain04/chain"
	"XianfengChain04/utils"
	"bytes"
	"crypto/sha256"
	"math/big"
)

const DIFFICULTY = 10 //难度系数


type PoW struct {
	Block chain.Block
	Target *big.Int
}

func (pow PoW) FindNonce() int64 {
	//给定一个nonce，获得区块的hash
	var nonce int64
	nonce = 0

	//无限循环
	for {
        hash := CalculateHash(pow.Block, nonce)

		//拿到系统的目标值
		target := pow.Target

		//比较大小

		result := bytes.Compare(hash[:], target.Bytes())
		//判断结果
		if result == -1 {
			return nonce
		}
		nonce++ //否则自增
	}
	return 0
}

/**
 *根据区块已有的信息和当前nonce的赋值，计算区块的hash
 */
func CalculateHash(block chain.Block, nonce int64) [32]byte {
	heightByte, _ := utils.Int2Byte(block.Height)
	versionByte, _ := utils.Int2Byte(block.Version)
	timeByte, _ := utils.Int2Byte(block.TimeStamp)
	nonceByte, _ := utils.Int2Byte(block.Nonce)

	blockByte := bytes.Join([][]byte{heightByte, versionByte, block.PrevHash[:], timeByte, nonceByte, block.Data}, []byte{})
	//计算系统的hash
	hash := sha256.Sum256(blockByte)
	return hash
}