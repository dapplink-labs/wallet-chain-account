package aptos

import (
	"encoding/base64"
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

type ChainAdaptor struct {
	aptosHttpClient *RestyClient
	aptosClient     *aptos.Client
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	rpcUrl := conf.WalletNode.Aptos.RPCs[0].RPCURL
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
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.SupportChainsResponse{
			Code:    common2.ReturnCode_ERROR,
			Msg:     msg,
			Support: false,
		}, nil
	}
	return &account.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

func (c *ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.ConvertAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	pubKeyHex := req.PublicKey

	pubKey := &crypto.Ed25519PublicKey{}
	err := pubKey.FromHex(pubKeyHex)
	if err != nil {
		log.Error("ConvertAddress pubKey FromHex fail", "err", err)
		return nil, fmt.Errorf("invalid pubKey: %v", err)
	}

	authKey := &crypto.AuthenticationKey{}
	authKey.FromPublicKey(pubKey)

	address := &aptos.AccountAddress{}
	address.FromAuthKey(authKey)

	return &account.ConvertAddressResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "convert address success",
		Address: address.String(),
	}, nil
}

func (c *ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	response := &account.ValidAddressResponse{
		Code:  common2.ReturnCode_ERROR,
		Msg:   "",
		Valid: false,
	}

	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		response.Valid = false
		response.Code = common2.ReturnCode_ERROR
		response.Msg = msg
		return response, nil
	}

	errTooShrot := "AccountAddress too short"
	errTooLong := "AccountAddress too long"

	address := req.Address
	aptosAccountAddress := &aptos.AccountAddress{}
	err := aptosAccountAddress.ParseStringRelaxed(address)
	if err != nil {
		switch err.Error() {
		case errTooShrot:
			response.Valid = false
			response.Code = common2.ReturnCode_ERROR
			response.Msg = errTooShrot
			return response, nil
		case errTooLong:
			response.Valid = false
			response.Code = common2.ReturnCode_ERROR
			response.Msg = errTooLong
			return response, nil
		default:
			tempErr := fmt.Errorf("invalid address format: %v", err)
			response.Valid = false
			response.Code = common2.ReturnCode_ERROR
			response.Msg = tempErr.Error()
			return response, nil
		}
	}

	cleanAddr := address
	if strings.HasPrefix(cleanAddr, "0x") {
		cleanAddr = cleanAddr[2:]
	}

	_, err = hex.DecodeString(cleanAddr)
	if err != nil {
		response.Valid = false
		response.Code = common2.ReturnCode_ERROR
		response.Msg = "address contains invalid hex characters"
		return response, nil
	}

	response.Valid = true
	response.Code = common2.ReturnCode_SUCCESS
	response.Msg = "ValidAddress success"
	return response, nil
}

func (c *ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, ""); !ok {
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	if req.Height == 0 {
		nodeInfo, err := c.aptosHttpClient.GetNodeInfo()
		if err != nil {
			log.Error("GetBlockByNumber GetNodeInfo fail", "err", err)
			return &account.BlockResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "GetBlockByNumber fail",
			}, nil
		}
		if req.ViewTx {

		}
		return &account.BlockResponse{
			Code:   common2.ReturnCode_SUCCESS,
			Msg:    "GetBlockByNumber GetNodeInfo success",
			Height: int64(nodeInfo.BlockHeight),
			// TODO: Transactionsdasda
			Transactions: nil,
		}, nil
	}

	blockResponse, err := c.aptosHttpClient.GetBlockByHeight(uint64(req.Height))
	if err != nil {
		log.Error("GetBlockByNumber GetBlockByHeight fail", "err", err)
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "GetBlockByNumber GetBlockByHeight fail",
		}, nil
	}
	if req.ViewTx {

	}
	return &account.BlockResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "GetBlockByNumber GetBlockByHeight success",
		Height: int64(blockResponse.BlockHeight),
		Hash:   blockResponse.BlockHash,
		// TODO: Transactionsdasda
		Transactions: nil,
	}, nil
}

