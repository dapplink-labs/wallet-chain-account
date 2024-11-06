package ethereum

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/status-im/keycard-go/hexutils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"

	account2 "github.com/dapplink-labs/chain-explorer-api/common/account"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

const ChainName = "Ethereum"

type ChainAdaptor struct {
	ethClient     EthClient
	ethDataClient *EthData
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	ethClient, err := DialEthClient(context.Background(), conf.WalletNode.Eth.RPCs[0].RPCURL)
	if err != nil {
		return nil, err
	}
	ethDataClient, err := NewEthDataClient(conf.WalletNode.Eth.DataApiUrl, conf.WalletNode.Eth.DataApiKey, time.Second*15)
	if err != nil {
		return nil, err
	}
	return &ChainAdaptor{
		ethClient:     ethClient,
		ethDataClient: ethDataClient,
	}, nil
}

func (c *ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	return &account.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

func (c *ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	publicKeyBytes, err := hex.DecodeString(req.PublicKey)
	if err != nil {
		return &account.ConvertAddressResponse{
			Code:    common2.ReturnCode_ERROR,
			Msg:     "convert address fail",
			Address: common.Address{}.String(),
		}, nil
	}
	addressCommon := common.BytesToAddress(crypto.Keccak256(publicKeyBytes[1:])[12:])
	return &account.ConvertAddressResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "convert address success",
		Address: addressCommon.String(),
	}, nil
}

func (c *ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	if len(req.Address) != 42 || !strings.HasPrefix(req.Address, "0x") {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_SUCCESS,
			Msg:   "invalid address",
			Valid: false,
		}, nil
	}
	ok := regexp.MustCompile("^[0-9a-fA-F]{40}$").MatchString(req.Address[2:])
	if ok {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_SUCCESS,
			Msg:   "valid address",
			Valid: true,
		}, nil
	} else {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_SUCCESS,
			Msg:   "invalid address",
			Valid: false,
		}, nil
	}
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	var blockNumber *big.Int
	if req.Height == 0 {
		blockNumber = nil // return latest block
	} else {
		blockNumber = big.NewInt(req.Height) // return special block by number
	}
	blockInfo, err := c.ethClient.BlockHeaderByNumber(blockNumber)
	if err != nil {
		log.Error("get latest block header fail", "err", err)
		return &account.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get latest block header fail",
		}, nil
	}
	blockHead := &account.BlockHeader{
		Hash:             blockInfo.Hash().String(),
		ParentHash:       blockInfo.ParentHash.String(),
		UncleHash:        blockInfo.UncleHash.String(),
		CoinBase:         blockInfo.Coinbase.String(),
		Root:             blockInfo.Root.String(),
		TxHash:           blockInfo.TxHash.String(),
		ReceiptHash:      blockInfo.ReceiptHash.String(),
		ParentBeaconRoot: blockInfo.ParentBeaconRoot.String(),
		Difficulty:       blockInfo.Difficulty.String(),
		Number:           blockInfo.Number.String(),
		GasLimit:         blockInfo.GasLimit,
		GasUsed:          blockInfo.GasUsed,
		Time:             blockInfo.Time,
		Extra:            string(blockInfo.Extra),
		MixDigest:        blockInfo.MixDigest.String(),
		Nonce:            strconv.FormatUint(blockInfo.Nonce.Uint64(), 10),
		BaseFee:          blockInfo.BaseFee.String(),
		WithdrawalsHash:  blockInfo.WithdrawalsHash.String(),
		BlobGasUsed:      *blockInfo.BlobGasUsed,
		ExcessBlobGas:    *blockInfo.ExcessBlobGas,
	}
	return &account.BlockHeaderResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "get latest block header success",
		BlockHeader: blockHead,
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	blockInfo, err := c.ethClient.BlockHeaderByHash(common.HexToHash(req.Hash))
	if err != nil {
		log.Error("get latest block header fail", "err", err)
		return &account.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get latest block header fail",
		}, nil
	}
	blockHeader := &account.BlockHeader{
		Hash:             blockInfo.Hash().String(),
		ParentHash:       blockInfo.ParentHash.String(),
		UncleHash:        blockInfo.UncleHash.String(),
		CoinBase:         blockInfo.Coinbase.String(),
		Root:             blockInfo.Root.String(),
		TxHash:           blockInfo.TxHash.String(),
		ReceiptHash:      blockInfo.ReceiptHash.String(),
		ParentBeaconRoot: blockInfo.ParentBeaconRoot.String(),
		Difficulty:       blockInfo.Difficulty.String(),
		Number:           blockInfo.Number.String(),
		GasLimit:         blockInfo.GasLimit,
		GasUsed:          blockInfo.GasUsed,
		Time:             blockInfo.Time,
		Extra:            string(blockInfo.Extra),
		MixDigest:        blockInfo.MixDigest.String(),
		Nonce:            strconv.FormatUint(blockInfo.Nonce.Uint64(), 10),
		BaseFee:          blockInfo.BaseFee.String(),
		WithdrawalsHash:  blockInfo.WithdrawalsHash.String(),
		BlobGasUsed:      *blockInfo.BlobGasUsed,
		ExcessBlobGas:    *blockInfo.ExcessBlobGas,
	}
	return &account.BlockHeaderResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "get latest block header success",
		BlockHeader: blockHeader,
	}, nil
}

