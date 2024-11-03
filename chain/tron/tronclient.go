package tron

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"io/ioutil"
	"net/http"
)

// TronClient 封装了与 TRON 区块链交互的 JSON-RPC 客户端
type TronClient struct {
	rpcURL string
}

// DialTronClient 初始化并返回一个 TronClient 实例
func DialTronClient(rpcURL string) *TronClient {
	return &TronClient{rpcURL: rpcURL}
}

// Post 使用 POST 方法发送 JSON-RPC 请求
func (client *TronClient) Post(method string, params interface{}) (json.RawMessage, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}
	resp, err := http.Post(client.rpcURL+"/jsonrpc", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("POST request failed: %v", err)
	}
	defer resp.Body.Close()
	return client.processResponse(resp)
}

// PostSolidity 使用 POST 方法发送 JSON-RPC 请求
func (client *TronClient) PostSolidity(method string, params map[string]interface{}) (json.RawMessage, error) {
	requestBody, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}
	resp, err := http.Post(client.rpcURL+"/walletsolidity/"+method, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("POST request failed: %v", err)
	}
	defer resp.Body.Close()
	return client.processResponse(resp)
}

// PostWallet 使用 POST 方法发送 JSON-RPC 请求
func (client *TronClient) PostWallet(method string, params map[string]interface{}) (json.RawMessage, error) {
	requestBody, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}
	resp, err := http.Post(client.rpcURL+"/wallet/"+method, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("POST request failed: %v", err)
	}
	defer resp.Body.Close()
	return client.processResponse(resp)
}

//// Get 使用 GET 方法发送 JSON-RPC 请求，参数以 JSON 格式放入请求体中
//func (client *TronClient) Get(method string, params interface{}) (json.RawMessage, error) {
//	// 构造 JSON-RPC 请求体
//	requestBody, err := json.Marshal(map[string]interface{}{
//		"jsonrpc": "2.0",
//		"method":  method,
//		"params":  params,
//		"id":      1,
//	})
//	if err != nil {
//		return nil, fmt.Errorf("failed to marshal request: %v", err)
//	}
//
//	// 创建 GET 请求，并将 JSON 放入请求体
//	req, err := http.NewRequest("GET", client.rpcURL+"/jsonrpc", bytes.NewBuffer(requestBody))
//	if err != nil {
//		return nil, fmt.Errorf("failed to create GET request: %v", err)
//	}
//	req.Header.Set("Content-Type", "application/json")
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		return nil, fmt.Errorf("GET request failed: %v", err)
//	}
//	defer resp.Body.Close()
//
//	return client.processResponse(resp)
//}

// GetSolidity 使用 GET 方法发送 JSON-RPC 请求，参数以 JSON 格式放入请求体中
func (client *TronClient) GetSolidity(method string, params map[string]interface{}) (json.RawMessage, error) {
	// 构造 JSON-RPC 请求体
	requestBody, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}
	// 创建 GET 请求，并将 JSON 放入请求体
	req, err := http.NewRequest("GET", client.rpcURL+"/walletsolidity/"+method, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET request failed: %v", err)
	}
	defer resp.Body.Close()

	return client.processResponse(resp)
}

// processResponse 处理通用的 HTTP 响应解析逻辑
func (client *TronClient) processResponse(resp *http.Response) (json.RawMessage, error) {
	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}
	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	// 检查是否是HTML响应（通常代表错误页面）
	if bytes.HasPrefix(body, []byte("<html>")) {
		return nil, errors.New("received unexpected HTML response, please check the URL or method")
	}
	// 解析 JSON 响应
	var rpcResponse json.RawMessage
	if err := json.Unmarshal(body, &rpcResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %v", err)
	}
	return rpcResponse, nil
}

// GetBlockByNumber 根据区块号获取区块信息
func (client *TronClient) GetBlockByNumber(blockNumber int64) (*Block, error) {
	params := []interface{}{blockNumber, false}
	result, err := client.Post("eth_getBlockByNumber", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by number: %v", err)
	}

	var blockInfo Block
	if err := json.Unmarshal(result, &blockInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block info: %v", err)
	}
	return &blockInfo, nil
}

