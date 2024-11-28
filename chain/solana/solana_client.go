package solana

import (
	"context"
	"errors"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	gresty "github.com/go-resty/resty/v2"
)

type SolClient interface {
	GetHealth() (string, error)

	GetAccountInfo(inputAddr string) (*AccountInfo, error)
	GetBalance(inputAddr string) (uint64, error)
	GetLatestBlockhash(commitmentType CommitmentType) (string, error)
	SendTransaction(
		signedTx string,
		config *SendTransactionRequest,
	) (string, error)
	SimulateTransaction(
		signedTx string,
		config *SimulateRequest,
	) (*SimulateResult, error)

	GetFeeForMessage(message string) (uint64, error)
	GetRecentPrioritizationFees() ([]PrioritizationFee, error)

	GetSlot(commitment CommitmentType) (uint64, error)
	GetBlocksWithLimit(startSlot uint64, limit uint64) ([]uint64, error)
	GetBlockBySlot(slot uint64, detailType TransactionDetailsType) (*BlockResult, error)

	GetTransaction(signature string) (*TransactionResult, error)
	GetTransactionRange(signatures []string) ([]*TransactionResult, error)
	GetTxForAddress(
		address string,
		commitment CommitmentType,
		limit uint64,
		beforeSignature string,
		untilSignature string,
	) ([]SignatureInfo, error)
}

type CommitmentType string

const (
	// Finalized Confirmed Processed
	// Finalized wait 32 slot
	Finalized CommitmentType = "finalized"
	// Confirmed wait 2-3 slot
	Confirmed CommitmentType = "confirmed"
	// Processed wait 0 slot
	Processed CommitmentType = "processed"
)

type TransactionDetailsType string

const (
	Full       TransactionDetailsType = "full"
	Accounts   TransactionDetailsType = "accounts"
	Signatures TransactionDetailsType = "signatures"
	None       TransactionDetailsType = "none"
)

const (
	HealthOk      = "ok"
	HealthBehind  = "behind"
	HealthUnknown = "unknown"
)

const (
	defaultRequestTimeout   = 30 * time.Second
	defaultRetryCount       = 3
	defaultRetryWaitTime    = 10 * time.Second
	defaultRetryMaxWaitTime = 30 * time.Second
	defaultWithDebug        = false

	blockLimit = 50_0000
)

type solclient struct {
	grestyClient *gresty.Client
}

var (
	errHTTPError       = errors.New("aptos http error")
	errInvalidAddress  = errors.New("invalid address")
	errInvalidResponse = errors.New("invalid response")
)

func NewSolHttpClient(baseUrl string) (SolClient, error) {
	return NewSolHttpClientAll(baseUrl, defaultWithDebug)
}

func NewSolHttpClientAll(baseUrl string, withDebug bool) (SolClient, error) {
	grestyClient := gresty.New()
	grestyClient.SetBaseURL(baseUrl)
	grestyClient.SetTimeout(defaultRequestTimeout)
	grestyClient.SetRetryCount(defaultRetryCount)
	grestyClient.SetRetryWaitTime(defaultRetryWaitTime)
	grestyClient.SetRetryMaxWaitTime(defaultRetryMaxWaitTime)
	grestyClient.SetDebug(withDebug)

	// Retry Condition
	//grestyClient.AddRetryCondition(func(r *gresty.Response, err error) bool {
	//	return err != nil || r.StatusCode() >= 500
	//})

	grestyClient.OnBeforeRequest(func(c *gresty.Client, r *gresty.Request) error {
		log.Printf("Making request to %s (Attempt %d)", r.URL, r.Attempt)
		return nil
	})

	grestyClient.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		attempt := r.Request.Attempt
		method := r.Request.Method
		url := r.Request.URL
		log.Printf("Response received: Method=%s, URL=%s, Status=%d, Attempt=%d",
			method, url, statusCode, attempt)

		if statusCode >= 400 {
			if statusCode == 404 {
				return fmt.Errorf("%d resource not found %s %s: %w",
					statusCode, method, url, errHTTPError)
			}
			if statusCode >= 500 {
				return fmt.Errorf("%d server error %s %s: %w",
					statusCode, method, url, errHTTPError)
			}
			return fmt.Errorf("%d cannot %s %s: %w",
				statusCode, method, url, errHTTPError)
		}
		return nil
	})
	return &solclient{grestyClient: grestyClient}, nil
}

