package chaindispatcher

import (
	"context"
	"runtime/debug"
	"strings"

	"github.com/dapplink-labs/wallet-chain-account/chain/arbitrum"
	"github.com/dapplink-labs/wallet-chain-account/chain/binance"
	"github.com/dapplink-labs/wallet-chain-account/chain/btt"
	"github.com/dapplink-labs/wallet-chain-account/chain/linea"
	"github.com/dapplink-labs/wallet-chain-account/chain/mantle"
	"github.com/dapplink-labs/wallet-chain-account/chain/optimism"
	"github.com/dapplink-labs/wallet-chain-account/chain/polygon"
	"github.com/dapplink-labs/wallet-chain-account/chain/scroll"
	"github.com/dapplink-labs/wallet-chain-account/chain/zksync"

	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/chain/aptos"
	"github.com/dapplink-labs/wallet-chain-account/chain/cosmos"
	"github.com/dapplink-labs/wallet-chain-account/chain/ethereum"
	"github.com/dapplink-labs/wallet-chain-account/chain/solana"
	"github.com/dapplink-labs/wallet-chain-account/chain/sui"
	"github.com/dapplink-labs/wallet-chain-account/chain/ton"
	"github.com/dapplink-labs/wallet-chain-account/chain/tron"
	"github.com/dapplink-labs/wallet-chain-account/chain/xlm"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

type CommonRequest interface {
	GetChain() string
}

type CommonReply = account.SupportChainsResponse

type ChainType = string

type ChainDispatcher struct {
	registry map[ChainType]chain.IChainAdaptor
}

func New(conf *config.Config) (*ChainDispatcher, error) {
	dispatcher := ChainDispatcher{
		registry: make(map[ChainType]chain.IChainAdaptor),
	}

	chainAdaptorFactoryMap := map[string]func(conf *config.Config) (chain.IChainAdaptor, error){
		strings.ToLower(ethereum.ChainName): ethereum.NewChainAdaptor,
		strings.ToLower(cosmos.ChainName):   cosmos.NewChainAdaptor,
		strings.ToLower(solana.ChainName):   solana.NewChainAdaptor,
		strings.ToLower(tron.ChainName):     tron.NewChainAdaptor,
		strings.ToLower(aptos.ChainName):    aptos.NewChainAdaptor,
		strings.ToLower(sui.ChainName):      sui.NewSuiAdaptor,
		strings.ToLower(ton.ChainName):      ton.NewChainAdaptor,
		strings.ToLower(polygon.ChainName):  polygon.NewChainAdaptor,
		strings.ToLower(zksync.ChainName):   zksync.NewChainAdaptor,
		strings.ToLower(arbitrum.ChainName): arbitrum.NewChainAdaptor,
		strings.ToLower(binance.ChainName):  binance.NewChainAdaptor,
		strings.ToLower(mantle.ChainName):   mantle.NewChainAdaptor,
		strings.ToLower(optimism.ChainName): optimism.NewChainAdaptor,
		strings.ToLower(linea.ChainName):    linea.NewChainAdaptor,
		strings.ToLower(scroll.ChainName):   scroll.NewChainAdaptor,
		strings.ToLower(xlm.ChainName):      xlm.NewChainAdaptor,
		strings.ToLower(btt.ChainName):      btt.NewChainAdaptor,
	}

	supportedChains := []string{
		strings.ToLower(ethereum.ChainName),
		strings.ToLower(cosmos.ChainName),
		strings.ToLower(solana.ChainName),
		strings.ToLower(tron.ChainName),
		strings.ToLower(sui.ChainName),
		strings.ToLower(ton.ChainName),
		strings.ToLower(aptos.ChainName),
		strings.ToLower(polygon.ChainName),
		strings.ToLower(arbitrum.ChainName),
		strings.ToLower(binance.ChainName),
		strings.ToLower(mantle.ChainName),
		strings.ToLower(optimism.ChainName),
		strings.ToLower(linea.ChainName),
		strings.ToLower(scroll.ChainName),
		strings.ToLower(xlm.ChainName),
		strings.ToLower(btt.ChainName),
		strings.ToLower(zksync.ChainName),
	}

	for _, c := range conf.Chains {
		chainName := strings.ToLower(c)
		if factory, ok := chainAdaptorFactoryMap[chainName]; ok {
			adaptor, err := factory(conf)
			if err != nil {
				log.Crit("failed to setup chain", "chain", chainName, "error", err)
			}
			dispatcher.registry[chainName] = adaptor
		} else {
			log.Error("unsupported chain", "chain", chainName, "supportedChains", supportedChains)
		}
	}
	return &dispatcher, nil
}

func (d *ChainDispatcher) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("panic error", "msg", e)
			log.Debug(string(debug.Stack()))
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	pos := strings.LastIndex(info.FullMethod, "/")
	method := info.FullMethod[pos+1:]

	chainName := req.(CommonRequest).GetChain()
	log.Info(method, "chain", chainName, "req", req)

	resp, err = handler(ctx, req)
	log.Debug("Finish handling", "resp", resp, "err", err)
	return
}

