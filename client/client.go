package client

import (
	"XianfengChain04/chain"
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
	case CREATEBLOCK:
		cmd.CreateBlock()
	case GETLASTBLOCK:
		cmd.GetLastBlock()
	case GETALLBLOKCS:
		cmd.GetAllBlocks()
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
	var gensis string
	generetesis.StringVar(&gensis,"gensis", "", "创世区块的数据")
	generetesis.Parse(os.Args[2:])

	fmt.Println("用户输入的自定义创世区块数据：", gensis)
	blockChain := cmd.Chain
	//1，先判断blockchain中是否已存在创世区块
	hashBig := new(big.Int)
	hashBig.SetBytes(blockChain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 1 {
		fmt.Println("创世区块已存在，不能重复生成")
		return
	}
	//2，调用方法实现创世区块的操作
	blockChain.CreateGensis([]byte(gensis))
	fmt.Println("已成功创建创世区块，并保存到文件中")
}

func (cmd *CmdClient) CreateBlock() {
	createBlock := flag.NewFlagSet(CREATEBLOCK, flag.ExitOnError)
	data := createBlock.String("data", "", "新建区块的自定义数据")

	if len(os.Args[2:]) > 2 {
		fmt.Println("CREATEBLOCK命令只支持data一个参数，请重试")
		return
	}

	//args := os.Args[2:]

	createBlock.Parse(os.Args[2:])
	//先判断是否已生成创世区块，如果没有创术区块则提示用户先创创世区块
	hashBig := new(big.Int)
	hashBig.SetBytes(cmd.Chain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 0 {//没有创世区块
		fmt.Println("That not a gensis block in blockchain, please use go run main.go generategensis command to create a gensis block first.")
		return
	}
	//生成一个新的区块，存到文件中
	blockChain := cmd.Chain
	err := blockChain.CreateNewBlock([]byte(*data))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("新区块创建成功，并成功写入文件")
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
	fmt.Printf("最新区块数据：%s\n", lasBlock.Data)
}

func (cmd *CmdClient) GetAllBlocks() {
	blocks, err := cmd.Chain.GetAllBlocks()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("恭喜，查询到所有区块数据")
	for _, block := range blocks {
		fmt.Printf("区块高度：%d，区块哈希：%x，区块数据：%s\n", block.Height, block.Hash, block.Data)
	}
}

func (cmd *CmdClient) Default() {
	fmt.Println("go run main.go: Unknow subcommand.")
	fmt.Println("Run 'go run main.go help' for usage.")
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
	fmt.Println("    createblock       this command used to create a new block, that can specified a data an argument named data.")
	fmt.Println("    getlastblock      get the lastest block data.")
	fmt.Println("    getallblock       return all blocks data to user.")
	fmt.Println("    help              use the command can print usage infomation.")
	fmt.Println()
	fmt.Println("Use go run main.go help [command] for more information about a command.")
}