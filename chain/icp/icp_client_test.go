package icp

import (
	"encoding/json"
	"testing"

	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/ethereum/go-ethereum/log"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
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
	rsp, err := adaptor.GetSupportChains(&account.SupportChainsRequest{
		Chain: ChainName,
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Info("isSupported:", rsp.GetSupport(), ", msg: ", rsp.GetMsg())
}

func TestChainAdaptor_ConvertAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.ConvertAddress(&account.ConvertAddressRequest{
		Chain:     ChainName,
		PublicKey: "fc4f8f805bd58e43fad80281a4b594e019fa18def2ce7a5bc083e048f8f60f36",
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Info("address:", rsp.GetAddress())
}

func TestChainAdaptor_ValidAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.ValidAddress(&account.ValidAddressRequest{
		Chain:   ChainName,
		Address: "9a2101940f1912e3cc2e64c78ec8ecd4547953c297b9e86f6dee89a0104189c6",
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Info("isValid:", rsp.GetValid())
}

func TestChainAdaptor_GetBlockByNumber(t *testing.T) {
}

func TestChainAdaptor_GetBlockByHash(t *testing.T) {
}

func TestChainAdaptor_GetAccount(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.GetAccount(&account.AccountRequest{
		Chain:   ChainName,
		Address: "9a2101940f1912e3cc2e64c78ec8ecd4547953c297b9e86f6dee89a0104189c6",
		Coin:    ChainName,
		//ContractAddress: "0x1Bdd8878252DaddD3Af2ba30628813271294eDc0",
	})
	if err != nil {
		t.Error(err)
	}
	log.Info("account number:", rsp.GetAccountNumber(), ", account balance:", rsp.GetBalance())
}

func TestChainAdaptor_GetFee(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.GetFee(&account.FeeRequest{
		Chain: ChainName,
	})
	if err != nil {
		t.Error(err)
	}
	log.Info("fast fee:", rsp.GetFastFee(), "normal fee:", rsp.GetNormalFee(), "slow fee:", rsp.GetSlowFee())
}

func TestChainAdaptor_GetTxByAddress(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	rsp, err := adaptor.GetTxByAddress(&account.TxAddressRequest{
		Chain:   ChainName,
		Address: "e9ed7def415a8c323953f578c9ce0decf0031cb5b2a1f88cf6f1d89af80ed43a",
		//ContractAddress: "0x1Bdd8878252DaddD3Af2ba30628813271294eDc0",
	})
	if err != nil {
		t.Error(err)
	}
	js, _ := json.Marshal(rsp)
	log.Info(string(js))
}
