package ton

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
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

func Test_GeneralTonWallet(t *testing.T) {
	conf, err := config.New("../../config.yml")
	if err != nil {
		log.Error("load config failed, error:", err)
		return
	}
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), conf.WalletNode.Ton.RpcUrl)
	if err != nil {
		log.Error("get config from ton url fail", "err", err)
	}

	client := liteclient.NewConnectionPool()
	// connect to mainnet lite server
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	//err := client.AddConnection(context.Background(), "135.181.140.212:13206", "K0t3+IWLOXHYMvMcrGZDPs+pn58a17LFbnXoQkKc2xw=")
	if err != nil {
		panic(err)
	}
	// initialize ton api lite connection wrapper
	api := ton.NewAPIClient(client).WithRetry()

	seed := wallet.NewSeed()
	// wallet version: V5R1Final  Version = 52
	w, err := wallet.FromSeed(api, seed, wallet.ConfigV5R1Final{
		NetworkGlobalID: wallet.MainnetGlobalID,
		Workchain:       0,
	}, true)

	if err != nil {
		t.Fatal(err)
	}
	addr := w.WalletAddress()

	// only test
	fmt.Println("seed:", seed)
	fmt.Println("testnet Address:", addr.Copy().Testnet(true))
	fmt.Println("mainnet Address:", addr.Copy().Testnet(false).Bounce(false))
	fmt.Println("encode privateKey:", hex.EncodeToString(w.PrivateKey()))
	fmt.Println("encode publicKey:", hex.EncodeToString(w.PrivateKey().Public().(ed25519.PublicKey)))
	// seed: [surround near bacon recycle bachelor thumb any spot wrong fancy attend fringe useful quarter luxury oil bench exist swift bread lucky pottery innocent fox]
	// testnet Address: 0QAHVc9OI1kqgh68auUYhe2YziHXDjw8fgla5sqw4kOszV9A
	// mainnet Address: UQAHVc9OI1kqgh68auUYhe2YziHXDjw8fgla5sqw4kOszeTK
	// encode privateKey: b0575c0125936c644c4619d51ab290a817bc62b078a77bb7c51b89c638e138aef7b418ebb5a0af910eab73aefd0c6cce147c6dea6dd87ec952eea75342f8a6a5
	// encode publicKey: f7b418ebb5a0af910eab73aefd0c6cce147c6dea6dd87ec952eea75342f8a6a5
}

func TestChainAdaptor_ConvertAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	// test account
	resp, err := adaptor.ConvertAddress(&account.ConvertAddressRequest{
		Chain:     ChainName,
		Network:   "mainnet",
		PublicKey: "f7b418ebb5a0af910eab73aefd0c6cce147c6dea6dd87ec952eea75342f8a6a5",
	})
	if err != nil {
		t.Error("convert address failed:", err)
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
		Address: "UQAHVc9OI1kqgh68auUYhe2YziHXDjw8fgla5sqw4kOszeTK",
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
		Chain:           ChainName,
		Network:         "mainnet",
		Address:         "0xD79053a14BC465d9C1434d4A4fAbdeA7b6a2A94b",
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
