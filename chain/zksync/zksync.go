package zksync

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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"

	account2 "github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/chain/evmbase"
	erc20base "github.com/dapplink-labs/wallet-chain-account/chain/evmbase"
	"github.com/dapplink-labs/wallet-chain-account/common/util"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

const ChainName = "Zksync"

type ChainAdaptor struct {
	ethClient     erc20base.EthClient
	ethDataClient *erc20base.EthData
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	ethClient, err := erc20base.DialEthClient(context.Background(), conf.WalletNode.Zksync.RpcUrl)
	if err != nil {
		return nil, err
	}
	ethDataClient, err := erc20base.NewEthDataClient(
		conf.WalletNode.Zksync.DataApiUrl,
		conf.WalletNode.Zksync.DataApiKey,
		time.Second*15,
	)
	if err != nil {
		return nil, err
	}
	return &ChainAdaptor{
		ethClient:     ethClient,
		ethDataClient: ethDataClient,
	}, nil

}

func (c ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	return &account.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

func (c ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	publicKeyBytes, err := hex.DecodeString(req.PublicKey)

	if err != nil {
		log.Error("decode public key failed")
		return &account.ConvertAddressResponse{
			Code:    common2.ReturnCode_ERROR,
			Msg:     "convert address fail",
			Address: common.Address{}.String(),
		}, err
	}

	address := common.BytesToAddress(crypto.Keccak256(publicKeyBytes[1:])[12:])
	return &account.ConvertAddressResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "convert address success",
		Address: address.String(),
	}, nil
}

func (c ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	// TODO: The `code` field is not showing in the response when it's set to ReturnCode_SUCCESS (0)
	// This is due to protobuf's default behavior of omitting fields with default values
	// To fix this, we need to:
	// 1. Modify the protobuf definition to force showing the field using `[(gogoproto.jsontag) = "code"]`
	// 2. Or handle this behavior in the client side by treating missing code as SUCCESS
	if len(req.Address) != 42 || !strings.HasPrefix(req.Address, "0x") {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_ERROR,
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
	}
	return &account.ValidAddressResponse{
		Code:  common2.ReturnCode_ERROR,
		Msg:   "invalid address format",
		Valid: false,
	}, nil
}

func (c ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	var blockNumber *big.Int
	if req.Height == 0 {
		blockNumber = nil // return latest block
	} else {
		blockNumber = big.NewInt(req.Height) // return special block by number
	}
	block, err := c.ethClient.BlockByNumber(blockNumber)
	if err != nil {
		log.Error("get block by number fail", "error", err)
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "block by number error",
		}, nil
	}
	height, _ := block.NumberUint64()
	txList := make([]*account.BlockInfoTransactionList, len(block.Transactions))
	for i, tx := range block.Transactions {
		txList[i] = &account.BlockInfoTransactionList{
			From:           tx.From,
			To:             tx.To,
			TokenAddress:   tx.To,
			ContractWallet: tx.To,
			Hash:           tx.Hash,
			Height:         height,
			Amount:         tx.Value,
		}
	}

	return &account.BlockResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "get latest block header success",
		Height:       int64(height),
		Hash:         block.Hash.String(),
		BaseFee:      block.BaseFee,
		Transactions: txList,
	}, nil
}

func (c ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	block, err := c.ethClient.BlockByHash(common.HexToHash(req.Hash))
	if err != nil {
		log.Error("get block by hash fail", "error", err)
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "block by hash error",
		}, nil
	}
	height, _ := block.NumberUint64()
	txList := make([]*account.BlockInfoTransactionList, len(block.Transactions))
	for i, tx := range block.Transactions {
		txList[i] = &account.BlockInfoTransactionList{
			From:   tx.From,
			To:     tx.To,
			Hash:   tx.Hash,
			Amount: tx.Value,
			Height: height,
		}
	}
	return &account.BlockResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "get block by hash success",
		Height:       int64(height),
		Hash:         block.Hash.String(),
		BaseFee:      block.BaseFee,
		Transactions: txList,
	}, nil
}

