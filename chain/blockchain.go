package chain

import (
	"errors"
	"github.com/bolt-master"
)

const BLOCKS = "blocks"
const LASTHASH = "lasthash"

/**
 *定义区块链结构体，该结构体用于们管理区块
 */
type BlockChain struct {
	//Blocks []Block
	DB *bolt.DB

	LastBlock Block//最新最后的区块
}

func CreateChain(db *bolt.DB) BlockChain {
	return BlockChain{
		DB: db,
	}
}

/**
 *创建一个区块链对象，包含一个创世区块
 */
func (chain *BlockChain) CreateGensis(data []byte) error {
	var err error
	//gensis持久化到db中去
	engine := chain.DB
	engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {//没有桶
			bucket, err = tx.CreateBucket([]byte(BLOCKS))
			if err != nil {
				return err
				//panic("操作区块存储文件失败，请重试！")
			}
		}
		//先查看
		lasthash := bucket.Get([]byte(LASTHASH))
		if len(lasthash) == 0 {
			gensis := CreateGenesis(data)
			genSerBytes, _ := gensis.Serialize()
			//bucket已经存在
			//key -> value
			//blockHash -> 区块序列化以后的数据
			bucket.Put(gensis.Hash[:], genSerBytes)//把创世区块保存到boltdb中去
			//使用一个标志，用来记录最新区块的哈希，以标明当前文件中存储到了最新的哪个区区块
			bucket.Put([]byte(LASTHASH), gensis.Hash[:])
			//把gensis赋值给chain.LastBlock
			chain.LastBlock = gensis
		}else {
			//从文件中读取出最新的区块并赋值给chain中的lastblock
			lastHash := bucket.Get([]byte(LASTHASH))
			lastBlockBytes := bucket.Get(lastHash)
			//把反序列化后的最新的区块赋值给chain中的lastblock
			chain.LastBlock, err = Deserialize(lastBlockBytes)

		}
		return nil
	})
	return err
}

/**
 *生成一个新区块
 */
func (chain *BlockChain) CreateNewBlock(data []byte) error {
	//目的：生成一个新区块，并存到bolt.DB文件中去
	//手段(步骤):
	//a，从文件中查到当前存储的最新区块数据
	lastBlock := chain.LastBlock
	//c，根据获取的最新区块生成一个新区块
	newBlock := NewBlock(lastBlock.Height, lastBlock.Hash, data)
	//d，将最新区块序列化，得到序列化数据
	var err error
	newBlockSerBytes, err := newBlock.Serialize()
	if err != nil {
		return err
	}
	//e，将序列化数据存储到文件，同时更新最新区块的标记lasthash，更新为最新区块的hash
	db := chain.DB
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			err = errors.New("区块数据库操作失败，请重试！")
			return err
		}
		//将新生成的区块保存到文件中去
		bucket.Put(newBlock.Hash[:], newBlockSerBytes)
		//更新最新区块的标记lasthash，更新为最新区块的hash
		bucket.Put([]byte(LASTHASH), newBlock.Hash[:])
		//更新内存中的blockchain的lastblock
		chain.LastBlock = newBlock
		return nil
	})
	return err
}

//获取最新的区块数据
func (chain *BlockChain) GetLastBlock() Block{
	return chain.LastBlock
}

//获取所有区块的数据
func (chain *BlockChain) GetAllBlocks() ([]Block, error) {
	//目的：获取所有的区块
	//手段（步骤）：
	//1，找到最后一个区块
	db := chain.DB
	var err error
	blocks := make([]Block, 0)
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			err = errors.New("区块数据库操作失败，请重试！")
			return err
		}
		var currentHash []byte
		currentHash = bucket.Get([]byte(LASTHASH))

		//2，根据最后一个区块依次往前找
		for  {
			currentBlockBytes := bucket.Get(currentHash)
			currentBlock, err := Deserialize(currentBlockBytes)
			if err != nil {
				break
			}

			//3，每次找到的区块放入到一个切片容器中
			blocks = append(blocks, currentBlock)

			//4，找到最开始的创世区块时，就结束了，不再找了
			if currentBlock.Height == 0 {
				break
			}
			currentHash = currentBlock.PrevHash[:]
		}
		return nil
	})
	return blocks, err
}