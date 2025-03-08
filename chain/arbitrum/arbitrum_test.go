package arbitrum

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/dapplink-labs/wallet-chain-account/chain/evmbase"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
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
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	adaptor.GetSupportChains(&account.SupportChainsRequest{
		Chain: ChainName,
	})
}

func TestChainAdaptor_ConvertAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.ConvertAddress(&account.ConvertAddressRequest{
		Chain:     ChainName,
		PublicKey: "048318535b54105d4a7aae60c08fc45f9687181b4fdfc625bd1a753fa7397fed753547f11ca8696646f2f3acb08e31016afac23e630c5d11f59f61fef57b0d2aa5",
	})

	log.Info("========", rsp.Address)

	js, _ := json.Marshal(rsp)

	log.Info(string(js))

}

func TestChainAdaptor_ValidAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.ValidAddress(&account.ValidAddressRequest{
		Chain:   ChainName,
		Address: "0x4740d7eE1bD4576aD962f2806b112998Cc3B72Fc",
	})
	if err != nil {
		t.Fatal(err)
	}
	js, _ := json.Marshal(rsp)
	log.Info(string(js))
}

func TestChainAdaptor_GetBlockByNumber(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.GetBlockByNumber(&account.BlockNumberRequest{
		Chain:  ChainName,
		Height: 305215338,
	})
	if err != nil {

	}
	js, _ := json.Marshal(rsp)
	log.Info(string(js))

}

func TestChainAdaptor_GetBlockByHash(t *testing.T) {
	adaptor, err := setup()
	if err != nil {

	}
	rsp, err := adaptor.GetBlockByHash(&account.BlockHashRequest{
		Chain: ChainName,
		Hash:  "0xb9e4fe3ee850aa8296dace7d6d70ccbe91a0a742b4c10d9fad756db58afe6425",
	})
	if err != nil {

	}
	js, _ := json.Marshal(rsp)

	log.Info(string(js))
}

func TestChainAdaptor_GetBlockHeaderByHash(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.GetBlockHeaderByHash(&account.BlockHeaderHashRequest{
		Chain: ChainName,
		Hash:  "0xcac38b7d9a8f27e1f96ad171d15fe219e2df2062cf52e3fa803907dafb0fabf0",
	})
	if err != nil {

	}
	js, _ := json.Marshal(rsp)
	log.Info(string(js))
}

func TestChainAdaptor_GetBlockHeaderByNumber(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.GetBlockHeaderByNumber(&account.BlockHeaderNumberRequest{
		Chain:  ChainName,
		Height: 305215338,
	})
	if err != nil {

	}
	js, _ := json.Marshal(rsp)
	log.Info(string(js))
}

func TestChainAdaptor_GetFee(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.GetFee(&account.FeeRequest{
		Chain: ChainName,
		RawTx: "",
	})
	if err != nil {

	}
	js, _ := json.Marshal(rsp)
	log.Info(string(js))
}

func TestChainAdaptor_GetTxByAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.GetTxByAddress(&account.TxAddressRequest{
		Chain:           ChainName,
		Address:         "0x4740d7ee1bd4576ad962f2806b112998cc3b72fc",
		ContractAddress: "0xfd086bc7cd5c481dcc9c85ebe478a1c0b69fcbb9",
	})
	if err != nil {
		t.Error(err)
	}
	js, _ := json.Marshal(rsp)
	log.Info(string(js))
}

func TestChainAdaptor_GetTxByHash(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.GetTxByHash(&account.TxHashRequest{
		Chain: ChainName,
		Hash:  "",
	})
	if err != nil {
		t.Error(err)
	}
	js, _ := json.Marshal(rsp)
	log.Info(string(js))
}

func TestChainAdaptor_GetBlockByRange(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.GetBlockByRange(&account.BlockByRangeRequest{
		Chain: ChainName,
		Start: "305216250",
		End:   "305216256",
	})
	if err != nil {
		t.Error(err)
	}
	js, _ := json.Marshal(rsp)
	log.Info(string(js))
}

func TestChainAdaptor_GetAccount(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.GetAccount(&account.AccountRequest{
		Chain:           ChainName,
		Address:         "0x4740d7eE1bD4576aD962f2806b112998Cc3B72Fc",
		ContractAddress: "0xf525E73bdeB4ac1b0e741aF3Ed8a8CBB43ab0756",
	})
	if err != nil {
		t.Error(err)
	}
	js, _ := json.Marshal(rsp)
	log.Info(string(js))

}

func TestChainAdaptor_BuildUnSignTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	baseTx := createTestBase64Tx("", 0, "", "")

	rsp, err := adaptor.BuildUnSignTransaction(&account.UnSignTransactionRequest{
		Chain:    ChainName,
		Base64Tx: baseTx,
	})
	if err != nil {
		t.Error(err)
	}
	js, _ := json.Marshal(rsp)
	log.Info(string(js))
}

