package tron

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58/base58"
	"strings"
)

// ComputeAddress 将公钥字节数组转换为 TRON 地址
func ComputeAddress(pubBytes []byte) (string, error) {
	// 如果公钥长度是 65 字节，移除第一个字节
	if len(pubBytes) == 65 {
		pubBytes = pubBytes[1:]
	}
	// 对公钥进行 Keccak256 哈希
	hash := crypto.Keccak256(pubBytes)
	// 转换哈希为十六进制字符串
	hashHex := hex.EncodeToString(hash)
	// 取哈希的第 24 位之后的字符，加上前缀 "41"
	addressHex := AddressPrefix + hashHex[24:]
	// 将十六进制字符串转换为字节数组
	addressBytes, err := hex.DecodeString(addressHex)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(addressBytes), nil
}

// ParseJSON 将 JSON 字符串解析为指定的结构体
func ParseJSON(data []byte, v interface{}) error {
	// 将 JSON 数据解析到 v 指向的结构体
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}
	return nil
}

// Base58ToHex 将 TRON 地址从 base58 转换为十六进制
func Base58ToHex(base58Address string) (string, error) {
	bytes, err := base58.Decode(base58Address)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// PadLeftZero 将十六进制字符串左侧填充零至指定长度
func PadLeftZero(hexStr string, length int) string {
	return strings.Repeat("0", length-len(hexStr)) + hexStr
}
