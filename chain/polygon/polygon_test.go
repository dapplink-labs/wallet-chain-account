package polygon

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dapplink-labs/wallet-chain-account/chain/evmbase"
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

	log.Info("ChainName========", ChainName)

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
		Height:  66444218,
	})
	log.Info("ChainName========", ChainName)
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
		Hash:    "0xbd56b33a34ce67fa1bee83da0c0135f16af5296b2d6ff97750f76f52c67eceb6",
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
		Address:         "0x67B94473D81D0cd00849D563C94d0432Ac988B49",
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

func TestChainAdaptor_GetTxByAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}
	resp, err := adaptor.GetTxByAddress(&account.TxAddressRequest{
		Chain:   ChainName,
		Network: "mainnet",
		Address: "0x8916B42a4DB16CA71080dBB0f3650162Ad1E7e3e",
	})
	log.Info("ChainName========", ChainName)
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
		Hash:    "0xc1ecf1f5412b5130f08070cebe83ff83e360c76a7df9ea6e397eb473c3da3282",
	})
	log.Info("ChainName========", ChainName)
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
		Start:   "66444215",
		End:     "66444218",
	})
	if err != nil {
		t.Error("get block by range failed:", err)
		return
	}
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.GetBlockHeader())
}

func createTestBase64Tx() string {

	testTx := evmbase.Eip1559DynamicFeeTx{
		Nonce:                10,
		FromAddress:          "0x67936bb11f8fd1d25da1f94e0aa51039409a7c97",
		ToAddress:            "0x85606bfea925e96285583fabadfcf9cb25bd6721",
		Amount:               "10000000000000000000",
		MaxPriorityFeePerGas: "36819559040",
		MaxFeePerGas:         "36819559041",
		Gas:                  2100,
		GasLimit:             21000,
		ChainId:              "137",
		ContractAddress:      "0x00",
		Signature:            "4f2f529b50b02d1f5b6bb14817cc1708943c2660d7938996b849f0ab3f186e0f6df57918924d23fe54c15602467d38c665aaabb2e54e877aa0452afc74983cd101",
	}

	jsonBytes, err := json.Marshal(testTx)
	if err != nil {
		panic(err)
	}

	base64Str := base64.StdEncoding.EncodeToString(jsonBytes)
	return base64Str
}
func TestChainAdaptor_BuildUnSignTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.BuildUnSignTransaction(&account.UnSignTransactionRequest{
		Chain:    ChainName,
		Network:  "mainnet",
		Base64Tx: createTestBase64Tx(),
	})
	if err != nil {
		log.Error("BuildUnSignTransaction failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.UnSignTx)
}

// 4f2f529b50b02d1f5b6bb14817cc1708943c2660d7938996b849f0ab3f186e0f6df57918924d23fe54c15602467d38c665aaabb2e54e877aa0452afc74983cd101
func TestChainAdaptor_BuildSignedTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.BuildSignedTransaction(&account.SignedTransactionRequest{
		Chain:     ChainName,
		Network:   "mainnet",
		Base64Tx:  createTestBase64Tx(),
		Signature: "0b8ea64c3d7cf8cb0008379afb7a1c79da73d92260ef3a374cefb3f5ae152a6339bac0e1eb80e86836db4fc444ba37ff9c4efb0368e95e16fc5439d8214c98fd00",
	})
	if err != nil {
		log.Error("TestChainAdaptor_BuildSignedTransaction failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.SignedTx)
}

func TestChainAdaptor_SendTx(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}
	RawTx := "0x02f87581890a8508929de2808508929de2818252089485606bfea925e96285583fabadfcf9cb25bd6721888ac7230489e8000080c080a00b8ea64c3d7cf8cb0008379afb7a1c79da73d92260ef3a374cefb3f5ae152a63a039bac0e1eb80e86836db4fc444ba37ff9c4efb0368e95e16fc5439d8214c98fd"
	log.Info("1:", RawTx)
	resp, err := adaptor.SendTx(&account.SendTxRequest{
		Chain:   ChainName,
		Network: "mainnet",
		RawTx:   RawTx,
	})
	if err != nil {
		log.Error("TestChainAdaptor_SendTx failed:", err)
		return
	}
	log.Info(resp.TxHash)
	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
}
func TestChainAdaptor_VerifySignedTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.VerifySignedTransaction(&account.VerifyTransactionRequest{
		Chain:     ChainName,
		Network:   "mainnet",
		Signature: "0x02f8738189808508929de2808508929de281809485606bfea925e96285583fabadfcf9cb25bd6721888ac7230489e8000080c001a04f2f529b50b02d1f5b6bb14817cc1708943c2660d7938996b849f0ab3f186e0fa06df57918924d23fe54c15602467d38c665aaabb2e54e877aa0452afc74983cd1",
	})
	if err != nil {
		log.Error("TestChainAdaptor_VerifySignedTransaction failed:", err)
		return
	}

	assert.Equal(t, common.ReturnCode_SUCCESS, resp.Code)
}
