package chain

import (
	"XianfengChain04/transaction"
	"XianfengChain04/wallet"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/bolt-master"
	"math/big"
)

const BLOCKS = "blocks"
const LASTHASH = "lasthash"

/**
 *定义区块链结构体，该结构体用于们管理区块
 */
type BlockChain struct {
	//Blocks []Block
	DB                 *bolt.DB
	LastBlock          Block//最新最后的区块
	IteratorBlockHash  [32]byte//表示当前迭代到了哪个区块，该变量用于记录迭代到的区块哈希
    Wallet             wallet.Wallet//引入wallet字段作为 blockchain的属性
}

func CreateChain(db *bolt.DB) (*BlockChain, error) {
	var lastBlock Block
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			bucket, _ = tx.CreateBucket([]byte(BLOCKS))
		}
		lastHash := bucket.Get([]byte(LASTHASH))
		if len(lastHash) <=  0 {
			return nil
		}
		lastBlockBytes := bucket.Get(lastHash)
		lastBlock, _ = Deserialize(lastBlockBytes)
		return nil
	})

	//创建或者加载wallet结构体对象
    wallet, err := wallet.LoadAddrAndKeyPairsFromDB(db)
    if err != nil {
    	return nil, err
	}

	blockChain := BlockChain{
		DB:                db,
		LastBlock:         lastBlock,
		IteratorBlockHash: lastBlock.Hash,
		Wallet:            *wallet,
	}
	return &blockChain, nil
}

/**
 *创建coinbase交易的方法
 */
func (chain *BlockChain) CreateCoinBase(addr string) error {
	//1，对用户传入的addr进行有效性检查
    isAddrValid := chain.Wallet.CheckAddress(addr)
    if !isAddrValid{
    	return errors.New("抱歉，地址不符合规范，请检查后重试")
	}

	//2，创建一笔coinbase交易
	coinbase, err := transaction.CreateCoinBase(addr)
	if err != nil {
		return err
	}
	//3，把coinbase交易存到区块中
	err = chain.CreateGensis([]transaction.Transaction{*coinbase})
	return err
}

/**
 *创建一个区块链对象，包含一个创世区块
 */
func (chain *BlockChain) CreateGensis(txs []transaction.Transaction) error {
	hashBig := new(big.Int)
	hashBig.SetBytes(chain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 1 {
		//最新区块哈希有值，则说明区块文件当中创世区块已存在
		return nil
	}

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
			gensis := CreateGenesis(txs)
			genSerBytes, _ := gensis.Serialize()
			//bucket已经存在
			//key -> value
			//blockHash -> 区块序列化以后的数据
			bucket.Put(gensis.Hash[:], genSerBytes)//把创世区块保存到boltdb中去
			//使用一个标志，用来记录最新区块的哈希，以标明当前文件中存储到了最新的哪个区区块
			bucket.Put([]byte(LASTHASH), gensis.Hash[:])
			//把gensis赋值给chain.LastBlock
			chain.LastBlock = gensis
			chain.IteratorBlockHash = gensis.Hash
			//fmt.Println("已成功创建创世区块，并写入文件中")
		}
		return nil
	})
	return err
}

/**
 *生成一个新区块
 */
func (chain *BlockChain) CreateNewBlock(txs []transaction.Transaction) error {
	//目的：生成一个新区块，并存到bolt.DB文件中去
	//手段(步骤):
	//a，从文件中查到当前存储的最新区块数据
	lastBlock := chain.LastBlock
	//c，根据获取的最新区块生成一个新区块
	newBlock := NewBlock(lastBlock.Height, lastBlock.Hash, txs)
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
		chain.IteratorBlockHash = newBlock.Hash
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

/**
 *该方法用于实现迭代器Iterator的HashNext方法，用于判断是否还有数据
 *如果有数据，返回true，否则返回false
 */
func (chain *BlockChain) HasNext() bool {
	db := chain.DB
	var  hasNext bool
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return errors.New("区块数据文件操作失败，请重试！")
		}
		//获取前一个区块的区块数据
		prevBlockBytes := bucket.Get(chain.IteratorBlockHash[:])
		//如果获取不到前一个区块的区块数据，说明前面没有区块了
		hasNext = len(prevBlockBytes) != 0
		return nil
	})
	return hasNext
}

