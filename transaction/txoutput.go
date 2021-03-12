package transaction

/**
 *定义交易输出的结构体
 */
type TxOutPut struct {
	Value     float64 //转账数量
	ScriptPub []byte  //锁定脚本
}
