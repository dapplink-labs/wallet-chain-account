package tezos

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/log"
	gresty "github.com/go-resty/resty/v2"
	"github.com/pkg/errors"

	"github.com/dapplink-labs/wallet-chain-account/config"
)

type TezosClient struct {
	client *gresty.Client
}

func NewTezosClients(conf *config.Config) (*TezosClient, error) {
	url := conf.WalletNode.Tezos.RpcUrl
	client := gresty.New()
	client.SetBaseURL(url)
	client.SetDebug(true)
	client.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			baseUrl := r.Request.URL
			body := r.Request.Body
			return fmt.Errorf("%d cannot %s %s %s: %w", statusCode, method, baseUrl, body, errBlockChainHTTPError)
		}
		return nil
	})
	log.Info("init tezos client", "url", url)
	return &TezosClient{
		client: client,
	}, nil
}

func (tc *TezosClient) GetNonce(address string) (*uint64, error) {
	res, err := tc.client.R().
		SetHeader("Content-Type", "application/json").
		Post("/chains/main/blocks/head/context/contracts/" + address + "/counter")
	if err != nil {
		return nil, err
	}
	nonce, ok := res.Result().(uint64)
	if !ok {
		return nil, errors.New("get nonce fail, ok is false")
	}
	return &nonce, nil
}

func (tc *TezosClient) estimateFee(tx []SendTransaction) (*big.Int, error) {
	var prTx []PreApplyTransaction
	_, err := tc.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(&tx).
		SetResult(&prTx).
		Post("/chains/main/blocks/head/helpers/preapply/operations")
	if err != nil {
		return nil, err
	}
	bigInt := new(big.Int)
	bigInt, success := bigInt.SetString(prTx[0].Contents[0].Metadata.OperationResult.ConsumedMilligas, 10)
	if !success {
		return nil, errors.New("string convert BigInt fail")
	}
	return bigInt, nil
}

func (tc *TezosClient) SendSignTransaction(boc string) (string, error) {
	var hash string
	_, err := tc.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`"%s"`, boc)).
		SetResult(&hash).
		Post("/injection/operation")
	if err != nil {
		return "", err
	}
	/*hash, ok := res.Result().(string)
	if !ok {
		return "", errors.New("send signed transaction fail, ok is false")
	}*/
	return hash, nil
}
