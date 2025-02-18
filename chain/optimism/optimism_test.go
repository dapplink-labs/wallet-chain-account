package optimism

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/chain/evmbase"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
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
		Height: 132043390,
	})
	if err != nil {

	}
	js, _ := json.Marshal(rsp)
	log.Info(string(js))

}

func TestChainAdaptor_GetBlockByHash(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.GetBlockByHash(&account.BlockHashRequest{
		Chain: ChainName,
		Hash:  "0x6f7e0f205950dce59896980b53c58d152baafcb085afe6ad0cd9fcd57efcd6b5",
	})
	if err != nil {
		t.Fatal(err)
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
		Hash:  "0x6f7e0f205950dce59896980b53c58d152baafcb085afe6ad0cd9fcd57efcd6b5",
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
		Height: 132043396,
	})
	if err != nil {

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
		Address:         "0x4740d7ee1bd4576ad962f2806b112998cc3b72fc",
		ContractAddress: "0xb12c13e66AdE1F72f71834f2FC5082Db8C091358",
	})
	if err != nil {
		t.Error(err)
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
		Address:         "0x4740d7eE1bD4576aD962f2806b112998Cc3B72Fc",
		ContractAddress: "0xb12c13e66AdE1F72f71834f2FC5082Db8C091358",
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
		Hash:  "0xa31ba448123cdf0a254cad8c5335bdf501023d21124a8d522cb7c91e75452aac",
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
		Start: "128544330",
		End:   "128544335",
	})
	if err != nil {
		t.Error(err)
	}
	js, _ := json.Marshal(rsp)
	log.Info(string(js))
}

func TestChainAdaptor_CreateUnSignTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	base64Tx := createTestBase64Tx("", 0, "", "")
	rsp, err := adaptor.CreateUnSignTransaction(&account.UnSignTransactionRequest{
		Chain:    ChainName,
		Base64Tx: base64Tx,
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

	txDataHash := "0xc1c234195ac9871215cd960190893c4b361699b207d4d546ad8d9de175633e08"
	privateKey := ""
	signature, err := signHash(txDataHash, privateKey)
	if err != nil {
		t.Error(err)
	}

	log.Info("signature:", signature)
	signBase64Tx := createTestBase64Tx(signature, 0, "", "")

	log.Info("signBase64Tx:", signBase64Tx)
	rsp, err := adaptor.BuildSignedTransaction(&account.SignedTransactionRequest{
		Chain:    ChainName,
		Base64Tx: signBase64Tx,
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

	rsp, err := adaptor.SendTx(&account.SendTxRequest{
		Chain: ChainName,
		RawTx: "",
	})
	if err != nil {
		t.Error(err)
	}
	js, _ := json.Marshal(rsp)
	log.Info(string(js))
}

func createTestBase64Tx(signature string, limit uint64, maxGas string, priorityGas string) string {

	if limit == 0 {
		limit = 10000000
	}
	if maxGas == "" {
		maxGas = "26000000000"
	}
	if priorityGas == "" {
		priorityGas = "20520000000"
	}

	testTx := evmbase.Eip1559DynamicFeeTx{
		Nonce:                5,
		FromAddress:          "0x4740d7eE1bD4576aD962f2806b112998Cc3B72Fc",
		ToAddress:            "0x8218a0F47F4c0dE0c1754f50874707cd6e7b2e5e",
		Amount:               "9000000000000000000",
		MaxPriorityFeePerGas: priorityGas,
		MaxFeePerGas:         maxGas,
		GasLimit:             limit,
		ChainId:              "10",
		ContractAddress:      "0xb12c13e66ade1f72f71834f2fc5082db8c091358",
		Signature:            signature,
	}

	jsonBytes, err := json.Marshal(testTx)
	if err != nil {
		panic(err)
	}

	base64Str := base64.StdEncoding.EncodeToString(jsonBytes)
	return base64Str
}

// {"msg":"get gas price success","slow_fee":"1001881|1000000","normal_fee":"1001881|1000000|*2","fast_fee":"1001881|1000000|*3"}
func TestChainAdaptor_SendTx2(t *testing.T) {

	privateKey := ""

	// 获取当前区块的 baseFee（即当前区块的基础费用）
	baseFee := new(big.Int).SetUint64(1001881)

	// 设置最大优先费用（maxPriorityFeePerGas）为 2 Gwei
	maxPriorityFeePerGas := new(big.Int).SetUint64(0.02 * params.GWei)

	// 计算 maxFeePerGas = baseFee + maxPriorityFeePerGas
	maxFeePerGas := new(big.Int).Add(baseFee, maxPriorityFeePerGas)

	limit := uint64(10000000)

	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	base64Tx := createTestBase64Tx("", limit, maxFeePerGas.String(), maxPriorityFeePerGas.String())
	rsp0, err := adaptor.CreateUnSignTransaction(&account.UnSignTransactionRequest{
		Chain:    ChainName,
		Base64Tx: base64Tx,
	})
	if err != nil {
		t.Error(err)
	}

	log.Info("hash =====  ", rsp0.UnSignTx)
	signature, err := signHash(rsp0.UnSignTx, privateKey)
	if err != nil {
		t.Error(err)
	}

	signBase64Tx := createTestBase64Tx(signature, limit, maxFeePerGas.String(), maxPriorityFeePerGas.String())

	rsp1, err := adaptor.BuildSignedTransaction(&account.SignedTransactionRequest{
		Chain:    ChainName,
		Base64Tx: signBase64Tx,
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
