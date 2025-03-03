package icp

import (
	"context"
	"fmt"
	"github.com/dapplink-labs/wallet-chain-account/common/helpers"
	"github.com/dapplink-labs/wallet-chain-account/common/retry"
	"net/http"
	"time"

	"github.com/coinbase/rosetta-sdk-go/client"
	"github.com/ethereum/go-ethereum/log"
)

const (
	userAgent           = "rosetta-sdk-go"
	defaultDialAttempts = 5
)

type IcpClient struct {
	apiClient *client.APIClient
}

func NewIcpClient(ctx context.Context, rpcUrl string, timeOut uint64) (*IcpClient, error) {
	bOff := retry.Exponential()
	icpClient, err := retry.Do(ctx, defaultDialAttempts, bOff, func() (*IcpClient, error) {
		if !helpers.IsURLAvailable(rpcUrl) {
			return nil, fmt.Errorf("address unavailable (%s)", rpcUrl)
		}
		configuration := client.NewConfiguration(rpcUrl, userAgent, &http.Client{
			Timeout: time.Duration(timeOut) * time.Second,
		})

		apiClient := client.NewAPIClient(configuration)
		return &IcpClient{
			apiClient: apiClient,
		}, nil
	})

	if err != nil {
		log.Error("New icp client failed:", err)
		return nil, err
	}
	return icpClient, nil
}
