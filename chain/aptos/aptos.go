package aptos

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/aptos-labs/aptos-go-sdk/crypto"
	"github.com/ethereum/go-ethereum/log"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

const ChainName = "Aptos"
const ResourceTypeAPT = "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>"

type ChainAdaptor struct {
	aptosHttpClient AptClient
	aptosClient     *aptos.Client
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	rpcUrl := conf.WalletNode.Aptos.RpcUrl
	apiKey := conf.WalletNode.Aptos.DataApiKey

	aptosHttpClient, err := NewAptosHttpClient(rpcUrl, apiKey)
	if err != nil {
		log.Error("NewChainAdaptor NewAptosHttpClient fail", "err", err)
		return nil, err
	}

	aptosConfNetWork := conf.NetWork
	newAptosClient, err := NewAptosClient(aptosConfNetWork)
	if err != nil {
		log.Error("NewChainAdaptor newAptosClient fail", "err", err)
		return nil, err
	}

	return &ChainAdaptor{
		aptosHttpClient: aptosHttpClient,
		aptosClient:     newAptosClient,
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
		return response, err
	}

	pubKeyHex := req.PublicKey
	if ok, msg := validateAptosPublicKey(req.PublicKey); !ok {
		err := fmt.Errorf("ConvertAddress validatePublicKey fail, err msg = %s", msg)
		log.Error("err", err)
		response.Msg = err.Error()
		return response, err
	}
	accountAddress, err := PubKeyHexToAccountAddress(pubKeyHex)
	if err != nil {
		err := fmt.Errorf("ConvertAddress PubKeyHexToAccountAddress failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "convert address success"
	response.Address = accountAddress.String()
	return response, nil
}

func (c *ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	response := &account.ValidAddressResponse{
		Code:  common2.ReturnCode_ERROR,
		Msg:   "",
		Valid: false,
	}

	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		err := fmt.Errorf("ConvertAddress validateChainAndNetwork failed: %s", msg)
		log.Error("err", err)
		response.Msg = err.Error()
		return response, err
	}

	errTooShrot := "AccountAddress too short"
	errTooLong := "AccountAddress too long"

	address := req.Address
	aptosAccountAddress := &aptos.AccountAddress{}
	err := aptosAccountAddress.ParseStringRelaxed(address)
	if err != nil {
		switch err.Error() {
		case errTooShrot:
			err := fmt.Errorf("ValidAddress ParseStringRelaxed errTooShrot failed: %w", err)
			log.Error("err", err)
			response.Msg = err.Error()
			return nil, err
		case errTooLong:
			err := fmt.Errorf("ValidAddress ParseStringRelaxed errTooLong failed: %w", err)
			log.Error("err", err)
			response.Msg = err.Error()
			return nil, err
		default:
			err := fmt.Errorf("ValidAddress ParseStringRelaxed default failed: %w", err)
			log.Error("err", err)
			response.Msg = err.Error()
			return nil, err
		}
	}

	cleanAddr := address
	if strings.HasPrefix(cleanAddr, "0x") {
		cleanAddr = cleanAddr[2:]
	}

	_, err = hex.DecodeString(cleanAddr)
	if err != nil {
		err := fmt.Errorf("ValidAddress DecodeString default failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	response.Valid = true
	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "ValidAddress success"
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
		err := fmt.Errorf("GetBlockByNumber validateChainAndNetwork default failed: %s", msg)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	if req.Height == 0 {
		nodeInfo, err := c.aptosHttpClient.GetNodeInfo()
		if err != nil {
			err := fmt.Errorf("GetBlockByNumber GetNodeInfo failed: %w", err)
			log.Error("err", err)
			response.Msg = err.Error()
			return nil, err
		}
		if req.ViewTx {

		}
		response.Code = common2.ReturnCode_SUCCESS
		response.Msg = "GetBlockByNumber success"
		response.Height = int64(nodeInfo.BlockHeight)
		// TODO: Transactionsdasda
		response.Transactions = nil
		return response, nil
	}

	blockResponse, err := c.aptosHttpClient.GetBlockByHeight(uint64(req.Height))
	if err != nil {
		err := fmt.Errorf("GetBlockByNumber GetBlockByHeight failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	if req.ViewTx {

	}
	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "GetBlockByNumber success"
	response.Height = int64(blockResponse.BlockHeight)
	response.Hash = blockResponse.BlockHash
	// TODO: Transactionsdasda
	response.Transactions = nil
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
	// no implement
	panic("implement me")
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	response := &account.BlockHeaderResponse{
		Code:        common2.ReturnCode_ERROR,
		Msg:         "",
		BlockHeader: nil,
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		response.Msg = msg
		err := fmt.Errorf("GetBlockHeaderByNumber validateChainAndNetwork fail, err msg = %s", msg)
		return response, err
	}
	if req.Height == 0 {
		nodeInfo, err := c.aptosHttpClient.GetNodeInfo()
		if err != nil {
			err := fmt.Errorf("GetBlockHeaderByNumber GetNodeInfo fail, err msg = %s", err.Error())
			log.Error("err", err)
			response.Msg = err.Error()
			return response, err
		}
		seconds := nodeInfo.LedgerTimestamp / 1_000_000
		blockHead := &account.BlockHeader{
			Number: strconv.FormatUint(nodeInfo.BlockHeight, 10),
			Time:   seconds,
		}
		response.BlockHeader = blockHead
		response.Code = common2.ReturnCode_SUCCESS
		response.Msg = "GetBlockHeaderByNumber success"
		return response, nil
	}

	blockResponse, err := c.aptosHttpClient.GetBlockByHeight(uint64(req.Height))
	if err != nil {
		err := fmt.Errorf("GetBlockHeaderByNumber GetBlockByHeight failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	blockHead := &account.BlockHeader{
		Hash:   blockResponse.BlockHash,
		Number: strconv.FormatUint(blockResponse.BlockHeight, 10),
		Time:   blockResponse.BlockTimestamp,
	}
	response.BlockHeader = blockHead
	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "GetBlockHeaderByNumber success"
	return response, nil
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
	// no implement
	panic("implement me")
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
	accountResponse, err := c.aptosHttpClient.GetAccount(req.Address)
	if err != nil {
		err := fmt.Errorf("GetAccount GetAccount failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	accountBalance, err := c.aptosHttpClient.GetAccountBalance(req.Address, ResourceTypeAPT)
	if err != nil {
		err := fmt.Errorf("GetAccount GetAccountBalance failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "GetAccount success"
	response.Sequence = strconv.FormatUint(accountResponse.SequenceNumber, 10)
	response.Network = req.Network
	response.Balance = strconv.FormatUint(accountBalance, 10)
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
	gasPriceResponse, err := c.aptosHttpClient.GetGasPrice()
	if err != nil {
		err := fmt.Errorf("GetFee GetGasPrice failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "GetFee success"
	response.SlowFee = strconv.FormatUint(gasPriceResponse.DeprioritizedGasEstimate, 10)
	response.NormalFee = strconv.FormatUint(gasPriceResponse.GasEstimate, 10)
	response.FastFee = strconv.FormatUint(gasPriceResponse.PrioritizedGasEstimate, 10)
	return response, nil
}

func (c *ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	response := &account.SendTxResponse{
		Code:   common2.ReturnCode_ERROR,
		Msg:    "",
		TxHash: "",
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		response.Msg = msg
		err := fmt.Errorf("SendTx validateChainAndNetwork fail, err msg = %s", msg)
		return response, err
	}
	rawTx := req.RawTx

	rawTxByteList, err := base64.StdEncoding.DecodeString(rawTx)
	if err != nil {
		err := fmt.Errorf("SendTx DecodeString rawTx failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	signedTxListDes := bcs.NewDeserializer(rawTxByteList)
	signedTx := &aptos.SignedTransaction{
		Transaction: &aptos.RawTransaction{},
		Authenticator: &aptos.TransactionAuthenticator{
			Variant: aptos.TransactionAuthenticatorEd25519,
			Auth: &crypto.Ed25519Authenticator{
				PubKey: &crypto.Ed25519PublicKey{},
				Sig:    &crypto.Ed25519Signature{},
			},
		},
	}
	signedTx.UnmarshalBCS(signedTxListDes)

	submitTransactionResponse, err := c.aptosClient.SubmitTransaction(signedTx)
	if err != nil {
		err := fmt.Errorf("SendTx SubmitTransaction rawTx failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "SubmitTransaction success"
	response.TxHash = submitTransactionResponse.Hash
	return response, nil
}

func (c *ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	response := &account.TxAddressResponse{
		Code: common2.ReturnCode_ERROR,
		Msg:  "",
		Tx:   nil,
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		response.Msg = msg
		err := fmt.Errorf("GetTxByAddress validateChainAndNetwork fail, err msg = %s", msg)
		return response, err
	}
	transactionsPtr, err := c.aptosHttpClient.GetTransactionByAddress(req.Address)
	if err != nil {
		err := fmt.Errorf("GetTxByAddress GetTransactionByAddress failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	if transactionsPtr == nil {
		err := fmt.Errorf("GetTxByAddress transactions is null: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	transactions := *transactionsPtr
	var txMessages []*account.TxMessage

	for _, tx := range transactions {
		var txStatus account.TxStatus
		if tx.Success {
			txStatus = account.TxStatus_Success
		} else {
			txStatus = account.TxStatus_Failed
		}

		feeStatement := GetFeeStatementFromEvents(tx.Events)
		var totalFee uint64
		if feeStatement != nil {
			totalFee = CalculateGasFee(tx.GasUnitPrice, feeStatement.TotalChargeGasUnits,
				feeStatement.StorageFeeOctas, feeStatement.StorageFeeRefundOctas)
		} else {
			totalFee = tx.GasUsed * tx.GasUnitPrice
		}
		fromAddr := &account.Address{
			Address: tx.Sender,
		}
		txMessage := &account.TxMessage{
			Hash:  tx.Hash,
			Froms: []*account.Address{fromAddr},
			//TODO to
			Tos: []*account.Address{},
			//TODO Value
			Values: []*account.Value{},
			Fee:    strconv.FormatUint(totalFee, 10),
			Status: txStatus,
			Type:   0,
			Height: strconv.FormatUint(tx.Version, 10),
			// ContractAddress:
			Datetime: strconv.FormatUint(tx.Timestamp, 10),
			Data:     convertExtraInfo(tx),
		}
		txMessages = append(txMessages, txMessage)
	}
	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "GetTxByAddress success"
	response.Tx = txMessages
	return response, nil
}

func (c *ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	response := &account.TxHashResponse{
		Code: common2.ReturnCode_ERROR,
		Msg:  "",
		Tx:   nil,
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		response.Msg = msg
		err := fmt.Errorf("GetTxByHash validateChainAndNetwork fail, err msg = %s", msg)
		return response, err
	}
	tx, err := c.aptosHttpClient.GetTransactionByHash(req.Hash)
	if err != nil {
		err := fmt.Errorf("GetTxByHash GetTransactionByHash failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	//var fromAddrList []*account.Address
	//var toAddrsList []*account.Address
	//var valueList []*account.Value
	var txStatus account.TxStatus
	if tx.Success {
		txStatus = account.TxStatus_Success
	} else {
		txStatus = account.TxStatus_Failed
	}

	feeStatement := GetFeeStatementFromEvents(tx.Events)
	totalFee := CalculateGasFee(tx.GasUnitPrice, feeStatement.TotalChargeGasUnits, feeStatement.StorageFeeOctas, feeStatement.StorageFeeRefundOctas)

	txMessage := &account.TxMessage{
		Hash:  tx.Hash,
		Index: uint32(tx.SequenceNumber),
		Froms: []*account.Address{{
			Address: tx.Sender,
		}},
		Tos: []*account.Address{{
			Address: tx.Payload.Arguments[0].(string),
		}},
		Values: []*account.Value{{
			Value: tx.Payload.Arguments[1].(string),
		}},
		Fee:    strconv.FormatUint(totalFee, 10),
		Status: txStatus,
		Type:   determineTransactionType(tx.Payload.Function),
		// Height: strconv.FormatUint(tx.Version, 10),
		// ContractAddress:
		Datetime: strconv.FormatUint(tx.Timestamp, 10),
		Data:     serializePayload(tx.Payload),
	}
	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "GetTxByHash success"
	response.Tx = txMessage
	return response, nil
}

func (c *ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	response := &account.BlockByRangeResponse{
		Code:        common2.ReturnCode_ERROR,
		Msg:         "",
		BlockHeader: nil,
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		response.Msg = msg
		err := fmt.Errorf("GetBlockByRange validateChainAndNetwork fail, err msg = %s", msg)
		return response, err
	}
	startBlock, err := strconv.ParseUint(req.Start, 10, 64)
	if err != nil {
		err := fmt.Errorf("GetBlockByRange startVersion failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	endBlock, err := strconv.ParseUint(req.End, 10, 64)
	if err != nil {
		err := fmt.Errorf("GetBlockByRange endVersion failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	if startBlock > endBlock {
		err := fmt.Errorf("GetBlockByRange startBlock > endBlock failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	response.BlockHeader = make([]*account.BlockHeader, 0, (endBlock-startBlock)*5)

	for height := startBlock; height <= endBlock; height++ {
		block, err := c.aptosHttpClient.GetBlockByHeight(height)
		if err != nil {
			err := fmt.Errorf("GetBlockByRange GetBlockByHeight failed at height %d: %w", height, err)
			log.Error("err", err)
			response.Msg = err.Error()
			return nil, err
		}
		startVersion := block.FirstVersion
		endVersion := block.LastVersion
		txs, err := c.aptosHttpClient.GetTransactionByVersionRange(startVersion, endVersion)
		if err != nil {
			err := fmt.Errorf("GetBlockByRange GetTransactionByVersionRange failed: %w", err)
			log.Error("err", err)
			response.Msg = err.Error()
			return nil, err
		}

		for _, tx := range txs {
			blockByVersion, err := c.aptosHttpClient.GetBlockByVersion(tx.Version)
			if err != nil {
				err := fmt.Errorf("GetBlockByRange GetBlockByVersion failed: %w", err)
				log.Error("err", err)
				response.Msg = err.Error()
				return nil, err
			}

			blockHeader := &account.BlockHeader{
				Hash: blockByVersion.BlockHash,
				//ParentHash:  tx.StateRootHash,
				//Root:        tx.StateRootHash,
				TxHash:      tx.Hash,
				ReceiptHash: tx.EventRootHash,
				Number:      strconv.FormatUint(blockByVersion.BlockHeight, 10),
				GasLimit:    tx.MaxGasAmount,
				GasUsed:     tx.GasUsed,
				Time:        blockByVersion.BlockTimestamp,
				Extra:       convertExtraInfo(tx),
				Nonce:       strconv.FormatUint(tx.SequenceNumber, 10),
			}
			response.BlockHeader = append(response.BlockHeader, blockHeader)
		}
	}

	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "GetBlockByRange success"
	return response, nil
}

func (c *ChainAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	response := &account.UnSignTransactionResponse{
		Code:     common2.ReturnCode_ERROR,
		Msg:      "",
		UnSignTx: "",
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		response.Msg = msg
		err := fmt.Errorf("CreateUnSignTransaction validateChainAndNetwork fail, err msg = %s", msg)
		return response, err
	}
	if req.Base64Tx == "" {
		response.Msg = "base64_tx cannot be empty"
		return response, errors.New(response.Msg)
	}
	txByteList, err := base64.StdEncoding.DecodeString(req.Base64Tx)
	if err != nil {
		err := fmt.Errorf("CreateUnSignTransaction failed to decode base64 transaction: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	var txRequest TransactionRequest
	if err := json.Unmarshal(txByteList, &txRequest); err != nil {
		err := fmt.Errorf("CreateUnSignTransaction failed to unmarshal transaction request: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	rawTransaction, err := ConvertToRawTransaction(&txRequest)
	if err != nil {
		err := fmt.Errorf("CreateUnSignTransaction failed to ConvertToRawTransaction: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	signingMessage, err := rawTransaction.SigningMessage()
	if err != nil {
		err := fmt.Errorf("CreateUnSignTransaction failed to SigningMessage: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	signingMessageHex := hex.EncodeToString(signingMessage)

	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "CreateUnSignTransaction success"
	response.UnSignTx = signingMessageHex
	return response, nil
}

func (c *ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	response := &account.SignedTransactionResponse{
		Code:     common2.ReturnCode_ERROR,
		Msg:      "",
		SignedTx: "",
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		response.Msg = msg
		err := fmt.Errorf("BuildSignedTransaction validateChainAndNetwork fail, err msg = %s", msg)
		return response, err
	}
	if req.Base64Tx == "" || req.Signature == "" || req.PublicKey == "" {
		err := fmt.Errorf("req.Base64Tx or req.Signature is empty")
		log.Error("BuildSignedTransaction Base64Tx or Signature or PublicKey is empty", "err", err)
		return nil, err
	}

	txByteList, err := base64.StdEncoding.DecodeString(req.Base64Tx)
	if err != nil {
		err := fmt.Errorf("BuildSignedTransaction DecodeString rawTx failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	var txRequest TransactionRequest
	if err := json.Unmarshal(txByteList, &txRequest); err != nil {
		err := fmt.Errorf("CreateUnSignTransaction failed to unmarshal transaction request: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	rawTransaction, err := ConvertToRawTransaction(&txRequest)
	if err != nil {
		err := fmt.Errorf("CreateUnSignTransaction failed to ConvertToRawTransaction: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	fmt.Printf("Transaction Details:\n")
	fmt.Printf("Sender: %s\n", rawTransaction.Sender.String())
	fmt.Printf("Sequence Number: %d\n", rawTransaction.SequenceNumber)
	fmt.Printf("Expiration: %d\n", rawTransaction.ExpirationTimestampSeconds)
	signingMessage, err := rawTransaction.SigningMessage()
	if err != nil {
		return nil, fmt.Errorf("get signing message failed: %w", err)
	}

	pubKeyHex := req.PublicKey
	ed25519PublicKey, err := PubKeyHexToPubKey(pubKeyHex)
	if err != nil {
		err := fmt.Errorf("BuildSignedTransaction PubKeyHexToPubKey failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	signatureByteList, err := hex.DecodeString(req.Signature)
	if err != nil {
		return nil, fmt.Errorf("decode signature failed: %w", err)
	}
	if len(signatureByteList) != ed25519.SignatureSize {
		return nil, fmt.Errorf("invalid signature length: got %d, want %d",
			len(signatureByteList), ed25519.SignatureSize)
	}
	fmt.Printf("Verification Details:\n")
	fmt.Printf("Public Key (hex): %x\n", ed25519PublicKey.Bytes())
	fmt.Printf("Signing Message (hex): %x\n", signingMessage)
	fmt.Printf("Signature (hex): %x\n", signatureByteList)
	if !ed25519.Verify(ed25519PublicKey.Bytes(), signingMessage, signatureByteList) {
		return nil, fmt.Errorf("signature verification failed")
	}

	signature := &crypto.Ed25519Signature{
		Inner: [ed25519.SignatureSize]byte{},
	}
	copy(signature.Inner[:], signatureByteList)

	authenticator := &aptos.TransactionAuthenticator{
		Variant: aptos.TransactionAuthenticatorEd25519,
		Auth: &crypto.Ed25519Authenticator{
			PubKey: ed25519PublicKey,
			Sig:    signature,
		},
	}

	signedTxn := &aptos.SignedTransaction{
		Transaction:   rawTransaction,
		Authenticator: authenticator,
	}

	//signedTxnJson, err := json.Marshal(signedTxn)
	//if err != nil {
	//	err := fmt.Errorf("BuildSignedTransaction signedTxn json failed: %w", err)
	//	log.Error("err", err)
	//	response.Msg = err.Error()
	//	return nil, err
	//}

	signedTxnByteList, err := bcs.Serialize(signedTxn)
	if err != nil {
		err := fmt.Errorf("BuildSignedTransaction signedTxn Serialize failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	signedTxBase64 := base64.StdEncoding.EncodeToString(signedTxnByteList)

	response.Code = common2.ReturnCode_SUCCESS
	response.SignedTx = signedTxBase64
	response.Msg = "BuildSignedTransaction success"
	return response, nil
}

func (c *ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	response := &account.DecodeTransactionResponse{
		Code:     common2.ReturnCode_ERROR,
		Msg:      "",
		Base64Tx: "",
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		response.Msg = msg
		err := fmt.Errorf("DecodeTransaction validateChainAndNetwork fail, err msg = %s", msg)
		return response, err
	}

	return &account.DecodeTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "verify tx success",
		Base64Tx: "0x000000",
	}, nil
}

func (c *ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	response := &account.VerifyTransactionResponse{
		Code:   common2.ReturnCode_ERROR,
		Verify: false,
		Msg:    "",
	}

	// 1. validateChainAndNetwork
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		response.Msg = msg
		return response, fmt.Errorf("validateChainAndNetwork fail: %s", msg)
	}

	if req.PublicKey == "" || req.Signature == "" {
		err := fmt.Errorf("PublicKey or Signature is empty")
		response.Msg = err.Error()
		return response, err
	}

	// 2. DecodeString Signature
	signedTxBytes, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		err = fmt.Errorf("decode signed transaction failed: %w", err)
		response.Msg = err.Error()
		return response, err
	}

	signedTx := &aptos.SignedTransaction{
		Transaction: &aptos.RawTransaction{},
		Authenticator: &aptos.TransactionAuthenticator{
			Variant: aptos.TransactionAuthenticatorEd25519,
			Auth: &crypto.Ed25519Authenticator{
				PubKey: &crypto.Ed25519PublicKey{},
				Sig:    &crypto.Ed25519Signature{},
			},
		},
	}
	des := bcs.NewDeserializer(signedTxBytes)
	signedTx.UnmarshalBCS(des)

	message, err := signedTx.Transaction.SigningMessage()
	if err != nil {
		err = fmt.Errorf("get signing message failed: %w", err)
		response.Msg = err.Error()
		return response, err
	}
	//fmt.Printf("Message to verify (hex): %x\n", message)

	if signedTx.Authenticator.Variant != aptos.TransactionAuthenticatorEd25519 {
		err = fmt.Errorf("invalid authenticator variant")
		response.Msg = err.Error()
		return response, err
	}
	//fmt.Printf("verifyResp Msg: %T\n", signedTx.Authenticator.Auth)
	auth, ok := signedTx.Authenticator.Auth.(*aptos.Ed25519TransactionAuthenticator)
	if !ok {
		err = fmt.Errorf("invalid authenticator type")
		response.Msg = err.Error()
		return response, err
	}
	ed25519Auth, ok := auth.Sender.Auth.(*crypto.Ed25519Authenticator)
	if !ok {
		err = fmt.Errorf("invalid authenticator type")
		response.Msg = err.Error()
		return response, err
	}

	inputPubKey, err := PubKeyHexToPubKey(req.PublicKey)
	if err != nil {
		err = fmt.Errorf("invalid input public key: %w", err)
		response.Msg = err.Error()
		return response, err
	}

	if !bytes.Equal(inputPubKey.Inner[:], ed25519Auth.PubKey.Inner[:]) {
		response.Code = common2.ReturnCode_ERROR
		response.Verify = false
		response.Msg = "public key mismatch"
		return response, nil
	}

	//fmt.Println("Verify message", hex.EncodeToString(message))
	//fmt.Println("Verify Sig", hex.EncodeToString(ed25519Auth.Sig.Bytes()))
	if !ed25519.Verify(ed25519Auth.PubKey.Inner[:], message, ed25519Auth.Sig.Inner[:]) {
		response.Code = common2.ReturnCode_ERROR
		response.Verify = false
		response.Msg = "invalid signature"
		return response, nil
	}

	response.Code = common2.ReturnCode_SUCCESS
	response.Verify = true
	response.Msg = "transaction verified"
	return response, nil
}

func (c *ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	response := &account.ExtraDataResponse{
		Code:  common2.ReturnCode_ERROR,
		Msg:   "",
		Value: "",
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		response.Msg = msg
		err := fmt.Errorf("GetExtraData validateChainAndNetwork fail, err msg = %s", msg)
		return response, err
	}
	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "get extra data success"
	response.Value = "no data"
	return response, nil
}

func CalculateGasFee(gasUnitPrice, totalChargeGasUnits, storageFeeOctas, storageRefundOctas uint64) uint64 {
	// calc base gas fee
	gasFee := gasUnitPrice * totalChargeGasUnits

	// Storage Fee
	netStorageFee := storageFeeOctas - storageRefundOctas

	// totalFee
	totalFee := gasFee + netStorageFee

	return totalFee
}

func GetFeeStatementFromEvents(events []Event) *FeeStatement {
	for _, event := range events {
		if event.Type == "0x1::transaction_fee::FeeStatement" {
			return &event.Data
		}
	}
	return nil
}

func convertExtraInfo(tx TransactionResponse) string {
	extraInfo := map[string]interface{}{
		"vm_status":             tx.VMStatus,
		"accumulator_root_hash": tx.AccumulatorRootHash,
		"changes":               tx.Changes,
		"signature":             tx.Signature,
		"events":                tx.Events,
		"payload":               tx.Payload,
		"success":               tx.Success,
	}

	extraJSON, err := json.Marshal(extraInfo)
	if err != nil {
		return ""
	}
	return string(extraJSON)
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

func determineTransactionType(function string) int32 {
	if strings.Contains(function, "transfer") {
		return 1
	}
	return 0
}

func serializePayload(payload Payload) string {
	data, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return string(data)
}

func validateAptosPublicKey(pubKey string) (bool, string) {
	if pubKey == "" {
		return false, "public key cannot be empty"
	}

	if !strings.HasPrefix(pubKey, "0x") {
		return false, "public key must start with 0x"
	}

	pubKeyWithoutPrefix := strings.TrimPrefix(pubKey, "0x")
	if len(pubKeyWithoutPrefix) != 64 { // Aptos 使用 32 字节的公钥，所以是 64 个十六进制字符
		return false, "invalid public key length, expected 32 bytes (64 hex chars)"
	}

	if _, err := hex.DecodeString(pubKeyWithoutPrefix); err != nil {
		return false, "invalid public key format: must be hex string"
	}

	return true, ""
}
