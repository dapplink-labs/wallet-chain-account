package aptos

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	gresty "github.com/go-resty/resty/v2"
)

type Client interface {
	GetNodeInfo() (*NodeInfo, error)

	GetAccount(inputAddr string) (*AccountResponse, error)
	GetGasPrice() (*EstimateGasPriceResponse, error)
	//SubmitTransaction(req *SubmitTransactionRequest) (*SubmitTransactionResponse, error)

	GetBlockByHeight(height uint64) (*BlockResponse, error)

	GetTransactionByHash(txHash string) (*TransactionResponse, error)
	GetTransactionByAddress(address string) (*[]TransactionResponse, error)
	GetTransactionByVersion(version string) (*TransactionResponse, error)
	GetTransactionByVersionRange(startVersion, endVersion uint64) ([]TransactionResponse, error)
}

var (
	errHTTPError       = errors.New("aptos http error")
	errInvalidAddress  = errors.New("invalid address")
	errInvalidResponse = errors.New("invalid response")
)

const (
	defaultRequestTimeout     = 10 * time.Second
	defaultRangRequestTimeout = 20 * time.Second
	defaultRetryCount         = 3
	defaultWithDebug          = false
	apikeyHeader              = "api-key"

	baseAPIPath = "/v1"

	// API Url
	pathNodeInfo = baseAPIPath + "/"

	pathGetSequence = baseAPIPath + "/accounts/%s"
	pathGasPrice    = baseAPIPath + "/estimate_gas_price"

	pathTransactions = baseAPIPath + "/transactions"
	pathTxByAddr     = baseAPIPath + "/accounts/%s/transactions"
	pathTxByHash     = baseAPIPath + "/transactions/by_hash/%s"
	pathTxByVersion  = baseAPIPath + "/transactions/by_version/%s"

	pathBlockByHeight = baseAPIPath + "/blocks/by_height/%s"
)

type RestyClient struct {
	client *gresty.Client
}

func NewAptosHttpClient(baseUrl, apiKey string) (*RestyClient, error) {
	return NewAptosHttpClientAll(baseUrl, apiKey, defaultWithDebug)
}

func NewAptosHttpClientAll(baseUrl, apiKey string, withDebug bool) (*RestyClient, error) {
	client := gresty.New()
	client.SetBaseURL(baseUrl)
	client.SetTimeout(defaultRequestTimeout)
	client.SetRetryCount(defaultRetryCount)
	client.SetDebug(withDebug)
	if apiKey != "" {
		client.SetHeader(apikeyHeader, apiKey)
	}

	client.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			url := r.Request.URL
			return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errHTTPError)
		}
		return nil
	})
	return &RestyClient{
		client: client,
	}, nil
}

func (c *RestyClient) GetNodeInfo() (*NodeInfo, error) {
	response := &NodeInfo{}
	resp, err := c.client.R().
		SetResult(response).
		Get(pathNodeInfo)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get node info: %w", errHTTPError)
	}

	return response, nil
}

// GetAccount info
func (c *RestyClient) GetAccount(inputAddr string) (*AccountResponse, error) {
	if !IsValidAddress(inputAddr) {
		return nil, fmt.Errorf("invalid address %s: %w", inputAddr, errInvalidAddress)
	}
	dealAddr := strings.TrimSpace(inputAddr)

	account := &AccountResponse{}
	resp, err := c.client.R().
		SetResult(account).
		Get(fmt.Sprintf(pathGetSequence, dealAddr))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get sequence number: %w", errHTTPError)
	}
	return account, nil
}

// GetGasPrice Get estimate gas price
func (c *RestyClient) GetGasPrice() (*EstimateGasPriceResponse, error) {
	gasPrice := &EstimateGasPriceResponse{}
	resp, err := c.client.R().
		SetResult(gasPrice).
		Get(fmt.Sprintf(pathGasPrice))

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to get gas price: %w", errHTTPError)
	}

	return gasPrice, nil
}

//func (c *RestyClient) SubmitTransaction(req *SubmitTransactionRequest) (*SubmitTransactionResponse, error) {
//	// check req
//	if err := ValidateSubmitTransaction(req); err != nil {
//		return nil, fmt.Errorf("invalid request: %w", err)
//	}
//
//	response := &SubmitTransactionResponse{}
//	resp, err := c.client.R().
//		SetBody(req).
//		SetResult(response).
//		Post(pathTransactions)
//
//	if err != nil {
//		return nil, fmt.Errorf("request failed: %w", err)
//	}
//
//	if resp.IsError() {
//		return nil, fmt.Errorf("failed to broadcast transaction: %w", errHTTPError)
//	}
//
//	return response, nil
//}

func (c *RestyClient) GetBlockByHeight(height uint64) (*BlockResponse, error) {
	if height < 0 {
		return nil, fmt.Errorf("invalid block height")
	}

	path := fmt.Sprintf(pathBlockByHeight, fmt.Sprint(height))

	response := &BlockResponse{}
	resp, err := c.client.R().
		SetResult(response).
		Get(path)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get block: %w", errHTTPError)
	}

	return response, nil
}

