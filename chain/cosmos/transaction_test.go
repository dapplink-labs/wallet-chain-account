package cosmos

import (
	"encoding/hex"
	"fmt"
	"testing"
)

// success
func TestCosmos_BuildUnSignTransaction(t *testing.T) {
	txStruct := &TxStructure{
		ChainId:       "cosmoshub-4",
		FromAddress:   "",
		ToAddress:     "",
		Amount:        10000,
		GasLimit:      137674,
		FeeAmount:     10000,
		Sequence:      5,
		AccountNumber: 3014650,
		Decimal:       6,
		Memo:          "10087",
		PubKey:        "",
	}
	txBytes, _ := BuildUnSignTransaction(txStruct)
	fmt.Printf("txBytes=%X\n", txBytes)
}

func TestCosmos_CreateSignTransaction(t *testing.T) {
	txStruct := &TxStructure{
		ChainId:       "cosmoshub-4",
		FromAddress:   "",
		ToAddress:     "",
		Amount:        10000,
		GasLimit:      137674,
		FeeAmount:     10000,
		Sequence:      5,
		AccountNumber: 3014650,
		Decimal:       6,
		Memo:          "10087",
		PubKey:        "",
	}
	signature := "4775e0b23e23c1cb44c75ce0f687e462408107a1a2365286857f750c1e475ead756bde5db2e0656235826c12eae2df1b991384f135295b992672c28d67fe176700"
	signBytes, _ := hex.DecodeString(signature)
	txBytes, _ := BuildSignTransaction(txStruct, signBytes)
	fmt.Printf("txBytes=%X\n", txBytes)
}
