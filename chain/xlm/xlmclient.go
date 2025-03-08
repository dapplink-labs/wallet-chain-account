package xlm

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/stellar/go/network"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/txnbuild"

	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

type XlmClient struct {
	client          *http.Client
	rpcURL          string
	data_api_rpcURL string
}

func NewXlmClients(conf *config.Config) (*XlmClient, error) {
	return &XlmClient{
		client:          &http.Client{},
		rpcURL:          conf.WalletNode.Xlm.RpcUrl,
		data_api_rpcURL: conf.WalletNode.Xlm.DataApiUrl,
	}, nil
}

func (xc *XlmClient) HttpProcess(requestData interface{}, rpcURL string, respData interface{}) error {
	// 将请求数据序列化为JSON，设置请求Http
	var req *http.Request
	var err error
	if requestData != nil {
		jsonData, err := json.Marshal(requestData)
		if err != nil {
			log.Fatalf("Failed to marshal request data: %v", err)
			return err
		}

		req, err = http.NewRequest(http.MethodPost, rpcURL, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
			return err
		}
	} else {
		req, err = http.NewRequest(http.MethodGet, rpcURL, nil)
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
			return err
		}
	}

	// 发送请求
	resp, err := xc.client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
		return err
	}
	defer resp.Body.Close()

	// 读取响应数据
	var content bytes.Buffer
	if _, err = io.Copy(&content, resp.Body); err != nil {
		log.Fatalf("Failed to reading response body: %v", err)
		return err
	}

	// 返回状态
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("response status %v %s, response body: %s", resp.StatusCode, resp.Status, content.String())
		return err
	}

	err = json.Unmarshal(content.Bytes(), &respData)
	if err != nil {
		log.Fatalf("Failed to read json response body: %v", err)
		return err
	}

	return nil
}

func (xc *XlmClient) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	isValidAddress := strkey.IsValidEd25519PublicKey(req.Address)
	if isValidAddress {
		return &account.ValidAddressResponse{
			Code:  common.ReturnCode_SUCCESS,
			Msg:   "ValidAddress Success",
			Valid: isValidAddress,
		}, nil
	} else {
		return &account.ValidAddressResponse{
			Code:  common.ReturnCode_SUCCESS,
			Msg:   "ValidAddress Failed",
			Valid: isValidAddress,
		}, nil
	}
}

func (xc *XlmClient) GetAccountInfo(addr string) (*account.AccountResponse, error) {
	var url = fmt.Sprintf("%s/accounts/%s", xc.data_api_rpcURL, addr)
	var result ResponseAccountInfo
	err := xc.HttpProcess(nil, url, &result)
	if err != nil {
		log.Fatalf("Failed to HttpProcess: %v", err)
		return &account.AccountResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "GetAccountInfo Failed",
		}, err
	}

	return &account.AccountResponse{
		Code:          common.ReturnCode_SUCCESS,
		Msg:           "GetAccountInfo Success",
		Network:       "mainnet",
		AccountNumber: result.AccountID,
		Sequence:      result.Sequence,
		Balance:       result.Balances[0].Balance,
	}, nil
}

func (xc *XlmClient) GetFee() (*account.FeeResponse, error) {
	// 创建请求数据
	requestData := RequestGetFeeStats{
		Jsonrpc: "2.0",
		ID:      8675309,
		Method:  "getFeeStats",
	}

	//创建接收数据
	var result ResponseGetFeeStats

	err := xc.HttpProcess(requestData, xc.rpcURL, &result)
	if err != nil {
		log.Fatalf("Failed to HttpProcess: %v", err)
		return &account.FeeResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "GetFee Failed",
		}, nil
	}

	return &account.FeeResponse{
		Code:      common.ReturnCode_SUCCESS,
		Msg:       "GetFee Success",
		SlowFee:   result.Result.InclusionFee.Min,
		NormalFee: result.Result.InclusionFee.Mode,
		FastFee:   result.Result.InclusionFee.Max,
	}, nil
}