func (c *ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, ""); !ok {
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	// no implement
	panic("implement me")
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	if req.Height == 0 {
		nodeInfo, err := c.aptosHttpClient.GetNodeInfo()
		if err != nil {
			log.Error("GetBlockHeaderByNumber GetNodeInfo fail", "err", err)
			return &account.BlockHeaderResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "GetBlockHeaderByNumber GetNodeInfo fail",
			}, nil
		}
		seconds := nodeInfo.LedgerTimestamp / 1_000_000
		blockHead := &account.BlockHeader{
			Number: strconv.FormatUint(nodeInfo.BlockHeight, 10),
			Time:   seconds,
		}
		return &account.BlockHeaderResponse{
			Code:        common2.ReturnCode_SUCCESS,
			Msg:         "GetBlockHeaderByNumber GetNodeInfo success",
			BlockHeader: blockHead,
		}, nil
	}

	blockResponse, err := c.aptosHttpClient.GetBlockByHeight(uint64(req.Height))
	if err != nil {
		log.Error("GetBlockHeaderByNumber GetBlockByHeight fail", "err", err)
		return &account.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "GetBlockHeaderByNumber GetBlockByHeight fail",
		}, nil
	}
	blockHead := &account.BlockHeader{
		Hash:   blockResponse.BlockHash,
		Number: strconv.FormatUint(blockResponse.BlockHeight, 10),
		Time:   blockResponse.BlockTimestamp,
	}
	return &account.BlockHeaderResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "GetBlockHeaderByNumber GetBlockByHeight success",
		BlockHeader: blockHead,
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	// no implement
	panic("implement me")
}

func (c *ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.AccountResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	accountResponse, err := c.aptosHttpClient.GetAccount(req.Address)
	if err != nil {
		log.Error("GetAccount fail", "err", err)
		return &account.AccountResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "GetAccount fail",
		}, nil
	}
	return &account.AccountResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "get account response success",
		Sequence: strconv.FormatUint(accountResponse.SequenceNumber, 10),
		Network:  req.Network,
	}, nil
}

func (c *ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.FeeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	response, err := c.aptosHttpClient.GetGasPrice()
	if err != nil {
		log.Error("GetFee fail", "err", err)
		return &account.FeeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "GetFee fail",
		}, nil
	}
	return &account.FeeResponse{
		Code:      common2.ReturnCode_SUCCESS,
		Msg:       "GetFee success",
		SlowFee:   strconv.FormatUint(response.DeprioritizedGasEstimate, 10),
		NormalFee: strconv.FormatUint(response.GasEstimate, 10),
		FastFee:   strconv.FormatUint(response.PrioritizedGasEstimate, 10),
	}, nil
}

func (c *ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.SendTxResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.TxAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	transactionsPtr, err := c.aptosHttpClient.GetTransactionByAddress(req.Address)
	if err != nil {
		log.Error("GetTxByAddress GetTransactionByAddress fail", "err", err)
		return &account.TxAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "GetTxByAddress GetTransactionByAddress fail",
		}, nil
	}
	if transactionsPtr == nil {
		return &account.TxAddressResponse{
			Code: common2.ReturnCode_SUCCESS,
			Msg:  "GetTxByAddress success but no transactions found",
			Tx:   []*account.TxMessage{},
		}, nil
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

	return &account.TxAddressResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "GetTxByAddress success",
		Tx:   txMessages,
	}, nil
}

func (c *ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	tx, err := c.aptosHttpClient.GetTransactionByHash(req.Hash)
	if err != nil {
		log.Error("GetTransactionByHash error", "err", err)
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "GetTransactionByHash error",
		}, nil
	}
	//TODO fromAddrList toAddrsList
	var fromAddrList []*account.Address
	var toAddrsList []*account.Address
	var valueList []*account.Value
	var txStatus account.TxStatus
	if tx.Success {
		txStatus = account.TxStatus_Success
	} else {
		txStatus = account.TxStatus_Failed
	}

	feeStatement := GetFeeStatementFromEvents(tx.Events)
	totalFee := CalculateGasFee(tx.GasUnitPrice, feeStatement.TotalChargeGasUnits, feeStatement.StorageFeeOctas, feeStatement.StorageFeeRefundOctas)

	return &account.TxHashResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "GetTransactionByHash success",
		Tx: &account.TxMessage{
			Hash:   tx.Hash,
			Froms:  fromAddrList,
			Tos:    toAddrsList,
			Values: valueList,
			Fee:    strconv.FormatUint(totalFee, 10),
			Status: txStatus,
			Type:   0,
			//Height:          tx.,
			//ContractAddress: tx.To().String(),
			//Data: hexutils.BytesToHex(tx.Data()),
		},
	}, nil
}