func (c ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	hash := common.HexToHash(req.Hash)
	var header *types.Header
	// zksync的header.Hash()得到的hash 和 block hash 不一致，使用原始hash进行对比
	err := c.ethClient.(interface{ GetRPC() evmbase.RPC }).GetRPC().CallContext(ctxwt, &header, "eth_getBlockByHash", hash, false)
	if err != nil {
		return &account.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block header by hash fail",
		}, nil
	}
	if header == nil {
		return &account.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "block header not found",
		}, nil
	}
	log.Debug("get block header by hash", "headerHash", header.Hash(), "blockHash", hash)

	blockHeader := &account.BlockHeader{
		Hash:        header.Hash().String(),
		ParentHash:  header.ParentHash.String(),
		UncleHash:   header.UncleHash.String(),
		CoinBase:    header.Coinbase.String(),
		Root:        header.Root.String(),
		TxHash:      header.TxHash.String(),
		ReceiptHash: header.ReceiptHash.String(),
		Difficulty:  header.Difficulty.String(),
		Number:      header.Number.String(),
		GasLimit:    header.GasLimit,
		GasUsed:     header.GasUsed,
		Time:        header.Time,
		Extra:       hex.EncodeToString(header.Extra),
		MixDigest:   header.MixDigest.String(),
		Nonce:       strconv.FormatUint(header.Nonce.Uint64(), 10),
		BaseFee:     header.BaseFee.String(),
	}

	// Handle optional fields
	if header.ParentBeaconRoot != nil {
		blockHeader.ParentBeaconRoot = header.ParentBeaconRoot.String()
	}
	if header.WithdrawalsHash != nil {
		blockHeader.WithdrawalsHash = header.WithdrawalsHash.String()
	}

	return &account.BlockHeaderResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "get block header by hash success",
		BlockHeader: blockHeader,
	}, nil
}

func (c ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	blockInfo, err := c.ethClient.BlockHeaderByNumber(big.NewInt(req.Height))
	if err != nil {
		log.Error("get block header by number fail", "error", err)
		return &account.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block header by number fail",
		}, nil
	}
	blockHeader := &account.BlockHeader{
		Hash:        blockInfo.Hash().String(),
		ParentHash:  blockInfo.ParentHash.String(),
		UncleHash:   blockInfo.UncleHash.String(),
		CoinBase:    blockInfo.Coinbase.String(),
		Root:        blockInfo.Root.String(),
		TxHash:      blockInfo.TxHash.String(),
		ReceiptHash: blockInfo.ReceiptHash.String(),
		Difficulty:  blockInfo.Difficulty.String(),
		Number:      blockInfo.Number.String(),
		GasLimit:    blockInfo.GasLimit,
		GasUsed:     blockInfo.GasUsed,
		Time:        blockInfo.Time,
		Extra:       hex.EncodeToString(blockInfo.Extra),
		MixDigest:   blockInfo.MixDigest.String(),
		Nonce:       strconv.FormatUint(blockInfo.Nonce.Uint64(), 10),
		BaseFee:     blockInfo.BaseFee.String(),
	}

	// Handle optional fields
	if blockInfo.ParentBeaconRoot != nil {
		blockHeader.ParentBeaconRoot = blockInfo.ParentBeaconRoot.String()
	}
	if blockInfo.WithdrawalsHash != nil {
		blockHeader.WithdrawalsHash = blockInfo.WithdrawalsHash.String()
	}
	if blockInfo.BlobGasUsed != nil {
		blockHeader.BlobGasUsed = *blockInfo.BlobGasUsed
	}
	if blockInfo.ExcessBlobGas != nil {
		blockHeader.ExcessBlobGas = *blockInfo.ExcessBlobGas
	}

	return &account.BlockHeaderResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "get block header by number success",
		BlockHeader: blockHeader,
	}, nil
}

