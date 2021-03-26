package transaction

import (
	"XianfengChain04/utils"
	"bytes"
)

/**
 *定义交易输出的结构体
 */
type TxOutPut struct {
	Value     float64 //转账数量
	//ScriptPub []byte  //锁定脚本
    PubkHash []byte //公钥哈希
}

/**
 *锁定一定数量的钱找到一个交易输出上
 */
func LockMoney2PubkHash(value float64, addr string) TxOutPut {
	//1，得到base58反编码以后的数据
	reAddr := utils.Decode(addr)
	//2，去除校验位，得到公钥hash
	pubkHash := reAddr[:len(reAddr) - 4]
	out := TxOutPut{
		Value:    value,
		PubkHash: pubkHash,
	}
	return out
}

/**
 *该函数用于验证某个交易输出是否属于某个地址
 */
func (output *TxOutPut) CheckPubKHashWithAddress(address string) bool {
	//1，base58反编码
	reAdd := utils.Decode(address)
	//2，截取
	pubKHash := reAdd[:len(reAdd) - 4]
	//3，比较给定address的公钥哈希与output的公钥哈希是否相等
	return bytes.Compare(pubKHash, output.PubkHash) == 0
}
