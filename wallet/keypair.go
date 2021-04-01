package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
)

/**
 *地址所对应的秘钥对（私钥 + 公钥），封装在一个自定义的结构体中
 */
type KeyPair struct {
	Priv *ecdsa.PrivateKey
	Pub []byte
}

/**
 *生成一对秘钥对
 */
func NewKeyPair() (*KeyPair, error) {
	curve := elliptic.P256()
	pri, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	pub := elliptic.Marshal(curve, pri.X, pri.Y)
	keyPair := KeyPair{
		Priv: pri,
		Pub:  pub,
	}
	return &keyPair, nil
}

/**
 *恢复公钥信息
 */
func RecoverPublicKey(curve elliptic.Curve, data []byte) ecdsa.PublicKey {
	x, y := elliptic.Unmarshal(curve, data)
	pub := ecdsa.PublicKey{curve, x, y}
	return pub
}

/**
 *将[]byte格式的签名数据转换为r和s的big.Int类型
 */
func ConverSignature(sign []byte) (r, s *big.Int) {
	rBig := new(big.Int)
	sBig := new(big.Int)
	rBig.SetBytes(sign[:len(sign)/2])
	sBig.SetBytes(sign[len(sign)/2:])

	return rBig, sBig
}