/**
 *该方法用于实现迭代器Iterator的Next方法，用于取出一个数据
 *此处，因为是区块数据集合，因此返回的数据类型是Block
 */
func (chain *BlockChain) Next() Block {
	db := chain.DB
	var iteratorBlock Block
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return errors.New("区块数据文件操作失败，请重试！")
		}
		//前一个区块的数据
		blockBytes := bucket.Get(chain.IteratorBlockHash[:])
		iteratorBlock, _ = Deserialize(blockBytes)
		//迭代到当前区块后，更新游标的区块内容
		chain.IteratorBlockHash = iteratorBlock.PrevHash
		return nil
	})
	return iteratorBlock
}

/**
 *该方法用于查询出指定地址的UTXO集合并返回
 */
func (chain *BlockChain) SerchDBUTXOs(addr string) ([]transaction.UTXO) {
	//花费记录的容器
	spend := make([]transaction.TxInput, 0)
	//收入记录的容器
	inCome := make([]transaction.UTXO, 0)
	//迭代遍历每一个区块
	for chain.HasNext() {
		block := chain.Next()
		//遍历区块中的交易
		for _, tx := range block.Transactions {
			//a、遍历每个交易的交易输入
			for _, input := range tx.Inputs {
				//找到了花费记录
				if !input.VertifyInputWithAddress(addr) {
					continue
				}
				spend = append(spend, input)
			}
			//b、遍历每个交易的交易输出:收入
			for index, output := range tx.Outputs {
				if !output.CheckPubKHashWithAddress(addr) {
					continue
				}
				utxo := transaction.NewUTXO(tx.TxHash, index, output)
				inCome = append(inCome, utxo)
			}
		}
	}
	//遍历收入集合和花费集合，把已花的剔除，找出为花费的记录
	utxos := make([]transaction.UTXO, 0)
	var isComeSpent bool
	for _, come := range inCome {
		isComeSpent = false
		for _, spen := range spend {
			if come.TxId == spen.TxId && come.Vout == spen.Vout {
				//说明该笔收入已被花费
				isComeSpent = true
				break
			}
		}
		if !isComeSpent {//当前遍历到的come没有被消费
			utxos = append(utxos, come)
		}
	}
	return utxos
}

/**
 *该方法用于实现地址余额的统计
 */
func (chain *BlockChain) GetBalance(addr string) (float64, error) {
	//1，检查地址的合法性
	isAddrValid := chain.Wallet.CheckAddress(addr)
	if !isAddrValid {
		return 0, errors.New("地址不符合规范，请检查后重试")
	}

	//2，获取地址的余额
	_, totaBalance := chain.GetUTXOsWithBalance(addr, []transaction.Transaction{})
	return totaBalance, nil
}

/**
 *该方法用于实现地址余额统计和地址所可以花费的utxo集合
 */
