package tron

import (
	"encoding/hex"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"math/big"
	"strings"
)

const (
	AddressPrefix = "41"
)

// Base58ToHex Convert TRON address from base58 to hexadecimal
func Base58ToHex(base58Address string) (string, error) {
	bytes, err := common.Decode(base58Address)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// PadLeftZero Fill the left side of the hexadecimal string with zero to the specified length
func PadLeftZero(hexStr string, length int) string {
	return strings.Repeat("0", length-len(hexStr)) + hexStr
}

// ParseTRC20TransferData Extract the 'to' address and 'amount' from ABI encoded data`
func ParseTRC20TransferData(data string) (string, *big.Int) {
	// Extract the receiving address (10-20 bytes, 2 characters per byte in hexadecimal, positions 20 to 40)
	toAddressHex := data[32:72]
	toAddress := address.HexToAddress(AddressPrefix + toAddressHex) // TRON addresses usually start with '41'
	valueHex := data[72:136]                                        // Get amount
	value := new(big.Int)
	value.SetString(valueHex, 16) // Parse hexadecimal to integer
	return toAddress.String(), value
}
