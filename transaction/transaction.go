package transaction

import (
	"XianfengChain04/utils"
	"XianfengChain04/wallet"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
)

const REWARSIXE = 50

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

/**
 *该函数用于定义一个coinbase交易，并返回该交易结构体
 */
func CreateCoinBase(addr string) (*Transaction, error) {

	output0 := LockMoney2PubkHash(REWARSIXE, addr)

	coinbase := Transaction{
		Outputs: []TxOutPut{output0},
	}
	coinbaseBytes, err := utils.Encoder(coinbase)
	if err != nil {
		return nil, err
	}
	coinbase.TxHash = sha256.Sum256(coinbaseBytes)

	return &coinbase, nil
}

/**
 *该函数用于构建一笔普通的交易，返回构建好的交易实例
 */
func CreateNewTransaction(utxos []UTXO, from string,pubk []byte, to string, amount float64) (*Transaction, error) {
	//1，构建inputs
	inputs := make([]TxInput, 0)//用于存放交易输入的容器
	var inputAmount float64//该变量用于记录转账发起者一共付了多少钱
    //input -> 交易输入：对某个交易的交易输出UTXO的引用
    for _, utxo := range utxos {
    	//1，根据from获取到对应的原始公钥
    	input := NewTxInput(utxo.TxId, utxo.Vout, pubk)
		inputAmount += utxo.Value
		//把构建好的input存入到交易输入容器中
		inputs = append(inputs, input)
	}

	//2，构建outputs
	outputs := make([]TxOutPut, 0)//用于存放交易输出的容器
	//构建转账接受者的交易输出
	outpus0 := LockMoney2PubkHash(amount, to)

	outputs = append(outputs, outpus0)//把第一个交易输出放入到专门存交易输出的容器中

	//判断是否需要找零，如果需要找零，则需要构建一个新的找零输出

    if inputAmount - amount > 0 {
    	output1 := LockMoney2PubkHash(inputAmount - amount, from)
		outputs = append(outputs, output1)
	}

	//3，构建transaction
	newTransaction := Transaction{
		Inputs:  inputs,
		Outputs: outputs,
	}
	//4，计算transaction的哈希，并赋值
	transactionBytes, err := utils.Encoder(newTransaction)
	newTransaction.TxHash = sha256.Sum256(transactionBytes)
	if err != nil {
		return nil, err
	}
	//5，将构建的transaction实例进行返回
	return &newTransaction, nil
}

/**
 *交易的签名验证方法，该方法返回一个bool值
 *返回true表示签名验证通过，返回false表示签名验证不通过
 */
func (tx *Transaction) VerifyTx(utxos []UTXO) (bool, error) {
	if tx.IsCoinbase() {//如果传入的交易是coinbase交易，不需要验签，直接返回true
		return true, nil
	}

	if len(tx.Inputs) != len(utxos) {
		//交易中的input与所引用的utxo个数不一致，直接返回false
		return false, errors.New("签名验证失败")
	}
	//签名：私钥，原始数据 -> Hash计算 -> hash值
    //验签：公钥，签名数据，原始数据
	//为了使交易对象的结构能够恢复为签名时的状态，需要将当前的交易做一次副本拷贝，并修改其中的字段
	txCopy := tx.CopyTx()
	for index, _ := range txCopy.Inputs {
		//①公钥：Input.PubK 原始公钥
		//②签名数据：Input.Sig 签名字段
		//③原始数据：交易的hash值

		//a，把input中的sign置空
		txCopy.Inputs[index].Sig = nil
		//b，把input中的PubK替换为所引用的utxo的pubkhash
		txCopy.Inputs[index].PubK = utxos[index].PubkHash
		//c，计算改造后的交易的hash
		txHash, err := txCopy.CalculateTxHash()
		if err != nil {
			return false, err
		}
		//调用系统api的签名验证方法，进行验签
		//公钥格式：[]byte -> PublicKey
		pub := wallet.RecoverPublicKey(elliptic.P256(), tx.Inputs[index].PubK)
		//签名格式：[]byte ->r,s *big.Int
		r, s := wallet.ConverSignature(tx.Inputs[index].Sig)
		isVerify := ecdsa.Verify(&pub, txHash, r, s)
		if !isVerify{
			return false, errors.New("验签失败")
		}
	}
	return true, nil
}

/**
 *对交易进行签名
 */
func (tx *Transaction) SignTx(priv *ecdsa.PrivateKey, utxos []UTXO) (error) {
	if tx.IsCoinbase() {//判断传入的交易是否是coinbase交易，是则直接返回
		return nil
	}

	var err error
	if len(tx.Inputs) != len(utxos) {
		err = errors.New("签名失败，请重试")
		return err
	}
	txCopy := tx.CopyTx()
    for i := 0; i < len(txCopy.Inputs); i++ {
    	//遍历得到每一笔消费input
    	//input := tx.Inputs[i]
    	//找到每一笔消费input所对应的花费utxo
    	utxo := utxos[i]
    	//得到utxo对应的锁定脚本
    	pubkHash := utxo.PubkHash
    	//将对应的input中的解锁脚本中的值设置为对应utxo中的锁定脚本
		txCopy.Inputs[i].PubK = pubkHash

    	txHash, err := txCopy.CalculateTxHash()
    	if err != nil {
    		return err
		}
		//真正的API签名调用
		r, s, err := ecdsa.Sign(rand.Reader, priv, txHash)
    	if err != nil {
    		return err
		}
		tx.Inputs[i].Sig = append(r.Bytes(), s.Bytes()...)
		txCopy.Inputs[i].PubK = nil
	}
	return err
}

/**
 *交易的序列化
 */
func (tx *Transaction) Serialize() ([]byte, error) {
	txBytes, err := utils.Encoder(tx)
	if err != nil {
		return nil, err
	}
	return txBytes, nil
}

/**
 *计算交易哈希值
 */
func (tx *Transaction) CalculateTxHash() ([]byte, error) {
	txBytes, err := tx.Serialize()
	if err != nil{
		return nil, err
	}
	return utils.Hash256(txBytes), nil
}

/**
 *拷贝交易对象实例
 */
func (tx Transaction)  CopyTx() Transaction {
	inputs := make([]TxInput, 0)
	for _, input := range tx.Inputs {
        txIn := TxInput{
			TxId: input.TxId,
			Vout: input.Vout,
			Sig:  input.Sig,
			PubK: input.PubK,
		}
		inputs = append(inputs, txIn)
	}

	outputs := make([]TxOutPut, 0)
	for _, output := range tx.Outputs {
		txOut := TxOutPut{
			Value:    output.Value,
			PubkHash: output.PubkHash,
		}
		outputs = append(outputs, txOut)
	}
	hash := tx.TxHash

	return Transaction{
		TxHash:  hash,
		Inputs:  inputs,
		Outputs: outputs,
	}
}

/**
 *该方法用于判断某个具体交易是否是coinbase交易
 *是返回true，不是返回false
 */
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 0 && len(tx.Outputs) == 1
}