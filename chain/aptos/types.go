package aptos

type AccountResponse struct {
	SequenceNumber    uint64 `json:"sequence_number,string"`
	AuthenticationKey string `json:"authentication_key"`
}

type EstimateGasPriceResponse struct {
	DeprioritizedGasEstimate uint64 `json:"deprioritized_gas_estimate"`
	GasEstimate              uint64 `json:"gas_estimate"`
	PrioritizedGasEstimate   uint64 `json:"prioritized_gas_estimate"`
}

type TransferRequest struct {
	FromAddress   string
	FromPublicKey string
	ToAddress     string
	Amount        uint64
	// TODO Need to support more currencies
	// CoinType    string
}

type TransactionRequest struct {
	Sender         string         `json:"Sender"`
	SequenceNumber uint64         `json:"SequenceNumber"`
	Payload        PayloadWrapper `json:"Payload"`

	MaxGasAmount               uint64 `json:"MaxGasAmount"`
	GasUnitPrice               uint64 `json:"GasUnitPrice"`
	ExpirationTimestampSeconds uint64 `json:"ExpirationTimestampSeconds"`

	ChainId uint8 `json:"ChainId"`
}

type PayloadWrapper struct {
	Payload PayloadFunction `json:"Payload"`
}

type PayloadFunction struct {
	Module   ModuleInfo `json:"Module"`
	Function string     `json:"Function"`
	ArgTypes []string   `json:"ArgTypes"`
	Args     []string   `json:"Args"`
}

type ModuleInfo struct {
	Address string `json:"Address"`
	Name    string `json:"Name"`
}

type TransactionResponse struct {
	Version             uint64    `json:"version,string"`
	Hash                string    `json:"hash"`
	StateChangeHash     string    `json:"state_change_hash"`
	EventRootHash       string    `json:"event_root_hash"`
	StateCheckpointHash *string   `json:"state_checkpoint_hash"`
	GasUsed             uint64    `json:"gas_used,string"`
	Success             bool      `json:"success"`
	VMStatus            string    `json:"vm_status"`
	AccumulatorRootHash string    `json:"accumulator_root_hash"`
	Changes             []Change  `json:"changes"`
	Sender              string    `json:"sender"`
	SequenceNumber      uint64    `json:"sequence_number,string"`
	MaxGasAmount        uint64    `json:"max_gas_amount,string"`
	GasUnitPrice        uint64    `json:"gas_unit_price,string"`
	ExpirationTimestamp uint64    `json:"expiration_timestamp_secs,string"`
	Payload             Payload   `json:"payload"`
	Signature           Signature `json:"signature"`
	Events              []Event   `json:"events"`
	Timestamp           uint64    `json:"timestamp,string"`
	Type                string    `json:"type"`
}

type Change struct {
	Address      string     `json:"address"`
	StateKeyHash string     `json:"state_key_hash"`
	Data         ChangeData `json:"data"`
	Type         string     `json:"type"`
}

type ChangeData struct {
	Type string        `json:"type"`
	Data MinerDataInfo `json:"data"`
}

type MinerDataInfo struct {
	Events       EventInfo `json:"events"`
	Genesis      string    `json:"genesis"`
	RewardEpochs []string  `json:"reward_epochs"`
	Rewards      []Reward  `json:"rewards"`
}

type EventInfo struct {
	Counter string   `json:"counter"`
	Guid    GuidInfo `json:"guid"`
}

type GuidInfo struct {
	ID IDInfo `json:"id"`
}

type IDInfo struct {
	Addr        string `json:"addr"`
	CreationNum string `json:"creation_num"`
}

type Reward struct {
	Eid    string `json:"eid"`
	Euid   string `json:"euid"`
	Share  string `json:"share"`
	Unlock string `json:"unlock"`
}

type Payload struct {
	Function      string   `json:"function"`
	TypeArguments []string `json:"type_arguments"`
	Arguments     []any    `json:"arguments"`
	Type          string   `json:"type"`
}