func (c *ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	block, err := c.ethClient.BlockByNumber(big.NewInt(req.Height))
	if err != nil {
		log.Error("block by number error", err)
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "block by number error",
		}, nil
	}
	var txListRet []*account.BlockInfoTransactionList
	for _, v := range block.Transactions {
		bitlItem := &account.BlockInfoTransactionList{
			From:   v.From,
			To:     v.To,
			Hash:   v.Hash,
			Amount: v.Value,
		}
		txListRet = append(txListRet, bitlItem)
	}
	return &account.BlockResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "block by number success",
		Height:       int64(block.Height),
		Hash:         block.Hash.String(),
		BaseFee:      block.BaseFee,
		Transactions: txListRet,
	}, nil
}

func (c *ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	block, err := c.ethClient.BlockByHash(common.HexToHash(req.Hash))
	if err != nil {
		log.Error("block by number error", err)
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "block by number error",
		}, nil
	}
	var txListRet []*account.BlockInfoTransactionList
	for _, v := range block.Transactions {
		bitlItem := &account.BlockInfoTransactionList{
			From:   v.From,
			To:     v.To,
			Hash:   v.Hash,
			Amount: v.Value,
		}
		txListRet = append(txListRet, bitlItem)
	}
	return &account.BlockResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "block by hash success",
		Height:       int64(block.Height),
		Hash:         block.Hash.String(),
		BaseFee:      block.BaseFee,
		Transactions: txListRet,
	}, nil
}

func (c *ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	nonceResult, err := c.ethClient.TxCountByAddress(common.HexToAddress(req.Address))
	if err != nil {
		log.Error("get nonce by address fail", "err", err)
		return &account.AccountResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get nonce by address fail",
		}, nil
	}
	balanceResult, err := c.ethDataClient.GetBalanceByAddress(req.ContractAddress, req.Address)
	if err != nil {
		return &account.AccountResponse{
			Code:    common2.ReturnCode_ERROR,
			Msg:     "get token balance fail",
			Balance: "0",
		}, err
	}
	log.Info("balance result", "balance=", balanceResult.Balance, "balanceStr=", balanceResult.BalanceStr)
	return &account.AccountResponse{
		Code:          common2.ReturnCode_SUCCESS,
		Msg:           "get account response success",
		AccountNumber: "0",
		Sequence:      nonceResult.String(),
		Balance:       balanceResult.Balance.Int().String(),
	}, nil
}

