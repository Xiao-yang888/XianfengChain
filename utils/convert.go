package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"golang.org/x/crypto/ripemd160"
)

/**
 *将int类型的数据转化为[]byte类型
 */
func Int2Byte(num int64) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	return buff.Bytes(), err
}

/**
 *gob编码序列化
 */
func Encoder(v interface{}) ([]byte, error) {
	buff := new(bytes.Buffer)
	encoder := gob.NewEncoder(buff)
	err := encoder.Encode(v)
	return buff.Bytes(), err
}

/**
 *gob反编码
 */
func Decodes(data []byte, v interface{}) (interface{}, error) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
    err := decoder.Decode(v)
    return v, err
}

/**
 *将json格式的数组转化为对应的字符串类型的切片
 */
func JSONArray2String(array string) ([]string, error) {
	var stringSlice []string
	err := json.Unmarshal([]byte(array), &stringSlice)
	return stringSlice, err
}

/**
 *将json格式的数组转化为对应的浮点型数据的切片
 */
func JSONArray2Float(array string) ([]float64, error) {
	var  floatSlice []float64
	err := json.Unmarshal([]byte(array), &floatSlice)
	return floatSlice, err
}

/**
 *sha256哈希计算
 */
func Hash256(data []byte) []byte {
	sha := sha256.New()
	sha.Write(data)
	return sha.Sum(nil)
}

/**
 *Ripemd160的哈希计算
 */
func HashRipemd160(data []byte) []byte {
	hash := ripemd160.New()
	hash.Write(data)
	return hash.Sum(nil)
}