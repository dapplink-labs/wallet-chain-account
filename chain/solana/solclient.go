package solana

import (
	"context"
	"encoding/base64"
	"log"
	"math/big"
	"strconv"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/dapplink-labs/wallet-chain-account/config"
)

type SolClient struct {
	RpcClient rpc.RpcClient
	Client    *client.Client
	SolConfig config.WalletNode
}

func NewSolClient(conf *config.Config) (*SolClient, error) {
	endpoint := conf.WalletNode.Sol.RpcUrl
	rpcClient := rpc.NewRpcClient(endpoint)
	clientNew := client.NewClient(endpoint)
	return &SolClient{
		RpcClient: rpcClient,
		Client:    clientNew,
		SolConfig: conf.WalletNode,
	}, nil
}

func (sol *SolClient) GetLatestBlockHeight() (int64, error) {
	res, err := sol.RpcClient.GetBlockHeight(context.Background())
	if err != nil {
		return 0, err
	}
	return int64(res.Result), nil
}

func (sol *SolClient) GetBalance(address string) (string, error) {
	balance, err := sol.Client.GetBalanceWithConfig(
		context.TODO(),
		address,
		client.GetBalanceConfig(rpc.GetBalanceConfig{
			Commitment: rpc.CommitmentProcessed,
		}),
	)
	if err != nil {
		return "", err
	}

	lamportsOnAccount := new(big.Float).SetUint64(balance)
	solBalance := new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(1000000000))

	return solBalance.String(), nil
}

func (sol *SolClient) GetTxByHash(hash string) (*TxMessage, error) {
	return nil, nil
}

func (sol *SolClient) SendTx(rawTx string) (string, error) {
	res, err := sol.RpcClient.SendTransactionWithConfig(
		context.Background(),
		base64.StdEncoding.EncodeToString([]byte(rawTx)),
		rpc.SendTransactionConfig{
			Encoding: rpc.SendTransactionConfigEncodingBase64,
		},
	)
	if err != nil {
		return "", err
	}
	if res.Error != nil {
		return "", res.Error
	}
	return res.Result, nil
}

func (sol *SolClient) GetNonce(NonceAccount string) (string, error) {
	nonce, err := sol.Client.GetNonceFromNonceAccount(context.Background(), NonceAccount)
	println("nonce:", nonce)
	if err != nil {
		log.Fatalf("failed to get nonce account, err: %v", err)
		return "", err
	}
	return nonce, nil
}

func (sol *SolClient) GetMinRent() (string, error) {
	bal, err := sol.RpcClient.GetMinimumBalanceForRentExemption(context.Background(), 100)
	if err != nil {
		log.Fatalf("failed to get GetMinimumBalanceForRentExemption , err: %v", err)
		return "", err
	}
	return strconv.FormatUint(bal.Result, 10), nil
}