func (s *solclient) GetHealth() (string, error) {
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getHealth",
		"params":  []interface{}{},
	}

	response := &GetHealthResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(response).
		Post("/")
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return "", fmt.Errorf("failed to get health: %w", errHTTPError)
	}

	if response.Error != nil {
		if response.Error.Code == -32005 {
			return HealthBehind, nil
		}
		return HealthUnknown, fmt.Errorf("RPC error: code=%d, message=%s",
			response.Error.Code,
			response.Error.Message,
		)
	}

	if response.Result == "" {
		return HealthUnknown, fmt.Errorf("invalid response: empty result")
	}

	switch response.Result {
	case HealthOk, HealthBehind:
		return response.Result, nil
	default:
		return HealthUnknown, fmt.Errorf("unknown health status: %s", response.Result)
	}
}

func (s *solclient) CreateNonceAccount(req *CreateNonceAccountRequest) (*CreateNonceAccountResponse, error) {
	// 1. 生成新的 nonce account 密钥对
	//nonceAccount := solana.NewWallet().PrivateKey

	// 2. 获取最新的 blockhash
	//blockhash, err := s.GetLatestBlockhash()
	//if err != nil {
	//	return nil, fmt.Errorf("get blockhash failed: %w", err)
	//}

	// 3. 获取最小租金
	//rentRequest := map[string]interface{}{
	//	"jsonrpc": "2.0",
	//	"id":      1,
	//	"method":  "getMinimumBalanceForRentExemption",
	//	"params": []interface{}{
	//		system.NonceAccountSize,
	//	},
	//}
	//
	//var rentResponse struct {
	//	Result uint64 `json:"result"`
	//}
	//
	//resp, err := s.grestyClient.R().
	//	SetBody(rentRequest).
	//	SetResult(&rentResponse).
	//	Post("/")
	//if err != nil {
	//	return nil, fmt.Errorf("get rent failed: %w", err)
	//}

	//rent := rentResponse.Result

	// 4. 构建创建 nonce account 的指令
	//createNonceAccIx := system.NewCreateNonceAccountInstruction(
	//	rent,
	//	system.NonceAccountSize,
	//	req.Payer.PublicKey(),
	//	nonceAccount.PublicKey(),
	//	req.Authority,
	//).Build()

	// 5. 构建交易
	//tx, err := solana.NewTransaction(
	//	[]solana.Instruction{
	//		createNonceAccIx,
	//	},
	//	blockhash.Result.Value.Blockhash,
	//	solana.TransactionPayer(req.Payer.PublicKey()),
	//)

	// 6. 签名交易
	//_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
	//	if key.Equals(req.Payer.PublicKey()) {
	//		return &req.Payer
	//	}
	//	if key.Equals(nonceAccount.PublicKey()) {
	//		return &nonceAccount
	//	}
	//	return nil
	//})
	//if err != nil {
	//	return nil, fmt.Errorf("sign transaction failed: %w", err)
	//}

	// 7. 发送交易
	//txHash, err := tx.Hash()
	//if err != nil {
	//	return nil, fmt.Errorf("get transaction hash failed: %w", err)
	//}
	//
	//sendRequest := map[string]interface{}{
	//	"jsonrpc": "2.0",
	//	"id":      1,
	//	"method":  "sendTransaction",
	//	"params": []interface{}{
	//		tx.MarshalBase64(),
	//		map[string]interface{}{
	//			"encoding": "base64",
	//		},
	//	},
	//}
	//
	//var sendResponse struct {
	//	Result string `json:"result"`
	//}
	//
	//resp, err = s.grestyClient.R().
	//	SetBody(sendRequest).
	//	SetResult(&sendResponse).
	//	Post("/")
	//if err != nil {
	//	return nil, fmt.Errorf("send transaction failed: %w", err)
	//}
	//
	//// 8. 等待交易确认
	//time.Sleep(time.Second * 2)

	// 9. 获取 nonce 值
	//nonceInfo, err := s.GetNonceAccount(nonceAccount.PublicKey().String())
	//if err != nil {
	//	return nil, fmt.Errorf("get nonce failed: %w", err)
	//}

	//return &CreateNonceAccountResponse{
	//	NonceAccount: nonceAccount.PublicKey(),
	//	//Nonce:        nonceInfo.Nonce.String(),
	//	Signature: sendResponse.Result,
	//}, nil
	return nil, nil
}

