package near

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/log"
	"github.com/test-go/testify/assert"

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
	adaptor, err := NewNearAdaptor(conf)
	if err != nil {
		log.Error("create chain adaptor failed, error:", err)
		return nil, err
	}
	return adaptor, nil
}

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
	t.Logf("响应: %s", resp.BlockHeader)

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
	t.Logf("响应: %s", resp.BlockHeader)
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
	t.Logf("响应: %s", resp.Transactions)
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
	t.Logf("响应: %s", resp.Tx)

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
	t.Logf("响应: %s", resp.GetBlockHeader())

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
	t.Logf("响应: %s", resp.GetTxHash())

}

func TestChainAdaptor_BuildSignedTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}
	sendTxRequest := &account.SignedTransactionRequest{
		Chain:    ChainName,
		Base64Tx: "ewogICJwcml2YXRlX2tleSI6ICJlZDI1NTE5OnlwWEo1YngyVDNzUXUxZXd3anVhNHphOHp1QXRnUnl1dEhKVHRocVJpUHlpZWJrQVRYcFl5Mkt3amRMcERjOFFnYlJRdXRydXlxb2FVcFU1d0VyZ2V3YiIsCiAgImFjY291bnRfaWQiOiAiMjBiOWJkZjMyZjc2OGFjNmU2ZmYzYzlhYjUxMmQ0YmQ3Zjk0ZGJjZjRlOWQxNWJiOGNkM2MzYjQwNjJkNTg1YSIsCiAgInJlY2VpdmVyX2lkIjogInJlY2VpdmVyLnRlc3RuZXQiLAogICJhbW91bnQiOiAiMTAwMDAwMDAwMDAwMDAwMDAwMDAwMCIsCiAgInB1YmxpY19rZXkiOiAiZWQyNTUxOTozQ2tLUjJlajJaWEVRaDd0WThia1ZrVnFpMnprdDMxc3ZhQTNNajN2M3BubSIKfQ==",
	}
	resp, err := adaptor.BuildSignedTransaction(sendTxRequest)
	if err != nil {
		t.Error("create build transaction failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	t.Logf("响应: %s", resp.GetSignedTx())
}

func TestChainAdaptor_CreateUnSignTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}
	sendTxRequest := &account.UnSignTransactionRequest{
		Chain:    ChainName,
		Base64Tx: "ewogICJwcml2YXRlX2tleSI6ICJlZDI1NTE5OnlwWEo1YngyVDNzUXUxZXd3anVhNHphOHp1QXRnUnl1dEhKVHRocVJpUHlpZWJrQVRYcFl5Mkt3amRMcERjOFFnYlJRdXRydXlxb2FVcFU1d0VyZ2V3YiIsCiAgImFjY291bnRfaWQiOiAiMjBiOWJkZjMyZjc2OGFjNmU2ZmYzYzlhYjUxMmQ0YmQ3Zjk0ZGJjZjRlOWQxNWJiOGNkM2MzYjQwNjJkNTg1YSIsCiAgInJlY2VpdmVyX2lkIjogInJlY2VpdmVyLnRlc3RuZXQiLAogICJhbW91bnQiOiAiMTAwMDAwMDAwMDAwMDAwMDAwMDAwMCIsCiAgInB1YmxpY19rZXkiOiAiZWQyNTUxOTozQ2tLUjJlajJaWEVRaDd0WThia1ZrVnFpMnprdDMxc3ZhQTNNajN2M3BubSIKfQ==",
	}
	resp, err := adaptor.CreateUnSignTransaction(sendTxRequest)
	if err != nil {
		t.Error("create unsigned transaction failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	t.Logf("响应: %s", resp.GetUnSignTx())
}

func TestChainAdaptor_GetTxByAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}
	sendTxRequest := &account.TxAddressRequest{
		Chain:   ChainName,
		Address: "20b9bdf32f768ac6e6ff3c9ab512d4bd7f94dbcf4e9d15bb8cd3c3b4062d585a",
	}
	resp, err := adaptor.GetTxByAddress(sendTxRequest)
	if err != nil {
		t.Error("create unsigned transaction failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	t.Logf("响应: %s", resp.GetTx())
}

func TestChainAdaptor_ConvertAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}
	sendTxRequest := &account.ConvertAddressRequest{
		Chain:     ChainName,
		PublicKey: "ed25519:3CkKR2ej2ZXEQh7tY8bkVkVqi2zkt31svaA3Mj3v3pnm",
	}
	resp, err := adaptor.ConvertAddress(sendTxRequest)
	if err != nil {
		t.Error("create unsigned transaction failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	t.Logf("响应: %s", resp.Address)
}
