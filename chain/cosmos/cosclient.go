package cosmos

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	authv1beta1 "cosmossdk.io/api/cosmos/auth/v1beta1"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	ed255192 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/dapplink-labs/wallet-chain-account/common/retry"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
)

const (
	defaultDialTimeout    = 30 * time.Second
	defaultDialAttempts   = 5
	defaultRequestTimeout = 10 * time.Second
)

type CosmosClient struct {
	context client.Context
	rpchttp *rpchttp.HTTP
	codec   *codec.ProtoCodec

	grpcConn        *grpc.ClientConn
	bankClient      banktypes.QueryClient
	txServiceClient sdktx.ServiceClient
	authClient      authv1beta1.QueryClient
	bankKeeper      keeper.BaseKeeper
}

type BlockResponse struct {
	BlockId struct {
		Hash string `json:"hash"`
	} `json:"block_id"`
	Block struct {
		Header struct {
			Height string `json:"height"`
			Time   string `json:"time"`
		} `json:"header"`
		Data struct {
			Txs []string `json:"txs"`
		} `json:"data"`
	} `json:"block"`
}

type DecodeTxRequest struct {
	Tx string `json:"tx_bytes"`
}

type DecodeTxResponse struct {
	Tx struct {
		Body struct {
			Messages []banktypes.MsgSend `json:"messages"`
		} `json:"body"`
	} `json:"tx"`
	Response struct {
		Height    string `json:"height"`
		Txhash    string `json:"txhash"`
		GasWanted string `json:"gas_wanted"`
		GasUsed   string `json:"gas_used"`
	} `json:"tx_response"`
}

type TxResponse struct {
	//Tx struct {
	//	Body struct {
	//		Messages []string `json:"messages"`
	//	} `json:"body"`
	//} `json:"tx"`
	Response struct {
		Height    string     `json:"height"`
		Txhash    string     `json:"txhash"`
		GasWanted string     `json:"gas_wanted"`
		GasUsed   string     `json:"gas_used"`
		Events    []*TxEvent `json:"events"`
		Timestamp string     `json:"timestamp"`
	} `json:"tx_response"`
}

type TxEvent struct {
	Type       string              `json:"type"`
	Attributes []*TxEventAttribute `json:"attributes"`
}

type TxEventAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func DialCosmosClient(ctx context.Context, nodeUrl string) (*CosmosClient, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultDialTimeout)
	defer cancel()

	bOff := retry.Exponential()
	connClient, err := retry.Do(ctx, defaultDialAttempts, bOff, func() (*CosmosClient, error) {
		//grpcConn, err := grpc.Dial(nodeUrl, grpc.WithInsecure())
		//if err != nil {
		//	log.Error("failed to retry dial nodeUrl (%s): %w", nodeUrl, err)
		//	return nil, err
		//}
		//return grpcConn, nil

		//if !IsURLAvailable(rpcUrl) {
		//	return nil, fmt.Errorf("address unavailable (%s)", rpcUrl)
		//}
		//
		ctx := client.Context{}
		conn, err := client.NewClientFromNode(nodeUrl)
		if err != nil {
			log.Error("failed to retry dial nodeUrl (%s): %w", nodeUrl, err)
			return nil, err
		}
		codec := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
		ctx = ctx.WithClient(conn).WithKeyring(keyring.NewInMemory(codec))
		return &CosmosClient{
			context:         ctx,
			rpchttp:         conn,
			codec:           codec,
			bankClient:      banktypes.NewQueryClient(ctx),
			txServiceClient: sdktx.NewServiceClient(ctx),
			authClient:      authv1beta1.NewQueryClient(ctx),
		}, nil

	})
	if err != nil {
		log.Error("failed to dial nodeUrl (%s): %w", nodeUrl, err)
		return nil, err
	}

	return connClient, nil
}

