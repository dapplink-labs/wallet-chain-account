package near

import (
	"github.com/dapplink-labs/wallet-chain-account/config"
	nearClient "github.com/eteu-technologies/near-api-go/pkg/client"
)

func NewNearClient(conf *config.Config) (*nearClient.Client, error) {
	newClient, err := nearClient.NewClient(conf.WalletNode.Near.RpcUrl)
	if err != nil {
		return nil, err
	}
	return &newClient, nil
}