func (xc *XlmClient) GetBlockByNumber(blockNumber *big.Int) (*account.BlockResponse, error) {
	/*
		https://horizon.stellar.org/ledgers/55983416/transactions
		这里只拿"交易简要信息"列表，作为示例

		需要再使用 "method": "getTransaction" 重新拿一次"交易详细信息"，暂不支持这种自动二次查询，请自行使用GetTxByHash进行查询。
	*/
	var url = fmt.Sprintf("%s/ledgers/%d/transactions", xc.data_api_rpcURL, blockNumber)
	//fmt.Println("url ", url)
	var result ResponseGetTransactionForLedgers
	err := xc.HttpProcess(nil, url, &result)
	if err != nil {
		log.Fatalf("Failed to HttpProcess: %v", err)
		return &account.BlockResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "GetBlockByNumber Failed",
		}, nil
	}

	var txListRet []*account.BlockInfoTransactionList
	for _, v := range result.Embedded.Records {
		bitlItem := &account.BlockInfoTransactionList{
			From:           v.Source_account,
			To:             "Not Support In This Function(Please Use GetTransactionByHash For Detail)",
			TokenAddress:   "",
			ContractWallet: "",
			Hash:           v.Hash,
			Height:         uint64(v.Ledger),
			Amount:         "Not Support In This Function(Please Use GetTransactionByHash For Detail)",
		}
		txListRet = append(txListRet, bitlItem)
	}

	var blockHash = ""
	if len(txListRet) > 0 {
		blockHash = result.Embedded.Records[0].Hash
	} else {
		blockHash = "nil"
	}

	return &account.BlockResponse{
		Code:         common.ReturnCode_SUCCESS,
		Msg:          "GetBlockByNumber Success",
		Height:       blockNumber.Int64(),
		Hash:         blockHash,
		BaseFee:      "not support",
		Transactions: txListRet,
	}, nil
}

func (xc *XlmClient) GetBlockHeaderByNumber(blockNumber *big.Int) (*account.BlockHeaderResponse, error) {
	var url = fmt.Sprintf("%s/ledgers/%d", xc.data_api_rpcURL, blockNumber)
	var result ResponseGetBlockHeader
	err := xc.HttpProcess(nil, url, &result)
	if err != nil {
		log.Fatalf("Failed to HttpProcess: %v", err)
		return &account.BlockHeaderResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "GetBlockHeaderByNumber Failed",
		}, nil
	}

	blockHeader := &account.BlockHeader{
		Hash:       result.Hash,
		ParentHash: result.PrevHash,
		CoinBase:   result.TotalCoins,
		Number:     strconv.Itoa(result.Sequence),
		BaseFee:    strconv.Itoa(result.BaseFeeInStroops),
	}
	return &account.BlockHeaderResponse{
		Code:        common.ReturnCode_SUCCESS,
		Msg:         "GetBlockHeaderByNumber Success",
		BlockHeader: blockHeader,
	}, nil
}

func (xc *XlmClient) GetTransactionByHash(txHash string) (*account.TxHashResponse, error) {
	var url = ""
	// 第一步：获取"交易体Part01"
	url = fmt.Sprintf("%s/transactions/%s", xc.data_api_rpcURL, txHash)
	var resultPart01 ResponseGetTransactionPart01
	err := xc.HttpProcess(nil, url, &resultPart01)
	if err != nil {
		log.Fatalf("(ResponseGetTransactionPart01) Failed to HttpProcess: %v", err)
		return &account.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "GetTransactionByHash Part01 Failed",
		}, nil
	}

	////测试打印原始Response
	//jsExpression, _ := json.MarshalIndent(resultPart01, "", "    ")
	//fmt.Println("Print Origin Json")
	//fmt.Println(string(jsExpression))

	// 第二步：获取"交易体效果"
	url = fmt.Sprintf("%s/transactions/%s/effects", xc.data_api_rpcURL, txHash)
	var resultEffect ResponseGetTransactionEffect
	err = xc.HttpProcess(nil, url, &resultEffect)
	if err != nil {
		log.Fatalf("(ResponseGetTransactionEffect) Failed to HttpProcess: %v", err)
		return &account.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "GetTransactionByHash Effect Failed",
		}, nil
	}

	if len(resultEffect.Embedded.Records) != 2 {
		log.Fatalf("Failed to HttpProcess: %v", err)
		return &account.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "GetTransactionByHash Failed [resultEffect.Embedded.Records] length != 2",
		}, nil
	}
	var txStatus account.TxStatus
	if resultPart01.Successful {
		txStatus = account.TxStatus_Success
	} else {
		txStatus = account.TxStatus_Failed
	}

	txID, err := strconv.ParseUint(resultPart01.ID, 10, 32)

	tx := &account.TxMessage{
		Hash:            resultPart01.Hash,
		Index:           uint32(txID),
		From:            resultEffect.Embedded.Records[1].Account,
		To:              resultEffect.Embedded.Records[0].Account,
		Value:           resultEffect.Embedded.Records[0].Amount,
		Fee:             resultPart01.FeeCharged,
		Status:          txStatus,
		Type:            0,
		Height:          strconv.FormatUint(uint64(resultPart01.Ledger), 10),
		ContractAddress: "Sorry, is currently not supported...",
	}

	return &account.TxHashResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "GetTransactionByHash Success",
		Tx:   tx,
	}, nil
}

