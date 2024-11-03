package ethereum

type TxStructure struct {
	ChainId         string `json:"chain_id"`
	Nonce           uint64 `json:"nonce"`
	GasPrice        string `json:"gas_price"`
	GasTipCap       string `json:"gas_tip_cap"`
	GasFeeCap       string `json:"gas_fee_cap"`
	Gas             uint64 `json:"gas"`
	ContractAddress string `json:"contract_address"`
	FromAddress     string `json:"from_address"`
	ToAddress       string `json:"to_address"`
	TokenId         string `json:"token_id"`
	Value           string `json:"value"`
}
