package utxoset

/**
 *定义一笔交易的消费记录
 */
type SpendRecord struct {
	TxId [32]byte
	Vout int
}

/**
 *结构体的返回txid的方法
 */
func (record SpendRecord) GetTxId() [32]byte {
	return record.TxId
}

/**
 *结构体的返回vout的方法
 */
func (record SpendRecord) GetVout() int {
	return record.Vout
}

/**
 *新建一笔消费记录的结构体实例，并将该实例返回
 */
func NewSpendRecord(txid [32]byte, vout int) SpendRecord {
	return SpendRecord{txid, vout}
}
