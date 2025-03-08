package binance

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

	"github.com/ethereum/go-ethereum"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/status-im/keycard-go/hexutils"

	account2 "github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/chain/evmbase"
	erc20Base "github.com/dapplink-labs/wallet-chain-account/chain/evmbase"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

const ChainName = "BscChain"

type ChainAdaptor struct {
	ethClient     erc20Base.EthClient
	ethDataClient *erc20Base.EthData
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	ethClient, err := erc20Base.DialEthClient(context.Background(), conf.WalletNode.Bsc.RpcUrl)
	if err != nil {
		return nil, err
	}
	ethDataClient, err := erc20Base.NewEthDataClient(conf.WalletNode.Bsc.DataApiUrl, conf.WalletNode.Bsc.DataApiKey, time.Duration(conf.WalletNode.Bsc.TimeOut))
	if err != nil {
		return nil, err
	}
	return &ChainAdaptor{
		ethClient:     ethClient,
		ethDataClient: ethDataClient,
	}, nil
}

func (c *ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {

	if req.Chain != ChainName {
		return &account.SupportChainsResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     "not Support Chain",
			Support: false,
		}, nil
	}

	return &account.SupportChainsResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "Support Chain",
		Support: true,
	}, nil
}
func (c *ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {

	publicKeyBytes, err := hex.DecodeString(req.PublicKey)
	if err != nil {
		log.Error("decode public key failed:", err)
		return &account.ConvertAddressResponse{
			Code:    common.ReturnCode_SUCCESS,
			Msg:     "convert address fail",
			Address: common2.Address{}.String(),
		}, nil
	}
	addressCommon := common2.BytesToAddress(crypto.Keccak256(publicKeyBytes[1:])[12:])
	return &account.ConvertAddressResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "convert address success",
		Address: addressCommon.String(),
	}, nil
}
func (c *ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	if len(req.Address) != 42 || !strings.HasPrefix(req.Address, "0x") {
		return &account.ValidAddressResponse{
			Code:  common.ReturnCode_SUCCESS,
			Msg:   "invalid address",
			Valid: false,
		}, nil
	}
	ok := regexp.MustCompile("^[0-9a-fA-F]{40}$").MatchString(req.Address[2:])
	if ok {
		return &account.ValidAddressResponse{
			Code:  common.ReturnCode_SUCCESS,
			Msg:   "valid address",
			Valid: true,
		}, nil
	} else {
		return &account.ValidAddressResponse{
			Code:  common.ReturnCode_SUCCESS,
			Msg:   "invalid address",
			Valid: false,
		}, nil
	}
}
func (c *ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	var blockNumber *big.Int
	if req.Height == 0 {
		blockNumber = nil // return latest block
	} else {
		blockNumber = big.NewInt(req.Height) // return special block by number
	}
	rsp, err := c.ethClient.BlockByNumber(blockNumber)
	if err != nil {
		log.Error("GetBlockByNumber fail:", err)
		return &account.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	height, _ := rsp.NumberUint64()
	var blockTxList []*account.BlockInfoTransactionList
	for _, tx := range rsp.Transactions {
		blockTxList = append(blockTxList, &account.BlockInfoTransactionList{
			From:   tx.From,
			To:     tx.To,
			Hash:   tx.Hash,
			Amount: tx.Value,
			Height: height,
		})
	}

	return &account.BlockResponse{
		Code:         common.ReturnCode_SUCCESS,
		Msg:          "GetBlockByNumber success",
		Height:       req.Height,
		Hash:         rsp.Hash.String(),
		BaseFee:      rsp.BaseFee,
		Transactions: blockTxList,
	}, nil
}
func (c *ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {

	rsp, err := c.ethClient.BlockByHash(common2.HexToHash(req.Hash))
	if err != nil {
		log.Error("GetBlockByHash fail:", err)
		return &account.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	height, _ := rsp.NumberUint64()

	var blockTxList []*account.BlockInfoTransactionList
	for _, tx := range rsp.Transactions {
		blockTxList = append(blockTxList, &account.BlockInfoTransactionList{
			From:   tx.From,
			To:     tx.To,
			Hash:   tx.Hash,
			Amount: tx.Value,
			Height: height,
		})
	}

	return &account.BlockResponse{
		Code:         common.ReturnCode_SUCCESS,
		Msg:          "GetBlockByNumber success",
		Height:       int64(height),
		Hash:         rsp.Hash.String(),
		BaseFee:      rsp.BaseFee,
		Transactions: blockTxList,
	}, nil
}
func (c *ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {

	rsp, err := c.ethClient.BlockHeaderByHash(common2.HexToHash(req.Hash))
	if err != nil {
		log.Error("GetBlockHeaderByHash fail:", err)
		return &account.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	blockHeader := &account.BlockHeader{
		Hash:        rsp.Hash().String(),
		ParentHash:  rsp.ParentHash.String(),
		UncleHash:   rsp.UncleHash.String(),
		CoinBase:    rsp.Coinbase.String(),
		Root:        rsp.Root.String(),
		TxHash:      rsp.TxHash.String(),
		ReceiptHash: rsp.ReceiptHash.String(),
		Difficulty:  rsp.Difficulty.String(),
		Number:      rsp.Number.String(),
		GasLimit:    rsp.GasLimit,
		GasUsed:     rsp.GasUsed,
		Time:        rsp.Time,
		Extra:       hex.EncodeToString(rsp.Extra),
		MixDigest:   rsp.MixDigest.String(),
		Nonce:       strconv.FormatUint(rsp.Nonce.Uint64(), 10),
		BaseFee:     rsp.BaseFee.String(),

		ParentBeaconRoot: getSafeHashString(rsp.ParentBeaconRoot),
		WithdrawalsHash:  getSafeHashString(rsp.WithdrawalsHash),
		BlobGasUsed:      getSafeUint64Ptr(rsp.BlobGasUsed),
		ExcessBlobGas:    getSafeUint64Ptr(rsp.ExcessBlobGas),
	}
	return &account.BlockHeaderResponse{
		Code:        common.ReturnCode_SUCCESS,
		Msg:         "get latest block header success",
		BlockHeader: blockHeader,
	}, nil
}
func (c *ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	var blockNumber *big.Int
	if req.Height == 0 {
		blockNumber = nil // return latest block
	} else {
		blockNumber = big.NewInt(req.Height) // return special block by number
	}
	rsp, err := c.ethClient.BlockHeaderByNumber(blockNumber)
	if err != nil {
		log.Error("GetBlockHeaderByNumber fail:", err)
		return &account.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}

	blockHeader := &account.BlockHeader{
		Hash:        rsp.Hash().String(),
		ParentHash:  rsp.ParentHash.String(),
		UncleHash:   rsp.UncleHash.String(),
		CoinBase:    rsp.Coinbase.String(),
		Root:        rsp.Root.String(),
		TxHash:      rsp.TxHash.String(),
		ReceiptHash: rsp.ReceiptHash.String(),
		Difficulty:  rsp.Difficulty.String(),
		Number:      rsp.Number.String(),
		GasLimit:    rsp.GasLimit,
		GasUsed:     rsp.GasUsed,
		Time:        rsp.Time,
		Extra:       hex.EncodeToString(rsp.Extra),
		MixDigest:   rsp.MixDigest.String(),
		Nonce:       strconv.FormatUint(rsp.Nonce.Uint64(), 10),
		BaseFee:     rsp.BaseFee.String(),

		ParentBeaconRoot: getSafeHashString(rsp.ParentBeaconRoot),
		WithdrawalsHash:  getSafeHashString(rsp.WithdrawalsHash),
		BlobGasUsed:      getSafeUint64Ptr(rsp.BlobGasUsed),
		ExcessBlobGas:    getSafeUint64Ptr(rsp.ExcessBlobGas),
	}
	return &account.BlockHeaderResponse{
		Code:        common.ReturnCode_SUCCESS,
		Msg:         "get latest block header success",
		BlockHeader: blockHeader,
	}, nil
}
func (c *ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {

	rspNonce, err := c.ethClient.TxCountByAddress(common2.HexToAddress(req.Address))
	if err != nil {
		log.Error("GetAccountByAddress fail:", err)
		return &account.AccountResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	sequence := strconv.FormatUint(uint64(rspNonce), 10)

	balanceRsp, err := c.ethDataClient.GetBalanceByAddress(req.ContractAddress, req.Address)
	if err != nil {
		log.Error("GetAccountByAddress fail:", err)
		return &account.AccountResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	balanceStr := "0"
	if balanceRsp.Balance != nil && balanceRsp.Balance.Int() != nil {
		balanceStr = balanceRsp.Balance.Int().String()
	}

	return &account.AccountResponse{
		Code:     common.ReturnCode_SUCCESS,
		Sequence: sequence,
		Balance:  balanceStr,
	}, nil
}
func (c *ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	gasPrice, err := c.ethClient.SuggestGasPrice()
	if err != nil {
		log.Error("get gas price failed", "err", err)
		return &account.FeeResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get suggest gas price fail",
		}, nil
	}
	gasTipCap, err := c.ethClient.SuggestGasTipCap()
	if err != nil {
		log.Error("get gas price failed", "err", err)
		return &account.FeeResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get suggest gas price fail",
		}, nil
	}
	return &account.FeeResponse{
		Code:      common.ReturnCode_SUCCESS,
		Msg:       "get gas price success",
		SlowFee:   gasPrice.String() + "|" + gasTipCap.String(),
		NormalFee: gasPrice.String() + "|" + gasTipCap.String() + "|" + "*2",
		FastFee:   gasPrice.String() + "|" + gasTipCap.String() + "|" + "*3",
	}, nil
}
func (c *ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	transaction, err := c.ethClient.SendRawTransaction(req.RawTx)
	if err != nil {
		return &account.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "Send tx error" + err.Error(),
		}, err
	}
	return &account.SendTxResponse{
		Code:   common.ReturnCode_SUCCESS,
		Msg:    "send tx success",
		TxHash: transaction.String(),
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
			Code: common.ReturnCode_ERROR,
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
			Code: common.ReturnCode_SUCCESS,
			Msg:  "get tx list success",
			Tx:   list,
		}, nil
	}
}
func (c *ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {

	rsp, err := c.ethClient.TxByHash(common2.HexToHash(req.Hash))
	if err != nil {
		return &account.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get tx by hash fail",
		}, err
	}

	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return &account.TxHashResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  "Ethereum Tx NotFound",
			}, nil
		}
		log.Error("get transaction error", "err", err)
		return &account.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "Ethereum Tx NotFound",
		}, nil
	}

	receipt, err := c.ethClient.TxReceiptByHash(common2.HexToHash(req.Hash))
	if err != nil {
		log.Error("get transaction receipt error", "err", err)
		return &account.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "Get transaction receipt error",
		}, nil
	}

	var beforeToAddress string
	var beforeTokenAddress string
	var beforeValue *big.Int

	code, err := c.ethClient.EthGetCode(common2.HexToAddress(rsp.To().String()))
	if err != nil {
		log.Info("Get account code fail", "err", err)
		return nil, err
	}

	if code == "contract" {
		inputData := hexutil.Encode(rsp.Data()[:])
		if len(inputData) >= 138 && inputData[:10] == "0xa9059cbb" {
			beforeToAddress = "0x" + inputData[34:74]
			trimHex := strings.TrimLeft(inputData[74:138], "0")
			rawValue, _ := hexutil.DecodeBig("0x" + trimHex)
			beforeTokenAddress = rsp.To().String()
			beforeValue = decimal.NewFromBigInt(rawValue, 0).BigInt()
		}
	} else {
		beforeToAddress = rsp.To().String()
		beforeTokenAddress = common2.Address{}.String()
		beforeValue = rsp.Value()
	}

	var txStatus account.TxStatus
	if receipt.Status == 1 {
		txStatus = account.TxStatus_Success
	} else {
		txStatus = account.TxStatus_Failed
	}

	tx := &account.TxMessage{
		Hash:            rsp.Hash().Hex(),
		Index:           uint32(receipt.TransactionIndex),
		From:            "",
		To:              beforeToAddress,
		Value:           beforeValue.String(),
		Fee:             rsp.GasFeeCap().String(),
		Status:          txStatus,
		Type:            0,
		Height:          receipt.BlockNumber.String(),
		ContractAddress: beforeTokenAddress,
		Data:            hexutils.BytesToHex(rsp.Data()),
	}

	return &account.TxHashResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "get tx by hash success",
		Tx:   tx,
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
			Code: common.ReturnCode_ERROR,
			Msg:  "get block range fail",
		}, err
	}
	blockHeaderList := make([]*account.BlockHeader, 0, len(blockRange))
	for _, block := range blockRange {
		blockItem := &account.BlockHeader{
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

			ParentBeaconRoot: getSafeHashString(block.ParentBeaconRoot),
			WithdrawalsHash:  getSafeHashString(block.WithdrawalsHash),
			BlobGasUsed:      getSafeUint64Ptr(block.BlobGasUsed),
			ExcessBlobGas:    getSafeUint64Ptr(block.ExcessBlobGas),
		}
		blockHeaderList = append(blockHeaderList, blockItem)
	}
	return &account.BlockByRangeResponse{
		Code:        common.ReturnCode_SUCCESS,
		Msg:         "get block range success",
		BlockHeader: blockHeaderList,
	}, nil
}
func (c *ChainAdaptor) BuildUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {

	dyTx, _, err := c.buildDynamicFeeTx(req.Base64Tx)
	if err != nil {
		log.Error("create unsign transaction fail", "err", err)
		return &account.UnSignTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "create unsign transaction fail",
		}, err
	}

	rwaTx, err := evmbase.CreateEip1559UnSignTx(dyTx, dyTx.ChainID)
	if err != nil {
		log.Error("create unsign transaction fail", "err", err)
		return &account.UnSignTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "create unsign transaction fail",
		}, err
	}

	return &account.UnSignTransactionResponse{
		Code:     common.ReturnCode_SUCCESS,
		Msg:      "create unsign transaction success",
		UnSignTx: rwaTx,
	}, nil
}
func (c *ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {

	dyTx, txEip1559Data, err := c.buildDynamicFeeTx(req.Base64Tx)
	if err != nil {
		log.Error("create unsign transaction fail", "err", err)
		return &account.SignedTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "create unsign transaction fail",
		}, err
	}

	// Decode signature and create signed transaction
	inputSignatureByteList, err := hex.DecodeString(txEip1559Data.Signature)
	if err != nil {
		log.Error("decode signature failed", "err", err)
		return &account.SignedTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "create unsign transaction fail",
		}, fmt.Errorf("invalid signature: %w", err)
	}

	signer, signedTx, rawTx, txHash, err := evmbase.CreateEip1559SignedTx(dyTx, inputSignatureByteList, dyTx.ChainID)
	if err != nil {
		log.Error("create signed tx fail", "err", err)
		return &account.SignedTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "create signed tx fail",
		}, fmt.Errorf("create signed tx fail: %w", err)
	}

	// Verify sender
	sender, err := types.Sender(signer, signedTx)
	if err != nil {
		log.Error("recover sender failed", "err", err)
		return &account.SignedTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "sender error",
		}, fmt.Errorf("recover sender failed: %w", err)
	}

	if strings.ToLower(sender.Hex()) != strings.ToLower(txEip1559Data.FromAddress) {
		log.Error("sender mismatch", "expected", txEip1559Data.FromAddress, "got", sender.Hex())
		return &account.SignedTransactionResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  "sign object error",
			}, fmt.Errorf("sender address mismatch: expected %s, got %s",
				txEip1559Data.FromAddress,
				sender.Hex(),
			)
	}
	return &account.SignedTransactionResponse{
		Code:     common.ReturnCode_SUCCESS,
		Msg:      txHash,
		SignedTx: rawTx,
	}, nil
}
func (c *ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	return &account.DecodeTransactionResponse{
		Code:     common.ReturnCode_SUCCESS,
		Msg:      "verify tx success",
		Base64Tx: "0x000000",
	}, nil
}
func (c *ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	return &account.VerifyTransactionResponse{
		Code:   common.ReturnCode_SUCCESS,
		Msg:    "verify tx success",
		Verify: true,
	}, nil
}
func (c *ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	return &account.ExtraDataResponse{
		Code:  common.ReturnCode_SUCCESS,
		Msg:   "get extra data success",
		Value: "not data",
	}, nil
}

func (c *ChainAdaptor) GetNftListByAddress(req *account.NftAddressRequest) (*account.NftAddressResponse, error) {
	panic("implement me")
}

func getSafeUint64Ptr(ptr *uint64) uint64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func getSafeHashString(hash *common2.Hash) string {
	if hash == nil {
		return common2.Hash{}.String()
	}
	return hash.String()
}

// buildDynamicFeeTx build eip1559 tx
func (c *ChainAdaptor) buildDynamicFeeTx(base64Tx string) (*types.DynamicFeeTx, *evmbase.Eip1559DynamicFeeTx, error) {
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
	toAddress := common2.HexToAddress(dynamicFeeTx.ToAddress)
	var finalToAddress common2.Address
	var finalAmount *big.Int
	var buildData []byte
	log.Info("contract address check", "contractAddress", dynamicFeeTx.ContractAddress, "isEthTransfer", isEthTransfer(&dynamicFeeTx))

	// 5. Handle contract interaction vs direct transfer
	if isEthTransfer(&dynamicFeeTx) {
		finalToAddress = toAddress
		finalAmount = amount
	} else {
		contractAddress := common2.HexToAddress(dynamicFeeTx.ContractAddress)
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
