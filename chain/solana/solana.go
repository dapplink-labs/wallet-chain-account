package solana

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/token"

	"github.com/ethereum/go-ethereum/log"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"

	account2 "github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

const ChainName = "Solana"

const (
	MaxBlockRange = 1000
)

type ChainAdaptor struct {
	solCli    SolClient
	sdkClient *rpc.Client
	solData   *SolData
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	rpcUrl := conf.WalletNode.Sol.RpcUrl

	solHttpCli, err := NewSolHttpClient(rpcUrl)
	if err != nil {
		return nil, err
	}
	dataApiUrl := conf.WalletNode.Sol.DataApiUrl
	dataApiKey := conf.WalletNode.Sol.DataApiKey
	dataApiTimeOut := conf.WalletNode.Sol.TimeOut
	solData, err := NewSolScanClient(dataApiUrl, dataApiKey, time.Duration(dataApiTimeOut))
	if err != nil {
		return nil, err
	}

	sdkClient := rpc.New(rpcUrl)

	return &ChainAdaptor{
		solCli:    solHttpCli,
		sdkClient: sdkClient,
		solData:   solData,
	}, nil
}

func (c *ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	response := &account.SupportChainsResponse{
		Code:    common2.ReturnCode_ERROR,
		Msg:     "",
		Support: false,
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		err := fmt.Errorf("GetSupportChains validateChainAndNetwork fail, err msg = %s", msg)
		log.Error("err", err)
		response.Msg = err.Error()
		return response, err
	}

	response.Msg = "Support this chain"
	response.Code = common2.ReturnCode_SUCCESS
	response.Support = true
	return response, nil
}

