package tron

import (
	"encoding/hex"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"math/big"
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

// ParseTRC20TransferData 从 ABI 编码的数据中提取 `to` 地址和 `amount`
func ParseTRC20TransferData(data string) (string, *big.Int) {
	// 提取接收地址（第10-20个字节，16进制每字节2字符，位置20到40）
	toAddressHex := data[32:72]
	toAddress := address.HexToAddress("41" + toAddressHex) // TRON 地址通常以 '41' 开头

	valueHex := data[72:136] // 获取金额
	value := new(big.Int)
	value.SetString(valueHex, 16) // 解析16进制为整数
	return toAddress.String(), value
}
