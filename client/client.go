package client

import (
	"XianfengChain04/chain"
	"XianfengChain04/utils"
	"flag"
	"fmt"
	"math/big"
	"os"
)

/**
 *该结构体定义了用于实现命令行参数解析的结构体
 */
type CmdClient struct {
    Chain chain.BlockChain
}

/**
 *该方法用于获取当前节点已经生成的地址列表
 */
func (cmd *CmdClient) ListAddress() {
	listAddress := flag.NewFlagSet(LISTADDRESS, flag.ExitOnError)
	listAddress.Parse(os.Args[2:])
    if len(os.Args[2:])	> 0 {
    	fmt.Println("无法解析参数，请检查后重试")
		return
	}
	addList, err := cmd.Chain.GetAddressList()
    if err != nil {
    	fmt.Println(err.Error())
		return
	}
	//如果本地节点暂时还没有地址，需要给用户提示
	if len(addList) == 0 {
		fmt.Println("暂无地址，可以使用go run main.go getnewaddress命令生成新地址")
		return
	}
	fmt.Println("获取地址列表成功，地址信息如下：")
	for index, add := range addList {
		fmt.Printf("[%d]:%s\n", index + 1, add)
	}
}

/**
 *定义新的方法，用于生成新的地址
 */
func (cmd *CmdClient) GetNewAddress() {
	getNewAddress := flag.NewFlagSet(GETNEWADDRESS, flag.ExitOnError)
	getNewAddress.Parse(os.Args[2:])
	if len(os.Args[2:]) > 0 {
		fmt.Println("抱歉，生成新地址功能无法解析参数，请重试")
		return
	}
	address, err := cmd.Chain.GetNewAddress()
	if err != nil {
		fmt.Println("生成新地址时遇到错误，请重试", err.Error())
		return
	}
	fmt.Println("生成新的地址：", address)
}

/**
 *该方法用于导出某个特定地址的私钥信息
 */
func (cmd *CmdClient) DumpPrivkey() {
	dumpPrivkey := flag.NewFlagSet(DUMPPRIVKEY, flag.ExitOnError)
	address := dumpPrivkey.String("address", "", "要导出的私钥的地址")
	dumpPrivkey.Parse(os.Args[2:])
	if len(os.Args[2:]) > 2 {
		fmt.Println("无法解析输入参数，请检查后重试")
		return
	}
	pri, err := cmd.Chain.DumpPrivkey(*address)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("私钥是：%x\n", pri.D.Bytes())
}

/**
 *client运行方法
 */
func (cmd *CmdClient) Run() {
	//处理用户没有输入任何命令参数的情况，打印输出说明书
	args := os.Args
	if len(args) == 1 {
		cmd.Help()
		return
	}

	//解析用户输入的第一个参数，作为功能命令进行解析
	switch os.Args[1] {
	case GENERATEGENSIS:
		cmd.GenerateGensis()
	case SENDTRANSACTION:
		cmd.SendTransaction()
	case GETBALANCE:
		cmd.GetBalance()
	case GETLASTBLOCK:
		cmd.GetLastBlock()
	case GETALLBLOKCS:
		cmd.GetAllBlocks()
	case GETNEWADDRESS:
		cmd.GetNewAddress()
	case LISTADDRESS:
		cmd.ListAddress()//获取所有的地址信息列表
	case DUMPPRIVKEY:
		cmd.DumpPrivkey()
	case SETCOINBASE:
		cmd.SetCoinbase()//设置挖矿矿工的地址
	case GETCOINBASE:
		cmd.GetCoinbase()//查看当前节点所设置的矿工地址
	case HELP:
		cmd.Help()
	default:
		cmd.Default()
	}
	//createBlock := flag.NewFlagSet("createblock", flag.ExitOnError)
	//data := createBlock.String("data", "默认值", "新区块的区块数据")
	//createBlock.Parse(os.Args[2:])
	//cmd.Chain.CreateNewBlock([]byte(*data))

}