func (c *CosmosClient) GetBalance(coin, addr string) (*sdk.Coin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	address, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, err
	}

	resp, err := c.bankClient.Balance(ctx, &banktypes.QueryBalanceRequest{
		Address: address.String(),
		Denom:   coin,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetBalance(), nil
}

func (c *CosmosClient) GetTxByHash(restURL, hash string) (*TxResponse, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	//defer cancel()

	//response, err := c.txServiceClient.GetTx(ctx, &sdktx.GetTxRequest{Hash: hash})
	//if err != nil {
	//	log.Error("failed to get block: %v", err)
	//	return nil, err
	//}

	resp, err := http.Get(fmt.Sprintf("%s/cosmos/tx/v1beta1/txs/%s", restURL, hash))
	if err != nil {
		log.Error("failed to get block: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read response body: %v", err)
		return nil, err
	}

	var txResponse TxResponse
	err = json.Unmarshal(body, &txResponse)
	if err != nil {
		log.Error("failed to unmarshal response: %v", err)
		return nil, err
	}

	return &txResponse, nil
}

func (c *CosmosClient) GetTxByEvent(event []string, page, limit uint64) (*sdktx.GetTxsEventResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 1000
	}

	if len(event) == 0 {
		return nil, errors.New("event cannot be empty")
	}
	eventTmp := make([]string, 0)
	for _, v := range event {
		_ = v
		eventTmp = append(eventTmp, fmt.Sprintf("message.sender='%s'", "cosmos1qd99t24whd3hfg22r53x8uw9ps3rrctwxqvn4m"))
	}

	tmp := &sdktx.GetTxsEventRequest{
		Events: eventTmp,
		Pagination: &query.PageRequest{
			Offset: page,
			Limit:  limit,
		},
	}

	return c.txServiceClient.GetTxsEvent(ctx, tmp)
}

func (c *CosmosClient) GetAccount(addr string) (*authv1beta1.QueryAccountResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	return c.authClient.Account(ctx, &authv1beta1.QueryAccountRequest{Address: addr})
}

func (c *CosmosClient) SendTx(fromAddr, toAddr, coin string, amount int64) (*banktypes.MsgSendResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()
	newCtx := sdk.UnwrapSDKContext(ctx)

	fromAddress, err := sdk.AccAddressFromBech32(fromAddr)
	if err != nil {
		return nil, err
	}
	toAddress, err := sdk.AccAddressFromBech32(toAddr)
	if err != nil {
		return nil, err
	}

	if c.bankKeeper.BlockedAddr(toAddress) {
		return nil, errors.New(fmt.Sprintf("%s is not allowed to receive funds", toAddress))
	}

	msg := banktypes.NewMsgSend(fromAddress, toAddress, sdk.NewCoins(sdk.NewInt64Coin(coin, amount)))
	if err := c.bankKeeper.SendCoins(newCtx, fromAddress, toAddress, msg.Amount); err != nil {
		return nil, err
	}

	//defer func() {
	//	for _, a := range msg.Amount {
	//		if a.Amount.IsInt64() {
	//			telemetry.SetGaugeWithLabels(
	//				[]string{"tx", "msg", "send"},
	//				float32(a.Amount.Int64()),
	//				[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
	//			)
	//		}
	//	}
	//}()

	return &banktypes.MsgSendResponse{}, nil
}

func (c *CosmosClient) GetAddressFromPubKey(key []byte) string {
	// todo check ed255192
	pub := ed255192.PubKey{Key: key}
	return pub.Address().String()
}

func (c *CosmosClient) BroadcastTx(txByte []byte) (*sdktx.BroadcastTxResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	return c.txServiceClient.BroadcastTx(ctx, &sdktx.BroadcastTxRequest{
		TxBytes: txByte,
		Mode:    sdktx.BroadcastMode_BROADCAST_MODE_SYNC,
	})
}

