package oasis

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

type OasisRosettaClient struct {
	apiClient *client.APIClient
}

func NewOasisRosettaClient(ctx context.Context, rpcUrl string, timeOut uint64) (*OasisRosettaClient, error) {
	bOff := retry.Exponential()
	oasisRosettaClient, err := retry.Do(ctx, defaultDialAttempts, bOff, func() (*OasisRosettaClient, error) {
		if !helpers.IsURLAvailable(rpcUrl) {
			return nil, fmt.Errorf("address unavailable (%s)", rpcUrl)
		}
		configuration := client.NewConfiguration(rpcUrl, userAgent, &http.Client{
			Timeout: time.Duration(timeOut) * time.Second,
		})

		apiClient := client.NewAPIClient(configuration)
		return &OasisRosettaClient{
			apiClient: apiClient,
		}, nil
	})

	if err != nil {
		log.Error("New oasis client failed:", err)
		return nil, err
	}
	return oasisRosettaClient, nil
}

func (oasis *OasisRosettaClient) ConstructionSubmit(ctx context.Context, signedTx string) (*types.TransactionIdentifierResponse, error) {
	response, _, err := oasis.apiClient.ConstructionAPI.ConstructionSubmit(ctx, &types.ConstructionSubmitRequest{
		NetworkIdentifier: &types.NetworkIdentifier{
			Blockchain: "Oasis",
			Network:    "bb3d748def55bdfb797a2ac53ee6ee141e54cd2ab2dc2375f4a0703a178e6e55",
		},
		SignedTransaction: signedTx,
	})

	if err != nil {
		log.Error("ConstructionSubmit failed:", err)
		return nil, err
	}
	return response, nil
}

func (oasis *OasisRosettaClient) ConstructionDerive(ctx context.Context, publicKey []byte, curveType types.CurveType) (*types.ConstructionDeriveResponse, error) {
	derive, _, err := oasis.apiClient.ConstructionAPI.ConstructionDerive(ctx, &types.ConstructionDeriveRequest{
		NetworkIdentifier: &types.NetworkIdentifier{
			Blockchain: "Oasis",
			Network:    "bb3d748def55bdfb797a2ac53ee6ee141e54cd2ab2dc2375f4a0703a178e6e55",
		},

		PublicKey: &types.PublicKey{
			Bytes:     publicKey,
			CurveType: curveType,
		},
	})

	if err != nil {
		log.Error("ConstructionDerive failed:", err)
		return nil, err
	}

	return derive, nil
}

func (oasis *OasisRosettaClient) FetchAccountBalances(ctx context.Context, address string) (*types.AccountBalanceResponse, error) {
	account, _, err := oasis.apiClient.AccountAPI.AccountBalance(ctx, &types.AccountBalanceRequest{
		NetworkIdentifier: &types.NetworkIdentifier{
			Blockchain: "Oasis",
			Network:    "bb3d748def55bdfb797a2ac53ee6ee141e54cd2ab2dc2375f4a0703a178e6e55",
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
