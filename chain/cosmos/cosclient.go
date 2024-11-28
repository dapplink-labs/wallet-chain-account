package cosmos

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	authv1beta1 "cosmossdk.io/api/cosmos/auth/v1beta1"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/query"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"

	"github.com/dapplink-labs/chain-explorer-api/explorer/oklink"
	"github.com/dapplink-labs/wallet-chain-account/common/helpers"
	"github.com/dapplink-labs/wallet-chain-account/common/retry"
	"github.com/dapplink-labs/wallet-chain-account/config"
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

	dataClient      *oklink.ChainExplorerAdaptor
	bankClient      banktypes.QueryClient
	txServiceClient sdktx.ServiceClient
	authClient      authv1beta1.QueryClient
	bankKeeper      keeper.BaseKeeper
}

func DialCosmosClient(ctx context.Context, conf *config.Config) (*CosmosClient, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultDialTimeout)
	defer cancel()

	nodeUrl := conf.WalletNode.Cosmos.RpcUrl
	dataApiKey := conf.WalletNode.Cosmos.DataApiKey
	dataApiUrl := conf.WalletNode.Cosmos.DataApiUrl
	timeOut := conf.WalletNode.Cosmos.TimeOut
	bOff := retry.Exponential()
	connClient, err := retry.Do(ctx, defaultDialAttempts, bOff, func() (*CosmosClient, error) {
		if !helpers.IsURLAvailable(nodeUrl) {
			return nil, fmt.Errorf("address unavailable (%s)", nodeUrl)
		}

		interfaceRegistry := codectypes.NewInterfaceRegistry()
		codecMarshaler := codec.NewProtoCodec(interfaceRegistry)
		txConfig := authtx.NewTxConfig(codecMarshaler, authtx.DefaultSignModes)
		clientCtx := client.Context{}.
			WithCodec(codecMarshaler).
			WithInterfaceRegistry(interfaceRegistry).
			WithTxConfig(txConfig)

		conn, err := client.NewClientFromNode(nodeUrl)
		if err != nil {
			log.Error("failed to retry dial nodeUrl (%s): %w", nodeUrl, err)
			return nil, err
		}
		clientCtx = clientCtx.WithClient(conn)

		dataClient, err := oklink.NewChainExplorerAdaptor(dataApiKey, dataApiUrl+"/", false, time.Duration(timeOut))
		if err != nil {
			log.Error("failed cosmos  new chain explorer adaptor", "err", err)
			return nil, err
		}
		return &CosmosClient{
			context:         clientCtx,
			rpchttp:         conn,
			codec:           codecMarshaler,
			dataClient:      dataClient,
			bankClient:      banktypes.NewQueryClient(clientCtx),
			txServiceClient: sdktx.NewServiceClient(clientCtx),
			authClient:      authv1beta1.NewQueryClient(clientCtx),
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

func (c *CosmosClient) GetTxByHash(hash string) (*ctypes.ResultTx, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		log.Error("failed to get block by hash: %v", err)
		return nil, err
	}

	return c.rpchttp.Tx(ctx, hashBytes, true)
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
	pubKey := &sdk256k1.PubKey{Key: key}
	accAddress := sdk.AccAddress(pubKey.Address())
	prefix := "cosmos"
	address, _ := bech32.ConvertAndEncode(prefix, accAddress)
	return address
}

func (c *CosmosClient) BroadcastTx(txBytes []byte) (*sdktx.BroadcastTxResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	//return c.context.BroadcastTxSync(txBytes)
	//return c.rpchttp.BroadcastTxCommit(ctx, txBytes)
	return c.txServiceClient.BroadcastTx(ctx, &sdktx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    sdktx.BroadcastMode_BROADCAST_MODE_SYNC,
	})
}

func (c *CosmosClient) GetBlock(height int64) (*ctypes.ResultBlock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	return c.rpchttp.Block(ctx, &height)
}

func (c *CosmosClient) GetBlockByHash(hash string) (*ctypes.ResultBlock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	return c.rpchttp.BlockByHash(ctx, hashBytes)
}

func (c *CosmosClient) TxDecode(txData []byte) (*sdktx.TxDecodeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	request := &sdktx.TxDecodeRequest{
		TxBytes: txData,
	}
	return c.txServiceClient.TxDecode(ctx, request)
}

func (c *CosmosClient) GetHeaderByHeight(height int64) (*ctypes.ResultHeader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	return c.rpchttp.Header(ctx, &height)
}

func (c *CosmosClient) GetHeaderByHash(hash string) (*ctypes.ResultHeader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	return c.rpchttp.HeaderByHash(ctx, hashBytes)
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

func (c *CosmosClient) Tx(hash string, prove bool) (*ctypes.ResultTx, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	return c.rpchttp.Tx(ctx, hashBytes, prove)
}

func (c *CosmosClient) GetTx(hash string) (*sdktx.GetTxResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	req := &sdktx.GetTxRequest{
		Hash: hash,
	}
	return c.txServiceClient.GetTx(ctx, req)
}

func (c *CosmosClient) Close() error {
	c.rpchttp.OnStop()
	return nil
}