func (c ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	nonceResult, err := c.ethClient.TxCountByAddress(common.HexToAddress(req.Address))
	if err != nil {
		log.Error("get nonce by address fail", "error", err)
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

	sequence := strconv.FormatUint(uint64(nonceResult), 10)

	return &account.AccountResponse{
		Code:          common2.ReturnCode_SUCCESS,
		Msg:           "get account response success",
		AccountNumber: "0",
		Sequence:      sequence,
		Balance:       balanceResult.BalanceStr,
	}, nil
}

func (c ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
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

func (c ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	transaction, err := c.ethClient.SendRawTransaction(req.RawTx)
	if err != nil {
		return &account.SendTxResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Send tx error" + err.Error(),
		}, err
	}
	return &account.SendTxResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "send tx success",
		TxHash: transaction.String(),
	}, nil
}

func (c ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
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
				Hash:            txs[i].TxId,
				To:              txs[i].To,
				From:            txs[i].From,
				Fee:             txs[i].TxId,
				Status:          account.TxStatus_Success,
				Value:           txs[i].Amount,
				Type:            1,
				Height:          txs[i].Height,
				ContractAddress: txs[i].TokenContractAddress,
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

func (c ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	log.Info("GetTxByHash request", "hash", req.Hash)

	// 直接通过 RPC 调用获取交易信息
	ctxwt, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var tx map[string]interface{}
	err := c.ethClient.(interface{ GetRPC() erc20base.RPC }).GetRPC().CallContext(ctxwt, &tx, "eth_getTransactionByHash", common.HexToHash(req.Hash))
	if err != nil {
		log.Error("RPC call error", "err", err)
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "RPC call error: " + err.Error(),
		}, nil
	}
	if tx == nil {
		log.Error("transaction not found")
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Zksync Tx NotFound",
		}, nil
	}

	log.Info("Raw transaction data", "tx", tx)

	// 获取交易回执
	var receipt map[string]interface{}
	err = c.ethClient.(interface{ GetRPC() erc20base.RPC }).GetRPC().CallContext(ctxwt, &receipt, "eth_getTransactionReceipt", common.HexToHash(req.Hash))
	if err != nil {
		log.Error("get receipt error", "err", err)
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Get receipt error: " + err.Error(),
		}, nil
	}
	if receipt == nil {
		log.Error("receipt not found")
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Receipt not found",
		}, nil
	}

	log.Info("Raw receipt data", "receipt", receipt)

	// 解析交易数据
	var (
		from            = tx["from"].(string)
		to              = tx["to"].(string)
		value           = tx["value"].(string)
		gasPrice        = tx["gasPrice"].(string)
		blockNumber     = receipt["blockNumber"].(string)
		status          = receipt["status"].(string)
		contractAddress = receipt["contractAddress"]
	)

	// 检查合约地址
	var toAddress string
	var tokenAddress string
	var txValue string

	if contractAddress != nil && contractAddress.(string) != "0x0000000000000000000000000000000000000000" {
		// 这是一个合约创建交易
		toAddress = contractAddress.(string)
		tokenAddress = contractAddress.(string)
		txValue = "0"
	} else {
		// 这是一个普通交易
		toAddress = to
		tokenAddress = "0x0000000000000000000000000000000000000000"
		txValue = value
	}

	// 确定交易状态
	var txStatus account.TxStatus
	if status == "0x1" {
		txStatus = account.TxStatus_Success
	} else {
		txStatus = account.TxStatus_Failed
	}

	return &account.TxHashResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get transaction success",
		Tx: &account.TxMessage{
			Hash:            req.Hash,
			From:            from,
			To:              toAddress,
			Value:           txValue,
			Fee:             gasPrice,
			Status:          txStatus,
			Height:          blockNumber,
			ContractAddress: tokenAddress,
		},
	}, nil
}