// GetBlockByHash 根据区块哈希获取区块信息
func (client *TronClient) GetBlockByHash(blockHash string) (*Block, error) {
	params := []interface{}{blockHash, false}
	result, err := client.Post("eth_getBlockByHash", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by hash: %v", err)

	}
	var blockInfo Block
	if err := json.Unmarshal(result, &blockInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block info: %v", err)
	}
	return &blockInfo, nil
}

// GetAccount 根据地址获取账户信息
func (client *TronClient) GetAccount(address string) (*Account, error) {
	params := map[string]interface{}{
		"address": address,
		"visible": true,
	}
	result, err := client.PostSolidity("getaccount", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %v", err)
	}
	var accountInfo Account
	if err := json.Unmarshal(result, &accountInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account info: %v", err)
	}
	return &accountInfo, nil
}

// GetTxByAddress 根据地址获取交易列表
func (client *TronClient) GetTxByAddress(address string) (interface{}, error) {
	params := map[string]interface{}{
		"address": address,
		"count":   10,
	}
	result, err := client.PostSolidity("gettransactionsbyaddress", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx by address: %v", err)
	}
	var txList interface{}
	if err := json.Unmarshal(result, &txList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tx list: %v", err)
	}
	return txList, nil
}

// GetTransactionByID 根据交易哈希获取交易信息
func (client *TronClient) GetTransactionByID(txHash string) (*Transaction, error) {
	params := map[string]interface{}{
		"hash": txHash,
	}
	result, err := client.PostSolidity("gettransactionbyid", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx by id: %v", err)
	}
	var txInfo Transaction
	if err := json.Unmarshal(result, &txInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tx info: %v", err)
	}
	return &txInfo, nil
}

// CreateTRXTransaction 创建待签名的 TRX 交易
func (client *TronClient) CreateTRXTransaction(from, to string, amount int64) (*UnSignTransaction, error) {
	params := map[string]interface{}{
		"owner_address": from,
		"to_address":    to,
		"amount":        amount,
	}

	// 调用 createtransaction API 创建 TRX 交易
	result, err := client.PostWallet("createtransaction", params)
	if err != nil {
		return nil, fmt.Errorf("failed to create TRX transaction: %v", err)
	}

	var txInfo UnSignTransaction
	if err := json.Unmarshal(result, &txInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TRX transaction info: %v", err)
	}
	return &txInfo, nil
}

// CreateTRC20Transaction 创建待签名的 TRC20 代币交易
func (client *TronClient) CreateTRC20Transaction(from, to, contractAddress string, amount int64) (*UnSignTransaction, error) {
	// 将地址转换为十六进制格式
	toHex, err := Base58ToHex(to)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to address: %v", err)
	}

	// 编码 TRC20 代币 transfer 函数参数
	toAddressHexPadded := PadLeftZero(toHex, 64)
	amountHex := fmt.Sprintf("%x", amount)
	amountHexPadded := PadLeftZero(amountHex, 64)
	parameter := toAddressHexPadded + amountHexPadded

	// 创建请求参数
	params := map[string]interface{}{
		"owner_address":     from,
		"contract_address":  contractAddress,
		"function_selector": "transfer(address,uint256)",
		"parameter":         parameter,
		"fee_limit":         100000000, // 手续费上限
	}

	// 调用 triggersmartcontract API 创建代币交易
	result, err := client.PostWallet("triggersmartcontract", params)
	if err != nil {
		return nil, fmt.Errorf("failed to create TRC20 transaction: %v", err)
	}

	var txInfo UnSignTransaction
	if err := json.Unmarshal(result, &txInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TRC20 transaction info: %v", err)
	}
	return &txInfo, nil
}

// BroadcastTransaction 广播交易
func (client *TronClient) BroadcastTransaction(raw string) (*BroadcastReturns, error) {
	jsonBytes, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		log.Error("decode string fail", "err", err)
		return nil, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		log.Error("parse json fail", "err", err)
		return nil, err
	}
	// 调用 broadcast transaction API 广播交易
	result, err := client.PostWallet("broadcasttransaction", data)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast transaction: %v", err)

	}
	var bt BroadcastReturns
	if err := json.Unmarshal(result, &bt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction info: %v", err)
	}
	return &bt, nil
}
