package consensus

import (
	"XianfengChain04/utils"
	"bytes"
	"crypto/sha256"
	"math/big"
)


//拿到区块的属性值
//1，通过结构体引用，引用block结构体，然后访问其属性，比如block，height
//2，接口


const DIFFICULTY = 10 //难度系数

type PoW struct {
	Block BlockInterface
	Target *big.Int
}

func (pow PoW) FindNonce() ([32]byte, int64) {
	//给定一个nonce，获得区块的hash
	var nonce int64
	nonce = 0

	//无限循环
	hashBig := new(big.Int)
	for {
        hash := CalculateHash(pow.Block, nonce)

		//拿到系统的目标值
		target := pow.Target

		//比较大小
        hashBig = hashBig.SetBytes(hash[:])
		result := hashBig.Cmp(target)
		//判断结果
		if result == -1 {
			return hash, nonce
		}
		nonce++ //否则自增
	}
}

/**
 *根据区块已有的信息和当前nonce的赋值，计算区块的hash
 */
func CalculateHash(block BlockInterface, nonce int64) [32]byte {
	heightByte, _ := utils.Int2Byte(block.GetHeight())
	versionByte, _ := utils.Int2Byte(block.GetVersion())
	timeByte, _ := utils.Int2Byte(block.GetTimeStamp())
	nonceByte, _ := utils.Int2Byte(nonce)

	prev := block.GetPrevHash()
	blockByte := bytes.Join([][]byte{heightByte,
		versionByte,
		prev[:],
		timeByte,
		nonceByte,
		block.GetData()}, []byte{})
	//计算系统的hash
	hash := sha256.Sum256(blockByte)
	return hash
}