package tron

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	defaultRequestTimeout = 10 * time.Second
	defaultRetryCount     = 3
)

// TronClient Define a Tron RPC client
type TronClient struct {
	rpc *resty.Client
}

// DialTronClient Initialize and return a TronClient instance
func DialTronClient(rpcURL, rpcUser, rpcPass string) *TronClient {
	client := resty.New()
	client.SetHeader(rpcUser, rpcPass)
	client.SetBaseURL(rpcURL)
	client.SetTimeout(defaultRequestTimeout)
	client.SetRetryCount(defaultRetryCount)

	return &TronClient{
		rpc: client,
	}
}

// JsonRpc Call JSON-RPC
func (client *TronClient) JsonRpc(method string, params interface{}, result interface{}) error {
	requestBody, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}
	_, err = client.rpc.R().SetBody(requestBody).SetResult(result).Post("/jsonrpc")
	return nil
}

// Solidity Call Solidity
func (client *TronClient) Solidity(method string, params interface{}, result interface{}) error {
	_, err := client.rpc.R().SetBody(params).SetResult(result).Post("/walletsolidity/" + method)
	return err
}

// Wallet Call Wallet
func (client *TronClient) Wallet(method string, params interface{}, result interface{}) error {
	_, err := client.rpc.R().SetBody(params).SetResult(result).Post("/wallet/" + method)
	return err
}

// GetBlockByNumber Obtain block information based on block number
func (client *TronClient) GetBlockByNumber(blockNumber int64) (*Block, error) {
	params := []interface{}{blockNumber, false}
	var response Response[Block]
	err := client.JsonRpc("eth_getBlockByNumber", params, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by number: %v", err)
	}
	return &response.Result, nil
}

// GetBlockByHash Obtain block information based on block hash
func (client *TronClient) GetBlockByHash(blockHash string) (*Block, error) {
	params := []interface{}{blockHash, false}
	var response Response[Block]
	err := client.JsonRpc("eth_getBlockByHash", params, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by hash: %v", err)

	}
	return &response.Result, nil
}

// GetAccount Get account information
func (client *TronClient) GetAccount(address string) (*Account, error) {
	params := map[string]interface{}{
		"address": address,
		"visible": true,
	}
	var accountInfo Account
	err := client.Solidity("getaccount", params, &accountInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %v", err)
	}
	return &accountInfo, nil
}

// GetTxByAddress Obtain transaction list based on address
func (client *TronClient) GetTxByAddress(address string) (interface{}, error) {
	params := map[string]interface{}{
		"address": address,
		"count":   10,
	}
	var txList []Transaction
	err := client.Solidity("gettransactionsbyaddress", params, &txList)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx by address: %v", err)
	}
	return txList, nil
}

// GetTransactionByID Obtain transaction information based on transaction hash
func (client *TronClient) GetTransactionByID(txHash string) (*Transaction, error) {
	params := map[string]interface{}{
		"value":   txHash,
		"visible": true,
	}
	var txInfo Transaction
	err := client.Wallet("gettransactionbyid", params, &txInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx by id: %v", err)
	}
	return &txInfo, nil
}

// CreateTRXTransaction Create TRX transaction to be signed
func (client *TronClient) CreateTRXTransaction(from, to string, amount int64) (*UnSignTransaction, error) {
	params := map[string]interface{}{
		"owner_address": from,
		"to_address":    to,
		"amount":        amount,
		"visible":       true,
	}
	var txInfo UnSignTransaction
	err := client.Wallet("createtransaction", params, &txInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to create TRX transaction: %v", err)
	}
	return &txInfo, nil
}

// CreateTRC20Transaction Create TRC20 token transaction to be signed
func (client *TronClient) CreateTRC20Transaction(from, to, contractAddress string, amount int64) (*UnSignTransaction, error) {
	// Convert the address to hexadecimal format
	toHex, err := Base58ToHex(to)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to address: %v", err)
	}

	// Encoding TRC20 token transfer function parameters
	toAddressHexPadded := PadLeftZero(toHex, 64)
	amountHex := fmt.Sprintf("%x", amount)
	amountHexPadded := PadLeftZero(amountHex, 64)
	parameter := toAddressHexPadded + amountHexPadded

	// Create request parameters
	params := map[string]interface{}{
		"owner_address":     from,
		"contract_address":  contractAddress,
		"function_selector": "transfer(address,uint256)",
		"parameter":         parameter,
		"fee_limit":         100000000, // 手续费上限
		"call_value":        amount,
		"visible":           true,
	}
	var txInfo UnSignTrc20Transaction
	err = client.Wallet("triggersmartcontract", params, &txInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to create TRC20 transaction: %v", err)
	}
	return &txInfo.Transaction, nil
}

// BroadcastTransaction Broadcast trading
func (client *TronClient) BroadcastTransaction(raw *SendTxReq) (*BroadcastReturns, error) {
	// Create request parameters
	params := map[string]interface{}{
		"raw_data":     raw.RawData,
		"raw_data_hex": raw.RawDataHex,
	}
	var bt BroadcastReturns
	// Call the broadcast transaction API to broadcast transactions
	err := client.Wallet("broadcasttransaction", params, &bt)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast transaction: %v", err)

	}
	return &bt, nil
}

type TRC20BalanceResponse struct {
	Data []struct {
		OwnerPermission       OwnerPermission     `json:"owner_permission"`
		AccountResource       AccountResource     `json:"account_resource"`
		ActivePermission      AccountResource     `json:"active_permission"`
		Address               string              `json:"address"`
		CreateTime            int64               `json:"create_time"`
		LatestOprationTime    int64               `json:"latest_opration_time"`
		FreeAssetNetUsageV2   FreeAssetNetUsageV2 `json:"free_asset_net_usageV2"`
		FreeNetUsage          int                 `json:"free_net_usage"`
		AssetV2               AssetV2             `json:"assetV2"`
		FrozenV2              FrozenV2            `json:"frozenV2"`
		Balance               int64               `json:"balance"`
		TRC20                 []map[string]string `json:"trc20"`
		LatestConsumeFreeTime int64               `json:"latest_consume_free_time"`
		NetWindowSize         int                 `json:"net_window_size"`
		NetWindowOptimized    bool                `json:"net_window_optimized"`
	} `json:"data"`
	Success bool `json:"success"`
	Meta    struct {
		At       int64 `json:"at"`
		PageSize int   `json:"page_size"`
	} `json:"meta"`
}

// GetTRC20Balance Query the TRC20 token balance at the specified address
func (client *TronClient) GetTRC20Balance(address, contractAddress string) (string, error) {
	var result TRC20BalanceResponse
	client.rpc.R().SetResult(&result).Get("/v1/accounts/" + address)
	// Check if the parsing is successful
	if !result.Success || len(result.Data) == 0 {
		return "0", fmt.Errorf("invalid response or empty data")
	}
	// Traverse the trc20 list to find matching contractAddress
	for _, trc20 := range result.Data[0].TRC20 {
		if balance, exists := trc20[contractAddress]; exists {
			return balance, nil
		}
	}
	return "0", fmt.Errorf("contract address not found")
}

// GetChainParameters Get chain parameters
func (client *TronClient) GetChainParameters() (*ChainParameters, error) {
	var result ChainParameters
	err := client.Wallet("getChainParameters", nil, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain parameters: %v", err)
	}
	return &result, nil
}
