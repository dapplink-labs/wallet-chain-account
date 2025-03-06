package icp

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/aviate-labs/agent-go/principal"
	"github.com/coinbase/rosetta-sdk-go/client"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/dapplink-labs/wallet-chain-account/common/helpers"
	"github.com/dapplink-labs/wallet-chain-account/common/retry"
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

func (c *Client) DeriveAddressByPublicKey(publicKey string) (*types.ConstructionDeriveResponse, error) {
	publicKeyBytes, _ := hex.DecodeString(publicKey)
	derive, _, err := c.apiClient.ConstructionAPI.ConstructionDerive(context.Background(), &types.ConstructionDeriveRequest{
		NetworkIdentifier: &types.NetworkIdentifier{
			Blockchain: rosettaChainName,
			Network:    rosettaChainNetwork,
		},
		PublicKey: &types.PublicKey{
			Bytes:     publicKeyBytes,
			CurveType: types.Edwards25519,
		},
	})
	if err != nil {
		log.Error("get balance err", "err", err)
		panic(err)
	}
	return derive, nil
}

func (c *Client) ValidAddress(address string) (bool, error) {
	isValid, err := principal.DecodeAccountID(address)
	if err != nil {
		log.Error("valid err", "err", err)
		panic(err)
	}
	return isValid.String() == address, nil
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

func (c *Client) BuildUnsignedTransaction(from string, to string, amount string) (string, error) {
	toAmount := amount
	fromAmount := "-" + amount

	network := &types.NetworkIdentifier{
		Blockchain: rosettaChainName,
		Network:    rosettaChainNetwork,
	}

	ops := []*types.Operation{
		{
			OperationIdentifier: &types.OperationIdentifier{Index: 0},
			Type:                "TRANSACTION",
			Account:             &types.AccountIdentifier{Address: from},
			Amount:              &types.Amount{Value: fromAmount, Currency: &types.Currency{Symbol: "ICP", Decimals: 8}},
		},
		{
			OperationIdentifier: &types.OperationIdentifier{Index: 1},
			Type:                "TRANSACTION",
			Account:             &types.AccountIdentifier{Address: to},
			Amount:              &types.Amount{Value: toAmount, Currency: &types.Currency{Symbol: "ICP", Decimals: 8}},
		},
		{
			OperationIdentifier: &types.OperationIdentifier{Index: 2},
			Type:                "FEE",
			Account:             &types.AccountIdentifier{Address: from},
			Amount:              &types.Amount{Value: "-10000", Currency: &types.Currency{Symbol: "ICP", Decimals: 8}},
		},
	}

	// 预处理
	preprocessResp, _, err := c.apiClient.ConstructionAPI.ConstructionPreprocess(context.Background(), &types.ConstructionPreprocessRequest{
		NetworkIdentifier: network,
		Operations:        ops,
	})
	if err != nil {
		log.Error("build tx err", "err", err)
		panic(err)
	}
	// 获取元数据
	metadataResp, _, err := c.apiClient.ConstructionAPI.ConstructionMetadata(context.Background(), &types.ConstructionMetadataRequest{
		NetworkIdentifier: network,
		Options:           preprocessResp.Options,
	})
	if err != nil {
		log.Error("build tx err", "err", err)
		panic(err)
	}
	// 生成签名 payload
	payloadsResp, _, err := c.apiClient.ConstructionAPI.ConstructionPayloads(context.Background(), &types.ConstructionPayloadsRequest{
		NetworkIdentifier: network,
		Operations:        ops,
		Metadata:          metadataResp.Metadata,
	})
	if err != nil {
		log.Error("build tx err", "err", err)
		panic(err)
	}

	return payloadsResp.UnsignedTransaction, nil
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
