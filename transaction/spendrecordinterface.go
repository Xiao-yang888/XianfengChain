package transaction

/**
 *消费记录的接口标准，提供两个方法，分别返回消费记录的txid和vout
 */
type SpendRecordInterface interface {
	GetTxId() [32]byte
	GetVout() int
}