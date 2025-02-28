package near

import (
	"encoding/json"
	"fmt"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/test-go/testify/assert"
	"testing"
)

func setup() (chain.IChainAdaptor, error) {
	conf, err := config.New("../../config.yml")
	if err != nil {
		log.Error("load config failed, error:", err)
		return nil, err
	}
	adaptor, err := NewNearAdaptor(conf)
	if err != nil {
		log.Error("create chain adaptor failed, error:", err)
		return nil, err
	}
	return adaptor, nil
}

//
//func TestChainAdaptor_ConvertAddress(t *testing.T) {
//	adaptor, err := setup()
//	if err != nil {
//		return
//	}
//
//	// test account
//	// privateKey: 861ae7df240f80e5492065dafeb6444bdbf2d55d01e1797d2abe0db0afd4f917
//	// publicKey: 02410c64fcd262512683b54576440e3d3033d825ef9f753b44c51ccdd70a7e90c3
//	resp, err := adaptor.ConvertAddress(&account.ConvertAddressRequest{
//		Chain:     ChainName,
//		Network:   "mainnet",
//		PublicKey: "048318535b54105d4a7aae60c08fc45f9687181b4fdfc625bd1a753fa7397fed753547f11ca8696646f2f3acb08e31016afac23e630c5d11f59f61fef57b0d2aa5",
//	})
//	if err != nil {
//		t.Error("convert address failed:", err)
//		return
//	}
//
//	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
//	fmt.Println(resp.Address)
//
//	respJson, _ := json.Marshal(resp)
//	t.Logf("响应: %s", respJson)
//}
//
//func TestChainAdaptor_ValidAddress(t *testing.T) {
//	adaptor, err := setup()
//	if err != nil {
//		return
//	}
//
//	resp, err := adaptor.ValidAddress(&account.ValidAddressRequest{
//		Chain:   ChainName,
//		Network: "mainnet",
//		Address: "0x8358d847Fc823097380c4996A3D3485D9D86941f",
//	})
//	if err != nil {
//		t.Error("valid address failed:", err)
//		return
//	}
//	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
//}

func TestChainAdaptor_GetBlockHeaderByNumber(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetBlockHeaderByNumber(&account.BlockHeaderNumberRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Height:  140083597,
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
		Hash:    "AkdCfheN2Q5ujei7qnCJRkhn33TrqXPyUmhte5epn6LY",
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
		Height: 140083597,
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
		Hash:   "AkdCfheN2Q5ujei7qnCJRkhn33TrqXPyUmhte5epn6LY",
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
		Chain:           ChainName,
		Network:         "mainnet",
		Address:         "20b9bdf32f768ac6e6ff3c9ab512d4bd7f94dbcf4e9d15bb8cd3c3b4062d585a",
		ContractAddress: "0x00",
	})
	if err != nil {
		t.Error("get account failed:", err)
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
		t.Error("get account failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	respJson, _ := json.Marshal(resp)
	t.Logf("响应: %s", respJson)
}

//func TestChainAdaptor_GetTxByAddress(t *testing.T) {
//	adaptor, err := setup()
//	if err != nil {
//		return
//	}
//
//	resp, err := adaptor.GetTxByAddress(&account.TxAddressRequest{
//		Chain:   ChainName,
//		Network: "mainnet",
//		Address: "20b9bdf32f768ac6e6ff3c9ab512d4bd7f94dbcf4e9d15bb8cd3c3b4062d585a",
//	})
//	if err != nil {
//		t.Error("get transaction by address failed:", err)
//		return
//	}
//
//	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
//	fmt.Println(resp.Tx)
//}

func TestChainAdaptor_GetTxByHash(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetTxByHash(&account.TxHashRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Hash:    "7ff3jB4jhFN3Pot4cX4uRiowtwQnqunPV3aRJhZNvwcN",
		Coin:    "users.kaiching",
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
		Start:   "140083597",
		End:     "140083599",
	})
	if err != nil {
		t.Error("get block by range failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.GetBlockHeader())
}

func TestChainAdaptor_SendTx(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}
	sendTxRequest := &account.SendTxRequest{
		Chain:   ChainName,
		Network: "mainnet",
		RawTx:   "QAAAADIwYjliZGYzMmY3NjhhYzZlNmZmM2M5YWI1MTJkNGJkN2Y5NGRiY2Y0ZTlkMTViYjhjZDNjM2I0MDYyZDU4NWEAILm98y92isbm/zyatRLUvX+U289OnRW7jNPDtAYtWFpFvQnBZ38AABAAAAByZWNlaXZlci50ZXN0bmV0/qNZXDvLObzPkHVht/hI4QThWo5Cld9SKwgqIxSpsGkBAAAAAwAAoN7Frck1NgAAAAAAAAAAw0MuGjNuAJZ+iWiXzVVCycYeQ8Vzb4BlYQiD7P55z5UywscYC53d2wTRZdslocF0TTDvAjVfBSXaowleDQNtCQ==",
	}
	resp, err := adaptor.SendTx(sendTxRequest)
	if err != nil {
		t.Error("get block by range failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
}
