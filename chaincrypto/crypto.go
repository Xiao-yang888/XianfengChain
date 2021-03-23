package chaincrypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

/**
 *使用密码学随机生成私钥：椭圆曲线数字签名算法ECDSA
 *elliptic curve digital signature algorihm
 *ECC：elliptic curve crypto
 *DES：data encryption stardand
 */
func NewPriKey(curve elliptic.Curve)(*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(curve, rand.Reader)
}

/**
 *根据私钥获得公钥
 */
func GetPub(curve elliptic.Curve, pri *ecdsa.PrivateKey) []byte {
	return elliptic.Marshal(curve, pri.X, pri.Y)
}