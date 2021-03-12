package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
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
func Decode(data []byte, v interface{}) (interface{}, error) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
    err := decoder.Decode(v)
    return v, err
}