func (c *ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.BlockByRangeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	startVersion, err := strconv.ParseUint(req.Start, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid start version: %w", err)
	}
	endVersion, err := strconv.ParseUint(req.End, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid end version: %w", err)
	}
	if startVersion > endVersion {
		return nil, fmt.Errorf("start version (%d) cannot be greater than end version (%d)", startVersion, endVersion)
	}
	txs, err := c.aptosHttpClient.GetTransactionByVersionRange(startVersion, endVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	response := &account.BlockByRangeResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "GetBlockByRange success",
		BlockHeader: make([]*account.BlockHeader, 0, len(txs)),
	}
	for _, tx := range txs {
		blockHeader := &account.BlockHeader{
			Hash: tx.Hash,
			//ParentHash:  tx.StateRootHash,
			//Root:        tx.StateRootHash,
			TxHash:      tx.Hash,
			ReceiptHash: tx.EventRootHash,
			//Number:      tx.Version,
			GasLimit: tx.MaxGasAmount,
			GasUsed:  tx.GasUsed,
			Time:     tx.Timestamp,
			Extra:    convertExtraInfo(tx),
			Nonce:    strconv.FormatUint(tx.SequenceNumber, 10),
		}

		response.BlockHeader = append(response.BlockHeader, blockHeader)
	}
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
		return response, nil
	}

	jsonBytes, err := base64.StdEncoding.DecodeString(req.Base64Tx)
	if err != nil {
		log.Error("CreateUnSignTransaction DecodeString fail", "err", err)
		return nil, err
	}

	var data TransferRequest
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		log.Error("CreateUnSignTransaction Unmarshal fail", "err", err)
		return nil, err
	}
	if data.Amount == 0 {
		return nil, fmt.Errorf("transfer amount must be greater than 0")
	}
	transferAmount := data.Amount
	fromAddress, err := AddressToAccountAddress(data.FromAddress)
	if err != nil {
		log.Error("FromAddress AddressToAccountAddress fail", "err", err)
		return nil, err
	}
	toAddress, err := AddressToAccountAddress(data.ToAddress)
	if err != nil {
		log.Error("ToAddress AddressToAccountAddress fail", "err", err)
		return nil, err
	}
	// TODO Need to support more coinType
	transferPayload, err := aptos.CoinTransferPayload(nil, toAddress, transferAmount)
	if err != nil {
		log.Error("aptos CoinTransferPayload fail", "err", err)
		return nil, err
	}

	rawTxn, err := c.aptosClient.BuildTransaction(
		fromAddress,
		aptos.TransactionPayload{Payload: transferPayload},
	)
	rawTxnBytes, err := bcs.Serialize(rawTxn)
	if err != nil {
		log.Error("rawTxn Serialize fail", "err", err)
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
		return response, nil
	}
	if req.Base64Tx == "" || req.Signature == "" {
		err := fmt.Errorf("req.Base64Tx or req.Signature is empty")
		log.Error("rawTxn Serialize fail", "err", err)
		return nil, err
	}

	rawTxBytes, err := base64.StdEncoding.DecodeString(req.Base64Tx)
	if err != nil {
		log.Error("rawTxn DecodeString fail", "err", err)
		return nil, err
	}
	authBytes, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		log.Error("signature DecodeString fail", "err", err)
		response.Msg = "req.Signature DecodeString error"
		return response, nil
	}

	des := bcs.NewDeserializer(rawTxBytes)
	rawTxn := &aptos.RawTransaction{}
	rawTxn.UnmarshalBCS(des)
	if des.Error() != nil {
		log.Error("DeserializeBCS fail", "err", des.Error())
		return nil, des.Error()
	}
	fmt.Printf("BuildSignedTransaction Sender: %s\n", rawTxn.Sender)
	fmt.Printf("BuildSignedTransaction Sequence Number: %d\n", rawTxn.SequenceNumber)
	fmt.Printf("BuildSignedTransaction Max Gas Amount: %d\n", rawTxn.MaxGasAmount)
	fmt.Printf("BuildSignedTransaction Gas Unit Price: %d\n", rawTxn.GasUnitPrice)
	fmt.Printf("BuildSignedTransaction Expiration: %d\n", rawTxn.ExpirationTimestampSeconds)
	fmt.Printf("BuildSignedTransaction Chain ID: %d\n", rawTxn.ChainId)

	authDes := bcs.NewDeserializer(authBytes)
	accountAuth := &crypto.AccountAuthenticator{}
	accountAuth.UnmarshalBCS(authDes)
	if authDes.Error() != nil {
		log.Error("Unmarshal authenticator fail", "err", authDes.Error())
		return nil, authDes.Error()
	}
	fmt.Printf("accountAuth PubKey: %s\n", accountAuth.PubKey().ToHex())
	fmt.Printf("accountAuth Signature: %s\n", accountAuth.Signature().ToHex())

	txnAuth, err := aptos.NewTransactionAuthenticator(accountAuth)
	if err != nil {
		log.Error("NewTransactionAuthenticator fail", "err", err)
		return nil, err
	}

	signedTxn := &aptos.SignedTransaction{
		Transaction:   rawTxn,
		Authenticator: txnAuth,
	}

	signedTxnJson, _ := json.Marshal(signedTxn)
	fmt.Printf("signedTxnJson: %s\n", signedTxnJson)

	signedTxnSer, err := bcs.Serialize(signedTxn)
	if err != nil {
		log.Error("signedTxn Serialize fail", "err", err)
		return nil, err
	}

	signedTxBase64 := base64.StdEncoding.EncodeToString(signedTxnSer)

	fmt.Printf("signedTxBase64: %s\n", signedTxBase64)

	response.Code = common2.ReturnCode_SUCCESS
	response.SignedTx = signedTxBase64
	response.Msg = "BuildSignedTransaction success"
	return response, nil
}

