package aptos

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
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

	jsonBytes, err := base64.StdEncoding.DecodeString(req.Base64Tx)
	if err != nil {
		err := fmt.Errorf("CreateUnSignTransaction DecodeString failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	var data TransferRequest
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		err := fmt.Errorf("CreateUnSignTransaction Unmarshal failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	if data.Amount == 0 {
		err := fmt.Errorf("CreateUnSignTransaction data.Amount == 0 failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	transferAmount := data.Amount
	fromAddress, err := AddressToAccountAddress(data.FromAddress)
	if err != nil {
		err := fmt.Errorf("CreateUnSignTransaction FromAddress failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	toAddress, err := AddressToAccountAddress(data.ToAddress)
	if err != nil {
		err := fmt.Errorf("CreateUnSignTransaction ToAddress failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	// TODO Need to support more coinType
	transferPayload, err := aptos.CoinTransferPayload(nil, toAddress, transferAmount)
	if err != nil {
		err := fmt.Errorf("CreateUnSignTransaction CoinTransferPayload failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	rawTxn, err := c.aptosClient.BuildTransaction(
		fromAddress,
		aptos.TransactionPayload{Payload: transferPayload},
	)
	rawTxnBytes, err := bcs.Serialize(rawTxn)
	if err != nil {
		err := fmt.Errorf("CreateUnSignTransaction rawTxn Serialize failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	base64Tx := base64.StdEncoding.EncodeToString(rawTxnBytes)

	response.Code = common2.ReturnCode_SUCCESS
	response.UnSignTx = base64Tx
	response.Msg = "CreateUnSignTransaction success"
	return response, err
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
	if req.Base64Tx == "" || req.Signature == "" {
		err := fmt.Errorf("req.Base64Tx or req.Signature is empty")
		log.Error("BuildSignedTransaction req.Base64Tx or req.Signature is empty", "err", err)
		return nil, err
	}

	rawTxBytes, err := base64.StdEncoding.DecodeString(req.Base64Tx)
	if err != nil {
		err := fmt.Errorf("BuildSignedTransaction DecodeString rawTx failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}
	authBytes, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		err := fmt.Errorf("BuildSignedTransaction DecodeString signature failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	des := bcs.NewDeserializer(rawTxBytes)
	rawTxn := &aptos.RawTransaction{}
	rawTxn.UnmarshalBCS(des)
	if des.Error() != nil {
		err := fmt.Errorf("BuildSignedTransaction rawTxn.UnmarshalBCS failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	authDes := bcs.NewDeserializer(authBytes)
	accountAuth := &crypto.AccountAuthenticator{}
	accountAuth.UnmarshalBCS(authDes)
	if authDes.Error() != nil {
		err := fmt.Errorf("BuildSignedTransaction accountAuth.UnmarshalBCS failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	txnAuth, err := aptos.NewTransactionAuthenticator(accountAuth)
	if err != nil {
		err := fmt.Errorf("BuildSignedTransaction NewTransactionAuthenticator failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	signedTxn := &aptos.SignedTransaction{
		Transaction:   rawTxn,
		Authenticator: txnAuth,
	}

	signedTxnSer, err := bcs.Serialize(signedTxn)
	if err != nil {
		err := fmt.Errorf("BuildSignedTransaction signedTxn Serialize failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	signedTxBase64 := base64.StdEncoding.EncodeToString(signedTxnSer)

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
		Msg:    "",
		Verify: false,
	}

	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		response.Msg = msg
		err := fmt.Errorf("VerifySignedTransaction validateChainAndNetwork fail, err msg = %s", msg)
		return response, err
	}

	signedTxBytes, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		err := fmt.Errorf("VerifySignedTransaction DecodeString Signature failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return nil, err
	}

	// 2. Deserializer
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
	signedTx.UnmarshalBCS(bcs.NewDeserializer(signedTxBytes))

	var messages []string
	isValid := true

	rawTxn, ok := signedTx.Transaction.(*aptos.RawTransaction)
	if !ok {
		err := fmt.Errorf("VerifySignedTransaction signedTx.Transaction rawTxn failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return response, err
	}
	signingMessage, err := rawTxn.SigningMessage()
	if err != nil {
		err := fmt.Errorf("VerifySignedTransaction rawTxn.SigningMessage failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return response, err
	}

	auth := signedTx.Authenticator.Auth
	if senderAuth, ok := auth.(*aptos.Ed25519TransactionAuthenticator); ok {
		if ed25519Auth, ok := senderAuth.Sender.Auth.(*crypto.Ed25519Authenticator); ok {
			// Step 1: Derive account address from the public key in signedTxn.Authenticator.Auth
			accountAddress, err := PubKeyToAccountAddress(ed25519Auth.PubKey)
			if err != nil {
				err := fmt.Errorf("VerifySignedTransaction PubKeyToAccountAddress(ed25519Auth.PubKey) failed: %w", err)
				log.Error("err", err)
				response.Msg = err.Error()
				return response, err
			}
			// Step 2: Verify if the derived address matches the sender address in raw transaction
			// This ensures the transaction is truly from the claimed sender
			// rawTxn.Sender from rawTxn, this is req.Sender.address
			// accountAddress from signedTxn.Authenticator.Auth, this is signedTxn.address
			if *accountAddress != rawTxn.Sender {
				isValid = false
				messages = append(messages, fmt.Sprintf("sender address mismatch\nexpected: %s\nactual: %s",
					accountAddress.String(), rawTxn.Sender.String()))
			}

			// Verify if the signature is valid for this transaction
			// ed25519Auth.PubKey: the public key of the signer
			// Parameters:
			// - signingMessage: the original transaction data to be signed
			// - ed25519Auth.Sig: the signature created by the private key
			if !ed25519Auth.PubKey.Verify(signingMessage, ed25519Auth.Sig) {
				isValid = false
				messages = append(messages, "invalid signature")
			}
		} else {
			err := fmt.Errorf("VerifySignedTransaction invalid ed25519 authenticator type failed: %w", err)
			log.Error("err", err)
			response.Msg = err.Error()
			return response, err
		}
	} else {
		err := fmt.Errorf("VerifySignedTransaction invalid sender authenticator type failed: %w", err)
		log.Error("err", err)
		response.Msg = err.Error()
		return response, err
	}

	// Step 4: Verify transaction basic parameters
	messages = append(messages, fmt.Sprintf("\n=== Transaction Parameters ==="+
		"\nSender Address: %s"+
		"\nSequence Number: %d"+
		"\nGas Limit: %d"+
		"\nGas Unit Price: %d"+
		"\nExpiration Time: %d"+
		"\nChain ID: %d",
		rawTxn.Sender.String(),
		rawTxn.SequenceNumber,
		rawTxn.MaxGasAmount,
		rawTxn.GasUnitPrice,
		rawTxn.ExpirationTimestampSeconds,
		rawTxn.ChainId))

	// Step 5: Verify transfer parameters if it's a transfer transaction
	if entryFunction, ok := rawTxn.Payload.Payload.(*aptos.EntryFunction); ok {
		messages = append(messages, "\n=== Transfer Parameters ===")

		// 1. Verify module and function name
		if entryFunction.Module.Name != "aptos_account" {
			isValid = false
			messages = append(messages, "invalid module name")
		}
		if entryFunction.Function != "transfer" {
			isValid = false
			messages = append(messages, "invalid function name")
		}

		// 2. Verify arguments length
		if len(entryFunction.Args) < 2 {
			//isValid = false
			messages = append(messages, "insufficient transfer arguments")
			//return isValid, strings.Join(messages, "\n")
			tempMessgae := strings.Join(messages, "\n")
			err := fmt.Errorf("VerifySignedTransaction %s : %w", tempMessgae, err)
			log.Error("err", err)
			response.Msg = err.Error()
			return response, err
		}

		// 3. Verify recipient address
		toAddrBytes := entryFunction.Args[0]
		if len(toAddrBytes) != 32 { // Aptos address length
			isValid = false
			messages = append(messages, "invalid recipient address length")
		}
		messages = append(messages, fmt.Sprintf("Recipient Address (bytes): %x", toAddrBytes))

		// 4. Verify transfer amount
		amountBytes := entryFunction.Args[1]
		if len(amountBytes) != 8 { // uint64 length
			isValid = false
			messages = append(messages, "invalid amount format")
		}
		amount := binary.LittleEndian.Uint64(amountBytes)

		// Check amount constraints
		if amount == 0 {
			isValid = false
			messages = append(messages, "transfer amount cannot be zero")
		}
		const MAX_TRANSFER_AMOUNT = uint64(1000000000000)
		if amount > MAX_TRANSFER_AMOUNT {
			isValid = false
			messages = append(messages, fmt.Sprintf("transfer amount exceeds maximum limit: %d", MAX_TRANSFER_AMOUNT))
		}
		messages = append(messages, fmt.Sprintf("Transfer Amount: %d", amount))

		// 5. Verify sender has sufficient balance
		// totalRequired := amount + (rawTxn.MaxGasAmount * rawTxn.GasUnitPrice)
		// if senderBalance < totalRequired {
		//     isValid = false
		//     messages = append(messages, "insufficient balance for transfer and gas")
		// }

		// 6. Verify recipient address is valid
		if bytes.Equal(toAddrBytes, rawTxn.Sender[:]) {
			isValid = false
			messages = append(messages, "cannot transfer to self")
		}
	}

	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "VerifySignedTransaction success"
	response.Verify = isValid
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
