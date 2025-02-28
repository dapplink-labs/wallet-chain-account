package near

import (
	"github.com/dapplink-labs/wallet-chain-account/config"
	nearClient "github.com/eteu-technologies/near-api-go/pkg/client"
)

//type NearClient struct {
//	client *nearClient.Client
//}

func NewNearClient(conf *config.Config) (*nearClient.Client, error) {
	newClient, err := nearClient.NewClient(conf.WalletNode.Near.RpcUrl)
	if err != nil {
		return nil, err
	}
	return &newClient, nil
}

/**

func NewSuiClient(conf *config.Config) (*SuiClient, error) {
	client := sui.NewSuiClient(conf.WalletNode.Sui.RpcUrl)
	return &SuiClient{client: client}, nil
}
type SuiClient struct {
	client sui.ISuiAPI
}

*/
