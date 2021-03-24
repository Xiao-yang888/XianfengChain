package chaincrypto

import (
	"XianfengChain04/utils"
	"bytes"
	"crypto/elliptic"
)

/**
 *新生成一个比特币的地址
 */
func NewAddress() (string, error) {
	curve := elliptic.P256()
	//1，生成私钥
	pri, err := NewPriKey(curve)
	if err != nil {
		return "", err
	}
	//2，获得对应公钥
	pub := GetPub(curve, pri)
	//3，对公钥进行sha256哈希
	pubHash := utils.Hash256(pub)
	//4，ripemd160计算
	ripemePub := utils.HashRipemd160(pubHash)
	//5，添加版本号
	versionPub := append([]byte{0x00}, ripemePub...)
	//6，双hash
	firstHash := utils.Hash256(versionPub)
	secondHash := utils.Hash256(firstHash)
	//7，截取前四个字节作为地址校验位
	check  := secondHash[:4]
	//8，拼接到versionPub的后面
	originAddress := append(versionPub, check...)
	//9，base58编码
	return base58.Encode(originAddress), nil
}

/**
 *该函数用于检查地址是否合法，如果符合地址规范，返回true
 *如果不符合，返回false
 */
func CheckAddress(addr string) bool {
	//1，使用base58对传入的地址进行解码
	reAddrBytes := base58.Decode(addr)//versionPub + check
	//2，取出校验位
	if len(reAddrBytes) < 4 {
		return false
	}
	reCheck := reAddrBytes[len(reAddrBytes) - 4:]
	//3，截取得到versionPubHash
	reVersionPubHash := reAddrBytes[:len(reAddrBytes) - 4]
	//4，对versionPub进行双哈希
	reFirstHash := utils.Hash256(reVersionPubHash)
	reSecondHash := utils.Hash256(reFirstHash)
	//5，对双哈希以后的内容进行前四个字节的截取
	check := reSecondHash[:4]
	return bytes.Compare(reCheck, check) == 0
}