func (s *solclient) GetAccountInfo(inputAddr string) (*AccountInfo, error) {
	dealAddr := strings.TrimSpace(inputAddr)
	if dealAddr == "" {
		return nil, fmt.Errorf("invalid input: empty address")
	}

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getAccountInfo",
		"params": []interface{}{
			dealAddr,
			map[string]string{
				"encoding": "base64",
			},
		},
	}
	response := &GetAccountInfoResponse{}
	resp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(response).
		Post("/")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get account info: %w", errHTTPError)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("RPC error: code=%d, message=%s",
			response.Error.Code,
			response.Error.Message,
		)
	}

	accountInfo := &response.Result.Value
	if accountInfo.Owner == "" {
		return nil, fmt.Errorf("invalid response: empty owner")
	}
	if len(accountInfo.Data) < 2 {
		return nil, fmt.Errorf("invalid response: missing data encoding")
	}
	if accountInfo.Data[1] != "base64" {
		return nil, fmt.Errorf("unexpected data encoding: %s", accountInfo.Data[1])
	}

	return accountInfo, nil
}

func (s *solclient) GetBalance(inputAddr string) (uint64, error) {
	dealAddr := strings.TrimSpace(inputAddr)
	if dealAddr == "" {
		return 0, fmt.Errorf("invalid input: empty address")
	}

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getBalance",
		"params": []interface{}{
			dealAddr,
		},
	}

	response := &GetBalanceResponse{}
	resp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(response).
		Post("/")
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return 0, fmt.Errorf("HTTP error: status=%d, body=%s",
			resp.StatusCode(),
			resp.String(),
		)
	}
	if response.Error != nil {
		return 0, fmt.Errorf("RPC error: code=%d, message=%s",
			response.Error.Code,
			response.Error.Message,
		)
	}
	if response.Result.Value == 0 {
		log.Printf("Warning: account balance is 0 for address: %s", dealAddr)
	}
	return response.Result.Value, nil
}

// GetBlockHeight get latest height
func (s *solclient) GetBlockHeight() (uint64, error) {
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getBlockHeight",
	}

	response := &BlockHeightResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(response).
		Post("/")

	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return 0, fmt.Errorf("failed to get block height: %w", errHTTPError)
	}
	if response.Error != nil {
		return 0, fmt.Errorf("RPC error: code=%d, message=%s",
			response.Error.Code,
			response.Error.Message,
		)
	}
	if response.Result == 0 {
		return 0, fmt.Errorf("invalid block height: got 0")
	}
	return response.Result, nil
}

// GetSlot get latest slot
func (s *solclient) GetSlot(commitment CommitmentType) (uint64, error) {
	if commitment == "" {
		return 0, fmt.Errorf("invalid input: empty commitment")
	}
	config := GetSlotRequest{
		Commitment: commitment,
	}

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getSlot",
		"params":  []interface{}{config},
	}

	response := &GetSlotResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(response).
		Post("/")

	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return 0, fmt.Errorf("failed to get slot: %w", errHTTPError)
	}

	if response.Error != nil {
		return 0, fmt.Errorf("RPC error: code=%d, message=%s",
			response.Error.Code,
			response.Error.Message,
		)
	}

	if response.Result == 0 {
		return 0, fmt.Errorf("invalid slot number: got 0")
	}

	return response.Result, nil
}

// GetBlocksWithLimit returns a list of confirmed blocks starting at the given slot with limit
func (s *solclient) GetBlocksWithLimit(startSlot uint64, limit uint64) ([]uint64, error) {
	if startSlot == 0 {
		return nil, fmt.Errorf("invalid input: start slot cannot be 0")
	}
	if limit == 0 {
		return nil, fmt.Errorf("invalid input: limit cannot be 0")
	}
	if limit > blockLimit {
		return nil, fmt.Errorf("limit must not exceed %d blocks", blockLimit)
	}

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getBlocksWithLimit",
		"params":  []uint64{startSlot, limit},
	}

	response := &GetBlocksWithLimitResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(response).
		Post("/")

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return nil, fmt.Errorf("failed to get blocks with limit: %w", errHTTPError)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("RPC error: code=%d, message=%s",
			response.Error.Code,
			response.Error.Message,
		)
	}

	if response.Result == nil {
		return []uint64{}, nil
	}

	if len(response.Result) == 0 {
		log.Printf("Warning: no blocks found for slot range %d to %d",
			startSlot, startSlot+limit-1)
	}

	if uint64(len(response.Result)) > limit {
		return nil, fmt.Errorf("received more blocks than requested limit: got %d, want <= %d",
			len(response.Result), limit)
	}

	return response.Result, nil
}