func (c *ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	gasPrice, err := c.ethClient.SuggestGasPrice()
	if err != nil {
		log.Error("get gas price failed", "err", err)
		return &account.FeeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get suggest gas price fail",
		}, nil
	}
	gasTipCap, err := c.ethClient.SuggestGasTipCap()
	if err != nil {
		log.Error("get gas price failed", "err", err)
		return &account.FeeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get suggest gas price fail",
		}, nil
	}
	return &account.FeeResponse{
		Code:      common2.ReturnCode_SUCCESS,
		Msg:       "get gas price success",
		SlowFee:   gasPrice.String() + "|" + gasTipCap.String(),
		NormalFee: gasPrice.String() + "|" + gasTipCap.String() + "|" + "*2",
		FastFee:   gasPrice.String() + "|" + gasTipCap.String() + "|" + "*3",
	}, nil
}

func (c *ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	err := c.ethClient.SendRawTransaction(req.RawTx)
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

func (c *ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	var resp *account2.TransactionResponse[account2.AccountTxResponse]
	var err error
	if req.ContractAddress != "0x00" && req.ContractAddress != "" {
		resp, err = c.ethDataClient.GetTxByAddress(uint64(req.Page), uint64(req.Pagesize), req.Address, "tokentx")
	} else {
		resp, err = c.ethDataClient.GetTxByAddress(uint64(req.Page), uint64(req.Pagesize), req.Address, "txlist")
	}
	if err != nil {
		log.Error("get GetTxByAddress error", "err", err)
		return &account.TxAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get tx list fail",
			Tx:   nil,
		}, err
	} else {
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
		fmt.Println("resp", resp)
		return &account.TxAddressResponse{
			Code: common2.ReturnCode_SUCCESS,
			Msg:  "get tx list success",
			Tx:   list,
		}, nil
	}
}

func (c *ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	tx, err := c.ethClient.TxByHash(common.HexToHash(req.Hash))
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return &account.TxHashResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "Ethereum Tx NotFound",
			}, nil
		}
		log.Error("get transaction error", "err", err)
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Ethereum Tx NotFound",
		}, nil
	}
	receipt, err := c.ethClient.TxReceiptByHash(common.HexToHash(req.Hash))
	if err != nil {
		log.Error("get transaction receipt error", "err", err)
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Get transaction receipt error",
		}, nil
	}
	var from_addrs []*account.Address
	var to_addrs []*account.Address
	var value_list []*account.Value
	from_addrs = append(from_addrs, &account.Address{Address: ""})
	to_addrs = append(to_addrs, &account.Address{Address: tx.To().Hex()})
	value_list = append(value_list, &account.Value{Value: tx.Value().String()})
	var txStatus account.TxStatus
	if receipt.Status == 1 {
		txStatus = account.TxStatus_Success
	} else {
		txStatus = account.TxStatus_Failed
	}
	return &account.TxHashResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get transaction success",
		Tx: &account.TxMessage{
			Hash:            tx.Hash().Hex(),
			Index:           uint32(receipt.TransactionIndex),
			Froms:           from_addrs,
			Tos:             to_addrs,
			Values:          value_list,
			Fee:             tx.GasFeeCap().String(),
			Status:          txStatus,
			Type:            0,
			Height:          receipt.BlockNumber.String(),
			ContractAddress: tx.To().String(),
			Data:            hexutils.BytesToHex(tx.Data()),
		},
	}, nil
}

