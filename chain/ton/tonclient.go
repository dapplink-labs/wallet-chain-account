package ton

import (
	"context"

	"github.com/ethereum/go-ethereum/log"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"

	"github.com/dapplink-labs/wallet-chain-account/config"
)

type TonClient struct {
	client *liteclient.ConnectionPool
	api    ton.APIClientWrapped
}

func NewTonClients(conf *config.Config) (*TonClient, error) {
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), conf.WalletNode.Ton.RpcUrl)
	if err != nil {
		log.Error("get config from ton url fail", "err", err)
		return nil, err
	}

	client := liteclient.NewConnectionPool()
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		log.Error("add connections from config fail", "err", err)
		return nil, err
	}
	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	return &TonClient{
		client: client,
		api:    api,
	}, nil
}

func (tc *TonClient) GetAccountInfo(addr string) (string, uint64, error) {
	block, err := tc.api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return "0", 0, err
	}
	acc, err := tc.api.GetAccount(context.Background(), block, address.MustParseAddr(addr))
	if err != nil {
		return "", 0, err
	}
	return acc.State.Balance.String(), acc.State.LastTransactionLT, nil
}
