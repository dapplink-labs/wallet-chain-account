package sui

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/log"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/dapplink-labs/wallet-chain-account/config"
)

type SuiClient struct {
	client sui.ISuiAPI
}

func (c SuiClient) GetGasPrice() (uint64, error) {
	ctx := context.Background()
	price, err := c.client.SuiXGetReferenceGasPrice(ctx)
	if err != nil {
		log.Error("get gas price fail", "err", err)
		panic(err)
	}
	return price, nil
}

func (c SuiClient) SendTx(txStr string) (*models.TxnMetaData, error) {
	ctx := context.Background()
	var req models.PublishRequest
	jsonErr := json.Unmarshal([]byte(txStr), &req)
	if jsonErr != nil {
		return nil, jsonErr
	}
	publish, err := c.client.Publish(ctx, req)
	if err != nil {
		log.Error("publish tx  fail", "err", err)
		panic(err)
	}
	return &publish, nil
}

func (c *SuiClient) GetAccountBalance(owner, coinType string) (models.CoinBalanceResponse, error) {
	ctx := context.Background()
	if coinType == "" {
		coinType = SuiCoinType
	}
	req := models.SuiXGetBalanceRequest{
		Owner:    owner,
		CoinType: coinType,
	}
	balance, err := c.client.SuiXGetBalance(ctx, req)
	if err != nil {
		log.Error("get balance fail", "err", err)
		panic(err)
	}
	return balance, nil
}

func (c SuiClient) GetTxListByAddress(address string, cursor string, limit uint32) (models.SuiXQueryTransactionBlocksResponse, error) {
	ctx := context.Background()
	req := models.SuiXQueryTransactionBlocksRequest{
		SuiTransactionBlockResponseQuery: models.SuiTransactionBlockResponseQuery{
			TransactionFilter: models.TransactionFilter{
				"FromAddress": address,
			},
			Options: models.SuiTransactionBlockOptions{
				ShowInput:          true,
				ShowRawInput:       true,
				ShowEffects:        true,
				ShowEvents:         true,
				ShowObjectChanges:  true,
				ShowBalanceChanges: true,
			},
		},
		Cursor:          cursor,
		Limit:           uint64(limit),
		DescendingOrder: false,
	}
	txList, err := c.client.SuiXQueryTransactionBlocks(ctx, req)
	if err != nil {
		log.Error("get tx list fail", "err", err)
		panic(err)
	}
	return txList, nil
}

func (c SuiClient) GetTxDetailByDigest(digest string) (models.SuiTransactionBlockResponse, error) {
	ctx := context.Background()
	req := models.SuiGetTransactionBlockRequest{
		Digest: digest,
		Options: models.SuiTransactionBlockOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowEvents:         true,
			ShowBalanceChanges: true,
			ShowObjectChanges:  true,
		},
	}
	txDetail, err := c.client.SuiGetTransactionBlock(ctx, req)
	if err != nil {
		log.Error("get tx detail fail", "err", err)
		panic(err)
	}
	return txDetail, nil
}

func NewSuiClient(conf *config.Config) (*SuiClient, error) {
	client := sui.NewSuiClient(conf.WalletNode.Sui.RpcUrl)
	return &SuiClient{client: client}, nil
}
