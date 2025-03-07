package xlm

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"github.com/ethereum/go-ethereum/log"
	"testing"
)

func setup() (adaptor chain.IChainAdaptor, err error) {
	conf, err := config.New("../../config.yml")
	if err != nil {
		log.Error("load config failed, error:", err)
		return nil, err
	}
	adaptor, err = NewChainAdaptor(conf)
	if err != nil {
		log.Error("create chain adaptor failed, error:", err)
		return nil, err
	}
	return adaptor, nil
}

func TestChainAdaptor_GetSupportChains(t *testing.T) {
	adaptor := ChainAdaptor{}

	req := &account.SupportChainsRequest{
		Chain:   ChainName,
		Network: "mainnet",
	}

	resp, err := adaptor.GetSupportChains(req)

	if err != nil {
		t.Errorf("GetSupportChains failed with error: %v", err)
		return
	}
	fmt.Printf("resp: %s\n", resp)

	if resp.Code != common2.ReturnCode_SUCCESS {
		t.Errorf("Expected success code, got %v", resp.Code)
		return
	}

	if !resp.Support {
		t.Error("Expected Support to be true")
		return
	}
}

func TestChainAdaptor_ValidAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
		return
	}
	rsp, err := adaptor.ValidAddress(&account.ValidAddressRequest{
		Chain:   ChainName,
		Address: "GAZEFFEFXCG2IOM7QBICUKCO3MOL3NVT3GORVBNS7TTETRQCQDYXPOQC",
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	jsExpression, _ := json.MarshalIndent(rsp, "", "    ")
	fmt.Println(string(jsExpression))
}

func TestChainAdaptor_GetAccount(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
		return
	}
	rsp, err := adaptor.GetAccount(&account.AccountRequest{
		Chain:   ChainName,
		Address: "GD6ARGMYT65UUC7FQDBK77GXMEONL44BL7E5G4WL2NDWMJ7NSWBUBYQQ",
	})
	if err != nil {
		t.Error(err)
		return
	}
	jsExpression, _ := json.MarshalIndent(rsp, "", "    ")
	fmt.Println(string(jsExpression))
}

