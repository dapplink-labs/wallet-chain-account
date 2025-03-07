package xlm

import (
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
)

const ChainName = "Xlm"

type ChainAdaptor struct {
	xlmClient *XlmClient
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	xlmClient, err := NewXlmClients(conf)
	if err != nil {
		log.Error("new xlm client fail", "err", err)
		return nil, err
	}

	return &ChainAdaptor{
		xlmClient: xlmClient,
	}, nil
}

func (c *ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	return &account.SupportChainsResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

func (c *ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	return &account.ConvertAddressResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "Do not support this rpc interface",
	}, nil
}

func (c *ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	resp, err := c.xlmClient.ValidAddress(req)
	if err != nil {
		return &account.ValidAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "ValidAddress Failed",
		}, err
	}

	return resp, err
}

func (c *ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	resp, err := c.xlmClient.GetAccountInfo(req.Address)
	if err != nil {
		return &account.AccountResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "GetAccount Failed",
		}, err
	}

	return resp, err
}

func (c *ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	resp, err := c.xlmClient.GetFee()
	if err != nil {
		return &account.FeeResponse{
			Code: common.ReturnCode_SUCCESS,
			Msg:  "GetFee Failed",
		}, err
	}

	return resp, err
}

func (c *ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	return &account.TxAddressResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "Do not support this rpc interface",
	}, nil
}

func (c *ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	resp, err := c.xlmClient.GetTransactionByHash(req.Hash)
	if err != nil {
		return &account.TxHashResponse{
			Code: common.ReturnCode_SUCCESS,
			Msg:  "GetTxByHash Failed",
		}, err
	}

	return resp, err
}

func (c *ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	blockNumber := big.NewInt(req.Height)
	resp, err := c.xlmClient.GetBlockByNumber(blockNumber)
	if err != nil {
		log.Error("GetBlockByNumber fail:", err)
		return &account.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	return resp, err
}

func (c *ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	return &account.BlockByRangeResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "Do not support this rpc interface",
	}, nil
}

func (c *ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	return &account.BlockResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "Do not support this rpc interface",
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	blockNumber := big.NewInt(req.Height)
	resp, err := c.xlmClient.GetBlockHeaderByNumber(blockNumber)
	if err != nil {
		log.Error("GetBlockHeaderByHash fail:", err)
		return &account.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	return resp, err
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	return &account.BlockHeaderResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "Do not support this rpc interface",
	}, nil
}

func (c *ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	resp, err := c.xlmClient.SendTx(req.RawTx)
	if err != nil {
		log.Error("SendTx fail:", err)
		return &account.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	return resp, err
}

func (c *ChainAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	resp, err := c.xlmClient.CreateUnsignTransaction(req)
	if err != nil {
		log.Error("CreateUnSignTransaction fail:", err)
		return &account.UnSignTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	return resp, err
}

func (c *ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	resp, err := c.xlmClient.SignedTransaction(req)
	if err != nil {
		log.Error("BuildSignedTransaction fail:", err)
		return &account.SignedTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	return resp, err
}

func (c *ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	return &account.DecodeTransactionResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "Do not support this rpc interface",
	}, nil
}

func (c *ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	return &account.VerifyTransactionResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "Do not support this rpc interface",
	}, nil
}

func (c *ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	return &account.ExtraDataResponse{
		Code:  common.ReturnCode_ERROR,
		Msg:   "Do not support this api",
		Value: req.Chain,
	}, nil
}
