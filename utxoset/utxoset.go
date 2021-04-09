package utxoset

import (
	"XianfengChain04/transaction"
	"XianfengChain04/utils"
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/boltdb/bolt"
)

const UTXOSET = "utxoset"//存放utxoset的桶名

/**
 *UTXO集合，表示用于优化代码结构，实现快速查询
 *map中的key为string类型，表示的是地址
 *map中的value为[]UTXO类型，表示的是某个地址对应的所有的可花费的UTXO
*/
type UTXOSet struct {
	//UTXOs map[string] []transaction.UTXO
	Engine *bolt.DB//bolt.db对象
}

/**
 *构建一个utxoset结构体实例并返回
 */
func NewUTXOSet(db *bolt.DB) UTXOSet {
	utxoset := UTXOSet{
		Engine: db,
	}
	return utxoset
}

/**
 *查询某个地址的可用的utxo的功能
 */
func (utxoset *UTXOSet) QueryUTXOsByAddress(address string) ([]transaction.UTXO, error) {
	//fmt.Println("查询某个地址可用的utxo")
	var utxos []transaction.UTXO
	var err error

	engine := utxoset.Engine
	engine.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(UTXOSET))
		if bucket != nil {
			return nil
		}
		utxosBytes := bucket.Get([]byte(address))
		//未查询到地址对应的utxo数据
		if len(utxosBytes) == 0 {
			return nil
		}
		//将得到的该地址的utxo数据反序列化
		decoder := gob.NewDecoder(bytes.NewReader(utxosBytes))
		err = decoder.Decode(&utxos)
		return err
	})
	return utxos, err
}

/**
 *把某个地址新产生的utxo保存起来
 */
func (utxuset *UTXOSet) AddUTXOsWithAddress(address string, utxos []transaction.UTXO) (bool, error) {
	//fmt.Println("把某个地址新产生的utxo保存起来")
	var err error
	engine := utxuset.Engine

	engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(UTXOSET))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(UTXOSET))
			if err != nil {
				return err
			}
		}

		//先看看utxoset桶中是否有address对应的数据
		var existUTXOs []transaction.UTXO
		utxosBytes := bucket.Get([]byte(address))
		if len(utxosBytes) != 0 {//如果长度不为0，表示之前已存在该地址对应的utxo数据
			decoder := gob.NewDecoder(bytes.NewReader(utxosBytes))
			err = decoder.Decode(&existUTXOs)
			if err != nil {
				return err
			}
		}

		sumUTXOs := make([]transaction.UTXO, 0)
		//把已有的utxos数据进行记录
		if len(existUTXOs) > 0 {
			sumUTXOs = append(sumUTXOs, existUTXOs...)
		}
		//把此次新增的utxo数据追加到bolt.DB中
		sumUTXOs = append(sumUTXOs, utxos...)

		//将utxos进行序列化
		utxosBytes, err := utils.Encoder(sumUTXOs)
		if err != nil {
			return err
		}
		//把序列化后的utxos的数据存到bolt.DB中
		err = bucket.Put([]byte(address), utxosBytes)
		return err
	})
	return err == nil, err
}

/**
 *把某个地址消费了的某些utxo删除掉
 */
func (utxoset *UTXOSet) DeleteUTXOsWithAddress(address string, records []SpendRecord) (bool, error) {
	//fmt.Println("把某个地址消费了的某些utxo删除掉")
	engine := utxoset.Engine
	var err error

	engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(UTXOSET))
		if bucket == nil {
			return nil
		}
		//utxoset桶存在，可以尝试删除操作
		utxosBytes := bucket.Get([]byte(address))
		if len(utxosBytes) == 0 {
			err = errors.New("未查到utxos记录，无法删除")
			return err
		}

		//解码decoder
		var existUTXOs []transaction.UTXO
		decoder := gob.NewDecoder(bytes.NewReader(utxosBytes))
		err = decoder.Decode(&existUTXOs)
		if err != nil {
			return err
		}

		//判断records与找到的utxos的关系：消费的钱必须是现存的utxos的子集
		isSub := IsSubUTXOs(existUTXOs, records)
		if !isSub {
			return errors.New("抱歉，余额不足")
		}

		var isSpent bool
		remainUTXOs := make([]transaction.UTXO, 0)//剩余的utxo容器
		for _, existUTXO := range existUTXOs {
			isSpent = false
			for _, record := range records {
				if existUTXO.EqualSpendRecord(record) {//true表示两个utxo相等
        			//说明existUTXO被花了
        			isSpent = true
				}
			}
			if !isSpent {//existUTXO没被花，可以留下来
				remainUTXOs = append(remainUTXOs, existUTXO)
			}
		}
		remainBytes, err := utils.Encoder(remainUTXOs)
		if err != nil {
			return err
		}
		//把剩余的utxos保存到utxoset桶中
		bucket.Put([]byte(address), remainBytes)
		return nil
	})
	return err == nil, err
}

/**
 *通过交易的消费记录获取对应的utxo
 */
func (utxoset *UTXOSet) GetUTXOsBySpendRecords(address string, records []SpendRecord) ([]transaction.UTXO, error) {
	var err error
	spentUTXOs := make([]transaction.UTXO, 0)

	engine := utxoset.Engine
	engine.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(UTXOSET))
		if bucket != nil {
			return nil
		}
		//桶存在，可以尝试获取utxo数据
		utxosBytes := bucket.Get([]byte(address))
		if len(utxosBytes) == 0 {
			return nil
		}
		var existUTXOs []transaction.UTXO
		decoder := gob.NewDecoder(bytes.NewReader(utxosBytes))
		err = decoder.Decode(&existUTXOs)
		if err != nil {
			return err
		}
		for _, record := range records {
			isContainer := false
			for _, utxo := range existUTXOs {
				if utxo.EqualSpendRecord(record) {//消费记录与当前的utxo匹配上了
					isContainer = true
					spentUTXOs = append(spentUTXOs, utxo)
				}
			}
			if !isContainer {
				err = errors.New("消费了本该不属于你的钱")
				return err
			}
		}
		return nil
	})
	return spentUTXOs ,err
}

/**
 *用于判断消费记录records是否是现存utxo的子集，如果是，返回true，不是子集，返回false
 */
func IsSubUTXOs(utxos []transaction.UTXO, records []SpendRecord) bool {
	for _, record := range records {
		isContainer := false
		for _, utxo := range utxos {
			if utxo.EqualSpendRecord(record) {
				isContainer = true
			}
		}
		if !isContainer {
			return false
		}
	}
	return true
}