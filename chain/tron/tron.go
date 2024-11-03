package tron

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	account2 "github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"github.com/ethereum/go-ethereum/log"
	"strconv"
	"strings"
	"time"
)

const (
	ChainName     = "Tron"
	AddressPrefix = "41"
	TronSymbol    = "TRX"
)

type ChainAdaptor struct {
	tronClient     *TronClient
	tronDataClient *TronData
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	tronClient := DialTronClient(conf.WalletNode.Tron.RPCs[0].RPCURL)
	tronDataClient, err := NewTronDataClient(conf.WalletNode.Tron.DataApiUrl, conf.WalletNode.Tron.DataApiKey, time.Second*15)
	if err != nil {
		return nil, err
	}
	return &ChainAdaptor{
		tronClient:     tronClient,
		tronDataClient: tronDataClient,
	}, nil
}

// GetSupportChains 返回是否支持链
func (c *ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	return &account.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

// ConvertAddress 将公钥转换为地址
func (c *ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	bytes, err := hex.DecodeString(req.PublicKey)
	if err != nil {
		return &account.ConvertAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}
	address, err := ComputeAddress(bytes)
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
func (c *ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
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
func (c *ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
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
func (c *ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
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
func (c *ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	return &account.BlockHeaderResponse{
		Code: common2.ReturnCode_ERROR,
		Msg:  "不支持啦！",
	}, nil
}

// GetBlockHeaderByNumber 根据区块高度获取区块头信息
func (c *ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	return &account.BlockHeaderResponse{
		Code: common2.ReturnCode_ERROR,
		Msg:  "不支持啦！",
	}, nil
}

// GetAccount 根据地址获取账户信息
func (c *ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
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
func (c *ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	gasPrice, err := c.tronDataClient.GetEstimateGasFee()
	if err != nil {
		return &account.FeeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}
	return &account.FeeResponse{
		Code:      common2.ReturnCode_SUCCESS,
		Msg:       "get gas price success",
		SlowFee:   gasPrice.SlowGasPrice,
		NormalFee: gasPrice.StandardGasPrice,
		FastFee:   gasPrice.BestTransactionFee,
	}, nil
}

// SendTx 发送交易
func (c *ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	_, err := c.tronClient.BroadcastTransaction(req.RawTx)
	if err != nil {
		return &account.SendTxResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Send tx error" + err.Error(),
		}, err
	}
	return &account.SendTxResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "send tx success",
	}, nil
}

// GetTxByAddress 根据地址获取交易
func (c *ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	var resp *account2.TransactionResponse[account2.AccountTxResponse]
	var err error
	if req.ContractAddress != "0x00" {
		resp, err = c.tronDataClient.GetTxByAddress(uint64(req.Page), uint64(req.Pagesize), req.Address, "token")
		if err != nil {
			return nil, err
		}
	} else {
		resp, err = c.tronDataClient.GetTxByAddress(uint64(req.Page), uint64(req.Pagesize), req.Address, "normal")
		if err != nil {
			return nil, err
		}
	}
	txs := resp.TransactionList
	list := make([]*account.TxMessage, 0, len(txs))
	for i := 0; i < len(txs); i++ {
		list = append(list, &account.TxMessage{
			Hash:   txs[i].TxId,
			Tos:    []*account.Address{{Address: txs[i].To}},
			Froms:  []*account.Address{{Address: txs[i].From}},
			Fee:    txs[i].TxId,
			Status: account.TxStatus_Success,
			Values: []*account.Value{{Value: txs[i].Amount}},
			Type:   1,
			Height: txs[i].Height,
		})
	}
	return &account.TxAddressResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get transactions by address success",
		Tx:   list,
	}, err
}

// GetTxByHash 根据交易hash获取交易
func (c *ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	resp, err := c.tronClient.GetTransactionByID(req.Hash)
	if err != nil {
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}
	fromAddress := resp.RawData.Contract[0].Parameter.Value.OwnerAddress
	toAddress := resp.RawData.Contract[0].Parameter.Value.ToAddress
	amount := resp.RawData.Contract[0].Parameter.Value.Amount

	return &account.TxHashResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get transactions by address success",
		Tx: &account.TxMessage{
			Hash:   resp.TxID,
			Tos:    []*account.Address{{Address: toAddress}},
			Froms:  []*account.Address{{Address: fromAddress}},
			Fee:    "",
			Status: account.TxStatus_Success,
			Values: []*account.Value{{Value: strconv.Itoa(amount)}},
			Type:   1,
			Height: "0",
		},
	}, nil
}

// GetBlockByRange 根据区块范围获取区块
func (c *ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	//TODO implement me
	panic("implement me")
}

// CreateUnSignTransaction 创建未签名的交易
func (c *ChainAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	jsonBytes, err := base64.StdEncoding.DecodeString(req.Base64Tx)
	if err != nil {
		log.Error("decode string fail", "err", err)
		return nil, err
	}
	var data TxStructure
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		log.Error("parse json fail", "err", err)
		return nil, err
	}
	var transaction *UnSignTransaction
	if data.ContractAddress == "" {
		transaction, err = c.tronClient.CreateTRXTransaction(data.FromAddress, data.ToAddress, data.Value)
	} else {
		transaction, err = c.tronClient.CreateTRC20Transaction(data.FromAddress, data.ToAddress, data.ContractAddress, data.Value)
	}
	if err != nil {
		return nil, err
	}
	return &account.UnSignTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "create un sign tx success",
		UnSignTx: transaction.RawDataHex,
	}, nil
}

// BuildSignedTransaction 创建签名的交易
func (c *ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

// DecodeTransaction 解码交易
func (c *ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

// VerifySignedTransaction 验证签名
func (c *ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

// GetExtraData 获取额外数据
func (c *ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	//TODO implement me
	panic("implement me")
}
