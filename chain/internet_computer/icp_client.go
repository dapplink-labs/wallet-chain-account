package internet_computer

import (
	"context"
	"fmt"
	"github.com/dapplink-labs/wallet-chain-account/common/helpers"
	"github.com/dapplink-labs/wallet-chain-account/common/retry"
	"net/http"
	"time"

	"github.com/coinbase/rosetta-sdk-go/client"
	"github.com/coinbase/rosetta-sdk-go/types"
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

func (icp *IcpClient) ConstructionDerive(ctx context.Context, publicKey []byte, curveType types.CurveType) (*types.ConstructionDeriveResponse, error) {
	derive, _, err := icp.apiClient.ConstructionAPI.ConstructionDerive(ctx, &types.ConstructionDeriveRequest{
		NetworkIdentifier: &types.NetworkIdentifier{
			Blockchain: "Internet Computer",
			Network:    "00000000000000020101",
		},

		PublicKey: &types.PublicKey{
			Bytes:     publicKey,
			CurveType: curveType,
		},
	})

	if err != nil {
		log.Error("Get block by number failed:", err)
		return nil, err
	}

	return derive, nil
}

func (icp *IcpClient) FetchBlocksByIndex(ctx context.Context, index int64) (*types.BlockResponse, error) {
	block, _, err := icp.apiClient.BlockAPI.Block(ctx, &types.BlockRequest{
		NetworkIdentifier: &types.NetworkIdentifier{
			Blockchain: "Internet Computer",
			Network:    "00000000000000020101",
		},

		BlockIdentifier: &types.PartialBlockIdentifier{
			Index: &index,
		},
	})
	if err != nil {
		log.Error("Get block by number failed:", err)
		return nil, err
	}

	return block, nil
}

func (icp *IcpClient) FetchBlocksByHash(ctx context.Context, hash string) (*types.BlockResponse, error) {
	block, _, err := icp.apiClient.BlockAPI.Block(ctx, &types.BlockRequest{
		NetworkIdentifier: &types.NetworkIdentifier{
			Blockchain: "Internet Computer",
			Network:    "00000000000000020101",
		},
		BlockIdentifier: &types.PartialBlockIdentifier{
			Hash: &hash,
		},
	})
	if err != nil {
		log.Error("Get block by number failed:", err)
		return nil, err
	}

	return block, nil
}

func (icp *IcpClient) FetchAccountBalances(ctx context.Context, address string) (*types.AccountBalanceResponse, error) {
	account, _, err := icp.apiClient.AccountAPI.AccountBalance(ctx, &types.AccountBalanceRequest{
		NetworkIdentifier: &types.NetworkIdentifier{
			Blockchain: "Internet Computer",
			Network:    "00000000000000020101",
		},
		AccountIdentifier: &types.AccountIdentifier{
			Address: address,
		},
	})

	if err != nil {
		log.Error("Get account failed:", err)
		return nil, err
	}
	return account, nil
}

func (icp *IcpClient) ConstructionSubmit(ctx context.Context, signedTx string) (*types.TransactionIdentifierResponse, error) {
	account, _, err := icp.apiClient.ConstructionAPI.ConstructionSubmit(ctx, &types.ConstructionSubmitRequest{
		NetworkIdentifier: &types.NetworkIdentifier{
			Blockchain: "Internet Computer",
			Network:    "00000000000000020101",
		},
		SignedTransaction: signedTx,
	})

	if err != nil {
		log.Error("Get account failed:", err)
		return nil, err
	}
	return account, nil
}
