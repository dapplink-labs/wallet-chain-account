package tezos

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
	"github.com/trilitech/tzgo/rpc"
	"time"
)

var errBlockChainHTTPError = errors.New("tezos blockchain http error")

type TezosDataClient struct {
	client *rpc.Client
}

func NewTezosDataClient(url string) (*TezosDataClient, error) {
	if url == "" {
		return nil, fmt.Errorf("tezos blockchain URL cannot be empty")
	}
	c, err := rpc.NewClient(url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errBlockChainHTTPError, err)
	}
	log.Info("create tezos data client", "url", url)
	return &TezosDataClient{
		client: c,
	}, nil
}

func (tdc *TezosDataClient) GetTxByTxHash(txHash string) (*Operations, error) {
	return nil, nil
}

func (tdc *TezosDataClient) GetTxByAddr(ctx context.Context, address string) (*[]Operations, error) {
	var Tx []Operations
	var res []Operations
	urlPath := "/v1/accounts/" + address + "/operations"
	err := tdc.client.Get(ctx, urlPath, &Tx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errBlockChainHTTPError, err)
	}
	for _, tx := range Tx {
		if tx.Type == "transaction" {
			res = append(res, tx)
		}
	}
	return &res, nil
}

func (tdc *TezosDataClient) GetBlockNumber(ctx context.Context) (int64, error) {
	var count float64
	log.Info("get tezos block number")
	err := tdc.client.Get(ctx, "/v1/blocks/count", &count)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", errBlockChainHTTPError, err)
	}
	return int64(count), nil
}

func (tdc *TezosDataClient) GetBlockByNumber(ctx context.Context, blockHigh interface{}) (*Block, error) {
	var count interface{}
	var res Block
	if blockHigh.(int64) == 0 {
		err := tdc.client.Get(ctx, "/v1/blocks/count", &count)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", errBlockChainHTTPError, err)
		}
		count = int64(count.(float64))
	} else {
		count = blockHigh.(int64)
	}
	fmt.Println("count: ", count)
	time.Sleep(1 * time.Second)
	err := tdc.client.Get(ctx, fmt.Sprintf("/v1/blocks/%d", count), &res)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errBlockChainHTTPError, err)
	}
	fmt.Println(res)
	return &res, nil
}

func (tdc *TezosDataClient) GetBlockByHash(ctx context.Context, hash string) (*Block, error) {
	var res Block
	err := tdc.client.Get(ctx, fmt.Sprintf("/v1/blocks/%s", hash), &res)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errBlockChainHTTPError, err)
	}
	return &res, nil
}

func (tdc *TezosDataClient) GetAccountInfo(ctx context.Context, Address string) (*Account, error) {
	var res Account
	err := tdc.client.Get(ctx, fmt.Sprintf("/v1/accounts/%s", Address), &res)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errBlockChainHTTPError, err)
	}
	return &res, nil
}

func (tdc *TezosDataClient) GetBalance(ctx context.Context, Address string) (uint64, error) {
	var res Account
	err := tdc.client.Get(ctx, fmt.Sprintf("/v1/accounts/%s", Address), &res)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", errBlockChainHTTPError, err)
	}
	return res.Balance, nil
}

/*func (tdc *TezosDataClient) PostSendTx(ctx context.Context, signTransactionBytes string) (string, error) {
	urlPath := "/v1/helpers/inject"
	var hash string
	err := tdc.client.Post(ctx, urlPath, signTransactionBytes, hash)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errBlockChainHTTPError, err)
	}
	return hash, nil
}
*/
