package flow

import (
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/onflow/flow-go-sdk/access/grpc"
	flow_http "github.com/onflow/flow-go-sdk/access/http"
)

type FlowClient struct {
	client     *grpc.Client
	httpClient *flow_http.Client
}

func NewFlowClient(conf *config.Config) (*FlowClient, error) {
	baseClient, err := grpc.NewClient(conf.WalletNode.Flow.RpcUrl)
	if err != nil {
		return nil, err
	}

	httpClient, err1 := flow_http.NewClient(conf.WalletNode.Flow.DataApiUrl)
	if err1 != nil {
		return nil, err1
	}
	return &FlowClient{
		client:     baseClient,
		httpClient: httpClient,
	}, nil
}