func (chain BlockChain) GetUTXOsWithBalance(addr string, txs []transaction.Transaction) ([]transaction.UTXO, float64) {
	dbUtxos := chain.SerchDBUTXOs(addr)
	//看一看是否已经花了某个bolt.db文件中的utxo，如果某个utxo被花掉了，应该剔除掉
	memSpends := make([]transaction.TxInput, 0)
	memInComes := make([]transaction.UTXO, 0)
	for _, tx := range txs {
		//a，遍历交易输入，把花的钱记录下来
		for _, input := range tx.Inputs {
			if input.VertifyInputWithAddress(addr) {
				memSpends = append(memSpends, input)
			}
		}
		//b，遍历交易输出，把收入的钱记录下来
		for outIndex, output := range tx.Outputs {
			if output.CheckPubKHashWithAddress(addr) {
				utxo := transaction.NewUTXO(tx.TxHash, outIndex, output)
				memInComes = append(memInComes, utxo)
			}
		}
	}

	utxos := make([]transaction.UTXO, 0)
	var isUTXOSpend bool
	for _, utxo := range dbUtxos {
		isUTXOSpend = false
		for _, spend := range memSpends {
			if utxo.IsUTXOSpend(spend) {
				isUTXOSpend = true
			}
		}
		if !isUTXOSpend {
			utxos = append(utxos, utxo)
		}
	}
	//把内存中的收入也加入到
	utxos = append(utxos, memInComes...)

	var totaBalance float64
	for _, utxo := range utxos{
		fmt.Println(utxo)
		totaBalance += utxo.Value
	}

	return utxos, totaBalance
}

/**
 *定义区块链的发送交易的功能
 */
func (chain *BlockChain) SendTransaction(froms []string, tos []string, amounts []float64) (error) {
	var err error
	//对所有的from和to进行合法性检查
	for i := 0; i < len(froms); i ++ {
		isFromValid := chain.Wallet.CheckAddress(froms[i])
		isToValid := chain.Wallet.CheckAddress(tos[i])
		if !isFromValid || !isToValid {
			return errors.New("地址不符合规范，请检查后重试")
		}
	}

	newTxs := make([]transaction.Transaction, 0)
	//遍历
	for from_index, from := range froms {
		utxos, totaBalance := chain.GetUTXOsWithBalance(from,newTxs)
		if totaBalance < amounts[from_index] {
			return errors.New(from + "余额不足，赶紧去搬砖挣钱")
		}
		totaBalance = 0
		var utxoNum int
		for index, utxo := range utxos {
			totaBalance += utxo.Value
			if totaBalance > amounts[from_index] {
				utxoNum = index
				break
			}
		}

		//获取from的原始公钥
		keyPair := chain.Wallet.Address[from]
		if len(keyPair.Pub) == 0 {
			err = errors.New("构建交易出现错误")
			break
		}
		//可花费的钱总额比要花费的钱数额大，才构建交易
		newTx, err := transaction.CreateNewTransaction(
			utxos[0:utxoNum +1],
			from,
			keyPair.Pub,
		tos[from_index],
			amounts[from_index])
		if err != nil {
			return err
		}
		//对构建的交易newTx进行签名
		newTxs = append(newTxs, *newTx)
	}
	err = chain.CreateNewBlock(newTxs)
	if err != nil {
		return err
	}
	return nil
}

/**
 *生成比特币地址的功能
 */
func (chain *BlockChain) GetNewAddress() (string, error) {
	return chain.Wallet.NewAddress()
}

/**
 *获取钱包中的地址列表
 */
func (chain *BlockChain) GetAddressList() ([]string, error) {
	if chain.Wallet.Address == nil {
		return nil, errors.New("暂无地址")
	}

	addList := make([]string, 0)
	for add, _ := range chain.Wallet.Address {
		addList = append(addList, add)
	}
	return addList, nil
}

/**
 *该方法用于导出某个特定地址的私钥
 */
func (chain *BlockChain) DumpPrivkey(addr string) (*ecdsa.PrivateKey, error) {
    //1，地址规范性检查
	isAddrValid := chain.Wallet.CheckAddress(addr)
	if !isAddrValid {
		return nil, errors.New("地址不符合规范，请重试")
	}
	//2，钱包为空
	if chain.Wallet.Address == nil {
    	return nil, errors.New("当前钱包为找到对应地址的私钥")
	}
	//3，到wallet中找addr的对应的keypair
	keyPair := chain.Wallet.Address[addr]
	if keyPair == nil {
		return nil, errors.New("当前钱包未找到对应的地址的私钥")
	}
	//4，找到了具体结果，将私钥返回
	return keyPair.Priv, nil
}