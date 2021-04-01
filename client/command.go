package client

const (
	GENERATEGENSIS = "generategensis"
    SENDTRANSACTION = "sendtransaction"
    GETBALANCE = "getbalance"
    GETLASTBLOCK = "getlastblock"
    GETALLBLOKCS = "getallblocks"
    GETNEWADDRESS = "getnewaddress"//生成新的比特币地址
    LISTADDRESS = "listaddress"//列出所有目前已经生成并管理的地址
    DUMPPRIVKEY = "dumpprivkey"//导出某个地址的私钥
    SETCOINBASE = "setcoinbase"//设置挖矿矿工的地址
    GETCOINBASE = "getcoinbase"//查看当前节点所设置的矿工地址
    HELP = "help"
)