func (c *ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.DecodeTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}

	return &account.DecodeTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "verify tx success",
		Base64Tx: "0x000000",
	}, nil
}

func (c *ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.VerifyTransactionResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}

	signedTxBytes, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		return nil, fmt.Errorf("decode signed transaction failed: %w", err)
	}

	// 2.
	signedTx := &aptos.SignedTransaction{}
	signedTx.UnmarshalBCS(bcs.NewDeserializer(signedTxBytes))

	// 3.
	//rawTx := signedTx.Transaction.(*aptos.RawTransaction)
	//signingMessage, err := rawTx.SigningMessage()
	//if err != nil {
	//	return nil, fmt.Errorf("SigningMessage failed: %w", err)
	//}
	//
	//// 4.
	//if ed25519Auth, ok := signedTx.Authenticator.Auth.(*crypto.Ed25519Authenticator); ok {
	//
	//pubKeyBytes := ed25519Auth.PubKey.Bytes()
	//authKey := crypto.HashSHA3256(append(pubKeyBytes, 0x00))
	//expectedSender := aptos.AccountAddress{}
	//copy(expectedSender[:], authKey)
	//
	//if expectedSender != rawTx.Sender {
	//	return &account.VerifyTransactionResponse{
	//		IsValid: false,
	//		Message: "sender address does not match public key",
	//	}, nil
	//}

	//
	//isValid := ed25519Auth.PubKey.Verify(
	//	signingMessage,
	//	ed25519Auth.Sig,
	//)

	//return &account.VerifyTransactionResponse{
	//	IsValid: isValid,
	//	Message: fmt.Sprintf("signature verification result: %v", isValid),
	//}, nil
	//}
	return &account.VerifyTransactionResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "verify tx success",
		Verify: true,
	}, nil
}

func (c *ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	if ok, msg := validateChainAndNetwork(req.Chain, req.Network); !ok {
		return &account.ExtraDataResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  msg,
		}, nil
	}
	return &account.ExtraDataResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "get extra data success",
		Value: "not data",
	}, nil
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