func (s *solclient) GetBlockBySlot(slot uint64, detailType TransactionDetailsType) (*BlockResult, error) {
	//s.grestyClient.SetTimeout(120 * time.Second)
	//defer s.grestyClient.SetTimeout(defaultRequestTimeout)

	config := GetBlockRequest{
		Commitment:                     Finalized,
		Encoding:                       "json",
		MaxSupportedTransactionVersion: 0,
		TransactionDetails:             string(detailType),
		Rewards:                        false,
	}

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getBlock",
		"params":  []interface{}{slot, config},
	}

	resp := &GetBlockResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(resp).
		Post("/")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return nil, fmt.Errorf("failed to get block: %w", errHTTPError)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("RPC error: (code: %d) %s", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result, nil
}

func (s *solclient) GetTransaction(signature string) (*TransactionResult, error) {
	signature = strings.TrimSpace(signature)
	if signature == "" {
		return nil, fmt.Errorf("invalid input: empty signature")
	}
	if len(signature) < 88 || len(signature) > 90 {
		return nil, fmt.Errorf("invalid signature length: expected 88-90 chars, got %d", len(signature))
	}
	config := map[string]interface{}{
		"encoding":                       "json",
		"commitment":                     Finalized,
		"maxSupportedTransactionVersion": 0,
	}

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getTransaction",
		"params":  []interface{}{signature, config},
	}

	response := &GetTransactionResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(response).
		Post("/")

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return nil, fmt.Errorf("failed to get transaction: %w", errHTTPError)
	}
	if response.Error != nil {
		if response.Error.Code == -32004 {
			return nil, fmt.Errorf("transaction not found: %s", signature)
		}
		return nil, fmt.Errorf("RPC error: code=%d, message=%s",
			response.Error.Code,
			response.Error.Message,
		)
	}
	if response.Result.Transaction.Signatures == nil {
		return nil, fmt.Errorf("invalid response: empty transaction data")
	}

	return &TransactionResult{
		Slot:        response.Result.Slot,
		Version:     response.Result.Version,
		BlockTime:   response.Result.BlockTime,
		Transaction: response.Result.Transaction,
		Meta:        response.Result.Meta,
	}, nil
}

func (s *solclient) GetTransactionRange(inputSignatureList []string) ([]*TransactionResult, error) {
	if len(inputSignatureList) == 0 {
		return nil, fmt.Errorf("empty signatures")
	}

	for i, sig := range inputSignatureList {
		inputSignatureList[i] = strings.TrimSpace(sig)
		if inputSignatureList[i] == "" {
			return nil, fmt.Errorf("invalid input: empty signature at index %d", i)
		}
		if len(inputSignatureList[i]) < 88 || len(inputSignatureList[i]) > 90 {
			return nil, fmt.Errorf("invalid signature length at index %d: expected 88-90 chars, got %d",
				i, len(inputSignatureList[i]))
		}
	}

	if len(inputSignatureList) == 1 {
		tx, err := s.GetTransaction(inputSignatureList[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get single transaction: %w", err)
		}
		return []*TransactionResult{tx}, nil
	}

	const (
		maxConcurrent   = 20
		requestInterval = 100 * time.Millisecond
		timeout         = 5 * time.Minute
		maxRetries      = 3
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultChannel := make(chan *TransactionResult, len(inputSignatureList))
	errorChannel := make(chan error, len(inputSignatureList))

	rateLimiter := time.NewTicker(requestInterval)
	defer rateLimiter.Stop()

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrent)

	for i, sig := range inputSignatureList {
		wg.Add(1)
		go func(index int, signature string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			var tx *TransactionResult
			var err error

			for retry := 0; retry < maxRetries; retry++ {
				select {
				case <-ctx.Done():
					errorChannel <- ctx.Err()
					return
				case <-rateLimiter.C:
					tx, err = s.GetTransaction(signature)
					if err == nil {
						resultChannel <- tx
						return
					}

					if retry < maxRetries-1 && strings.Contains(err.Error(), "request failed") {
						time.Sleep(time.Second * time.Duration(retry+1))
						continue
					}

					errorChannel <- fmt.Errorf("failed to get transaction %s: %w", signature, err)
					return
				}
			}

			if err != nil {
				errorChannel <- fmt.Errorf("max retries exceeded for %s: %w", signature, err)
			}
		}(i, sig)
	}

	wg.Wait()
	close(resultChannel)
	close(errorChannel)

	var errorList []string
	for err := range errorChannel {
		if err != nil {
			errorList = append(errorList, err.Error())
		}
	}
	if len(errorList) > 0 {
		return nil, fmt.Errorf("multiple errors occurred: %s", strings.Join(errorList, "; "))
	}

	validResults := make([]*TransactionResult, 0, len(resultChannel))
	for result := range resultChannel {
		if result != nil && result.Transaction.Signatures != nil {
			validResults = append(validResults, result)
		} else {
			log.Println("Skipping invalid transaction", "result", result)
		}
	}

	if len(validResults) == 0 {
		return nil, fmt.Errorf("no valid transactions found")
	}

	return validResults, nil
}

func (s *solclient) GetFeeForMessage(message string) (uint64, error) {
	if message == "" {
		return 0, fmt.Errorf("invalid input: empty message")
	}
	config := GetFeeForMessageRequest{
		Commitment: string(Finalized),
	}

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getFeeForMessage",
		"params":  []interface{}{message, config},
	}

	resp := &GetFeeForMessageResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(resp).
		Post("/")
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return 0, fmt.Errorf("failed to get fee for message: %w", errHTTPError)
	}
	if resp.Error != nil {
		return 0, fmt.Errorf("RPC error: code=%d, message=%s",
			resp.Error.Code,
			resp.Error.Message,
		)
	}
	if resp.Result.Value == nil {
		return 0, fmt.Errorf("invalid message or unable to estimate fee")
	}

	return *resp.Result.Value, nil
}

