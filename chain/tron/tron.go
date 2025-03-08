package tron

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcutil/base58"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/dapplink-labs/wallet-chain-account/rpc/common"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/fbsobreira/gotron-sdk/pkg/address"

	"math/big"
	"strconv"
	"time"
)

const (
	ChainName = "Tron"
)

type ChainAdaptor struct {
	tronClient     *TronClient
	tronDataClient *TronData
}

func (c *ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) BuildUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	rpc := conf.WalletNode.Tron
	tronClient := DialTronClient(rpc.RpcUrl, rpc.RpcUser, rpc.RpcPass)
	tronDataClient, err := NewTronDataClient(conf.WalletNode.Tron.DataApiUrl, conf.WalletNode.Tron.DataApiKey, time.Second*15)
	if err != nil {
		return nil, err
	}
	return &ChainAdaptor{
		tronClient:     tronClient,
		tronDataClient: tronDataClient,
	}, nil
}

// GetSupportChains Return whether the chain is supported
func (c *ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	return &account.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

// ConvertAddress Convert public key to address
func (c *ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	// Decoding hexadecimal strings into byte arrays
	pubKeyBytes, err := hex.DecodeString(req.PublicKey)
	// Parse byte array into public key
	pubKey, _ := btcec.ParsePubKey(pubKeyBytes)
	// Convert public key to address
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

// ValidAddress verify address
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

// GetBlockByNumber Obtain block information based on block height
func (c *ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	// 获取区块数据
	rsp, err := c.tronClient.GetBlockByNumber(req.Height)
	if err != nil {
		return &account.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	// 解析交易列表
	transactions := make([]*account.BlockInfoTransactionList, 0)
	for _, tx := range rsp.Transactions {
		// 跳过非转账交易
		if len(tx.RawData.Contract) == 0 {
			continue
		}

		contract := tx.RawData.Contract[0]
		txInfo := &account.BlockInfoTransactionList{
			Hash:   tx.TxID,
			Height: uint64(rsp.BlockHeader.RawData.Number),
		}

		// 根据不同合约类型处理
		switch contract.Type {
		case "TransferContract":
			// TRX 转账
			value := contract.Parameter.Value
			txInfo.From = value.OwnerAddress
			txInfo.To = value.ToAddress
			txInfo.Amount = strconv.FormatInt(value.Amount, 10)
			// TRX 转账不需要设置 TokenAddress

		case "TriggerSmartContract":
			// TRC20 代币转账
			value := contract.Parameter.Value
			txInfo.From = value.OwnerAddress
			txInfo.TokenAddress = value.ContractAddress

			// 解析 data 字段获取接收地址和金额
			if len(value.Data) >= 138 { // 0x + 方法(8) + 参数1(64) + 参数2(64)
				// data 格式: 0xa9059cbb + 接收地址(32字节) + 金额(32字节)
				// 跳过方法ID(4字节)和填充的0
				to := "41" + value.Data[32:72] // 添加 TRON 地址前缀
				txInfo.To = to

				// 解析金额(去除前导0)
				amountHex := strings.TrimLeft(value.Data[72:136], "0")
				if amountHex == "" {
					amountHex = "0"
				}
				amount, _ := new(big.Int).SetString(amountHex, 16)
				txInfo.Amount = amount.String()
			}

		case "TransferAssetContract":
			// TRC10 代币转账
			value := contract.Parameter.Value
			txInfo.From = value.OwnerAddress
			txInfo.To = value.ToAddress
			txInfo.Amount = strconv.FormatInt(value.Amount, 10)
			txInfo.TokenAddress = value.AssetName // TRC10 使用资产名称作为标识
		}

		// 只添加解析成功的交易
		if txInfo.From != "" && txInfo.To != "" {
			transactions = append(transactions, txInfo)
		}
	}

	return &account.BlockResponse{
		Code:         common.ReturnCode_SUCCESS,
		Height:       rsp.BlockHeader.RawData.Number,
		Hash:         rsp.BlockID,
		BaseFee:      "", // TRON 没有 BaseFee 概念
		Transactions: transactions,
	}, nil
}

// GetBlockByHash Obtain block information based on block hash
func (c *ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	// 获取区块数据
	rsp, err := c.tronClient.GetBlockByNumber(req.Hash)
	if err != nil {
		return &account.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	// 解析交易列表
	transactions := make([]*account.BlockInfoTransactionList, 0)
	for _, tx := range rsp.Transactions {
		// 跳过非转账交易
		if len(tx.RawData.Contract) == 0 {
			continue
		}

		contract := tx.RawData.Contract[0]
		txInfo := &account.BlockInfoTransactionList{
			Hash:   tx.TxID,
			Height: uint64(rsp.BlockHeader.RawData.Number),
		}

		// 根据不同合约类型处理
		switch contract.Type {
		case "TransferContract":
			// TRX 转账
			value := contract.Parameter.Value
			txInfo.From = value.OwnerAddress
			txInfo.To = value.ToAddress
			txInfo.Amount = strconv.FormatInt(value.Amount, 10)
			// TRX 转账不需要设置 TokenAddress

		case "TriggerSmartContract":
			// TRC20 代币转账
			value := contract.Parameter.Value
			txInfo.From = value.OwnerAddress
			txInfo.TokenAddress = value.ContractAddress

			// 解析 data 字段获取接收地址和金额
			if len(value.Data) >= 138 { // 0x + 方法(8) + 参数1(64) + 参数2(64)
				// data 格式: 0xa9059cbb + 接收地址(32字节) + 金额(32字节)
				// 跳过方法ID(4字节)和填充的0
				to := "41" + value.Data[32:72] // 添加 TRON 地址前缀
				txInfo.To = to

				// 解析金额(去除前导0)
				amountHex := strings.TrimLeft(value.Data[72:136], "0")
				if amountHex == "" {
					amountHex = "0"
				}
				amount, _ := new(big.Int).SetString(amountHex, 16)
				txInfo.Amount = amount.String()
			}

		case "TransferAssetContract":
			// TRC10 代币转账
			value := contract.Parameter.Value
			txInfo.From = value.OwnerAddress
			txInfo.To = value.ToAddress
			txInfo.Amount = strconv.FormatInt(value.Amount, 10)
			txInfo.TokenAddress = value.AssetName // TRC10 使用资产名称作为标识
		}

		// 只添加解析成功的交易
		if txInfo.From != "" && txInfo.To != "" {
			transactions = append(transactions, txInfo)
		}
	}

	return &account.BlockResponse{
		Code:         common.ReturnCode_SUCCESS,
		Height:       rsp.BlockHeader.RawData.Number,
		Hash:         rsp.BlockID,
		BaseFee:      "", // TRON 没有 BaseFee 概念
		Transactions: transactions,
	}, nil
}

// GetBlockHeaderByHash Obtain block header information based on block hash
func (c *ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {

	rsp, err := c.tronClient.GetBlockHeaderByHash(req.Hash)
	if err != nil {
		log.Error("GetBlockHeaderByHash fail:", err)
		return &account.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	blockHeader := &account.BlockHeader{
		Hash:        rsp.BlockID,
		ParentHash:  rsp.BlockHeader.RawData.ParentHash,
		UncleHash:   "",
		CoinBase:    "",
		Root:        "",
		TxHash:      "",
		ReceiptHash: "",
		Difficulty:  "",
		Number:      strconv.FormatUint(uint64(rsp.BlockHeader.RawData.Number), 10),
		GasLimit:    0,
		GasUsed:     0,
		Time:        uint64(rsp.BlockHeader.RawData.Timestamp),
		Extra:       "",
		MixDigest:   "",
		Nonce:       "",
		BaseFee:     "",

		ParentBeaconRoot: "",
		WithdrawalsHash:  "",
		BlobGasUsed:      0,
		ExcessBlobGas:    0,
	}
	return &account.BlockHeaderResponse{
		Code:        common.ReturnCode_SUCCESS,
		Msg:         "get latest block header success",
		BlockHeader: blockHeader,
	}, nil

}

// GetBlockHeaderByNumber Obtain block header information based on block height
func (c *ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {

	rsp, err := c.tronClient.GetBlockHeaderByNumber(req.Height)
	if err != nil {
		log.Error("GetBlockHeaderByHash fail:", err)
		return &account.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	blockHeader := &account.BlockHeader{
		Hash:        rsp.BlockID,
		ParentHash:  rsp.BlockHeader.RawData.ParentHash,
		UncleHash:   "",
		CoinBase:    "",
		Root:        "",
		TxHash:      "",
		ReceiptHash: "",
		Difficulty:  "",
		Number:      strconv.FormatUint(uint64(rsp.BlockHeader.RawData.Number), 10),
		GasLimit:    0,
		GasUsed:     0,
		Time:        uint64(rsp.BlockHeader.RawData.Timestamp),
		Extra:       "",
		MixDigest:   "",
		Nonce:       "",
		BaseFee:     "",

		ParentBeaconRoot: "",
		WithdrawalsHash:  "",
		BlobGasUsed:      0,
		ExcessBlobGas:    0,
	}
	return &account.BlockHeaderResponse{
		Code:        common.ReturnCode_SUCCESS,
		Msg:         "get latest block header success",
		BlockHeader: blockHeader,
	}, nil
}

// GetAccount Obtain account information based on the address
func (c *ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {

	info, err := c.tronClient.GetBalance(req.Address)
	if err != nil {
		return &account.AccountResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}
	return &account.AccountResponse{
		Code:          common2.ReturnCode_SUCCESS,
		Msg:           "get account response success",
		AccountNumber: "0",
		Sequence:      "0",
		Balance:       strconv.FormatInt(info.Balance, 10),
	}, nil

}

// GetFee Obtain transaction fees
func (c *ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {

	return &account.FeeResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "does not currently support",
	}, nil
}

func (c *ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {

	return &account.TxAddressResponse{
		Code: common2.ReturnCode_ERROR,
		Msg:  "not support",
	}, nil
}

// GetTxByHash Obtain transactions based on transaction hash
func (c *ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	// 获取交易信息
	resp, err := c.tronClient.GetTransactionByHash(req.Hash)
	if err != nil {
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Get transaction error: " + err.Error(),
		}, nil
	}

	// 检查合约数据
	if len(resp.RawData.Contract) == 0 {
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Invalid transaction: no contract data",
		}, nil
	}

	contract := resp.RawData.Contract[0]
	value := contract.Parameter.Value

	// 准备交易数据
	var fromAddr = value.OwnerAddress
	var toAddr = value.ToAddress
	var amountStr = strconv.FormatInt(value.Amount, 10) // 保留原来的赋值
	var contractAddress string

	// 根据合约类型处理
	switch contract.Type {
	case "TransferContract":
		// TRX 转账
		contractAddress = ""
	case "TriggerSmartContract":
		// TRC20 代币转账
		contractAddress = value.ContractAddress
		// 解析接收地址 (跳过方法ID 8字节)
		toAddr = "41" + value.Data[32:72] // 添加 41 前缀，取完整的地址部分

		// 解析金额
		amountHex := value.Data[72:136] // 取完整的金额部分
		amountHex = strings.TrimLeft(amountHex, "0")
		if amountHex == "" {
			amountHex = "0"
		}
		bigIntAmount, _ := new(big.Int).SetString(amountHex, 16)
		amountStr = bigIntAmount.String()

	case "TransferAssetContract":
		// TRC10 代币转账
		contractAddress = value.AssetName
	}

	txStatus := account.TxStatus_Success
	if len(resp.Ret) > 0 && resp.Ret[0].ContractRet != "SUCCESS" {
		txStatus = account.TxStatus_Failed
	}

	return &account.TxHashResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get transaction success",
		Tx: &account.TxMessage{
			Hash:            resp.TxID,
			Index:           0, // Tron 不使用此字段
			From:            fromAddr,
			To:              toAddr,
			Value:           amountStr,
			Fee:             strconv.FormatInt(resp.RawData.FeeLimit, 10),
			Status:          txStatus,
			Type:            0,
			Height:          "", // 需要单独查询区块高度
			ContractAddress: contractAddress,
			Data:            resp.RawDataHex,
		},
	}, nil
}

// GetBlockByRange Obtain blocks based on their scope
func (c *ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	//// Convert the starting and ending blocks to big Int
	//startBlock := new(big.Int)
	//endBlock := new(big.Int)
	//if _, ok := startBlock.SetString(req.Start, 10); !ok {
	//	return nil, fmt.Errorf("invalid start block number: %s", req.Start)
	//}
	//if _, ok := endBlock.SetString(req.End, 10); !ok {
	//	return nil, fmt.Errorf("invalid end block number: %s", req.End)
	//}
	//
	//// Ensure that the starting block number is less than or equal to the ending block number
	//if startBlock.Cmp(endBlock) > 0 {
	//	return &account.BlockByRangeResponse{
	//		Code: common2.ReturnCode_ERROR,
	//		Msg:  "start block number must be less than or equal to end block number",
	//	}, nil
	//}
	//
	//// Pre allocated slice length
	//blockHeaderList := make([]*account.BlockHeader, 0, endBlock.Int64()-startBlock.Int64()+1)
	//
	//// Loop to obtain block data
	//for i := startBlock.Int64(); i <= endBlock.Int64(); i++ {
	//	block, err := c.tronClient.GetBlockByNumber(i)
	//	if err != nil {
	//		return &account.BlockByRangeResponse{
	//			Code: common2.ReturnCode_ERROR,
	//			Msg:  fmt.Sprintf("failed to get block %d: %v", i, err),
	//		}, nil
	//	}
	//
	//	// Add the obtained block data to blockHeaderList
	//	blockHeaderList = append(blockHeaderList, &account.BlockHeader{
	//		ParentHash: block.ParentHash,
	//		Difficulty: block.Difficulty,
	//		Number:     block.Number,
	//		Nonce:      block.Nonce,
	//	})
	//}

	// Return successful response
	return &account.BlockByRangeResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "not support get block by range ",
		//BlockHeader: blockHeaderList,
	}, nil
}

//
//// BuildUnSignTransaction Create unsigned transactions
//func (c *ChainAdaptor) BuildUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
//	jsonBytes, err := base64.StdEncoding.DecodeString(req.Base64Tx)
//	if err != nil {
//		log.Error("decode string fail", "err", err)
//		return nil, err
//	}
//	var data TxStructure
//	if err := json.Unmarshal(jsonBytes, &data); err != nil {
//		log.Error("parse json fail", "err", err)
//		return nil, err
//	}
//	var transaction *UnSignTransaction
//	if data.ContractAddress == "" {
//		transaction, err = c.tronClient.CreateTRXTransaction(data.FromAddress, data.ToAddress, data.Value)
//	} else {
//		transaction, err = c.tronClient.CreateTRC20Transaction(data.FromAddress, data.ToAddress, data.ContractAddress, data.Value)
//	}
//	if err != nil {
//		return nil, err
//	}
//	return &account.UnSignTransactionResponse{
//		Code:     common2.ReturnCode_SUCCESS,
//		Msg:      "create un sign tx success",
//		UnSignTx: transaction.RawDataHex,
//	}, nil
//}
//
//// BuildSignedTransaction Create a signed transaction
//func (c *ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
//	return &account.SignedTransactionResponse{
//		Code: common2.ReturnCode_ERROR,
//		Msg:  "not support build signed transaction",
//	}, nil
//}
//
//// DecodeTransaction Decoding transactions
//func (c *ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
//	return &account.DecodeTransactionResponse{
//		Code:     common2.ReturnCode_SUCCESS,
//		Msg:      "decode tx success",
//		Base64Tx: "0x000000",
//	}, nil
//}
//
//// VerifySignedTransaction verify signature
//func (c *ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
//	return &account.VerifyTransactionResponse{
//		Code:   common2.ReturnCode_SUCCESS,
//		Msg:    "verify tx success",
//		Verify: true,
//	}, nil
//}

// GetExtraData Obtain additional data
func (c *ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	return &account.ExtraDataResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "get extra data success",
		Value: "not data",
	}, nil
}

func (c *ChainAdaptor) GetNftListByAddress(req *account.NftAddressRequest) (*account.NftAddressResponse, error) {
	panic("implement me")
}

//
//// 检查交易确认状态
//func (c *ChainAdaptor) CheckTxConfirmations(txHash string) (bool, int64, error) {
//	// 获取交易信息
//	tx, err := c.tronClient.GetTransactionByHash(txHash)
//	if err != nil {
//		return false, 0, err
//	}
//
//	// 获取当前区块
//	currentBlock, err := c.tronClient.GetLatestBlock()
//	if err != nil {
//		return false, 0, err
//	}
//
//	// 计算确认数
//	confirmations := currentBlock.Number - tx.BlockNumber
//
//	// 根据交易类型和金额判断所需确认数
//	requiredConfirmations := TRON_NORMAL_CONFIRMATIONS
//	if isContractTransaction(tx) || isExchangeTransaction(tx) {
//		requiredConfirmations = TRON_EXCHANGE_CONFIRMATIONS
//	} else if isLargeTransaction(tx) {
//		requiredConfirmations = TRON_LARGE_CONFIRMATIONS
//	}
//
//	return confirmations >= requiredConfirmations, confirmations, nil
//}

func parseBlockTimestamp(timestampStr string) (uint64, error) {
	timestamp, err := strconv.ParseUint(timestampStr, 10, 64)
	if err != nil {
		// 尝试十六进制解析
		if strings.HasPrefix(timestampStr, "0x") {
			timestamp, err = strconv.ParseUint(strings.TrimPrefix(timestampStr, "0x"), 16, 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse timestamp %s: %v", timestampStr, err)
			}
		} else {
			return 0, fmt.Errorf("invalid timestamp format %s: %v", timestampStr, err)
		}
	}
	return timestamp, nil
}

func ParseHexto10(a string) (b string) {
	hexStr := strings.TrimPrefix(a, "0x")
	decimal, _ := strconv.ParseInt(hexStr, 16, 64)
	return strconv.FormatInt(decimal, 10)
}

func HexToUint64(hexStr string) (uint64, error) {
	// 空字符串检查
	if hexStr == "" {
		return 0, fmt.Errorf("empty hex string")
	}

	// 去掉 0x 前缀
	hexStr = strings.TrimPrefix(hexStr, "0x")

	// 空字符串检查（去掉前缀后）
	if hexStr == "" {
		return 0, nil
	}

	// 转换为 uint64
	result, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert hex to uint64: %v", err)
	}

	return result, nil
}

func FormatTronAddress(address string) string {
	// 如果是 Base58 地址，转换为 Hex
	if strings.HasPrefix(address, "T") {
		// 需要先解码 Base58，然后转为 Hex
		return "0x" + hex.EncodeToString(base58.Decode(address))
	}

	// 如果是 Hex 地址，确保有 0x 前缀
	if !strings.HasPrefix(address, "0x") {
		return "0x" + address
	}

	return address
}
