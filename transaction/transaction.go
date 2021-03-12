package transaction

/**
 *定义交易的结构体
 */
type Transaction struct {
	//交易哈希
	TxHash [32]byte
	//交易输入
	Inputs  []TxInput
	//交易输出
	Outputs []TxOutPut
}
