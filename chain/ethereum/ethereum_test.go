package ethereum

import (
	"fmt"
	"github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"testing"

	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
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
		PublicKey: "02410c64fcd262512683b54576440e3d3033d825ef9f753b44c51ccdd70a7e90c3",
	})
	if err != nil {
		t.Error("convert address failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.Address)
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
		t.Error("valid address failed:", err)
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
		t.Error("get block header by number failed:", err)
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
		t.Error("get block header by hash failed:", err)
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
		t.Error("get block by number failed:", err)
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
		t.Error("get block by hash failed:", err)
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
		Address: "0x8358d847Fc823097380c4996A3D3485D9D86941f",
	})
	if err != nil {
		t.Error("get account failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp)
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
		t.Error("get account failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp)
}

func TestChainAdaptor_GetTxByAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetTxByAddress(&account.TxAddressRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Address: "0x0B70B578aBd96AAb5e80D24D1f3C28DbdE14356a",
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
		Hash:    "0xa35cdf2fc94521b344ea02591eda0b767be2fc8620bb4ed8adabba49fcc0c89a",
	})
	if err != nil {
		t.Error("get transaction by address failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.Tx)
}

func TestChainAdaptor_GetBlockByRange(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetBlockByRange(&account.BlockByRangeRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Start:   "21118661",
		End:     "21118662",
	})
	if err != nil {
		t.Error("get block by range failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.GetBlockHeader())
}