func TestChainAdaptor_GetBlockByNumber(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
		return
	}
	rsp, err := adaptor.GetBlockByNumber(&account.BlockNumberRequest{
		Chain:  ChainName,
		Height: 56013360,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	jsExpression, _ := json.MarshalIndent(rsp, "", "    ")
	fmt.Println(string(jsExpression))
}

func TestChainAdaptor_GetBlockHeaderByNumber(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
		return
	}
	rsp, err := adaptor.GetBlockHeaderByNumber(&account.BlockHeaderNumberRequest{
		Chain:  ChainName,
		Height: 56013360,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	jsExpression, _ := json.MarshalIndent(rsp, "", "    ")
	fmt.Println(string(jsExpression))
}

func TestChainAdaptor_GetFee(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
		return
	}
	rsp, err := adaptor.GetFee(&account.FeeRequest{
		Chain:   ChainName,
		Network: "mainnet",
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	jsExpression, _ := json.MarshalIndent(rsp, "", "    ")
	fmt.Println(string(jsExpression))
}

func TestChainAdaptor_GetTransactionByHash(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
		return
	}

	/*
		// 普通转账（且不携带多种效果，也就是只有native转账，当前仅支持主网的rpc，测试网暂不支持，因为结构有差异）
		b4b8618b6b0782eb4b138ae8ea068af988fa6ca50760311a5ece4c771c7c9457

		// 调用合约（暂不支持吧，Json结构定义，有点麻烦）
		e7780551a2a2371d067b4160f56fde9847abb837ec85168fcf40bb01ebd99db7
	*/

	rsp, err := adaptor.GetTxByHash(&account.TxHashRequest{
		Chain: ChainName,
		Hash:  "47476b985c63ee571505048c179a79226e0968ca35dca0f0c9a58968bddafc6b",
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	jsExpression, _ := json.MarshalIndent(rsp, "", "    ")
	fmt.Println(string(jsExpression))
}

func TestChainAdaptor_CreateUnSignTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
		return
	}

	// 创建请求数据
	addrFrom := "GDYDI34YZQCP7WU726B626KTJIE6COPSXMWL2VLR7KNRHMNB6HLNFUJB"
	addrTo := "GD6ARGMYT65UUC7FQDBK77GXMEONL44BL7E5G4WL2NDWMJ7NSWBUBYQQ"
	sequenceFrom := 239763383908302853
	amount := "0.123"

	requestData := RequestCreateUnsignTransaction{
		AddrFrom:     addrFrom,
		AddrTo:       addrTo,
		SequenceFrom: int64(sequenceFrom),
		Amount:       amount,
	}

	jsonBytes, err := json.Marshal(requestData)
	if err != nil {
		t.Fatal(err)
		return
	}

	base64Encoded := base64.StdEncoding.EncodeToString(jsonBytes)
	fmt.Println("json: ", base64Encoded)
	// eyJhZGRyRnJvbSI6IkdEWURJMzRZWlFDUDdXVTcyNkI2MjZLVEpJRTZDT1BTWE1XTDJWTFI3S05SSE1OQjZITE5GVUpCIiwiYWRkclRvIjoiR0Q2QVJHTVlUNjVVVUM3RlFEQks3N0dYTUVPTkw0NEJMN0U1RzRXTDJORFdNSjdOU1dCVUJZUVEiLCJzZXF1ZW5jZUZyb20iOjIzOTc2MzM4MzkwODMwMjg1MywiYW1vdW50IjoiMC4xMjMifQ==

	rsp, err := adaptor.CreateUnSignTransaction(&account.UnSignTransactionRequest{
		Chain:    ChainName,
		Network:  "mainnet",
		Base64Tx: base64Encoded,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	jsExpression, _ := json.MarshalIndent(rsp, "", "    ")
	fmt.Println(string(jsExpression))

	// AAAAAgAAAADwNG+YzAT/2p/Xg+15U0oJ4Tnyuyy9VXH6mxOxofHW0gAAAGQDU8+HAAAABgAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAABAAAAAPA0b5jMBP/an9eD7XlTSgnhOfK7LL1VcfqbE7Gh8dbSAAAAAQAAAAD8CJmYn7tKC+WAwq/812Ec1fOBX8nTcsvTR2Yn7ZWDQAAAAAAAAAAAABLEsAAAAAAAAAAA
	// 成功后，把“未签名的交易值” 放入 "标准签名机"，返回出"待发送的交易值"
	// 签名机代码，我就不开放了。
}

func TestChainAdaptor_BuildSignedTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
		return
	}

	rsp, err := adaptor.BuildSignedTransaction(&account.SignedTransactionRequest{
		Chain:     ChainName,
		Network:   "mainnet",
		Base64Tx:  "AAAAAgAAAADwNG+YzAT/2p/Xg+15U0oJ4Tnyuyy9VXH6mxOxofHW0gAAAGQDU8+HAAAABQAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAABAAAAAPA0b5jMBP/an9eD7XlTSgnhOfK7LL1VcfqbE7Gh8dbSAAAAAQAAAAD8CJmYn7tKC+WAwq/812Ec1fOBX8nTcsvTR2Yn7ZWDQAAAAAAAAAAAABLEsAAAAAAAAAAA",
		Signature: "oCFwF8LbKGTC3CpMy+W+eMQbJm4zLouG2oJOY0geHGi/LZnYCgemuNL85IySg48frjRF3VRPT5gPtwv5LuvtBA==",
		PublicKey: "GDYDI34YZQCP7WU726B626KTJIE6COPSXMWL2VLR7KNRHMNB6HLNFUJB",
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	jsExpression, _ := json.MarshalIndent(rsp, "", "    ")
	fmt.Println(string(jsExpression))
}

func TestChainAdaptor_SendTx(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
		return
	}

	sendRawValue := "AAAAAgAAAADwNG+YzAT/2p/Xg+15U0oJ4Tnyuyy9VXH6mxOxofHW0gAAAGQDU8+HAAAABQAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAABAAAAAPA0b5jMBP/an9eD7XlTSgnhOfK7LL1VcfqbE7Gh8dbSAAAAAQAAAAD8CJmYn7tKC+WAwq/812Ec1fOBX8nTcsvTR2Yn7ZWDQAAAAAAAAAAAABLEsAAAAAAAAAABofHW0gAAAECgIXAXwtsoZMLcKkzL5b54xBsmbjMui4bagk5jSB4caL8tmdgKB6a40vzkjJKDjx+uNEXdVE9PmA+3C/ku6+0E"

	rsp, err := adaptor.SendTx(&account.SendTxRequest{
		Chain:   ChainName,
		Network: "mainnet",
		RawTx:   sendRawValue,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	jsExpression, _ := json.MarshalIndent(rsp, "", "    ")
	fmt.Println(string(jsExpression))
}