func (c ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	startBlock := new(big.Int)
	endBlock := new(big.Int)
	startBlock.SetString(req.Start, 10)
	endBlock.SetString(req.End, 10)
	blockRange, err := c.ethClient.BlockHeadersByRange(startBlock, endBlock, 324)
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
			Hash:        block.Hash().String(),
			ParentHash:  block.ParentHash.String(),
			UncleHash:   block.UncleHash.String(),
			CoinBase:    block.Coinbase.String(),
			Root:        block.Root.String(),
			TxHash:      block.TxHash.String(),
			ReceiptHash: block.ReceiptHash.String(),
			Difficulty:  block.Difficulty.String(),
			Number:      block.Number.String(),
			GasLimit:    block.GasLimit,
			GasUsed:     block.GasUsed,
			Time:        block.Time,
			Extra:       hex.EncodeToString(block.Extra),
			MixDigest:   block.MixDigest.String(),
			Nonce:       strconv.FormatUint(block.Nonce.Uint64(), 10),
			BaseFee:     block.BaseFee.String(),
		}

		// Handle optional fields
		if block.ParentBeaconRoot != nil {
			blockItem.ParentBeaconRoot = block.ParentBeaconRoot.String()
		}
		if block.WithdrawalsHash != nil {
			blockItem.WithdrawalsHash = block.WithdrawalsHash.String()
		}
		if block.BlobGasUsed != nil {
			blockItem.BlobGasUsed = *block.BlobGasUsed
		}
		if block.ExcessBlobGas != nil {
			blockItem.ExcessBlobGas = *block.ExcessBlobGas
		}

		blockHeaderList = append(blockHeaderList, blockItem)
	}
	return &account.BlockByRangeResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "get block range success",
		BlockHeader: blockHeaderList,
	}, nil
}

func (c ChainAdaptor) BuildUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	response := &account.UnSignTransactionResponse{
		Code: common2.ReturnCode_ERROR,
	}

	dFeeTx, _, err := c.buildDynamicFeeTx(req.Base64Tx)
	if err != nil {
		return nil, err
	}

	log.Info("zksync BuildUnSignTransaction", "dFeeTx", util.ToJSONString(dFeeTx))

	// Create unsigned transaction
	rawTx, err := evmbase.CreateEip1559UnSignTx(dFeeTx, dFeeTx.ChainID)
	if err != nil {
		log.Error("create un sign tx fail", "err", err)
		response.Msg = "get un sign tx fail"
		return response, nil
	}

	log.Info("zksync BuildUnSignTransaction", "rawTx", rawTx)
	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "create un sign tx success"
	response.UnSignTx = rawTx
	return response, nil
}

func (c ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	response := &account.SignedTransactionResponse{
		Code: common2.ReturnCode_ERROR,
	}

	dFeeTx, dynamicFeeTx, err := c.buildDynamicFeeTx(req.Base64Tx)
	if err != nil {
		log.Error("buildDynamicFeeTx failed", "err", err)
		return nil, err
	}

	log.Info("zksync BuildSignedTransaction", "dFeeTx", util.ToJSONString(dFeeTx))
	log.Info("zksync BuildSignedTransaction", "dynamicFeeTx", util.ToJSONString(dynamicFeeTx))
	log.Info("zksync BuildSignedTransaction", "req.Signature", req.Signature)

	// Decode signature and create signed transaction
	inputSignatureByteList, err := hex.DecodeString(req.Signature)
	if err != nil {
		log.Error("decode signature failed", "err", err)
		return nil, fmt.Errorf("invalid signature: %w", err)
	}

	signer, signedTx, rawTx, txHash, err := evmbase.CreateEip1559SignedTx(dFeeTx, inputSignatureByteList, dFeeTx.ChainID)
	if err != nil {
		log.Error("create signed tx fail", "err", err)
		return nil, fmt.Errorf("create signed tx fail: %w", err)
	}

	log.Info("zksync BuildSignedTransaction", "rawTx", rawTx)

	// Verify sender
	sender, err := types.Sender(signer, signedTx)
	if err != nil {
		log.Error("recover sender failed", "err", err)
		return nil, fmt.Errorf("recover sender failed: %w", err)
	}

	if sender.Hex() != dynamicFeeTx.FromAddress {
		log.Error("sender mismatch",
			"expected", dynamicFeeTx.FromAddress,
			"got", sender.Hex(),
		)
		return nil, fmt.Errorf("sender address mismatch: expected %s, got %s",
			dynamicFeeTx.FromAddress,
			sender.Hex(),
		)
	}

	log.Info("zksync BuildSignedTransaction", "sender", sender.Hex())

	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = txHash
	response.SignedTx = rawTx
	return response, nil
}