func TestChainAdaptor_BuildSignedTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	txDataHash := "0x717760a880f5b90b9f58b3522942452db93b3321e6ce6d7a9b3eac6f20127e95"
	privateKey := ""
	signature, err := signHash(txDataHash, privateKey)
	if err != nil {
		t.Error(err)
	}

	baseTx := createTestBase64Tx(signature, 0, "", "")
	rsp, err := adaptor.BuildSignedTransaction(&account.SignedTransactionRequest{
		Chain:    ChainName,
		Base64Tx: baseTx,
	})
	if err != nil {
		t.Error(err)
	}
	js, _ := json.Marshal(rsp)
	log.Info(string(js))

}

func TestChainAdaptor_SendTx(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rawTx := "0x02f8b282a4b1278464409c40849fdb664082520894f525e73bdeb4ac1b0e741af3ed8a8cbb43ab075680b845a9059cbb2a0000000000000000000000008218a0f47f4c0de0c1754f50874707cd6e7b2e5e00000000000000000000000000000000000000000000004eea56726e9f5e5ca4c080a0018ce3908327a15a0f2a13203ea7800625631609bf0fcc5a2974aef37dc880cfa0082a6a05844718ed39c3322f7769d3c0bfa1d0a8cc0106ebeec9cc600b1b5ece"
	rsp, err := adaptor.SendTx(&account.SendTxRequest{
		Chain: ChainName,
		RawTx: rawTx,
	})
	if err != nil {
		t.Error(err)
	}
	js, _ := json.Marshal(rsp)
	log.Info(string(js))
}

func createTestBase64Tx(signature string, limit uint64, maxGas string, priorityGas string) string {

	if limit == 0 {
		limit = 21000
	}
	if maxGas == "" {
		maxGas = "26000000000"
	}
	if priorityGas == "" {
		priorityGas = "20520000000"
	}

	testTx := evmbase.Eip1559DynamicFeeTx{
		Nonce:                50,
		FromAddress:          "0x4740d7eE1bD4576aD962f2806b112998Cc3B72Fc",
		ToAddress:            "0x8218a0F47F4c0dE0c1754f50874707cd6e7b2e5e",
		Amount:               "88000000000000000000",
		MaxPriorityFeePerGas: priorityGas,
		MaxFeePerGas:         maxGas,
		GasLimit:             limit,
		ChainId:              "42161",
		ContractAddress:      "0xf525E73bdeB4ac1b0e741aF3Ed8a8CBB43ab0756",
		Signature:            signature,
	}

	jsonBytes, err := json.Marshal(testTx)
	if err != nil {
		panic(err)
	}

	base64Str := base64.StdEncoding.EncodeToString(jsonBytes)
	return base64Str
}

func TestChainAdaptor_SendTx2(t *testing.T) {

	privateKey := ""

	// 获取当前区块的 baseFee（即当前区块的基础费用）
	baseFee := new(big.Int).SetUint64(10000000)

	// 设置最大优先费用（maxPriorityFeePerGas）为 2 Gwei
	maxPriorityFeePerGas := new(big.Int).SetUint64(2 * params.GWei)

	// 计算 maxFeePerGas = baseFee + maxPriorityFeePerGas
	maxFeePerGas := new(big.Int).Add(baseFee, maxPriorityFeePerGas)

	limit := uint64(100000)

	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	baseTx := createTestBase64Tx("", limit, maxFeePerGas.String(), maxPriorityFeePerGas.String())
	rsp0, err := adaptor.BuildUnSignTransaction(&account.UnSignTransactionRequest{
		Chain:    ChainName,
		Base64Tx: baseTx,
	})

	signature, err := signHash(rsp0.UnSignTx, privateKey)
	if err != nil {
		t.Error(err)
	}

	signBaseTx := createTestBase64Tx(signature, limit, maxFeePerGas.String(), maxPriorityFeePerGas.String())
	rsp1, err := adaptor.BuildSignedTransaction(&account.SignedTransactionRequest{
		Chain:    ChainName,
		Base64Tx: signBaseTx,
	})
	if err != nil {
		t.Error(err)
	}

	log.Info(rsp1.SignedTx)
	rsp2, err := adaptor.SendTx(&account.SendTxRequest{
		Chain: ChainName,
		RawTx: rsp1.SignedTx,
	})
	if err != nil {
		t.Error(err)
	}

	js, _ := json.Marshal(rsp2)
	log.Info(string(js))

}

func signHash(hash string, privateKey string) (signature string, err error) {
	bytes, err := hexutil.Decode(hash)
	if err != nil {
		panic(err)
	}
	prkByte, err := hex.DecodeString(privateKey)
	if err != nil {
		panic(err)
	}
	prk, err := crypto.ToECDSA(prkByte)
	if err != nil {
		panic(err)
	}
	sig, err := crypto.Sign(bytes, prk)
	if err != nil {
		panic(err)
	}
	signature = hex.EncodeToString(sig)
	return signature, nil
}
