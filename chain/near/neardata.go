package near

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Transaction 定义交易数据结构

// Transaction represents a single transaction.
type Transaction2 struct {
	ID                           string      `json:"id"`
	SignerAccountID              string      `json:"signer_account_id"`
	ReceiverAccountID            string      `json:"receiver_account_id"`
	TransactionHash              string      `json:"transaction_hash"`
	IncludedInBlockHash          string      `json:"included_in_block_hash"`
	BlockTimestamp               string      `json:"block_timestamp"`
	ReceiptConversionTokensBurnt string      `json:"receipt_conversion_tokens_burnt"`
	Block                        Block       `json:"block"`
	Actions                      []Action2   `json:"actions"`
	ActionsAgg                   ActionsAgg  `json:"actions_agg"`
	Outcomes                     Outcomes    `json:"outcomes"`
	OutcomesAgg                  OutcomesAgg `json:"outcomes_agg"`
}

// Block represents the block information.
type Block struct {
	BlockHeight int `json:"block_height"`
}

// Action represents a single action within a transaction.
type Action2 struct {
	Action  string      `json:"action"`
	Method  interface{} `json:"method"`
	Deposit float64     `json:"deposit"`
	Fee     float64     `json:"fee"`
	Args    interface{} `json:"args"`
}

// ActionsAgg represents aggregated actions information.
type ActionsAgg struct {
	Deposit float64 `json:"deposit"`
}

// Outcomes represents the outcomes of a transaction.
type Outcomes struct {
	Status bool `json:"status"`
}

// OutcomesAgg represents aggregated outcomes information.
type OutcomesAgg struct {
	TransactionFee float64 `json:"transaction_fee"`
}

// TransactionsResponse represents the response containing multiple transactions.
type TransactionsResponse struct {
	Transactions []Transaction2 `json:"txns"`
}

// GetAccountTransactions 完整请求封装
func getAccountTransactions(accountID string) ([]Transaction2, error) {
	// 构建请求
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.nearblocks.io/v1/account/%s/txns-only", accountID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("构建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "NEAR-Client/1.0")

	// 发送请求（带超时控制）
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回异常状态码: %d", resp.StatusCode)
	}

	// 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		Transactions []Transaction2 `json:"txns"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}

	return result.Transactions, nil
}
