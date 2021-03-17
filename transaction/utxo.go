package transaction

/**
 *
 */
type UTXO struct {
	TxId [32]byte //该笔收入来自哪个交易
	Vout int  //该笔收入来自哪个交易输出
    TxOutPut  //该笔收入的面额和拥有者
}
