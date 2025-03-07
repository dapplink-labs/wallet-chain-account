package btt

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
		PublicKey: "041b91a40686b23eb5f42fb40c98c6a16ce95d999928b4fbf2c6dc7274abedc97a9404a1e95af3983eeef6e02af4f7cca3d4a3a91152dba520e3f08a9f8e6543a1",
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
		Address: "0xdda22000e1bcc0c70c8b1947ce7074df1dc5b80b",
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
		Height: 49309386,
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
		Hash:  "0xeda79567b1b15e638bdb0469864c40ed4a3b80df65f4b721610a9d448b9b9346",
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
		Hash:  "0xeda79567b1b15e638bdb0469864c40ed4a3b80df65f4b721610a9d448b9b9346",
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
		Height: 49309386,
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
		Address:         "0xdda22000e1bcc0c70c8b1947ce7074df1dc5b80b",
		ContractAddress: "0x1Bdd8878252DaddD3Af2ba30628813271294eDc0",
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
		Address:         "0xDDA22000e1bCC0c70C8b1947CE7074df1DC5B80B",
		ContractAddress: "0x1Bdd8878252DaddD3Af2ba30628813271294eDc0",
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
		Hash:  "0xfe66799cd6de5b8a6a9657bf91cb64101d8c0f511b52ab644b43bb92688d2a26",
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
		Start: "75827131",
		End:   "75827133",
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

	txDataHash := "0xb419c2ec2cc8297a42f81dd4ebec5697bcc4ae992606c7eb689913382b8f2813"
	privateKey := ""
	signature, err := signHash(txDataHash, privateKey)
	if err != nil {
		t.Error(err)
	}

	signBase64Tx := createTestBase64Tx(signature, 0, "", "")

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
		RawTx: "0x02f8b4821388018504c7165a0085060db88400830186a0941bdd8878252daddd3af2ba30628813271294edc080b844a9059cbb0000000000000000000000008218a0f47f4c0de0c1754f50874707cd6e7b2e5e000000000000000000000000000000000000000000000000016345785d8a0000c001a08649363f6ecab9c46f47bfb718ca2d90c4ebd836d4bd1bddb0cbafa62ea76beba00a3fa998a3274eef0925eb076c36017c65c0e03ad7f037b768a3040b22ae4b0e",
	})
	if err != nil {
		t.Error(err)
	}
	js, _ := json.Marshal(rsp)
	log.Info(string(js))
}

func createTestBase64Tx(signature string, limit uint64, maxGas string, priorityGas string) string {

	if limit == 0 {
		limit = 100000
	}
	if maxGas == "" {
		maxGas = "26000000000"
	}
	if priorityGas == "" {
		priorityGas = "20520000000"
	}

	testTx := evmbase.Eip1559DynamicFeeTx{
		Nonce:                1,
		FromAddress:          "0xDDA22000e1bCC0c70C8b1947CE7074df1DC5B80B",
		ToAddress:            "0x8218a0F47F4c0dE0c1754f50874707cd6e7b2e5e",
		Amount:               "100000000000000000",
		MaxPriorityFeePerGas: priorityGas,
		MaxFeePerGas:         maxGas,
		GasLimit:             limit,
		ChainId:              "5000",
		ContractAddress:      "0x1Bdd8878252DaddD3Af2ba30628813271294eDc0",
		Signature:            signature,
	}

	jsonBytes, err := json.Marshal(testTx)
	if err != nil {
		panic(err)
	}

	base64Str := base64.StdEncoding.EncodeToString(jsonBytes)
	return base64Str
}

// {"msg":"get gas price success","slow_fee":"20000000|0","normal_fee":"20000000|0|*2","fast_fee":"20000000|0|*3"}
func TestChainAdaptor_SendTx2(t *testing.T) {

	privateKey := "8191d4626b096f3b7dcf90d71931011de7750f7bbf1684792f3d91d93c5926e3"

	// 获取当前区块的 baseFee（即当前区块的基础费用）
	baseFee := new(big.Int).SetUint64(20000000)

	// 设置最大优先费用（maxPriorityFeePerGas）为 2 Gwei
	maxPriorityFeePerGas := new(big.Int).SetUint64(2 * params.GWei)

	// 计算 maxFeePerGas = baseFee + maxPriorityFeePerGas
	maxFeePerGas := new(big.Int).Add(baseFee, maxPriorityFeePerGas)

	limit := uint64(250226535)

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
