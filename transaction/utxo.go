package transaction

import (
	"XianfengChain04/utils"
	"XianfengChain04/wallet"
	"bytes"
)

/**
 *定义结构体UTXO，表示未花费的交易输出
 */
type UTXO struct {
	TxId [32]byte //该笔收入来自哪个交易
	Vout int  //该笔收入来自哪个交易输出
    TxOutPut  //该笔收入的面额和拥有者
}

/**
 *实例化一个UTXO的结构体实例
 */
func NewUTXO(txid [32]byte, vout int, out TxOutPut) UTXO {
	return UTXO{
		TxId:     txid,
		Vout:     vout,
		TxOutPut: out,
	}
}

/**
 *判定某个utxo是否被某个交易引用进而被消费了
 */
func (utxo *UTXO) IsUTXOSpend(spend TxInput) bool {
	//判断txid是否一致
	equalTxId := bytes.Compare(utxo.TxId[:], spend.TxId[:]) == 0
	//判断索引的下标是q否一致
	equalVout := utxo.Vout == spend.Vout
	//判断消费者是否一致
	//utxo.punkHash：公钥哈希
	//spend.Punk：原始公钥
	hash := utils.Hash256(spend.PubK)
	ripemd160 := utils.HashRipemd160(hash)
	pubkHash := append([]byte{wallet.VERSION}, ripemd160...)
	equalConsumer := bytes.Compare(pubkHash, utxo.PubkHash) == 0
	//三个条件同时满足
	return equalTxId && equalVout && equalConsumer
}