func (c *ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	startBlock := new(big.Int)
	endBlock := new(big.Int)
	startBlock.SetString(req.Start, 10)
	endBlock.SetString(req.End, 10)
	blockRange, err := c.ethClient.BlockHeadersByRange(startBlock, endBlock, 1)
	if err != nil {
		log.Error("get block range fail", "err", err)
		return &account.BlockByRangeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block range fail",
		}, err
	}
	blockHeaderList := make([]*account.BlockHeader, 0, len(blockRange))
	for _, block := range blockRange {
		blockItem := &account.BlockHeader{
			ParentHash:       block.ParentHash.String(),
			UncleHash:        block.UncleHash.String(),
			CoinBase:         block.Coinbase.String(),
			Root:             block.Root.String(),
			TxHash:           block.TxHash.String(),
			ReceiptHash:      block.ReceiptHash.String(),
			ParentBeaconRoot: block.ParentBeaconRoot.String(),
			Difficulty:       block.Difficulty.String(),
			Number:           block.Number.String(),
			GasLimit:         block.GasLimit,
			GasUsed:          block.GasUsed,
			Time:             block.Time,
			Extra:            string(block.Extra),
			MixDigest:        block.MixDigest.String(),
			Nonce:            strconv.FormatUint(block.Nonce.Uint64(), 10),
			BaseFee:          block.BaseFee.String(),
			WithdrawalsHash:  block.WithdrawalsHash.String(),
			BlobGasUsed:      *block.BlobGasUsed,
			ExcessBlobGas:    *block.ExcessBlobGas,
		}
		blockHeaderList = append(blockHeaderList, blockItem)
	}
	return &account.BlockByRangeResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "get block range success",
		BlockHeader: blockHeaderList,
	}, nil
}

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
	var amount *big.Int
	chainID := new(big.Int)
	gasTipCap := new(big.Int)
	maxFeePerGas := new(big.Int)
	valueData := new(big.Int)
	chainID.SetString(data.ChainId, 10)
	gasTipCap.SetString(data.GasTipCap, 10)
	maxFeePerGas.SetString(data.GasFeeCap, 10)
	valueData.SetString(data.Value, 10)
	toAddress := common.HexToAddress(data.ToAddress)
	tokenAddress := common.HexToAddress(data.ContractAddress)
	var buildData []byte
	if data.ContractAddress != "0x00" {
		buildData = BuildErc20Data(toAddress, valueData)
		toAddress = tokenAddress
		amount = big.NewInt(0)
	} else {
		toAddress = toAddress
		amount = valueData
	}
	dFeeTx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     data.Nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: maxFeePerGas,
		Gas:       data.Gas,
		To:        &toAddress,
		Value:     amount,
		Data:      buildData,
	}
	rawTx, err := CreateEip1559UnSignTx(dFeeTx, chainID)
	if err != nil {
		log.Error("create un sign tx fail", "err", err)
		return &account.UnSignTransactionResponse{
			Code:     common2.ReturnCode_ERROR,
			Msg:      "create un sign tx fail",
			UnSignTx: "",
		}, err
	}
	return &account.UnSignTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "create un sign tx success",
		UnSignTx: rawTx,
	}, nil
}

func (c *ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
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
	var amount *big.Int
	chainID := new(big.Int)
	gasTipCap := new(big.Int)
	maxFeePerGas := new(big.Int)
	valueData := new(big.Int)
	chainID.SetString(data.ChainId, 10)
	gasTipCap.SetString(data.GasTipCap, 10)
	maxFeePerGas.SetString(data.GasFeeCap, 10)
	valueData.SetString(data.Value, 10)
	toAddress := common.HexToAddress(data.ToAddress)
	tokenAddress := common.HexToAddress(data.ContractAddress)
	var buildData []byte
	if data.ContractAddress != "0x00" {
		buildData = BuildErc20Data(toAddress, valueData)
		toAddress = tokenAddress
		amount = big.NewInt(0)
	} else {
		toAddress = toAddress
		amount = valueData
	}
	dFeeTx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     data.Nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: maxFeePerGas,
		Gas:       data.Gas,
		To:        &toAddress,
		Value:     amount,
		Data:      buildData,
	}
	rawTx, txHash, err := CreateEip1559SignedTx(dFeeTx, []byte(req.Signature), chainID)
	if err != nil {
		log.Error("create un sign tx fail", "err", err)
		return &account.SignedTransactionResponse{
			Code:     common2.ReturnCode_ERROR,
			Msg:      "create un sign tx fail",
			SignedTx: "",
		}, err
	}
	return &account.SignedTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      txHash,
		SignedTx: rawTx,
	}, nil
}

func (c *ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	return &account.DecodeTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "verify tx success",
		Base64Tx: "0x000000",
	}, nil
}

func (c *ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	return &account.VerifyTransactionResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "verify tx success",
		Verify: true,
	}, nil
}

func (c *ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	return &account.ExtraDataResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "get extra data success",
		Value: "not data",
	}, nil
}