func (cmd *CmdClient) GenerateGensis() {
	//命令参数集合
	generetesis := flag.NewFlagSet(GENERATEGENSIS, flag.ExitOnError)
	//解析参数
	var addr string
	generetesis.StringVar(&addr,"address", "", "用户指定的矿工的地址")
	generetesis.Parse(os.Args[2:])

	fmt.Println("用户输入的自定义创世区块数据：", addr)
	blockChain := cmd.Chain
	//1，先判断blockchain中是否已存在创世区块
	hashBig := new(big.Int)
	hashBig.SetBytes(blockChain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 1 {
		fmt.Println("创世区块已存在，不能重复生成")
		return
	}

    err := blockChain.CreateCoinBase(addr)
    if err != nil {
    	fmt.Println("抱歉，创建coinbase交易遇到错误：", err.Error())
		return
	}
	fmt.Println("恭喜！生成了一笔COINBASE交易，奖励已到账")
}

/**
 *用户发起交易
 */
func (cmd *CmdClient) SendTransaction() {
	createBlock := flag.NewFlagSet(SENDTRANSACTION, flag.ExitOnError)
	from := createBlock.String("from", "", "交易发起人")
    to := createBlock.String("to", "", "交易接受者地址")
    amount := createBlock.String("amount", "", "转账的数量")

    if len(os.Args[2:]) > 6 {
		fmt.Println("SENDTRANSACTION命令只支持三个参数和参数值，请重试")
		return
	}

	createBlock.Parse(os.Args[2:])

    fromSlice, err := utils.JSONArray2String(*from)
    if err != nil {
    	fmt.Println("抱歉，参数格式不正确，清检查后重试！")
		return
    }

    toSlice, err := utils.JSONArray2String(*to)
    if err != nil {
		fmt.Println("抱歉，参数格式不正确，清检查后重试！")
    	return
	}

    amountSlice, err := utils.JSONArray2Float(*amount)
    if err != nil {
		fmt.Println("抱歉，参数格式不正确，清检查后重试！")
		return
	}


    //先看看参数个数是否一样
    fromLen := len(fromSlice)
    toLen := len(toSlice)
    amountLen := len(amountSlice)
    if fromLen != toLen || fromLen != amountLen || toLen != amountLen {
    	fmt.Println("参数个数不一致，请检查参数后重试")
		return
	}

	//先判断是否已生成创世区块，如果没有创术区块则提示用户先创创世区块
	hashBig := new(big.Int)
	hashBig.SetBytes(cmd.Chain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 0 {//没有创世区块
		fmt.Println("That not a gensis block in blockchain, please use go run main.go generategensis command to create a gensis block first.")
		return
	}

	err = cmd.Chain.SendTransaction(fromSlice, toSlice, amountSlice)
	if err != nil {
		fmt.Println("抱歉，发送交易出现错误:", err.Error())
		return
	}
	fmt.Println("交易发送成功")
}

/**
 *获取地址的余额
 */
func (Cmd *CmdClient) GetBalance() {
	getbalance := flag.NewFlagSet(GETBALANCE, flag.ExitOnError)
	var addr string
	getbalance.StringVar(&addr, "address", "", "用户的地址")
	getbalance.Parse(os.Args[2:])

	blockChain := Cmd.Chain
	//先判断是否有创世区块
    hashBig := new(big.Int)
	hashBig.SetBytes(blockChain.LastBlock.Hash[:])
    if hashBig.Cmp(big.NewInt(0)) == 0 {//没有创世区块
    	fmt.Println("抱歉，该网络链暂未存在，无法查询")
		return
	}

	//调用余额查询功能
	balance, err := blockChain.GetBalance(addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
    fmt.Printf("地址%s的余额是：%f\n", addr, balance)
}

func (cmd *CmdClient) GetLastBlock() {
	blockChain := cmd.Chain
	lasBlock := blockChain.GetLastBlock()
	//判断是否为空
	hashBig := new(big.Int)
	hashBig.SetBytes(lasBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 0 {
		fmt.Println("抱歉，当前暂无最新区块")
		return
	}
	fmt.Println("恭喜获取到最新区块")
	fmt.Printf("区块高度：%d\n", lasBlock.Height)
	fmt.Printf("区块哈希：%f\n", lasBlock.Hash)
	for index, tx := range lasBlock.Transactions {
		fmt.Printf("最新区块交易：%d, 交易：%v\n", index, tx)
	}

}

func (cmd *CmdClient) GetAllBlocks() {
	blocks, err := cmd.Chain.GetAllBlocks()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("恭喜，查询到所有区块数据")
	for _, block := range blocks {
		fmt.Printf("区块高度：%d，区块哈希：%x\n", block.Height, block.Hash)
		fmt.Print("区块中国的交易信息：\n")
		for index, tx := range block.Transactions {
			fmt.Printf("     第%d笔交易，交易hash：%x\n", index, tx.TxHash)
		    for inputIndex, input := range tx.Inputs {
		    	fmt.Printf("           第%d笔交易输入,花了%x的%d的钱\n", inputIndex, input.TxId, input.Vout)
			}
			for outputIndex, output := range tx.Outputs {
				fmt.Printf("      第%d笔交易输出，实现收入%f\n", outputIndex, output.Value)
			}
		}
		fmt.Println()
	}
}

func (cmd *CmdClient) Default() {
	fmt.Println("go run main.go: Unknow subcommand.")
	fmt.Println("Run 'go run main.go help' for usage.")
}

/**
 *设置矿工地址功能
 */
func (cmd *CmdClient) SetCoinbase() {
	setCoinbase := flag.NewFlagSet(SETCOINBASE, flag.ExitOnError)
	address := setCoinbase.String("address", "", "用户自定义设置的矿工地址")
    setCoinbase.Parse(os.Args[2:])
	if len(os.Args[2:]) > 2 {
		fmt.Println("参数无法解析，请重试")
		return
	}
	cmd.Chain.Setcoinbase(*address)
}

/**
 *获取当前节点的coinbase矿工地址
 */
func (cmd *CmdClient) GetCoinbase() {
	getCoinbase := flag.NewFlagSet(GETCOINBASE, flag.ExitOnError)
	getCoinbase.Parse(os.Args[2:])
	if len(os.Args[2:]) > 0 {
		fmt.Println("无法解析参数，请重试")
		return
	}
	miner := cmd.Chain.GetCoinbase()
	if len(miner) == 0 {
		fmt.Println("暂未设置coinbase矿工地址")
		return
	}
	fmt.Println("coinbase矿工地址：", miner)
}

/**
 *该方法用于打印输出项目的使用和说明信息，相当于项目的帮助文档和说明书
 */
func (cmd *CmdClient) Help() {
	fmt.Println("-------Welcome to XianfengCHAIN04 project-------")
	fmt.Println("XianfengChain04 is a custom blockchain project, the project plan to bulid a very simple public chain")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("go run main.go command [arguments]")
	fmt.Println()
	fmt.Println("AVAILABLE COMMANDS")
	fmt.Println("    generategensis    use the command can create a gensis block and save to the boltdb file. use the gensis argument to set the custom data.")
	fmt.Println("    sendtransaction   this command used to send a new transaction, that can specified a data an argument named data.")
	fmt.Println("    getbalance        this is a comand that can get the balance of specified address.")
	fmt.Println("    getlastblock      get the lastest block data.")
	fmt.Println("    getallblock       return all blocks data to user.")
	fmt.Println("    getnewaddress     this command use to create a new address by bition algorithm.")
	fmt.Println("    help              use the command can print usage infomation.")
	fmt.Println()
	fmt.Println("Use go run main.go help [command] for more information about a command.")
}