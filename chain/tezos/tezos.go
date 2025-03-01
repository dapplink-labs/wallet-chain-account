package tezos

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/log"
	"github.com/trilitech/tzgo/tezos"
	"strconv"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

const ChainName = "Tezos"

type ChainAdaptor struct {
	tezosClient     *TezosClient
	tezosDataClient *TezosDataClient
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	tezosClient, err := NewTezosClients(conf)
	if err != nil {
		log.Error("new tezos client fail", "err", err)
		return nil, err
	}

	tezosDataClient, err := NewTezosDataClient(conf.WalletNode.Tezos.DataApiUrl)
	if err != nil {
		log.Error("new tezos data client fail", "err", err)
		return nil, err
	}
	log.Info("new tezos chain adaptor success")
	return &ChainAdaptor{
		tezosClient:     tezosClient,
		tezosDataClient: tezosDataClient,
	}, nil
}

func (c ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	return &account.SupportChainsResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

func (c ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	// 解析公钥
	pubKey, err := tezos.ParseKey(req.PublicKey)
	if err != nil {
		return &account.ConvertAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}
	address := pubKey.Address().String()
	return &account.ConvertAddressResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "convert address successs",
		Address: address,
	}, nil
}

func (c ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	_, err := tezos.ParseAddress(req.Address)
	if err != nil {
		return &account.ValidAddressResponse{
			Code:  common.ReturnCode_ERROR,
			Msg:   "convert address fail",
			Valid: false,
		}, nil
	}
	return &account.ValidAddressResponse{
		Code:  common.ReturnCode_SUCCESS,
		Msg:   "convert address successs",
		Valid: err == nil,
	}, nil
}

func (c ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	block, err := c.tezosDataClient.GetBlockByNumber(context.Background(), req.Height)
	if err != nil {
		return &account.BlockResponse{
			Code:   common.ReturnCode_ERROR,
			Msg:    "get block number fail:" + err.Error(),
			Height: 0,
		}, nil
	}
	return &account.BlockResponse{
		Code:   common.ReturnCode_SUCCESS,
		Msg:    "get block number success",
		Height: int64(block.Level),
		Hash:   block.Hash,
	}, nil
}

func (c ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	block, err := c.tezosDataClient.GetBlockByHash(context.Background(), req.Hash)
	if err != nil {
		return &account.BlockResponse{
			Code:   common.ReturnCode_ERROR,
			Msg:    "get block number fail:" + err.Error(),
			Height: 0,
		}, nil
	}
	return &account.BlockResponse{
		Code:   common.ReturnCode_SUCCESS,
		Msg:    "get block number success",
		Height: int64(block.Level),
		Hash:   block.Hash,
	}, nil
}

func (c ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	return &account.BlockHeaderResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this rpc interface",
	}, nil
}

func (c ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	return &account.BlockHeaderResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this rpc interface",
	}, nil
}

func (c ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	accountInfo, err := c.tezosDataClient.GetAccountInfo(context.Background(), req.Address)
	if err != nil {
		return &account.AccountResponse{
			Code:          common.ReturnCode_ERROR,
			Msg:           "get account fail:" + err.Error(),
			AccountNumber: "",
			Sequence:      "",
			Balance:       "",
		}, nil
	} else {
		return &account.AccountResponse{
			Code:          common.ReturnCode_SUCCESS,
			Msg:           "get account success",
			AccountNumber: strconv.FormatUint(accountInfo.Id, 10),
			Sequence:      strconv.FormatUint(accountInfo.Counter, 10),
			Balance:       strconv.FormatUint(accountInfo.Balance, 10),
		}, nil
	}
}

func (c ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	var rawTransaction []SendTransaction
	err := json.Unmarshal([]byte(req.RawTx), &rawTransaction)
	if err != nil {
		return &account.FeeResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get transaction fail:" + err.Error(),
		}, nil
	}
	fee, err := c.tezosClient.estimateFee(rawTransaction)
	if err != nil {
		return &account.FeeResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get fee fail:" + err.Error(),
		}, nil
	}
	return &account.FeeResponse{
		Code:      common.ReturnCode_SUCCESS,
		Msg:       "get fee success",
		NormalFee: fee.String(),
	}, nil
}

func (c ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	hash, err := c.tezosClient.SendSignTransaction(req.RawTx)
	if err != nil {
		return &account.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "send raw transaction fail:" + err.Error(),
		}, nil
	}
	return &account.SendTxResponse{
		Code:   common.ReturnCode_SUCCESS,
		Msg:    "send raw transaction success",
		TxHash: hash,
	}, nil
}

func (c ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	op, err := c.tezosDataClient.GetTxByAddr(context.Background(), req.Address)
	if err != nil {
		return &account.TxAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get transaction by address fail:" + err.Error(),
		}, nil
	}
	var txMs []*account.TxMessage
	for _, v := range *op {
		txMs = append(txMs, &account.TxMessage{
			Hash:  v.Hash,
			Index: uint32(v.Counter),
			Froms: []*account.Address{
				{
					Address: v.Sender.Address,
				},
			},
			Tos: []*account.Address{
				{
					Address: v.Target.Address,
				},
			},
			Values: []*account.Value{
				{
					Value: strconv.FormatUint(v.Amount, 10),
				},
			},
			Status:   account.TxStatus_Success,
			Height:   strconv.FormatUint(v.Level, 10),
			Datetime: v.TimeStamp,
			Fee:      strconv.FormatUint((v.BakerFee + v.StorageFee + v.AllocationFee), 10),
		})
	}
	return &account.TxAddressResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get transaction by address success",
		Tx:   txMs,
	}, nil
}

func (c ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	return &account.TxHashResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this rpc interface",
	}, nil
}

func (c ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	return &account.BlockByRangeResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this rpc interface",
	}, nil
}

func (c ChainAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	return &account.UnSignTransactionResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this rpc interface",
	}, nil
}

func (c ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	return &account.SignedTransactionResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this rpc interface",
	}, nil
}

func (c ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	return &account.DecodeTransactionResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this rpc interface",
	}, nil
}

func (c ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	return &account.VerifyTransactionResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this rpc interface",
	}, nil
}

func (c ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	return &account.ExtraDataResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this rpc interface",
	}, nil
}
