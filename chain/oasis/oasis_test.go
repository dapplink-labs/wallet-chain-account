package oasis

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

func setup() (chain.IChainAdaptor, error) {
	conf, err := config.New("../../config.yml")
	if err != nil {
		log.Error("load config failed, error:", err)
		return nil, err
	}
	adaptor, err := NewChainAdaptor(conf)
	if err != nil {
		log.Error("create chain adaptor failed, error:", err)
		return nil, err
	}
	return adaptor, nil
}

func TestChainAdaptor_GetSupportChains(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := adaptor.GetSupportChains(&account.SupportChainsRequest{
		Chain: ChainName,
	})
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
}

func TestChainAdaptor_ConvertAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}
	resp, err := adaptor.ConvertAddress(&account.ConvertAddressRequest{
		Chain:     ChainName,
		Network:   "mainnet",
		PublicKey: "587c480305e3b3a1b69b95f7f1c2b2390d2e79d7465280808b0915801f269084",
	})
	if err != nil {
		t.Error("convert address failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	assert.Equal(t, "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr", resp.Address)
	fmt.Println(resp.Address)
	respJson, _ := json.Marshal(resp)
	t.Logf("响应: %s", respJson)
}

func TestChainAdaptor_ValidAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.ValidAddress(&account.ValidAddressRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Address: "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr",
	})
	if err != nil {
		t.Error("valid address failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
}

func TestChainAdaptor_GetBlockByNumber(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetBlockByNumber(&account.BlockNumberRequest{
		Chain:  ChainName,
		Height: 23359474,
		//Height: 0,
		ViewTx: true,
	})
	if err != nil {
		t.Error("get block by number failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.Transactions)
}

func TestChainAdaptor_GetAccount(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetAccount(&account.AccountRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Address: "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr",
	})
	if err != nil {
		t.Error("get account failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	respJson, _ := json.Marshal(resp)
	t.Logf("响应: %s", respJson)
}

func TestChainAdaptor_GetTxByAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetTxByAddress(&account.TxAddressRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Address: "oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr",
	})
	if err != nil {
		t.Error("get transaction by address failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.Tx)
}

func TestChainAdaptor_GetTxByHash(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetTxByHash(&account.TxHashRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Hash:    "ad367bf760255afce848769a9cc601bea7fe16b32a808e24a247b82472123dc3",
	})
	if err != nil {
		t.Error("get transaction by address failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.Tx)
}

func TestChainAdaptor_SendTx(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.SendTx(&account.SendTxRequest{
		Chain:   ChainName,
		Network: "mainnet",
		RawTx:   "omlzaWduYXR1cmWiaXNpZ25hdHVyZVhAQj8E0PzyVF47cqC1FPgqC3rKKAg8ALarR2XtH8Irw7BshvMyP8TdjlgNs7hI/BZ7p2WHDo+I4rhZjr5LM39TCGpwdWJsaWNfa2V5WCBYfEgDBeOzobablffxwrI5DS5510ZSgICLCRWAHyaQhHN1bnRydXN0ZWRfcmF3X3ZhbHVlWGCkY2ZlZaJjZ2FzGQYhZmFtb3VudEEAZGJvZHmiYnRvVQDFGO5FeVUsepbIutjehiwlwJMjrmZhbW91bnREAvrwgGVub25jZRBmbWV0aG9kcHN0YWtpbmcuVHJhbnNmZXI=",
	})
	if err != nil {
		t.Error("get block by range failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp)
}