func (c ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	// 1. 解码原始交易数据
	tx := common.FromHex(req.RawTx)
	if len(tx) == 0 {
		return &account.DecodeTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "invalid raw transaction",
		}, nil
	}

	// 2. 尝试解析交易
	parsedTx := new(types.Transaction)
	err := parsedTx.UnmarshalBinary(tx)
	if err != nil {
		return &account.DecodeTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "failed to parse transaction: " + err.Error(),
		}, nil
	}

	// 3. 构建交易数据结构
	txData := &evmbase.Eip1559DynamicFeeTx{
		ChainId:              parsedTx.ChainId().String(),
		Nonce:                uint64(parsedTx.Nonce()),
		MaxPriorityFeePerGas: parsedTx.GasTipCap().String(),
		MaxFeePerGas:         parsedTx.GasFeeCap().String(),
		GasLimit:             parsedTx.Gas(),
		ToAddress:            parsedTx.To().Hex(),
		Amount:               parsedTx.Value().String(),
	}

	// 4. 转换为 JSON 然后 Base64 编码
	txJson, err := json.Marshal(txData)
	if err != nil {
		return &account.DecodeTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "failed to encode transaction data: " + err.Error(),
		}, nil
	}

	return &account.DecodeTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "decode transaction success",
		Base64Tx: base64.StdEncoding.EncodeToString(txJson),
	}, nil
}

// VerifyTransactionRequestWrapper wraps the VerifyTransactionRequest with additional fields
type VerifyTransactionRequestWrapper struct {
	*account.VerifyTransactionRequest
	TxHash string
}

// VerifySignedTransaction verifies the signature of a transaction
// TODO: Currently the signature field combines both txHash and signature in format "txHash:signature"
// This is a temporary solution and should be improved by:
//  1. Modifying the protobuf definition of VerifyTransactionRequest to add a separate txHash field
//  2. Update the proto file to include:
//     message VerifyTransactionRequest {
//     string public_key = 1;
//     string signature = 2;
//     string tx_hash = 3;  // New field
//     }
//  3. This will make the API more clear and maintainable
func (c ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	if req == nil {
		return &account.VerifyTransactionResponse{
			Code:   common2.ReturnCode_ERROR,
			Msg:    "request is nil",
			Verify: false,
		}, nil
	}

	if req.PublicKey == "" {
		return &account.VerifyTransactionResponse{
			Code:   common2.ReturnCode_ERROR,
			Msg:    "public key is empty",
			Verify: false,
		}, nil
	}

	if req.Signature == "" {
		return &account.VerifyTransactionResponse{
			Code:   common2.ReturnCode_ERROR,
			Msg:    "signature is empty",
			Verify: false,
		}, nil
	}

	// Split the signature field to get the transaction hash and signature
	parts := strings.Split(req.Signature, ":")
	if len(parts) != 2 {
		return &account.VerifyTransactionResponse{
			Code:   common2.ReturnCode_ERROR,
			Msg:    "invalid signature format, expected 'txHash:signature'",
			Verify: false,
		}, nil
	}

	txHash := parts[0]
	signature := parts[1]

	// Convert public key from hex to bytes
	publicKeyBytes, err := hex.DecodeString(req.PublicKey)
	if err != nil {
		return &account.VerifyTransactionResponse{
			Code:   common2.ReturnCode_ERROR,
			Msg:    "failed to decode public key: " + err.Error(),
			Verify: false,
		}, nil
	}

	// Convert transaction hash from hex to bytes
	txHashBytes, err := hex.DecodeString(strings.TrimPrefix(txHash, "0x"))
	if err != nil {
		return &account.VerifyTransactionResponse{
			Code:   common2.ReturnCode_ERROR,
			Msg:    "failed to decode transaction hash: " + err.Error(),
			Verify: false,
		}, nil
	}

	// Convert signature from hex to bytes
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return &account.VerifyTransactionResponse{
			Code:   common2.ReturnCode_ERROR,
			Msg:    "failed to decode signature: " + err.Error(),
			Verify: false,
		}, nil
	}

	// Verify the signature using the first 64 bytes (R,S) of the signature
	verified := crypto.VerifySignature(publicKeyBytes, txHashBytes, signatureBytes[:64])

	return &account.VerifyTransactionResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "verify transaction success",
		Verify: verified,
	}, nil
}

