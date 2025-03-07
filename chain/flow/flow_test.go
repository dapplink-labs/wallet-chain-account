package flow

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/ethereum/go-ethereum/log"
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

func TestFun1(t *testing.T) {

	adaptor, err := setup()
	if err != nil {
		log.Error("setup failed, error:", err)
	}
	res, err := adaptor.GetBlockByNumber(&account.BlockNumberRequest{Chain: ChainName, Height: 105019513})
	fmt.Println(res)
}

func TestFun2(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		log.Error("setup failed, error:", err)
	}
	res, err := adaptor.GetAccount(&account.AccountRequest{
		Address: "e40303d90745ff5b",
	})
	if err != nil {
		fmt.Println("get account failed, error:", err)
	}
	fmt.Println(res)
}

// 构建哈希
func TestFun3(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		log.Error("setup failed, error:", err)
	}
	message := flowMessage{
		Script:           TransferTpl,
		Argument:         []string{"0.010000", "0x694b6cbf9e36770a"},
		GasLimit:         9999,
		Payer:            "0xe40303d90745ff5b",
		ReferenceBlockId: "3dc4bb6f526ea0c8357e31f35eccc14cebbaf2f7baa75464d1f65fe606c8b1a6",
	}
	fmt.Println(message)
	messageJson, err := json.Marshal(message)
	if err != nil {
		fmt.Println("json marshal err:", err)
	}
	res, err := adaptor.CreateUnSignTransaction(&account.UnSignTransactionRequest{
		Base64Tx: base64.StdEncoding.EncodeToString(messageJson),
	})
	if err != nil {
		fmt.Println("create unsign transaction failed, error:", err)
	}
	fmt.Println("输出res", res)
}

// build完整交易
func TestFun4(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		log.Error("setup failed, error:", err)
	}
	message := flowMessage{
		Script:           TransferTpl,
		Argument:         []string{"0.010000", "0x694b6cbf9e36770a"},
		GasLimit:         9999,
		Payer:            "0xe40303d90745ff5b",
		ReferenceBlockId: "3dc4bb6f526ea0c8357e31f35eccc14cebbaf2f7baa75464d1f65fe606c8b1a6",
	}
	messageJson, err := json.Marshal(message)
	if err != nil {
		fmt.Println("json marshal err:", err)
	}
	base64Msg := base64.StdEncoding.EncodeToString(messageJson)
	hash, _ := GetMassageHash(base64Msg)
	decodeString, err := hex.DecodeString(hash)
	fmt.Println("打印hash decode", decodeString)
	//模拟签名
	sign, err := Sign(decodeString, "5e7163b759d07980dc88cd4d7aba91583d075871c7ffd9f1a9ef11c9645a682b", "seck256k1")
	if err != nil {
		fmt.Println("sign err:", err)
	}
	fmt.Println("打印签名", sign)
	res, err := adaptor.BuildSignedTransaction(&account.SignedTransactionRequest{
		Base64Tx:  base64.StdEncoding.EncodeToString(messageJson),
		Signature: sign,
	})
	if err != nil {
		fmt.Println("build transaction failed, error:", err)
	}
	fmt.Println("输出res", res)

}
