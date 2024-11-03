package tron

type Block struct {
	BaseFeePerGas    string        `json:"baseFeePerGas"`
	Difficulty       string        `json:"difficulty"`
	ExtraData        string        `json:"extraData"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Hash             string        `json:"hash"`
	LogsBloom        string        `json:"logsBloom"`
	Miner            string        `json:"miner"`
	MixHash          string        `json:"mixHash"`
	Nonce            string        `json:"nonce"`
	Number           string        `json:"number"`
	ParentHash       string        `json:"parentHash"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	Size             string        `json:"size"`
	StateRoot        string        `json:"stateRoot"`
	Timestamp        string        `json:"timestamp"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	Transactions     []string      `json:"transactions"`
	TransactionsRoot string        `json:"transactionsRoot"`
	Uncles           []interface{} `json:"uncles"`
}

type Account struct {
	Address             string                `json:"address"`
	Balance             int                   `json:"balance"`
	CreateTime          int64                 `json:"create_time"`
	LatestConsumeTime   int64                 `json:"latest_consume_time"`
	NetWindowSize       int                   `json:"net_window_size"`
	NetWindowOptimized  bool                  `json:"net_window_optimized"`
	AccountResource     AccountResource       `json:"account_resource"`
	OwnerPermission     OwnerPermission       `json:"owner_permission"`
	ActivePermission    []ActivePermission    `json:"active_permission"`
	FrozenV2            []FrozenV2            `json:"frozenV2"`
	AssetV2             []AssetV2             `json:"assetV2"`
	FreeAssetNetUsageV2 []FreeAssetNetUsageV2 `json:"free_asset_net_usageV2"`
	AssetOptimized      bool                  `json:"asset_optimized"`
}

type AccountResource struct {
	LatestConsumeTimeForEnergy                int64 `json:"latest_consume_time_for_energy"`
	EnergyWindowSize                          int   `json:"energy_window_size"`
	AcquiredDelegatedFrozenV2BalanceForEnergy int   `json:"acquired_delegated_frozenV2_balance_for_energy"`
	EnergyWindowOptimized                     bool  `json:"energy_window_optimized"`
}

type Keys struct {
	Address string `json:"address"`
	Weight  int    `json:"weight"`
}

type OwnerPermission struct {
	PermissionName string `json:"permission_name"`
	Threshold      int    `json:"threshold"`
	Keys           []Keys `json:"keys"`
}

type ActivePermission struct {
	Type           string `json:"type"`
	ID             int    `json:"id"`
	PermissionName string `json:"permission_name"`
	Threshold      int    `json:"threshold"`
	Operations     string `json:"operations"`
	Keys           []Keys `json:"keys"`
}

type FrozenV2 struct {
	Type string `json:"type,omitempty"`
}

type AssetV2 struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

type FreeAssetNetUsageV2 struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}