func (d *ChainDispatcher) preHandler(req interface{}) (resp *CommonReply, chainName string) {
	chainName = strings.ToLower(req.(CommonRequest).GetChain())
	log.Debug("chain", chainName, "req", req)
	if _, ok := d.registry[chainName]; !ok {
		return &CommonReply{
			Code:    common.ReturnCode_ERROR,
			Msg:     config.UnsupportedOperation,
			Support: false,
		}, chainName
	}
	return nil, chainName
}

func (d *ChainDispatcher) GetSupportChains(ctx context.Context, request *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.SupportChainsResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  config.UnsupportedOperation,
		}, nil
	}
	return d.registry[chainName].GetSupportChains(request)
}

func (d *ChainDispatcher) ConvertAddress(ctx context.Context, request *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.ConvertAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "covert address fail at pre handle",
		}, nil
	}
	return d.registry[chainName].ConvertAddress(request)
}

func (d *ChainDispatcher) ValidAddress(ctx context.Context, request *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.ValidAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "valid address error at pre handle",
		}, nil
	}
	return d.registry[chainName].ValidAddress(request)
}

func (d *ChainDispatcher) GetBlockByNumber(ctx context.Context, request *account.BlockNumberRequest) (*account.BlockResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block by number fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetBlockByNumber(request)
}

func (d *ChainDispatcher) GetBlockByHash(ctx context.Context, request *account.BlockHashRequest) (*account.BlockResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block by hash fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetBlockByHash(request)
}

func (d *ChainDispatcher) GetBlockHeaderByHash(ctx context.Context, request *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block header by hash fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetBlockHeaderByHash(request)
}

func (d *ChainDispatcher) GetBlockHeaderByNumber(ctx context.Context, request *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block header by number fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetBlockHeaderByNumber(request)
}

func (d *ChainDispatcher) GetBlockHeaderByRange(ctx context.Context, request *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.BlockByRangeResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get block range header fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetBlockByRange(request)
}

func (d *ChainDispatcher) GetAccount(ctx context.Context, request *account.AccountRequest) (*account.AccountResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.AccountResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get account information fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetAccount(request)
}

func (d *ChainDispatcher) GetFee(ctx context.Context, request *account.FeeRequest) (*account.FeeResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.FeeResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get fee fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetFee(request)
}

func (d *ChainDispatcher) SendTx(ctx context.Context, request *account.SendTxRequest) (*account.SendTxResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "send tx fail at pre handle",
		}, nil
	}
	return d.registry[chainName].SendTx(request)
}

func (d *ChainDispatcher) GetTxByAddress(ctx context.Context, request *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get tx by address fail pre handle",
		}, nil
	}
	return d.registry[chainName].GetTxByAddress(request)
}

func (d *ChainDispatcher) GetTxByHash(ctx context.Context, request *account.TxHashRequest) (*account.TxHashResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get tx by hash fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetTxByHash(request)
}

func (d *ChainDispatcher) GetBlockByRange(ctx context.Context, request *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.BlockByRangeResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get blcok by range fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetBlockByRange(request)
}

func (d *ChainDispatcher) BuildUnSignTransaction(ctx context.Context, request *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.UnSignTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get un sign tx fail at pre handle",
		}, nil
	}
	return d.registry[chainName].BuildUnSignTransaction(request)
}

func (d *ChainDispatcher) BuildSignedTransaction(ctx context.Context, request *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.SignedTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "signed tx fail at pre handle",
		}, nil
	}
	return d.registry[chainName].BuildSignedTransaction(request)
}

func (d *ChainDispatcher) DecodeTransaction(ctx context.Context, request *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.DecodeTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "decode tx fail at pre handle",
		}, nil
	}
	return d.registry[chainName].DecodeTransaction(request)
}

func (d *ChainDispatcher) VerifySignedTransaction(ctx context.Context, request *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.VerifyTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "verify tx fail at pre handle",
		}, nil
	}
	return d.registry[chainName].VerifySignedTransaction(request)
}

func (d *ChainDispatcher) GetExtraData(ctx context.Context, request *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	resp, chainName := d.preHandler(request)
	if resp != nil {
		return &account.ExtraDataResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get extra data fail at pre handle",
		}, nil
	}
	return d.registry[chainName].GetExtraData(request)
}

func (d *ChainDispatcher) GetNftListByAddress(ctx context.Context, request *account.NftAddressRequest) (*account.NftAddressResponse, error) {
	panic("implement me")
}

func (d *ChainDispatcher) GetNftCollection(ctx context.Context, request *account.NftCollectionRequest) (*account.NftCollectionResponse, error) {
	panic("implement me")
}

func (d *ChainDispatcher) GetNftDetail(ctx context.Context, request *account.NftDetailRequest) (*account.NftDetailResponse, error) {
	panic("implement me")
}

func (d *ChainDispatcher) GetNftHolderList(ctx context.Context, request *account.NftHolderListRequest) (*account.NftHolderListResponse, error) {
	panic("implement me")
}

func (d *ChainDispatcher) GetNftTradeHistory(ctx context.Context, request *account.NftTradeHistoryRequest) (*account.NftTradeHistoryResponse, error) {
	panic("implement me")
}

func (d *ChainDispatcher) GetAddressNftTradeHistory(ctx context.Context, request *account.AddressNftTradeHistoryRequest) (*account.AddressNftTradeHistoryResponse, error) {
	panic("implement me")
}
