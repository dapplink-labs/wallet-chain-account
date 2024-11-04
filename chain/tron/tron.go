package tron

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	account2 "github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"math/big"
	"strconv"
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
	// 将十六进制字符串解码为字节数组
	pubKeyBytes, err := hex.DecodeString(req.PublicKey)
	// 将字节数组解析为公钥
	pubKey, _ := btcec.ParsePubKey(pubKeyBytes)
	// 将公钥转换为地址
	addr := address.PubkeyToAddress(*pubKey.ToECDSA())

	if err != nil {
		return &account.ConvertAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	} else {
		return &account.ConvertAddressResponse{
			Code:    common2.ReturnCode_SUCCESS,
			Msg:     "convert address successs",
			Address: addr.String(),
		}, nil
	}
}

// ValidAddress 验证地址
func (c *ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	_, err := address.Base58ToAddress(req.Address)
	if err != nil {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_ERROR,
			Msg:   "convert address error",
			Valid: false,
		}, nil
	}
	return &account.ValidAddressResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "convert address success",
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
	height, err := strconv.ParseInt(block.Number[2:], 16, 64)
	if err != nil {
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}
	if req.ViewTx {
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
			Height:       height,
			Transactions: txListRet,
		}, nil
	} else {
		return &account.BlockResponse{
			Code:    common2.ReturnCode_SUCCESS,
			Msg:     "block by number success",
			Hash:    block.Hash,
			BaseFee: block.BaseFeePerGas,
			Height:  height,
		}, nil
	}
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
	height, err := strconv.ParseInt(block.Number[2:], 16, 64)
	if err != nil {
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}
	if req.ViewTx {
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
			Height:       height,
			Transactions: txListRet,
		}, nil
	} else {

		return &account.BlockResponse{
			Code:    common2.ReturnCode_SUCCESS,
			Msg:     "block by number success",
			Hash:    block.Hash,
			BaseFee: block.BaseFeePerGas,
			Height:  height,
		}, nil
	}
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

	if req.Coin == TronSymbol || req.Coin == "" {
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
	jsonBytes, err := base64.StdEncoding.DecodeString(req.RawTx)
	if err != nil {
		log.Error("decode string fail", "err", err)
		return nil, err
	}
	var data SendTxReq
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		log.Error("parse json fail", "err", err)
		return nil, err
	}
	tx, err := c.tronClient.BroadcastTransaction(&data)
	if err != nil {
		return &account.SendTxResponse{
			Code:   common2.ReturnCode_ERROR,
			Msg:    "Send tx error" + err.Error(),
			TxHash: tx.Txid,
		}, err
	}
	return &account.SendTxResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "send tx success",
		TxHash: tx.Txid,
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
	// 调用 tronClient 的 GetTransactionByID 方法获取交易
	resp, err := c.tronClient.GetTransactionByID(req.Hash)
	if err != nil {
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  fmt.Sprintf("failed to get transaction: %v", err),
		}, nil
	}
	// 检查 resp 是否为 nil
	if resp == nil {
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "transaction response is nil",
		}, nil
	}
	// 检查 resp.RawData 是否为 nil
	if len(resp.RawData.Contract) == 0 {
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "transaction raw data or contract is nil",
		}, nil
	}
	// 获取交易的基本信息
	fromAddress := resp.RawData.Contract[0].Parameter.Value.OwnerAddress
	toAddress := resp.RawData.Contract[0].Parameter.Value.ToAddress
	amount := resp.RawData.Contract[0].Parameter.Value.Amount

	// 返回成功的交易信息
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
			//Height: "0",
		},
	}, nil
}

// GetBlockByRange 根据区块范围获取区块
func (c *ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	// 将起始区块和结束区块转换为 big.Int
	startBlock := new(big.Int)
	endBlock := new(big.Int)
	if _, ok := startBlock.SetString(req.Start, 10); !ok {
		return nil, fmt.Errorf("invalid start block number: %s", req.Start)
	}
	if _, ok := endBlock.SetString(req.End, 10); !ok {
		return nil, fmt.Errorf("invalid end block number: %s", req.End)
	}

	// 确保起始区块号小于等于结束区块号
	if startBlock.Cmp(endBlock) > 0 {
		return &account.BlockByRangeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "start block number must be less than or equal to end block number",
		}, nil
	}

	// 预分配 slice 长度
	blockHeaderList := make([]*account.BlockHeader, 0, endBlock.Int64()-startBlock.Int64()+1)

	// 循环获取区块数据
	for i := startBlock.Int64(); i <= endBlock.Int64(); i++ {
		block, err := c.tronClient.GetBlockByNumber(i)
		if err != nil {
			return &account.BlockByRangeResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  fmt.Sprintf("failed to get block %d: %v", i, err),
			}, nil
		}

		// 将获取的区块数据添加到 blockHeaderList
		blockHeaderList = append(blockHeaderList, &account.BlockHeader{
			ParentHash: block.ParentHash,
			Difficulty: block.Difficulty,
			Number:     block.Number,
			Nonce:      block.Nonce,
		})
	}

	// 返回成功响应
	return &account.BlockByRangeResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "get block by range success",
		BlockHeader: blockHeaderList,
	}, nil
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
	return &account.DecodeTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "verify tx success",
		Base64Tx: "0x000000",
	}, nil
}

// VerifySignedTransaction 验证签名
func (c *ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	return &account.VerifyTransactionResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "verify tx success",
		Verify: true,
	}, nil
}

// GetExtraData 获取额外数据
func (c *ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	return &account.ExtraDataResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "get extra data success",
		Value: "not data",
	}, nil
}
