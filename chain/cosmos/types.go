package cosmos

type NativeBalanceResponse struct {
	Code     string `json:"code"`
	Msg      string `json:"msg"`
	Response struct {
		Address                string `json:"address"`
		AvailableBalance       string `json:"availableBalance"`
		Delegated              string `json:"delegated"`
		DelegatedReward        string `json:"delegatedReward"`
		RewardRecipientAddress string `json:"rewardRecipientAddress"`
		Unbonding              string `json:"unbonding"`
		Symbol                 string `json:"symbol"`
		Commission             string `json:"commission"`
		Incentive              string `json:"incentive"`
		EthereumCoChainAddress string `json:"ethereumCoChainAddress"`
	} `json:"data"`
}

type NormalTransactionResponse struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data TransactionData `json:"data"`
}

type TransactionData struct {
	Page            string              `json:"page"`
	Limit           string              `json:"limit"`
	TotalPage       string              `json:"totalPage"`
	TransactionList []NormalTransaction `json:"transactionList"`
}

type NormalTransaction struct {
	Symbol          string   `json:"symbol"`
	TxId            string   `json:"txId"`
	BlockHash       string   `json:"blockHash"`
	Height          string   `json:"height"`
	TransactionTime string   `json:"transactionTime"`
	From            []string `json:"from"`
	To              []string `json:"to"`
	TxFee           string   `json:"txFee"`
	GasLimit        string   `json:"gasLimit"`
	GasUsed         string   `json:"gasUsed"`
	Type            []string `json:"type"`
	Value           string   `json:"value"`
	State           string   `json:"state"`
}

type TxStructure struct {
	ChainId         string `json:"chainId"`
	FromAddress     string `json:"from_address"`
	ToAddress       string `json:"to_address"`
	Amount          int64  `json:"amount"`
	ContractAddress string `json:"contract_address"`
	GasLimit        uint64 `json:"gas_limit"`
	FeeAmount       int64  `json:"free_amount"`
	Memo            string `json:"memo"`
	Decimal         int    `json:"decimal"`
	Sequence        uint64 `json:"sequence"`
	AccountNumber   uint64 `json:"account_number"`
	PubKey          string `json:"pub_key"`
}

type NativeBlockResponse struct {
	Code     string      `json:"code"`
	Msg      string      `json:"msg"`
	Response []BlockData `json:"data"`
}

type BlockData struct {
	ChainFullName  string `json:"chainFullName"`
	ChainShortName string `json:"chainShortName"`
	Hash           string `json:"hash"`
	Height         string `json:"height"`
	Validator      string `json:"validator"`
	BlockTime      string `json:"blockTime"`
	TxnCount       string `json:"txnCount"`
	Round          string `json:"round"`
	MineReward     string `json:"mineReward"`
	GasUsed        string `json:"gasUsed"`
	GasLimit       string `json:"gasLimit"`
	TotalFee       string `json:"totalFee"`
}
