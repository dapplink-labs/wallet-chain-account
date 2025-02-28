package oasis

import "time"

type BlockInfo struct {
	TotalCount   int `json:"total_count"`
	Transactions []struct {
		Block int `json:"block"`
		Body  struct {
			Amount string `json:"amount"`
			To     string `json:"to"`
		} `json:"body"`
		Fee      string `json:"fee"`
		GasLimit string `json:"gas_limit"`
		Hash     string `json:"hash"`
		Index    int    `json:"index"`
		Method   string `json:"method"`
		Nonce    int    `json:"nonce"`
		Sender   string `json:"sender"`
		Success  bool   `json:"success"`
	} `json:"transactions"`
}

type LatestBlock struct {
	LatestBlock int64 `json:"latest_block"`
}

type TxInfo struct {
	IsTotalCountClipped bool `json:"is_total_count_clipped"`
	TotalCount          int  `json:"total_count"`
	Transactions        []struct {
		Block int `json:"block"`
		Body  struct {
			Amount string `json:"amount"`
			To     string `json:"to"`
		} `json:"body"`
		Fee       string    `json:"fee"`
		GasLimit  string    `json:"gas_limit"`
		Hash      string    `json:"hash"`
		Index     int       `json:"index"`
		Method    string    `json:"method"`
		Nonce     int       `json:"nonce"`
		Sender    string    `json:"sender"`
		Success   bool      `json:"success"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"transactions"`
}