func (s *solclient) GetRecentPrioritizationFees() ([]PrioritizationFee, error) {
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getRecentPrioritizationFees",
		"params":  []interface{}{},
	}

	resp := &getRecentPrioritizationFeesResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(resp).
		Post("/")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return nil, fmt.Errorf("failed to get prioritization fees: %w", errHTTPError)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("RPC error: code=%d, message=%s",
			resp.Error.Code,
			resp.Error.Message,
		)
	}
	if resp.Result == nil {
		return []PrioritizationFee{}, nil
	}
	if len(resp.Result) == 0 {
		return []PrioritizationFee{}, nil
	}

	return resp.Result, nil
}

func GetSuggestedPriorityFee(fees []PrioritizationFee) uint64 {
	if len(fees) == 0 {
		return 0
	}

	priorityFees := make([]uint64, len(fees))
	for i, fee := range fees {
		priorityFees[i] = fee.PrioritizationFee
	}

	sort.Slice(priorityFees, func(i, j int) bool {
		return priorityFees[i] < priorityFees[j]
	})

	index := int(float64(len(priorityFees)) * 0.75)
	return priorityFees[index]
}

func (s *solclient) GetTxForAddress(
	address string,
	commitment CommitmentType,
	limit uint64,
	beforeSignature string,
	untilSignature string,
) ([]SignatureInfo, error) {
	if address == "" {
		return nil, fmt.Errorf("empty address")
	}

	config := &GetSignaturesRequest{
		Commitment: string(commitment),
		Limit:      limit,
		Before:     beforeSignature,
		Until:      untilSignature,
	}

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getSignaturesForAddress",
		"params":  []interface{}{address, config},
	}

	resp := &GetSignaturesResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(resp).
		Post("/")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return nil, fmt.Errorf("failed to get signatures: %w", errHTTPError)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("RPC error: code=%d, message=%s",
			resp.Error.Code,
			resp.Error.Message,
		)
	}

	if resp.Result == nil {
		return []SignatureInfo{}, nil
	}

	return resp.Result, nil
}

func (s *solclient) SimulateTransaction(
	signedTx string,
	config *SimulateRequest,
) (*SimulateResult, error) {
	if signedTx == "" {
		return nil, fmt.Errorf("invalid input: empty transaction")
	}
	if config == nil {
		config = &SimulateRequest{
			Commitment: string(Finalized),
			Encoding:   "base64",
		}
	}

	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "simulateTransaction",
		"params": []interface{}{
			signedTx,
			config,
		},
	}

	resp := &SimulateTransactionResponse{}
	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(resp).
		Post("/")
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return nil, fmt.Errorf("failed to simulate transaction: %w", errHTTPError)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("RPC error: code=%d, message=%s",
			resp.Error.Code,
			resp.Error.Message,
		)
	}
	if resp.Result.Err != nil {
		return nil, fmt.Errorf("simulation failed: %v", resp.Result.Err)
	}
	if resp.Result.UnitsConsumed == 0 && len(resp.Result.Logs) == 0 {
		return nil, fmt.Errorf("empty simulation result")
	}
	return &resp.Result, nil
}

