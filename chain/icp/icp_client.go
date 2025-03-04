package icp

import (
	"context"
	"fmt"
	"github.com/coinbase/rosetta-sdk-go/types"
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
	rosettaChainName    = "Internet Computer"
	rosettaChainNetwork = "00000000000000020101"
)

type Client struct {
	apiClient *client.APIClient
}

func (c *Client) GetAccountBalance(address string) (*types.AccountBalanceResponse, error) {
	balance, _, err := c.apiClient.AccountAPI.AccountBalance(context.Background(), &types.AccountBalanceRequest{
		NetworkIdentifier: &types.NetworkIdentifier{
			Blockchain: rosettaChainName,
			Network:    rosettaChainNetwork,
		},
		AccountIdentifier: &types.AccountIdentifier{
			Address: address,
		},
	})
	if err != nil {
		log.Error("get balance err", "err", err)
		panic(err)
	}
	return balance, nil
}

func (c *Client) GetBlockByNumber(number int64) (*types.BlockResponse, error) {
	block, _, err := c.apiClient.BlockAPI.Block(context.Background(), &types.BlockRequest{
		NetworkIdentifier: &types.NetworkIdentifier{
			Blockchain: rosettaChainName,
			Network:    rosettaChainNetwork,
		},
		BlockIdentifier: &types.PartialBlockIdentifier{
			Index: &number,
		},
	})
	if err != nil {
		log.Error("get block err", "err", err)
		panic(err)
	}
	return block, nil
}

func (c *Client) GetBlockByHash(hash string) (*types.BlockResponse, error) {
	block, _, err := c.apiClient.BlockAPI.Block(context.Background(), &types.BlockRequest{
		NetworkIdentifier: &types.NetworkIdentifier{
			Blockchain: rosettaChainName,
			Network:    rosettaChainNetwork,
		},
		BlockIdentifier: &types.PartialBlockIdentifier{
			Hash: &hash,
		},
	})
	if err != nil {
		log.Error("get block err", "err", err)
		panic(err)
	}
	return block, nil
}

func (c *Client) GetTxByHash(hash string) (*types.SearchTransactionsResponse, error) {
	tx, _, err := c.apiClient.SearchAPI.SearchTransactions(context.Background(), &types.SearchTransactionsRequest{
		NetworkIdentifier: &types.NetworkIdentifier{
			Blockchain: rosettaChainName,
			Network:    rosettaChainNetwork,
		},
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: hash,
		},
	})
	if err != nil {
		log.Error("get tx err", "err", err)
		panic(err)
	}
	return tx, nil
}

func (c *Client) GetTxByAddress(address string) (*types.SearchTransactionsResponse, error) {
	tx, _, err := c.apiClient.SearchAPI.SearchTransactions(context.Background(), &types.SearchTransactionsRequest{
		NetworkIdentifier: &types.NetworkIdentifier{
			Blockchain: rosettaChainName,
			Network:    rosettaChainNetwork,
		},
		AccountIdentifier: &types.AccountIdentifier{
			Address: address,
		},
	})
	if err != nil {
		log.Error("get tx err", "err", err)
		panic(err)
	}
	return tx, nil
}

func NewIcpClient(ctx context.Context, rpcUrl string, timeOut uint64) (*Client, error) {
	bOff := retry.Exponential()
	icpClient, err := retry.Do(ctx, defaultDialAttempts, bOff, func() (*Client, error) {
		if !helpers.IsURLAvailable(rpcUrl) {
			return nil, fmt.Errorf("address unavailable (%s)", rpcUrl)
		}
		configuration := client.NewConfiguration(rpcUrl, userAgent, &http.Client{
			Timeout: time.Duration(timeOut) * time.Second,
		})

		apiClient := client.NewAPIClient(configuration)
		return &Client{
			apiClient: apiClient,
		}, nil
	})

	if err != nil {
		log.Error("New icp client failed:", err)
		return nil, err
	}
	return icpClient, nil
}