func (xc *XlmClient) CreateUnsignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(req.Base64Tx)
	if err != nil {
		log.Fatalf("CreateUnsignTransaction Failed %v", err)
		return &account.UnSignTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "CreateUnsignTransaction Failed",
		}, nil
	}

	// 将 JSON 字节数组反序列化为结构体
	var decodedTx RequestCreateUnsignTransaction
	err = json.Unmarshal(decodedBytes, &decodedTx)
	if err != nil {
		log.Fatalf("CreateUnsignTransaction Failed %v", err)
		return &account.UnSignTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "CreateUnsignTransaction Failed",
		}, nil
	}

	// 开始构建交易
	myOperate := txnbuild.Payment{
		SourceAccount: decodedTx.AddrFrom,
		Destination:   decodedTx.AddrTo,
		Amount:        decodedTx.Amount,
		Asset:         txnbuild.NativeAsset{},
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: decodedTx.AddrFrom,
				Sequence:  int64(decodedTx.SequenceFrom),
			},
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&myOperate},
			BaseFee:              txnbuild.MinBaseFee,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()}, // Use a real timeout in production!
		},
	)
	if err != nil {
		log.Fatalf("CreateUnsignTransaction Failed %v", err)
		return &account.UnSignTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "CreateUnsignTransaction Failed",
		}, nil
	}

	base64_tx, err := tx.Base64()
	if err != nil {
		log.Fatalf("CreateUnsignTransaction Failed %v", err)
		return &account.UnSignTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "CreateUnsignTransaction Failed",
		}, nil
	}

	return &account.UnSignTransactionResponse{
		Code:     common.ReturnCode_SUCCESS,
		Msg:      "CreateUnsignTransaction Success",
		UnSignTx: base64_tx,
	}, nil
}

func (xc *XlmClient) SignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	tx_Generic_UnSign, err := txnbuild.TransactionFromXDR(req.Base64Tx)
	if err != nil {
		log.Fatalf("SignedTransaction: %v", err)
		return &account.SignedTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "SignedTransaction Failed",
		}, nil
	}

	tx_UnSign, success := tx_Generic_UnSign.Transaction()
	if success != true {
		return &account.SignedTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "SignedTransaction Failed",
		}, nil
	}

	if tx_UnSign.SourceAccount().AccountID != req.PublicKey {
		return &account.SignedTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "SignedTransaction Failed",
		}, nil
	}

	tx_UnSign_New, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount: &txnbuild.SimpleAccount{
				AccountID: tx_UnSign.SourceAccount().AccountID,
				Sequence:  tx_UnSign.SequenceNumber() - 1,
			},
			IncrementSequenceNum: true,
			Operations:           tx_UnSign.Operations(),
			BaseFee:              tx_UnSign.BaseFee(),
			Preconditions:        txnbuild.Preconditions{TimeBounds: tx_UnSign.Timebounds()},
		},
	)
	if err != nil {
		log.Fatalf("SignedTransaction: %v", err)
		return &account.SignedTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "SignedTransaction Failed",
		}, nil
	}

	tx_UnSign_New, err = tx_UnSign_New.AddSignatureBase64(network.PublicNetworkPassphrase, req.PublicKey, req.Signature)
	if err != nil {
		log.Fatalf("SignedTransaction: %v", err)
		return &account.SignedTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "SignedTransaction Failed",
		}, nil
	}

	base64_AfterSign, err := tx_UnSign_New.Base64()
	if err != nil {
		log.Fatalf("SignedTransaction: %v", err)
		return &account.SignedTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "SignedTransaction Failed",
		}, nil
	}

	return &account.SignedTransactionResponse{
		Code:     common.ReturnCode_SUCCESS,
		Msg:      "SignedTransaction Success",
		SignedTx: base64_AfterSign,
	}, nil
}

func (xc *XlmClient) SendTx(rawTx string) (*account.SendTxResponse, error) {
	// 创建请求数据
	requestData := RequestSendTransaction{
		Jsonrpc: "2.0",
		ID:      2,
		Method:  "sendTransaction",
		Params: struct {
			Transaction string `json:"transaction"`
			XdrFormat   string `json:"xdrFormat"`
		}{
			Transaction: rawTx,
			XdrFormat:   "base64",
		},
	}

	//创建接收数据
	var result ResponseSendTransaction
	err := xc.HttpProcess(requestData, xc.rpcURL, &result)
	if err != nil {
		log.Fatalf("Failed to HttpProcess: %v", err)
		return &account.SendTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "SendTransactionByHash Failed",
		}, nil
	}

	//fmt.Println("SendTransactionByHash Status", result.Result.Status)
	if result.Result.Status == "PENDING" {
		return &account.SendTxResponse{
			Code:   common.ReturnCode_SUCCESS,
			Msg:    "SendTx " + result.Result.Status,
			TxHash: result.Result.Hash,
		}, nil
	} else {
		return &account.SendTxResponse{
			Code:   common.ReturnCode_ERROR,
			Msg:    "SendTx " + result.Result.Status,
			TxHash: result.Result.Hash,
		}, nil
	}
}
