package tron

import (
	"encoding/hex"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"strconv"
	"strings"
)

const (
	ChainName     = "Tron"
	AddressPrefix = "41"
	TronSymbol    = "TRX"
)

type ChainAdaptor struct {
	tronClient *TronClient
}

// GetSupportChains 返回是否支持链
func (c ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	return &account.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

// ConvertAddress 将公钥转换为地址
func (c ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	address, err := ComputeAddress(req.PublicKey)
	if err != nil {
		return &account.ConvertAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	} else {
		return &account.ConvertAddressResponse{
			Code:    common2.ReturnCode_SUCCESS,
			Msg:     "convert address successs",
			Address: address,
		}, nil
	}
}

// ValidAddress 验证地址
func (c ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	// 1. 检查地址长度，TRON 地址应为 42 个字符（16 进制表示）并以 "41" 开头
	if len(req.Address) != 42 || !strings.HasPrefix(req.Address, AddressPrefix) {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_SUCCESS,
			Msg:   "convert address successs",
			Valid: false,
		}, nil
	}
	// 2. 检查地址是否为有效的十六进制字符串
	_, err := hex.DecodeString(req.Address)
	if err != nil {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_SUCCESS,
			Msg:   "convert address successs",
			Valid: false,
		}, nil
	}
	return &account.ValidAddressResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "convert address successs",
		Valid: true,
	}, nil
}

// GetBlockByNumber 根据区块高度获取区块信息
func (c ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	block, err := c.tronClient.GetBlockByNumber(req.Height)
	if err != nil {
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}
	txListRet := make([]*account.BlockInfoTransactionList, 0, len(block.Transactions))
	for _, v := range block.Transactions {
		bitlItem := &account.BlockInfoTransactionList{
			Hash: v,
		}
		txListRet = append(txListRet, bitlItem)
	}
	return &account.BlockResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "block by number success",
		Hash:         block.Hash,
		BaseFee:      block.BaseFeePerGas,
		Transactions: txListRet,
	}, nil
}

// GetBlockByHash 根据区块hash获取区块信息
func (c ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	block, err := c.tronClient.GetBlockByHash(req.Hash)
	if err != nil {
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}
	txListRet := make([]*account.BlockInfoTransactionList, 0, len(block.Transactions))
	for _, v := range block.Transactions {
		bitlItem := &account.BlockInfoTransactionList{
			Hash: v,
		}
		txListRet = append(txListRet, bitlItem)
	}
	return &account.BlockResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "block by number success",
		Hash:         block.Hash,
		BaseFee:      block.BaseFeePerGas,
		Transactions: txListRet,
	}, nil
}

// GetBlockHeaderByHash 根据区块hash获取区块头信息
func (c ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	return &account.BlockHeaderResponse{
		Code: common2.ReturnCode_ERROR,
		Msg:  "不支持啦！",
	}, nil
}

// GetBlockHeaderByNumber 根据区块高度获取区块头信息
func (c ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	return &account.BlockHeaderResponse{
		Code: common2.ReturnCode_ERROR,
		Msg:  "不支持啦！",
	}, nil
}

// GetAccount 根据地址获取账户信息
func (c ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	info, err := c.tronClient.GetAccount(req.Address)
	if err != nil {
		return &account.AccountResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}

	if req.Coin == TronSymbol {
		return &account.AccountResponse{
			Code:          common2.ReturnCode_SUCCESS,
			Msg:           "get account response success",
			AccountNumber: "0",
			//Sequence:      info.Nonce,
			Balance: strconv.Itoa(info.Balance),
		}, nil
	} else {
		for _, v := range info.AssetV2 {
			if v.Key == req.Coin {
				return &account.AccountResponse{
					Code:          common2.ReturnCode_SUCCESS,
					Msg:           "get account response success",
					AccountNumber: "0",
					//Sequence:      info.Nonce,
					Balance: strconv.FormatInt(v.Value, 10),
				}, nil
			}
		}
		return &account.AccountResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "不支持的币种",
		}, nil
	}
}

// GetFee 获取手续费
func (c ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {

}

func (c ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	tronClient := DialTronClient(conf.WalletNode.Tron.RPCs[0].RPCURL)

	return &ChainAdaptor{
		tronClient: tronClient,
	}, nil
}
