package tron

import (
	"encoding/hex"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"strings"
)

// Base58ToHex 将 TRON 地址从 base58 转换为十六进制
func Base58ToHex(base58Address string) (string, error) {
	bytes, err := common.Decode(base58Address)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// PadLeftZero 将十六进制字符串左侧填充零至指定长度
func PadLeftZero(hexStr string, length int) string {
	return strings.Repeat("0", length-len(hexStr)) + hexStr
}