func (c *ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	response := &account.ConvertAddressResponse{
		Code:    common2.ReturnCode_ERROR,
		Msg:     "",
		Address: "",
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		err := fmt.Errorf("ConvertAddress validateChainAndNetwork fail, err msg = %s", msg)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	pubKeyHex := req.PublicKey
	if ok, msg := validatePublicKey(pubKeyHex); !ok {
		err := fmt.Errorf("ConvertAddress validatePublicKey fail, err msg = %s", msg)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	accountAddress, err := PubKeyHexToAddress(pubKeyHex)
	if err != nil {
		err := fmt.Errorf("ConvertAddress PubKeyHexToAddress failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "convert address success"
	response.Address = accountAddress
	return response, nil
}

func (c *ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	response := &account.ValidAddressResponse{
		Code:  common2.ReturnCode_ERROR,
		Msg:   "",
		Valid: false,
	}

	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		err := fmt.Errorf("ValidAddress validateChainAndNetwork failed: %s", msg)
		log.Error("err", err)
		response.Msg = err.Error()
		return response, err
	}
	address := req.Address
	if len(address) == 0 {
		err := fmt.Errorf("ValidAddress address is empty")
		log.Error("err", err)
		response.Msg = err.Error()
		return response, err
	}
	if len(address) != 43 && len(address) != 44 {
		err := fmt.Errorf("invalid Solana address length: expected 43 or 44 characters, got %d", len(address))
		response.Msg = err.Error()
		return response, err
	}
	response.Code = common2.ReturnCode_SUCCESS
	response.Valid = true
	return response, nil
}

func (c *ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	response := &account.BlockResponse{
		Code:         common2.ReturnCode_ERROR,
		Msg:          "",
		Height:       0,
		Hash:         "",
		BaseFee:      "",
		Transactions: nil,
	}
	if ok, msg := validateChainAndNetwork(req.Chain, ""); !ok {
		err := fmt.Errorf("GetBlockByNumber validateChainAndNetwork failed: %s", msg)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	resultSlot := uint64(req.Height)
	if req.Height == 0 {
		latestSlot, err := c.solCli.GetSlot(Finalized)
		if err != nil {
			err := fmt.Errorf("GetBlockByNumber GetSlot failed: %w", err)
			log.Error("err", err)
			response.Msg = err.Error()
			return nil, err
		}
		resultSlot = latestSlot
	}

	blockResult := &BlockResult{}
	if req.ViewTx {
		tempBlockBySlot, err := c.solCli.GetBlockBySlot(resultSlot, Signatures)
		if err != nil {
			err := fmt.Errorf("GetBlockByNumber GetBlockBySlot failed: %w", err)
			log.Error("err", err)
			response.Msg = err.Error()
			return nil, err
		}
		blockResult = tempBlockBySlot
	} else {
		tempBlockBySlot, err := c.solCli.GetBlockBySlot(resultSlot, None)
		if err != nil {
			err := fmt.Errorf("GetBlockByNumber GetBlockBySlot failed: %w", err)
			log.Error("err", err)
			response.Msg = err.Error()
			return nil, err
		}
		blockResult = tempBlockBySlot
	}

	response.Hash = blockResult.BlockHash
	response.Height = int64(resultSlot)
	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "GetBlockByNumber success"
	if req.ViewTx {
		response.Transactions = make([]*account.BlockInfoTransactionList, 0, len(blockResult.Signatures))
		for _, signature := range blockResult.Signatures {
			txInfo := &account.BlockInfoTransactionList{
				Hash: signature,
			}
			response.Transactions = append(response.Transactions, txInfo)
		}
	}
	return response, nil
}

func (c *ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	response := &account.BlockResponse{
		Code:         common2.ReturnCode_ERROR,
		Msg:          "",
		Height:       0,
		Hash:         "",
		BaseFee:      "",
		Transactions: nil,
	}
	if ok, msg := validateChainAndNetwork(req.Chain, ""); !ok {
		response.Msg = msg
		err := fmt.Errorf("GetBlockByHash validateChainAndNetwork fail, err msg = %s", msg)
		return response, err
	}
	panic("implement me")
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	response := &account.BlockHeaderResponse{
		Code:        common2.ReturnCode_ERROR,
		Msg:         "",
		BlockHeader: nil,
	}

	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		response.Msg = msg
		err := fmt.Errorf("GetBlockHeaderByHash validateChainAndNetwork fail, err msg = %s", msg)
		return response, err
	}
	panic("implement me")
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	response := &account.BlockHeaderResponse{
		Code:        common2.ReturnCode_ERROR,
		Msg:         "",
		BlockHeader: nil,
	}
	if ok, msg := validateChainAndNetwork(req.Chain, ""); !ok {
		err := fmt.Errorf("GetBlockHeaderByNumber validateChainAndNetwork failed: %s", msg)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	resultSlot := uint64(req.Height)
	if req.Height == 0 {
		latestSlot, err := c.solCli.GetSlot(Finalized)
		if err != nil {
			err := fmt.Errorf("GetBlockHeaderByNumber GetSlot failed: %w", err)
			log.Error("err", err)
			response.Msg = err.Error()
			return nil, err
		}
		resultSlot = latestSlot
	}

	blockResult, err := c.solCli.GetBlockBySlot(resultSlot, None)
	if err != nil {
		err := fmt.Errorf("GetBlockHeaderByNumber GetBlockBySlot failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	blockHead := &account.BlockHeader{
		Hash:       blockResult.BlockHash,
		Number:     strconv.FormatUint(resultSlot, 10),
		ParentHash: blockResult.PreviousBlockhash,
		Time:       uint64(blockResult.BlockTime),
	}

	response.BlockHeader = blockHead
	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "GetBlockHeaderByNumber success"
	return response, nil
}

func (c *ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	response := &account.AccountResponse{
		Code:          common2.ReturnCode_ERROR,
		Msg:           "",
		Network:       "",
		AccountNumber: "",
		Sequence:      "",
		Balance:       "",
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		response.Msg = msg
		err := fmt.Errorf("GetAccount validateChainAndNetwork fail, err msg = %s", msg)
		return response, err
	}
	accountInfoResp, err := c.solCli.GetAccountInfo(req.Address)

	if err != nil {
		err := fmt.Errorf("GetAccount GetAccountInfo failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	latestBlockhashResponse, err := c.solCli.GetLatestBlockhash(Finalized)
	if err != nil {
		err := fmt.Errorf("GetAccount GetLatestBlockhash failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "GetAccount success"
	response.Sequence = latestBlockhashResponse
	response.Network = req.Network
	response.Balance = strconv.FormatUint(accountInfoResp.Lamports, 10)
	return response, nil
}

func (c *ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	response := &account.FeeResponse{
		Code:      common2.ReturnCode_ERROR,
		Msg:       "",
		SlowFee:   "",
		NormalFee: "",
		FastFee:   "",
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		response.Msg = msg
		err := fmt.Errorf("GetFee validateChainAndNetwork fail, err msg = %s", msg)
		return response, err
	}
	baseFee, err := c.solCli.GetFeeForMessage(req.RawTx)
	if err != nil {
		err := fmt.Errorf("GetFee GetFeeForMessage failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	priorityFees, err := c.solCli.GetRecentPrioritizationFees()
	if err != nil {
		err := fmt.Errorf("GetFee GetRecentPrioritizationFees failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	priorityFee := GetSuggestedPriorityFee(priorityFees)
	slowFee := baseFee + uint64(float64(priorityFee)*0.75)
	normalFee := baseFee + priorityFee
	fastFee := baseFee + uint64(float64(priorityFee)*1.25)

	response.SlowFee = strconv.FormatUint(slowFee, 10)
	response.NormalFee = strconv.FormatUint(normalFee, 10)
	response.FastFee = strconv.FormatUint(fastFee, 10)

	return response, nil
}

func (c *ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	if req.RawTx == "" {
		return &account.SendTxResponse{
			Code:   common2.ReturnCode_ERROR,
			Msg:    "invalid input: empty transaction",
			TxHash: "",
		}, nil
	}
	log.Info("2:", req.RawTx)
	// Send the transaction
	txHash, err := c.solCli.SendTransaction(req.RawTx, nil)
	if err != nil {
		log.Error("Failed to send transaction", "err", err)
		return &account.SendTxResponse{
			Code:   common2.ReturnCode_ERROR,
			Msg:    "failed to send transaction",
			TxHash: "",
		}, err
	}

	return &account.SendTxResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "transaction sent successfully",
		TxHash: txHash,
	}, nil
}

func (c *ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	var resp *account2.TransactionResponse[account2.AccountTxResponse]
	var err error
	fmt.Println("req.ContractAddress", req.ContractAddress)
	if req.ContractAddress != "0x00" && req.ContractAddress != "" {
		log.Info("Spl token transfer record")
		resp, err = c.solData.GetTxByAddress(uint64(req.Page), uint64(req.Pagesize), req.Address, "spl")
	} else {
		log.Info("Sol transfer record")
		resp, err = c.solData.GetTxByAddress(uint64(req.Page), uint64(req.Pagesize), req.Address, "sol")
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
	response := &account.TxHashResponse{
		Code: common2.ReturnCode_ERROR,
		Msg:  "",
		Tx:   nil,
	}

	if err := validateRequest(req); err != nil {
		response.Msg = err.Error()
		return response, err
	}

	txResult, err := c.solCli.GetTransaction(req.Hash)
	if err != nil {
		response.Msg = err.Error()
		log.Error("GetTransaction failed", "error", err)
		return response, err
	}

	tx, err := buildTxMessage(txResult)
	if err != nil {
		response.Msg = err.Error()
		return response, err
	}

	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "success"
	response.Tx = tx

	return response, nil
}

func validateRequest(req *account.TxHashRequest) error {
	if req == nil {
		return fmt.Errorf("invalid request: request is nil")
	}
	if req.Hash == "" {
		return fmt.Errorf("invalid request: empty transaction hash")
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return fmt.Errorf("invalid chain or network: %s", msg)
	}
	return nil
}

func buildTxMessage(txResult *TransactionResult) (*account.TxMessage, error) {
	if txResult == nil {
		return nil, fmt.Errorf("empty transaction result")
	}

	if len(txResult.Transaction.Signatures) == 0 {
		return nil, fmt.Errorf("invalid transaction: no signatures")
	}
	if len(txResult.Transaction.Message.AccountKeys) == 0 {
		return nil, fmt.Errorf("invalid transaction: no account keys")
	}

	tx := &account.TxMessage{
		Hash:   txResult.Transaction.Signatures[0],
		Height: strconv.FormatUint(txResult.Slot, 10),
		Fee:    strconv.FormatUint(txResult.Meta.Fee, 10),
	}

	if txResult.Meta.Err != nil {
		tx.Status = account.TxStatus_Failed
	} else {
		tx.Status = account.TxStatus_Success
	}

	if txResult.BlockTime != nil {
		tx.Datetime = time.Unix(*txResult.BlockTime, 0).Format(time.RFC3339)
	}

	tx.Froms = []*account.Address{{
		Address: txResult.Transaction.Message.AccountKeys[0],
	}}

	tx.Tos = make([]*account.Address, 0)
	tx.Values = make([]*account.Value, 0)

	if err := processInstructions(txResult, tx); err != nil {
		return nil, fmt.Errorf("failed to process instructions: %w", err)
	}

	return tx, nil
}

func processInstructions(txResult *TransactionResult, tx *account.TxMessage) error {
	for i, inst := range txResult.Transaction.Message.Instructions {
		if inst.ProgramIdIndex >= len(txResult.Transaction.Message.AccountKeys) {
			log.Warn("Invalid program ID index", "instruction", i)
			continue
		}

		if txResult.Transaction.Message.AccountKeys[inst.ProgramIdIndex] != "11111111111111111111111111111111" {
			continue
		}

		if len(inst.Accounts) < 2 {
			log.Warn("Invalid accounts length", "instruction", i)
			continue
		}

		toIndex := inst.Accounts[1]
		if toIndex >= len(txResult.Transaction.Message.AccountKeys) {
			log.Warn("Invalid to account index", "instruction", i)
			continue
		}

		toAddr := txResult.Transaction.Message.AccountKeys[toIndex]
		tx.Tos = append(tx.Tos, &account.Address{Address: toAddr})

		if err := calculateAmount(txResult, toIndex, tx); err != nil {
			log.Warn("Failed to calculate amount", "error", err)
			continue
		}
	}

	return nil
}

func calculateAmount(txResult *TransactionResult, toIndex int, tx *account.TxMessage) error {
	if toIndex >= len(txResult.Meta.PostBalances) || toIndex >= len(txResult.Meta.PreBalances) {
		return fmt.Errorf("invalid balance index: %d", toIndex)
	}

	amount := txResult.Meta.PostBalances[toIndex] - txResult.Meta.PreBalances[toIndex]
	tx.Values = append(tx.Values, &account.Value{
		Value: strconv.FormatUint(amount, 10),
	})

	return nil
}

func (c *ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	response := &account.BlockByRangeResponse{
		Code:        common2.ReturnCode_ERROR,
		Msg:         "",
		BlockHeader: nil,
	}
	if err := validateBlockRangeRequest(req); err != nil {
		response.Msg = err.Error()
		return response, err
	}
	startSlot, _ := strconv.ParseUint(req.Start, 10, 64)
	endSlot, _ := strconv.ParseUint(req.End, 10, 64)

	for slot := startSlot; slot <= endSlot; slot++ {
		blockResult, err := c.solCli.GetBlockBySlot(slot, Signatures)
		if err != nil {
			if len(response.BlockHeader) > 0 {
				response.Code = common2.ReturnCode_SUCCESS
				response.Msg = fmt.Sprintf("partial success, stopped at slot %d: %v", slot, err)
				return response, nil
			}
			response.Msg = fmt.Sprintf("failed to get signatures for slot %d: %v", slot, err)
			return response, err
		}

		if len(blockResult.Signatures) == 0 {
			continue
		}

		txResults, err := c.solCli.GetTransactionRange(blockResult.Signatures)
		if err != nil {
			if len(response.BlockHeader) > 0 {
				response.Code = common2.ReturnCode_SUCCESS
				response.Msg = fmt.Sprintf("partial success, stopped at slot %d: %v", slot, err)
				return response, nil
			}
			response.Msg = fmt.Sprintf("failed to get transactions for slot %d: %v", slot, err)
			return response, err
		}

		block, err := organizeTransactionsByBlock(txResults)
		if err != nil {
			if len(response.BlockHeader) > 0 {
				response.Code = common2.ReturnCode_SUCCESS
				response.Msg = fmt.Sprintf("partial success, stopped at slot %d: %v", slot, err)
				return response, nil
			}
			response.Msg = fmt.Sprintf("failed to organize transactions for slot %d: %v", slot, err)
			return response, err
		}

		if len(block) > 0 {
			response.BlockHeader = append(response.BlockHeader, block...)
		}
	}

	if len(response.BlockHeader) == 0 {
		response.Code = common2.ReturnCode_SUCCESS
		response.Msg = "no transactions found in range"
		return response, nil
	}

	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "success"
	return response, nil
}

func validateBlockRangeRequest(req *account.BlockByRangeRequest) error {
	if req == nil {
		return fmt.Errorf("invalid request: request is nil")
	}
	startSlot, err := strconv.ParseUint(req.Start, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid start height format: %s", err)
	}
	endSlot, err := strconv.ParseUint(req.End, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid end height format: %s", err)
	}

	if startSlot > endSlot {
		return fmt.Errorf("invalid height range: start height greater than end height")
	}

	if startSlot-endSlot > MaxBlockRange {
		return fmt.Errorf("invalid range: exceeds maximum allowed range of %d", MaxBlockRange)
	}

	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return fmt.Errorf("invalid chain or network: %s", msg)
	}

	return nil
}

func organizeTransactionsByBlock(txResults []*TransactionResult) ([]*account.BlockHeader, error) {
	if len(txResults) == 0 {
		return nil, nil
	}

	blockMap := make(map[uint64]*account.BlockHeader)

	for _, txResult := range txResults {
		if txResult == nil {
			continue
		}

		slot := txResult.Slot

		block, exists := blockMap[slot]
		if !exists {
			block = &account.BlockHeader{
				Number: strconv.FormatUint(slot, 10),
			}

			if txResult.BlockTime != nil {
				block.Time = uint64(*txResult.BlockTime)
			}

			if len(txResult.Transaction.Signatures) > 0 {
				block.Hash = txResult.Transaction.Signatures[0]
			}

			txHashes := make([]string, 0)
			for _, sig := range txResult.Transaction.Signatures {
				txHashes = append(txHashes, sig)
			}
			block.TxHash = strings.Join(txHashes, ",")

			block.GasUsed = txResult.Meta.ComputeUnitsConsumed

			blockMap[slot] = block
		} else {
			if len(txResult.Transaction.Signatures) > 0 {
				if block.TxHash != "" {
					block.TxHash += "," + txResult.Transaction.Signatures[0]
				} else {
					block.TxHash = txResult.Transaction.Signatures[0]
				}
			}

			block.GasUsed += txResult.Meta.ComputeUnitsConsumed
		}
	}

	blocks := make([]*account.BlockHeader, 0, len(blockMap))
	for _, block := range blockMap {
		blocks = append(blocks, block)
	}

	sort.Slice(blocks, func(i, j int) bool {
		heightI, _ := strconv.ParseUint(blocks[i].Number, 10, 64)
		heightJ, _ := strconv.ParseUint(blocks[j].Number, 10, 64)
		return heightI < heightJ
	})

	return blocks, nil
}

func (c *ChainAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	// Decode the base64 transaction string
	jsonBytes, err := base64.StdEncoding.DecodeString(req.Base64Tx)
	if err != nil {
		log.Error("Failed to decode base64 string", "err", err)
		return nil, err
	}

	// Unmarshal JSON into TxStructure
	var data TxStructure
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		log.Error("Failed to parse JSON", "err", err)
		return nil, err
	}

	// Parse the value from string to float
	valueFloat, err := strconv.ParseFloat(data.Value, 64)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse value: %w", err)
	}
	value := uint64(valueFloat * 1000000000)

	// Convert from address to public key
	fromPubkey, err := solana.PublicKeyFromBase58(data.FromAddress)
	if err != nil {
		return nil, err
	}

	// Convert to address to public key
	toPubkey, err := solana.PublicKeyFromBase58(data.ToAddress)
	if err != nil {
		return nil, err
	}

	var tx *solana.Transaction
	if isSOLTransfer(data.ContractAddress) {
		// Create a new SOL transfer transaction
		tx, err = solana.NewTransaction(
			[]solana.Instruction{
				system.NewTransferInstruction(
					value,
					fromPubkey,
					toPubkey,
				).Build(),
			},
			solana.MustHashFromBase58(data.Nonce),
			solana.TransactionPayer(fromPubkey),
		)
	} else {
		// Handle SPL token transfer
		mintPubkey := solana.MustPublicKeyFromBase58(data.ContractAddress)

		fromTokenAccount, _, err := solana.FindAssociatedTokenAddress(
			fromPubkey,
			mintPubkey,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to find associated token address: %w", err)
		}

		toTokenAccount, _, err := solana.FindAssociatedTokenAddress(
			toPubkey,
			mintPubkey,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to find associated token address: %w", err)
		}

		//tokenInfo, err := c.sdkClient.GetTokenSupply(context.Background(), mintPubkey, rpc.CommitmentFinalized)
		tokenInfo, err := GetTokenSupply(c.sdkClient, mintPubkey)
		if err != nil {
			return nil, fmt.Errorf("Failed to get token info: %w", err)
		}
		decimals := tokenInfo.Value.Decimals

		valueFloat, err := strconv.ParseFloat(data.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse value: %w", err)
		}
		actualValue := uint64(valueFloat * math.Pow10(int(decimals)))

		transferInstruction := token.NewTransferInstruction(
			actualValue,
			fromTokenAccount,
			toTokenAccount,
			fromPubkey,
			[]solana.PublicKey{},
		).Build()

		//accountInfo, err := c.sdkClient.GetAccountInfo(context.Background(), toTokenAccount)
		accountInfo, err := GetAccountInfo(c.sdkClient, toTokenAccount)

		if err != nil || accountInfo.Value == nil {
			// Create associated token account if it doesn't exist
			createATAInstruction := associatedtokenaccount.NewCreateInstruction(
				fromPubkey,
				toPubkey,
				mintPubkey,
			).Build()

			tx, err = solana.NewTransaction(
				[]solana.Instruction{createATAInstruction, transferInstruction},
				solana.MustHashFromBase58(data.Nonce),
				solana.TransactionPayer(fromPubkey),
			)
		} else {
			// Directly create transfer transaction
			tx, err = solana.NewTransaction(
				[]solana.Instruction{transferInstruction},
				solana.MustHashFromBase58(data.Nonce),
				solana.TransactionPayer(fromPubkey),
			)
		}
	}

	// Log the transaction details
	log.Info("Transaction:", tx.String())

	// Serialize the transaction message
	txm, _ := tx.Message.MarshalBinary()
	signingMessageHex := hex.EncodeToString(txm)

	// Return the unsigned transaction response
	return &account.UnSignTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "Successfully created unsigned transaction",
		UnSignTx: signingMessageHex,
	}, nil
}
func (c ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	// Decode the base64 transaction string
	jsonBytes, err := base64.StdEncoding.DecodeString(req.Base64Tx)
	if err != nil {
		log.Error("Failed to decode base64 string", "err", err)
		return nil, err
	}

	// Unmarshal JSON into TxStructure
	var data TxStructure
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		log.Error("Failed to parse JSON", "err", err)
		return nil, err
	}

	// Parse the value from string to float
	valueFloat, err := strconv.ParseFloat(data.Value, 64)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse value: %w", err)
	}
	value := uint64(valueFloat * 1000000000)

	// Convert from address to public key
	fromPubkey, err := solana.PublicKeyFromBase58(data.FromAddress)
	if err != nil {
		return nil, err
	}

	// Convert to address to public key
	toPubkey, err := solana.PublicKeyFromBase58(data.ToAddress)
	if err != nil {
		return nil, err
	}

	var tx *solana.Transaction
	if isSOLTransfer(data.ContractAddress) {
		// Create a new SOL transfer transaction
		tx, err = solana.NewTransaction(
			[]solana.Instruction{
				system.NewTransferInstruction(
					value,
					fromPubkey,
					toPubkey,
				).Build(),
			},
			solana.MustHashFromBase58(data.Nonce),
			solana.TransactionPayer(fromPubkey),
		)
	} else {
		// Handle SPL token transfer
		mintPubkey := solana.MustPublicKeyFromBase58(data.ContractAddress)

		fromTokenAccount, _, err := solana.FindAssociatedTokenAddress(
			fromPubkey,
			mintPubkey,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to find associated token address: %w", err)
		}

		toTokenAccount, _, err := solana.FindAssociatedTokenAddress(
			toPubkey,
			mintPubkey,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to find associated token address: %w", err)
		}

		//tokenInfo, err := c.sdkClient.GetTokenSupply(context.Background(), mintPubkey, rpc.CommitmentFinalized)
		tokenInfo, err := GetTokenSupply(c.sdkClient, mintPubkey)
		if err != nil {
			return nil, fmt.Errorf("Failed to get token info: %w", err)
		}
		decimals := tokenInfo.Value.Decimals

		valueFloat, err := strconv.ParseFloat(data.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse value: %w", err)
		}
		actualValue := uint64(valueFloat * math.Pow10(int(decimals)))

		transferInstruction := token.NewTransferInstruction(
			actualValue,
			fromTokenAccount,
			toTokenAccount,
			fromPubkey,
			[]solana.PublicKey{},
		).Build()
		//accountInfo, err := c.sdkClient.GetAccountInfo(context.Background(), toTokenAccount)
		accountInfo, err := GetAccountInfo(c.sdkClient, toTokenAccount)

		if err != nil || accountInfo.Value == nil {
			// Create associated token account if it doesn't exist
			createATAInstruction := associatedtokenaccount.NewCreateInstruction(
				fromPubkey,
				toPubkey,
				mintPubkey,
			).Build()

			tx, err = solana.NewTransaction(
				[]solana.Instruction{createATAInstruction, transferInstruction},
				solana.MustHashFromBase58(data.Nonce),
				solana.TransactionPayer(fromPubkey),
			)
		} else {
			// Directly create transfer transaction
			tx, err = solana.NewTransaction(
				[]solana.Instruction{transferInstruction},
				solana.MustHashFromBase58(data.Nonce),
				solana.TransactionPayer(fromPubkey),
			)
		}
	}

	// Ensure the Signatures slice is initialized
	if len(tx.Signatures) == 0 {
		tx.Signatures = make([]solana.Signature, 1)
	}

	// Decode the signature from hex
	signatureBytes, err := hex.DecodeString(data.Signature)
	if err != nil {
		log.Error("Failed to decode hex signature", "err", err)
	}

	// Verify the signature length
	if len(signatureBytes) != 64 {
		log.Error("Invalid signature length", "length", len(signatureBytes))
	}

	// Convert to Solana Signature
	var solanaSig solana.Signature
	copy(solanaSig[:], signatureBytes)

	// Set the signature
	tx.Signatures[0] = solanaSig

	// Dump the transaction for debugging
	spew.Dump(tx)
	if err := tx.VerifySignatures(); err != nil {
		log.Info("Invalid signatures", "err", err)
	}

	// Serialize the transaction
	serializedTx, err := tx.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("Failed to serialize transaction: %w", err)
	}

	// Encode the serialized transaction to base58
	base58Tx := base58.Encode(serializedTx)
	//base64Tx := base64.StdEncoding.EncodeToString(serializedTx)
	return &account.SignedTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "Successfully created signed transaction",
		SignedTx: base58Tx,
	}, nil
}

func (c *ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {

	txBytes, err := base58.Decode(req.Signature)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction: %w", err)
	}

	tx, err := solana.TransactionFromBytes(txBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %w", err)
	}

	if err := tx.VerifySignatures(); err != nil {
		log.Info("Invalid signatures", "err", err)
		return &account.VerifyTransactionResponse{
			Code:   common2.ReturnCode_ERROR,
			Msg:    "invalid signature",
			Verify: false,
		}, nil
	}

	return &account.VerifyTransactionResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "valid signature",
		Verify: true,
	}, nil
}

func (c *ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	//TODO implement me
	panic("implement me")
}
func isSOLTransfer(coinAddress string) bool {

	return coinAddress == "" ||
		coinAddress == "So11111111111111111111111111111111111111112"
}
func getPrivateKey(keyStr string) (solana.PrivateKey, error) {
	// base58
	if prikey, err := solana.PrivateKeyFromBase58(keyStr); err == nil {
		return prikey, nil
	}

	// 16
	if bytes, err := hex.DecodeString(keyStr); err == nil {
		return solana.PrivateKey(bytes), nil
	}

	return nil, fmt.Errorf("invalid private key format")
}

func validateChainAndNetwork(chain, network string) (bool, string) {
	if chain != ChainName {
		return false, "invalid chain"
	}
	//if network != NetworkMainnet && network != NetworkTestnet {
	//	return false, "invalid network"
	//}
	return true, ""
}

func validatePublicKey(pubKey string) (bool, string) {
	if pubKey == "" {
		return false, "public key cannot be empty"
	}
	pubKeyWithoutPrefix := strings.TrimPrefix(pubKey, "0x")

	if len(pubKeyWithoutPrefix) != 64 {
		return false, "invalid public key length"
	}
	if _, err := hex.DecodeString(pubKeyWithoutPrefix); err != nil {
		return false, "invalid public key format: must be hex string"
	}

	return true, ""
}
