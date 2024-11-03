package tron

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

// Get 使用 GET 方法发送 JSON-RPC 请求，参数以 JSON 格式放入请求体中
func (client *TronClient) Get(method string, params interface{}) (json.RawMessage, error) {
	// 构造 JSON-RPC 请求体
	requestBody, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// 创建 GET 请求，并将 JSON 放入请求体
	req, err := http.NewRequest("GET", client.rpcURL+"/jsonrpc", bytes.NewBuffer(requestBody))
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

// GetTransactionByHash 根据交易哈希获取交易信息
func (client *TronClient) GetTransactionByHash(txHash string) (interface{}, error) {
	params := []interface{}{txHash}
	result, err := client.Get("eth_getTransactionByHash", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by hash: %v", err)
	}

	var transactionInfo interface{}
	if err := json.Unmarshal(result, &transactionInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction info: %v", err)
	}
	return transactionInfo, nil
}

// GetBalance 获取指定地址的余额
func (client *TronClient) GetBalance(address string) (interface{}, error) {
	params := []interface{}{address, "latest"}
	result, err := client.Get("eth_getBalance", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %v", err)
	}

	var balance interface{}
	if err := json.Unmarshal(result, &balance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal balance: %v", err)
	}
	return balance, nil
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
