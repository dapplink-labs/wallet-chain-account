package zksync

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/chain/evmbase"
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
		return
	}

	resp, err := adaptor.GetSupportChains(&account.SupportChainsRequest{
		Chain:   ChainName,
		Network: "mainnet",
	})
	if err != nil {
		t.Error("get support chains failed:", err)
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
}

func TestChainAdaptor_ConvertAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	publicKey := "048846b3ce4376e8d58c83c1c6420a784caa675d7f26c496f499585d09891af8fc9167a4b658b57b28211783cdee651caa8b5341b753fa39c995317670123f12d8"
	resp, err := adaptor.ConvertAddress(&account.ConvertAddressRequest{
		Chain:     ChainName,
		PublicKey: publicKey,
		Network:   "mainnet",
	})
	if err != nil {
		t.Error("convert address failed:", err)
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
}

func TestChainAdaptor_ValidAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.ValidAddress(&account.ValidAddressRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Address: "0x8916B42a4DB16CA71080dBB0f3650162Ad1E7e3e",
	})
	if err != nil {
		t.Error("valid address failed:", err)
	}
	t.Logf("Code：%s", resp.Code)
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	assert.True(t, resp.Valid)

	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
}

func TestChainAdaptor_GetBlockByNumber(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetBlockByNumber(&account.BlockNumberRequest{
		Chain:  ChainName,
		Height: 57458640, // 0x36f77d0
	})
	if err != nil {
		t.Error("get block by number failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
}

func TestChainAdaptor_GetBlockByHash(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetBlockByHash(&account.BlockHashRequest{
		Chain:  ChainName,
		Hash:   "0x604e1f80266e3306e8b9b093d2ae2c5dba9e45318eca968a8bac8d46a9f53283",
		ViewTx: true,
	})
	if err != nil {
		t.Error("get block by hash failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
}

func TestChainAdaptor_GetBlockHeaderByHash(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetBlockHeaderByHash(&account.BlockHeaderHashRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Hash:    "0x604e1f80266e3306e8b9b093d2ae2c5dba9e45318eca968a8bac8d46a9f53283",
	})
	if err != nil {
		t.Error("get block header by hash failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
}

func TestChainAdaptor_GetBlockHeaderByNumber(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetBlockHeaderByNumber(&account.BlockHeaderNumberRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Height:  57458640, // 0x36f77d0
	})
	if err != nil {
		t.Error("get block header by number failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
}

func TestChainAdaptor_GetAccount(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetAccount(&account.AccountRequest{
		Chain:           ChainName,
		Network:         "mainnet",
		Address:         "0x000002c34bAE6DD7BeC72AcbA6aAAC1e01A359De",
		ContractAddress: "0x00",
	})
	if err != nil {
		t.Error("get account failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
}

func TestChainAdaptor_GetAccount_By_ContractAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetAccount(&account.AccountRequest{
		Chain:           ChainName,
		Network:         "mainnet",
		Address:         "0x000002c34bAE6DD7BeC72AcbA6aAAC1e01A359De",
		ContractAddress: "0x145e082e384a9fc86e95eea805dc9012f1b76cb7",
	})
	if err != nil {
		t.Error("get account failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
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
		t.Error("get fee failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
}

func TestChainAdaptor_SendTx(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.SendTx(&account.SendTxRequest{
		Chain:   ChainName,
		Network: "mainnet",
		RawTx:   "",
	})
	if err != nil {
		t.Error("send tx failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
}

func TestChainAdaptor_GetTxByAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetTxByAddress(&account.TxAddressRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Address: "0x000002c34bAE6DD7BeC72AcbA6aAAC1e01A359De",
	})
	if err != nil {
		t.Error("get transaction by address failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
}

func TestChainAdaptor_GetTxByAddress_By_ContractAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetTxByAddress(&account.TxAddressRequest{
		Chain:           ChainName,
		Network:         "mainnet",
		Address:         "0x000002c34bAE6DD7BeC72AcbA6aAAC1e01A359De",
		ContractAddress: "0x145e082e384a9fc86e95eea805dc9012f1b76cb7",
	})
	if err != nil {
		t.Error("get transaction by address failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
}

func TestChainAdaptor_GetTxByHash(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetTxByHash(&account.TxHashRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Hash:    "0xa5d66082c85a722424675105002724f2e8c442281daf1b82ca22136f1a242342",
	})
	if err != nil {
		t.Error("get transaction by hash failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
}

func TestChainAdaptor_GetBlockByRange(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.GetBlockByRange(&account.BlockByRangeRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Start:   "57458640", // 0x36f77d0
		End:     "57458642", // 0x36f77d2
	})
	if err != nil {
		t.Error("get block by range failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
}

// 创建辅助函数生成测试交易数据
func createTestTransactionBase64() string {
	testTx := &evmbase.Eip1559DynamicFeeTx{
		ChainId:              "324",
		Nonce:                1,
		MaxPriorityFeePerGas: "1000000000",
		MaxFeePerGas:         "20000000000",
		GasLimit:             21000,
		FromAddress:          "0x82565b64e8063674CAea7003979280f4dbC3aAE7",
		ToAddress:            "0x8916B42a4DB16CA71080dBB0f3650162Ad1E7e3e",
		Amount:               "1000000000000000000",
		ContractAddress:      "0x00",
	}
	testTxJson, _ := json.Marshal(testTx)
	return base64.StdEncoding.EncodeToString(testTxJson)
}

func TestChainAdaptor_BuildUnSignTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.BuildUnSignTransaction(&account.UnSignTransactionRequest{
		Chain:    ChainName,
		Network:  "mainnet",
		Base64Tx: createTestTransactionBase64(),
	})
	if err != nil {
		t.Error("build unsigned transaction failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
}

func TestChainAdaptor_BuildSignedTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	Signature := "52cf3aa0a66dfe64b6ec18f0bef0e0c90371fc5117c808a024d2c56db5e690f91af6509a8a619438b5babe2c352d2fc20fbf62bffffe72538e0eaa466ad327d601"

	resp, err := adaptor.BuildSignedTransaction(&account.SignedTransactionRequest{
		Chain:     ChainName,
		Network:   "mainnet",
		Base64Tx:  createTestTransactionBase64(),
		Signature: Signature,
	})
	if err != nil {
		t.Error("build signed transaction failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
}

func TestChainAdaptor_DecodeTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	// 已经签名的交易，从BuildSignedTransaction得到
	RawTx := "0x02f87582014401843b9aca008504a817c800825208948916b42a4db16ca71080dbb0f3650162ad1e7e3e880de0b6b3a764000080c001a052cf3aa0a66dfe64b6ec18f0bef0e0c90371fc5117c808a024d2c56db5e690f9a01af6509a8a619438b5babe2c352d2fc20fbf62bffffe72538e0eaa466ad327d6"
	resp, err := adaptor.DecodeTransaction(&account.DecodeTransactionRequest{
		Chain:   ChainName,
		Network: "mainnet",
		RawTx:   RawTx,
	})
	if err != nil {
		t.Error("decode transaction failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)

	// 3. 解码返回的base64数据并验证
	txBytes, err := base64.StdEncoding.DecodeString(resp.Base64Tx)
	if err != nil {
		t.Error("decode base64 failed:", err)
		return
	}

	var decodedTx evmbase.Eip1559DynamicFeeTx
	err = json.Unmarshal(txBytes, &decodedTx)
	if err != nil {
		t.Error("unmarshal transaction failed:", err)
		return
	}

	// 4. 验证解码后的交易数据
	assert.Equal(t, "324", decodedTx.ChainId)
	assert.Equal(t, "0x8916B42a4DB16CA71080dBB0f3650162Ad1E7e3e", decodedTx.ToAddress)
	assert.Equal(t, "1000000000000000000", decodedTx.Amount) // 1 ETH
	assert.Equal(t, uint64(21000), decodedTx.GasLimit)
	assert.Equal(t, "1000000000", decodedTx.MaxPriorityFeePerGas) // 1 Gwei
	assert.Equal(t, "20000000000", decodedTx.MaxFeePerGas)        // 20 Gwei

	resJson, _ := json.Marshal(decodedTx)
	t.Logf("响应：%s", resJson)
}

func TestChainAdaptor_VerifySignedTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	// 1. 测试用的公钥、签名和交易哈希
	publicKey := "048846b3ce4376e8d58c83c1c6420a784caa675d7f26c496f499585d09891af8fc9167a4b658b57b28211783cdee651caa8b5341b753fa39c995317670123f12d8"
	txHash := "0x6e959617f7fdfff5379834171f28680021219479bb189a51c312a7f584224269"
	signature := "52cf3aa0a66dfe64b6ec18f0bef0e0c90371fc5117c808a024d2c56db5e690f91af6509a8a619438b5babe2c352d2fc20fbf62bffffe72538e0eaa466ad327d601"
	// Combine txHash and signature in the format "txHash:signature"
	combinedSignature := txHash + ":" + signature

	// 2. 验证签名
	resp, err := adaptor.VerifySignedTransaction(&account.VerifyTransactionRequest{
		Chain:     ChainName,
		Network:   "mainnet",
		PublicKey: publicKey,
		Signature: combinedSignature,
	})
	if err != nil {
		t.Error("verify signed transaction failed:", err)
		return
	}
	t.Logf("resp: %+v", resp)
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	assert.True(t, resp.Verify)

	// 3. 测试无效的公钥
	resp, err = adaptor.VerifySignedTransaction(&account.VerifyTransactionRequest{
		Chain:     ChainName,
		Network:   "mainnet",
		PublicKey: "invalid_public_key",
		Signature: combinedSignature,
	})
	if err != nil {
		t.Error("verify signed transaction with invalid public key failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_ERROR, resp.Code)
	assert.False(t, resp.Verify)

	// 4. 测试无效的签名格式
	resp, err = adaptor.VerifySignedTransaction(&account.VerifyTransactionRequest{
		Chain:     ChainName,
		Network:   "mainnet",
		PublicKey: publicKey,
		Signature: "invalid_signature_format",
	})
	if err != nil {
		t.Error("verify signed transaction with invalid signature format failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_ERROR, resp.Code)
	assert.False(t, resp.Verify)

	// 5. 测试无效的交易哈希
	resp, err = adaptor.VerifySignedTransaction(&account.VerifyTransactionRequest{
		Chain:     ChainName,
		Network:   "mainnet",
		PublicKey: publicKey,
		Signature: "invalid_tx_hash:" + signature,
	})
	if err != nil {
		t.Error("verify signed transaction with invalid tx hash failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_ERROR, resp.Code)
	assert.False(t, resp.Verify)

	// 6. 测试缺少参数的情况
	resp, err = adaptor.VerifySignedTransaction(&account.VerifyTransactionRequest{
		Chain:   ChainName,
		Network: "mainnet",
	})
	if err != nil {
		t.Error("verify signed transaction with missing parameters failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_ERROR, resp.Code)
	assert.False(t, resp.Verify)

	resJson, _ := json.Marshal(resp)
	t.Logf("响应：%s", resJson)
}

func TestChainAdaptor_GetNftListByAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	// This function is not implemented yet, so we expect it to panic
	assert.Panics(t, func() {
		adaptor.GetNftListByAddress(&account.NftAddressRequest{
			Chain:   ChainName,
			Network: "mainnet",
		})
	})
}
