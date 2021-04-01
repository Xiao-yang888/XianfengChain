package transaction

import (
	"XianfengChain04/utils"
	"XianfengChain04/wallet"
	"bytes"
)

/**
 *定义交易输入的结构体
 */
type TxInput struct {
	TxId      [32]byte  //该字段确定引用自哪笔交易
	Vout      int       //该字段确定引用自该交易的哪笔输出
	//ScritpSig []byte    //该字段表示使用交易输出的证明，解锁脚本
    Sig []byte // 签名
	PubK []byte// 原始公钥
}

/**
 *该函数用于生成一个交易输入案例，即一笔新的花费
 */
func NewTxInput(txid [32]byte, vout int, pubk []byte) TxInput {
	input := TxInput{
		TxId: txid,
		Vout: vout,
		PubK: pubk,
	}
	return input
}

/**
 *验证某个TxInput是否是某个特定address的消费
 */
func (input *TxInput) VertifyInputWithAddress(address string) bool {
	//已知：input消费记录，address某个特定地址
	//目标：判定input是否是address的消费
	//txinput中包含原始公钥pubk字段
	//根据之前的知识：原始公钥pubk与address之间有特定的关系
	//pubk -> 经过n个步骤的变换计算 -> address
	//思路①：pubk变换计算得到一个addr，与address进行比较

	//思路②：pubk变换为pubkHash，address变换为pubkHash1，比较两个hash
	pubk := input.PubK
	hash256 := utils.Hash256(pubk)
	ripemd160 := utils.HashRipemd160(hash256)
	pubkHash := append([]byte{wallet.VERSION}, ripemd160...)
	//变换address为punkhash1
	reAddress := utils.Decode(address)
	rePubkHash := reAddress[:len(reAddress) - 4]
	return bytes.Compare(pubkHash, rePubkHash) == 0
}