// https://cosmos-rest.publicnode.com/cosmos/tx/v1beta1/txs/block/22879895
func (c *CosmosClient) GetBlock(restURL string, height int64) (*BlockResponse, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	//defer cancel()
	//
	//txsRequest := &sdktx.GetBlockWithTxsRequest{
	//	Height: height,
	//}
	//result, err := c.txServiceClient.GetBlockWithTxs(ctx, txsRequest)
	//if err != nil {
	//	return nil, err
	//}
	//decodedHash, err := base64.StdEncoding.DecodeString(string(result.Block.Data.Txs[0]))
	//if err != nil {
	//	log.Error("解码失败: %v", err)
	//}
	//fmt.Printf("解码前的二进制数据: %x\n", string(result.Block.Data.Txs[0]))
	//fmt.Printf("解码后的二进制数据: %x\n", string(decodedHash))
	//
	//return nil, nil

	resp, err := http.Get(fmt.Sprintf("%s/cosmos/tx/v1beta1/txs/block/%d", restURL, height))
	if err != nil {
		log.Error("failed to get block: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read response body: %v", err)
		return nil, err
	}

	var blockResp BlockResponse
	err = json.Unmarshal(body, &blockResp)
	if err != nil {
		log.Error("failed to unmarshal response: %v", err)
		return nil, err
	}
	decodedHash, err := base64.StdEncoding.DecodeString(blockResp.BlockId.Hash)
	if err != nil {
		log.Error("解码失败: %v", err)
	}
	fmt.Printf("解码前的二进制数据: %x\n", blockResp.BlockId.Hash)
	fmt.Printf("解码后的二进制数据: %x\n", string(decodedHash))
	blockResp.BlockId.Hash = fmt.Sprintf("%x", string(decodedHash))
	return &blockResp, nil
}

func (c *CosmosClient) DecodeBlockTx(restURL string, block *BlockResponse) ([]*account.BlockInfoTransactionList, error) {
	var blockTransactions []*account.BlockInfoTransactionList

	for _, txData := range block.Block.Data.Txs {
		// 创建解码交易请求
		decodeTxReq := DecodeTxRequest{Tx: txData}
		jsonReq, err := json.Marshal(decodeTxReq)
		if err != nil {
			log.Error("failed to marshal decode tx request: %v", err)
			return nil, nil
		}

		// 发送 POST 请求到 tx/v1beta1/decode 端点
		resp, err := http.Post(restURL+"/cosmos/tx/v1beta1/decode", "application/json", bytes.NewBuffer(jsonReq))
		if err != nil {
			log.Error("failed to send decode tx request: %v", err)
			return nil, nil
		}
		// 读取响应
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error("failed to read decode tx response: %v", err)
			return nil, nil
		}
		resp.Body.Close()
		var decodeTxResp DecodeTxResponse
		err = json.Unmarshal(body, &decodeTxResp)
		if err != nil {
			log.Error("failed to unmarshal decode tx response: %v", err)
			return nil, nil
		}
		for _, msg := range decodeTxResp.Tx.Body.Messages {
			log.Info("Sender: %s", msg.FromAddress)
			log.Info("Receiver: %s", msg.ToAddress)
			log.Info("Amount: %s", msg.Amount.String())

			blockTransaction := &account.BlockInfoTransactionList{
				From:   msg.FromAddress,
				To:     msg.ToAddress,
				Hash:   decodeTxResp.Response.Txhash,
				Amount: msg.Amount.String(),
			}
			blockTransactions = append(blockTransactions, blockTransaction)
		}
	}
	return blockTransactions, nil
}

func (c *CosmosClient) GetHeaderByHeight(height *int64) (*ctypes.ResultHeader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	result, err := c.rpchttp.Header(ctx, height)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CosmosClient) GetHeaderByHash(hash []byte) (*ctypes.ResultHeader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	result, err := c.rpchttp.HeaderByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CosmosClient) GetBlockByHash(hash []byte) (*ctypes.ResultBlock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	result, err := c.rpchttp.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CosmosClient) BlockchainInfo(minHeight, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	result, err := c.rpchttp.BlockchainInfo(ctx, minHeight, maxHeight)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CosmosClient) Tx(hash []byte, prove bool) (*ctypes.ResultTx, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	result, err := c.rpchttp.Tx(ctx, hash, prove)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CosmosClient) Close() error {
	//c.rpcHttp.OnStop()
	//return c.grpcConn.Close()
	return nil
}
