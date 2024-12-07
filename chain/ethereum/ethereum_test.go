package ethereum

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

func TestChainAdaptor_ConvertAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	// test account
	// privateKey: 861ae7df240f80e5492065dafeb6444bdbf2d55d01e1797d2abe0db0afd4f917
	// publicKey: 02410c64fcd262512683b54576440e3d3033d825ef9f753b44c51ccdd70a7e90c3
	resp, err := adaptor.ConvertAddress(&account.ConvertAddressRequest{
		Chain:     ChainName,
		Network:   "mainnet",
		PublicKey: "048318535b54105d4a7aae60c08fc45f9687181b4fdfc625bd1a753fa7397fed753547f11ca8696646f2f3acb08e31016afac23e630c5d11f59f61fef57b0d2aa5",
	})
	if err != nil {
		log.Error("convert address failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
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
		Address: "0x8358d847Fc823097380c4996A3D3485D9D86941f",
	})
	if err != nil {
		log.Error("valid address failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
}

func TestChainAdaptor_GetBlockHeaderByNumber(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetBlockHeaderByNumber(&account.BlockHeaderNumberRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Height:  21118661,
	})
	if err != nil {
		log.Error("get block header by number failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.BlockHeader)
}

func TestChainAdaptor_GetBlockHeaderByHash(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetBlockHeaderByHash(&account.BlockHeaderHashRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Hash:    "0x17933cce37211452df901718afd30e4fe013b67c0d262dffd2eb5a3f1b091431",
	})
	if err != nil {
		log.Error("get block header by hash failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.BlockHeader)
}

func TestChainAdaptor_GetBlockByNumber(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetBlockByNumber(&account.BlockNumberRequest{
		Chain:  ChainName,
		Height: 21118661,
		ViewTx: true,
	})
	if err != nil {
		log.Error("get block by number failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.Transactions)
}

func TestChainAdaptor_GetBlockByHash(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetBlockByHash(&account.BlockHashRequest{
		Chain:  ChainName,
		Hash:   "0x17933cce37211452df901718afd30e4fe013b67c0d262dffd2eb5a3f1b091431",
		ViewTx: true,
	})
	if err != nil {
		log.Error("get block by hash failed:", err)
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
		Chain:           ChainName,
		Network:         "mainnet",
		Address:         "0xD79053a14BC465d9C1434d4A4fAbdeA7b6a2A94b",
		ContractAddress: "0x00",
	})
	if err != nil {
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	respJson, _ := json.Marshal(resp)
	t.Logf("响应: %s", respJson)
}

func TestChainAdaptor_GetFee(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetFee(&account.FeeRequest{
		Chain:   ChainName,
		Network: "mainnet",
	})
	if err != nil {
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	respJson, _ := json.Marshal(resp)
	t.Logf("响应: %s", respJson)
}