func (c ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	return &account.ExtraDataResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "get extra data success",
		Value: "not data",
	}, nil
}

func (c ChainAdaptor) GetNftListByAddress(req *account.NftAddressRequest) (*account.NftAddressResponse, error) {
	panic("implement me")
}

// buildDynamicFeeTx build eip1559 tx
func (c ChainAdaptor) buildDynamicFeeTx(base64Tx string) (*types.DynamicFeeTx, *evmbase.Eip1559DynamicFeeTx, error) {
	// Decode base64 string
	txReqJsonByte, err := base64.StdEncoding.DecodeString(base64Tx)
	if err != nil {
		log.Error("decode string fail", "err", err)
		return nil, nil, err
	}

	var dynamicFeeTx evmbase.Eip1559DynamicFeeTx
	if err := json.Unmarshal(txReqJsonByte, &dynamicFeeTx); err != nil {
		log.Error("parse json fail", "err", err)
		return nil, nil, err
	}

	chainID := new(big.Int)
	maxPriorityFeePerGas := new(big.Int)
	maxFeePerGas := new(big.Int)
	amount := new(big.Int)

	if _, ok := chainID.SetString(dynamicFeeTx.ChainId, 10); !ok {
		return nil, nil, fmt.Errorf("invalid chain ID: %s", dynamicFeeTx.ChainId)
	}
	if _, ok := maxPriorityFeePerGas.SetString(dynamicFeeTx.MaxPriorityFeePerGas, 10); !ok {
		return nil, nil, fmt.Errorf("invalid max priority fee: %s", dynamicFeeTx.MaxPriorityFeePerGas)
	}
	if _, ok := maxFeePerGas.SetString(dynamicFeeTx.MaxFeePerGas, 10); !ok {
		return nil, nil, fmt.Errorf("invalid max fee: %s", dynamicFeeTx.MaxFeePerGas)
	}
	if _, ok := amount.SetString(dynamicFeeTx.Amount, 10); !ok {
		return nil, nil, fmt.Errorf("invalid amount: %s", dynamicFeeTx.Amount)
	}

	// 4. Handle addresses and data
	toAddress := common.HexToAddress(dynamicFeeTx.ToAddress)
	var finalToAddress common.Address
	var finalAmount *big.Int
	var buildData []byte
	log.Info("contract address check", "contractAddress", dynamicFeeTx.ContractAddress, "isEthTransfer", isEthTransfer(&dynamicFeeTx))

	// 5. Handle contract interaction vs direct transfer
	if isEthTransfer(&dynamicFeeTx) {
		finalToAddress = toAddress
		finalAmount = amount
	} else {
		contractAddress := common.HexToAddress(dynamicFeeTx.ContractAddress)
		buildData = evmbase.BuildErc20Data(toAddress, amount)
		finalToAddress = contractAddress
		finalAmount = big.NewInt(0)
	}

	// 6. Create dynamic fee transaction
	dFeeTx := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     dynamicFeeTx.Nonce,
		GasTipCap: maxPriorityFeePerGas,
		GasFeeCap: maxFeePerGas,
		Gas:       dynamicFeeTx.GasLimit,
		To:        &finalToAddress,
		Value:     finalAmount,
		Data:      buildData,
	}

	return dFeeTx, &dynamicFeeTx, nil
}

// 判断是否为 ETH 转账
func isEthTransfer(tx *evmbase.Eip1559DynamicFeeTx) bool {
	// 检查合约地址是否为空或零地址
	if tx.ContractAddress == "" ||
		tx.ContractAddress == "0x0000000000000000000000000000000000000000" ||
		tx.ContractAddress == "0x00" {
		return true
	}
	return false
}
