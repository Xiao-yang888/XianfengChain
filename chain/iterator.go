package chain

/**
 *定义迭代器的接口标准，该迭代器通过总结分析，迭代有两个功能
         ①判断容器中是否还有数据
         ②从容器中取出一个数据
 */
type Iterator interface {
	HasNext() bool//判断容器中是否还有数据
	Next()    Block//如果容器中有数据，则取出包含的一个数据区块
}