func validateSimulateResponse(resp *SimulateTransactionResponse) error {
	if resp == nil {
		return fmt.Errorf("empty response")
	}

	if resp.Error != nil {
		return fmt.Errorf("RPC error: code=%d, message=%s",
			resp.Error.Code,
			resp.Error.Message,
		)
	}

	if resp.Jsonrpc != "2.0" {
		return fmt.Errorf("invalid jsonrpc version: %s", resp.Jsonrpc)
	}

	if resp.Result.Err != nil {
		return fmt.Errorf("simulation failed: %v", resp.Result.Err)
	}

	if resp.Result.UnitsConsumed == 0 {
		return fmt.Errorf("invalid units consumed: 0")
	}

	if len(resp.Result.Accounts) > 0 {
		for i, account := range resp.Result.Accounts {
			if account.Owner == "" {
				return fmt.Errorf("invalid account owner at index %d", i)
			}
		}
	}

	if resp.Result.ReturnData != nil {
		if resp.Result.ReturnData.ProgramId == "" {
			return fmt.Errorf("invalid return data: empty program id")
		}
	}

	if len(resp.Result.InnerInstructions) > 0 {
		for i, inner := range resp.Result.InnerInstructions {
			if len(inner.Instructions) == 0 {
				return fmt.Errorf("empty instructions for inner instruction at index %d", i)
			}
		}
	}

	return nil
}

func (s *solclient) SendTransaction(
	signedTx string,
	config *SendTransactionRequest,
) (string, error) {
	if signedTx == "" {
		return "", fmt.Errorf("invalid input: empty transaction")
	}
	if config == nil {
		config = &SendTransactionRequest{
			Commitment: string(Finalized),
			Encoding:   "base58",
		}
	}
	fmt.Println("3:", signedTx)
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "sendTransaction",
		"params": []interface{}{
			signedTx,
			config,
		},
	}

	resp := &SendTransactionResponse{}

	httpResp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(&resp).
		Post("/")
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	if httpResp.IsError() {
		return "", fmt.Errorf("failed to send transaction: %w", errHTTPError)
	}

	if resp.Error != nil {
		return "", fmt.Errorf("RPC error: code=%d, message=%s",
			resp.Error.Code,
			resp.Error.Message,
		)
	}

	if resp.Result == "" {
		return "", fmt.Errorf("empty transaction signature returned")
	}

	return resp.Result, nil
}

func (s *solclient) GetLatestBlockhash(commitmentType CommitmentType) (string, error) {
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getLatestBlockhash",
		"params": []interface{}{
			map[string]string{
				"commitment": string(commitmentType),
			},
		},
	}

	response := &GetLatestBlockhashResponse{}
	resp, err := s.grestyClient.R().
		SetBody(requestBody).
		SetResult(response).
		Post("/")
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return "", fmt.Errorf("failed to get latest blockhash: %w", errHTTPError)
	}

	if response.Error != nil {
		return "", fmt.Errorf("RPC error: code=%d, message=%s",
			response.Error.Code,
			response.Error.Message,
		)
	}

	blockhash := response.Result.Value.Blockhash
	if blockhash == "" {
		return "", fmt.Errorf("invalid blockhash response: empty blockhash")
	}

	return blockhash, nil
}

// GetAccountInfo retrieves account information for a given token account
func GetAccountInfo(sdkClient *rpc.Client, tokenAccount solana.PublicKey) (*rpc.GetAccountInfoResult, error) {
	accountInfo, err := sdkClient.GetAccountInfo(context.Background(), tokenAccount)
	if err != nil {
		log.Println("Failed to get account info", "err", err)
		return nil, err
	}
	return accountInfo, nil
}

// GetTokenSupply retrieves the token supply for a given mint public key
func GetTokenSupply(sdkClient *rpc.Client, mintPubkey solana.PublicKey) (*rpc.GetTokenSupplyResult, error) {
	tokenInfo, err := sdkClient.GetTokenSupply(context.Background(), mintPubkey, rpc.CommitmentFinalized)
	if err != nil {
		log.Println("Failed to get token supply", "err", err)
		return nil, err
	}
	return tokenInfo, nil
}