type Signature struct {
	PublicKey any    `json:"public_key"`
	Signature any    `json:"signature"`
	Type      string `json:"type"`
}

type Event struct {
	Guid           EventGuid    `json:"guid"`
	SequenceNumber string       `json:"sequence_number"`
	Type           string       `json:"type"`
	Data           FeeStatement `json:"data"`
}

type EventGuid struct {
	CreationNumber string `json:"creation_number"`
	AccountAddress string `json:"account_address"`
}

type FeeStatement struct {
	ExecutionGasUnits     uint64 `json:"execution_gas_units,string"`
	IoGasUnits            uint64 `json:"io_gas_units,string"`
	StorageFeeOctas       uint64 `json:"storage_fee_octas,string"`
	StorageFeeRefundOctas uint64 `json:"storage_fee_refund_octas,string"`
	TotalChargeGasUnits   uint64 `json:"total_charge_gas_units,string"`
}

type BlockResponse struct {
	BlockHeight    uint64                `json:"block_height,string"`
	BlockHash      string                `json:"block_hash"`
	BlockTimestamp uint64                `json:"block_timestamp,string"`
	FirstVersion   uint64                `json:"first_version,string"`
	LastVersion    uint64                `json:"last_version,string"`
	Transactions   []TransactionResponse `json:"transactions"`
}

type NodeInfo struct {
	// Chain ID of the current chain
	ChainID             uint8  `json:"chain_id"`
	Epoch               uint64 `json:"epoch,string"`
	LedgerVersion       uint64 `json:"ledger_version,string"`
	OldestLedgerVersion uint64 `json:"oldest_ledger_version,string"`
	LedgerTimestamp     uint64 `json:"ledger_timestamp,string"`
	NodeRole            string `json:"node_role"`
	OldestBlockHeight   uint64 `json:"oldest_block_height,string"`
	BlockHeight         uint64 `json:"block_height,string"`
	GitHash             string `json:"git_hash"`
}

type AccountBalanceResponse struct {
	Type string `json:"type"`
	Data struct {
		Coin struct {
			Value string `json:"value"`
		} `json:"coin"`
	} `json:"data"`
}

type SubmitTransactionRequest struct {
	ChainId                 uint8                `json:"chain_id"`
	Sender                  string               `json:"sender"`
	SequenceNumber          uint64               `json:"sequence_number"`
	MaxGasAmount            uint64               `json:"max_gas_amount"`
	GasUnitPrice            uint64               `json:"gas_unit_price"`
	ExpirationTimestampSecs uint64               `json:"expiration_timestamp_secs"`
	Payload                 TransactionPayload   `json:"payload"`
	Signature               TransactionSignature `json:"signature"`
}

type PayloadType string

const (
	EntryFunctionPayload PayloadType = "entry_function_payload"
	ScriptPayload        PayloadType = "script_payload"
	ModuleBundlePayload  PayloadType = "module_bundle_payload"
)

type TransactionPayload struct {
	Type          PayloadType   `json:"type"`
	Function      string        `json:"function"`
	TypeArguments []string      `json:"type_arguments"`
	Arguments     []interface{} `json:"arguments"`
}
type SignatureType string

const (
	Ed25519Signature SignatureType = "ed25519_signature"
	MultiSignature   SignatureType = "multi_signature"
	MultiEd25519     SignatureType = "multi_ed25519"
)

type TransactionSignature struct {
	Type      SignatureType `json:"type"`
	PublicKey string        `json:"public_key"`
	Signature string        `json:"signature"`
}

type SubmitTransactionResponse struct {
	Hash                    string               `json:"hash"`
	Sender                  string               `json:"sender"`
	SequenceNumber          string               `json:"sequence_number"`
	MaxGasAmount            string               `json:"max_gas_amount"`
	GasUnitPrice            string               `json:"gas_unit_price"`
	ExpirationTimestampSecs string               `json:"expiration_timestamp_secs"`
	Payload                 TransactionPayload   `json:"payload"`
	Signature               TransactionSignature `json:"signature"`
}
