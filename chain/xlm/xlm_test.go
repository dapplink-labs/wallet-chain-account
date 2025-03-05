package xlm

import (
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

func TestChainAdaptor_SendTx(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
		return
	}

	// 就不构建交易了，测试项目构建一笔，拿到这里来
	sendRawValue := "AAAAAgAAAAAyQpSFuI2kOZ+AUCooTtscvbaz2Z0ahbL85knGAoDxdwAAAGQAEyKQAAAAAwAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAABAAAAADJClIW4jaQ5n4BQKihO2xy9trPZnRqFsvzmScYCgPF3AAAAAQAAAAAMEVFY7uijL1jbeT6If7G72cl2DOevB6xkbATCD/L4cAAAAAAAAAAAByfRUAAAAAAAAAABAoDxdwAAAEDzeLuHFx4p3MFhNh/YbD5qwfws/BsfLz9kBr7zMYr8h+TeV6kTXrd7rxKcxK8D/JxUapOtjKPdPLQOYJQ3PzQM"

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