func (c *RestyClient) GetTransactionByAddress(inputAddr string) (*[]TransactionResponse, error) {
	if !IsValidAddress(inputAddr) {
		return nil, fmt.Errorf("invalid address %s: %w", inputAddr, errInvalidAddress)
	}
	path := fmt.Sprintf(pathTxByAddr, inputAddr)

	var response []TransactionResponse
	resp, err := c.client.R().
		SetResult(&response).
		Get(path)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get transactions: %w", errHTTPError)
	}
	return &response, nil
}

func (c *RestyClient) GetTransactionByHash(txHash string) (*TransactionResponse, error) {
	if txHash == "" {
		return nil, fmt.Errorf("transaction hash cannot be empty")
	}
	path := fmt.Sprintf(pathTxByHash, txHash)

	response := &TransactionResponse{}
	resp, err := c.client.R().
		SetResult(response).
		Get(path)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get transaction: %w", errHTTPError)
	}

	return response, nil
}

func (c *RestyClient) GetTransactionByVersion(version string) (*TransactionResponse, error) {
	if version == "" {
		return nil, fmt.Errorf("transaction version cannot be empty")
	}
	path := fmt.Sprintf(pathTxByVersion, version)

	response := &TransactionResponse{}
	resp, err := c.client.R().
		SetResult(response).
		Get(path)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get transaction by version: %w", errHTTPError)
	}

	return response, nil
}

func (c *RestyClient) GetTransactionByVersionRange(startVersion, endVersion uint64) ([]TransactionResponse, error) {
	if startVersion > endVersion {
		return nil, fmt.Errorf("start version (%d) cannot be greater than end version (%d)", startVersion, endVersion)
	}
	// Handle single version case
	if startVersion == endVersion {
		tx, err := c.GetTransactionByVersion(fmt.Sprint(startVersion))
		if err != nil {
			return nil, err
		}
		return []TransactionResponse{*tx}, nil
	}

	// Calculate total transactions to fetch
	count := endVersion - startVersion + 1
	transactions := make([]TransactionResponse, count)
	var wg sync.WaitGroup

	// Use smaller batch size for concurrent requests
	const groupSize = 20
	numGroups := (int(count)-1)/groupSize + 1
	errChan := make(chan error, numGroups)

	ctx, cancel := context.WithTimeout(context.Background(), defaultRangRequestTimeout)
	defer cancel()

	rateLimiter := time.NewTicker(100 * time.Millisecond)
	defer rateLimiter.Stop()

	for i := 0; i < int(count); i += groupSize {
		start := i
		end := i + groupSize - 1
		if end >= int(count) {
			end = int(count) - 1
		}
		wg.Add(1)

		go func(start, end int) {
			defer wg.Done()

			for j := start; j <= end; j++ {
				select {
				case <-ctx.Done():
					errChan <- ctx.Err()
					return
				case <-rateLimiter.C:
					version := startVersion + uint64(j)
					tx, err := c.GetTransactionByVersion(fmt.Sprint(version))
					if err != nil {
						errChan <- fmt.Errorf("failed to get transaction at version %d: %w", version, err)
						return
					}
					transactions[j] = *tx
				}
			}
		}(start, end)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return transactions, nil
}

func IsValidAddress(inputAddr string) bool {
	if len(inputAddr) == 0 {
		return false
	}
	// Prefix 0x
	if !strings.HasPrefix(inputAddr, "0x") {
		return false
	}
	// Trim Prefix 0x
	trimPrefix0xAddr := strings.TrimPrefix(inputAddr, "0x")

	// white space
	trimmedAddr := strings.TrimSpace(trimPrefix0xAddr)
	if len(trimmedAddr) == 0 {
		return false
	}

	isAllZeros := true
	for _, c := range trimPrefix0xAddr {
		if c != '0' {
			isAllZeros = false
			break
		}
	}
	if isAllZeros {
		return false
	}

	// check all hex str
	_, err := hex.DecodeString(trimPrefix0xAddr)
	return err == nil
}

//func ValidateSubmitTransaction(req *SubmitTransactionRequest) error {
//	if req == nil {
//		return errors.New("request cannot be nil")
//	}
//
//	if !IsValidAddress(req.Sender) {
//		return fmt.Errorf("invalid sender address: %s", req.Sender)
//	}
//
//	// require req
//	if req.SequenceNumber == 0 {
//		return errors.New("sequence number is required")
//	}
//	if req.MaxGasAmount == 0 {
//		return errors.New("max gas amount is required")
//	}
//	if req.GasUnitPrice == 0 {
//		return errors.New("gas unit price is required")
//	}
//	if req.ExpirationTimestampSecs == 0 {
//		return errors.New("expiration timestamp is required")
//	}
//
//	return nil
//}
