package solana

import "github.com/gagliardetto/solana-go"

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GetHealthResponse struct {
	Jsonrpc string    `json:"jsonrpc"`
	ID      int       `json:"id"`
	Result  string    `json:"result"`
	Error   *RPCError `json:"error,omitempty"`
}

type CreateNonceAccountRequest struct {
	// payer privateKey
	Payer solana.PrivateKey
	// nonce account Auth, PublicKey
	Authority solana.PublicKey
}

type CreateNonceAccountResponse struct {
	// nonce account PublicKey
	NonceAccount solana.PublicKey
	// nonce
	Nonce string
	// nonce
	Signature string
}

type GetAccountInfoResponse struct {
	JsonRPC string    `json:"jsonrpc"`
	ID      int       `json:"id"`
	Error   *RPCError `json:"error,omitempty"`
	Result  struct {
		Context struct {
			// now slot
			Slot uint64 `json:"slot"`
		} `json:"context"`
		Value AccountInfo `json:"value"`
	} `json:"result"`
}

type AccountInfo struct {
	// account now balance
	Lamports uint64 `json:"lamports"`
	Owner    string `json:"owner"`
	// slice index = 0, data
	// slice index = 1, encode = base58, and other
	Data       []string `json:"data"`
	Executable bool     `json:"executable"`
	RentEpoch  uint64   `json:"rentEpoch"`
	Space      uint64   `json:"space"`
}

type GetBalanceResponse struct {
	JsonRPC string    `json:"jsonrpc"`
	ID      int       `json:"id"`
	Error   *RPCError `json:"error,omitempty"`
	Result  struct {
		Context struct {
			Slot uint64 `json:"slot"`
		} `json:"context"`
		Value uint64 `json:"value"`
	} `json:"result"`
}

type BlockHeightResponse struct {
	JsonRPC string    `json:"jsonrpc"`
	ID      int       `json:"id"`
	Error   *RPCError `json:"error,omitempty"`
	// block height
	Result uint64 `json:"result"`
}

type GetSlotRequest struct {
	Commitment CommitmentType `json:"commitment,omitempty"`
	//MinContextSlot uint64     `json:"minContextSlot,omitempty"`
}

type GetSlotResponse struct {
	JsonRPC string    `json:"jsonrpc"`
	ID      int       `json:"id"`
	Error   *RPCError `json:"error,omitempty"`
	// slot
	Result uint64 `json:"result"`
}

// GetBlocksWithLimitResponse represents the response structure
type GetBlocksWithLimitResponse struct {
	JsonRPC string    `json:"jsonrpc"`
	ID      int       `json:"id"`
	Error   *RPCError `json:"error,omitempty"`
	Result  []uint64  `json:"result"`
}

type GetBlockRequest struct {
	// slot status
	// Finalized Confirmed Processed
	Commitment CommitmentType `json:"commitment,omitempty"`
	// "json", "jsonParsed", "base58", "base64"
	Encoding string `json:"encoding"`
	// max version
	// Legacy = 0, no other version
	MaxSupportedTransactionVersion int `json:"maxSupportedTransactionVersion"`
	// "full", "accounts", "signatures", "none"
	TransactionDetails string `json:"transactionDetails"`
	// contain rewards
	Rewards bool `json:"rewards"`
}

type GetBlockResponse struct {
	JsonRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Error   *RPCError   `json:"error,omitempty"`
	Result  BlockResult `json:"result"`
}

type BlockResult struct {
	ParentSlot        uint64              `json:"parentSlot"`
	BlockTime         int64               `json:"blockTime"`
	BlockHeight       uint64              `json:"blockHeight"`
	BlockHash         string              `json:"blockhash"`
	PreviousBlockhash string              `json:"previousBlockhash"`
	Signatures        []string            `json:"signatures"`
	Transactions      []TransactionDetail `json:"transactions"`
}

type TransactionDetail struct {
	Signature string           `json:"signature"`
	Slot      uint64           `json:"slot"`
	BlockTime int64            `json:"blockTime"`
	Meta      *TransactionMeta `json:"meta"`
	// "version": "legacy"   or   "version": 0
	Version         any         `json:"version"`
	Message         interface{} `json:"message"` // 使用 interface{} 因为可能为 null
	RecentBlockhash string      `json:"recentBlockhash"`
}

type LoadedAddresses struct {
	Readonly []string `json:"readonly"`
	Writable []string `json:"writable"`
}

type TransactionVersion struct {
	value interface{}
}

type Status struct {
	Ok interface{} `json:"Ok"`
}

//type GetTransactionRequest struct {
//	Encoding   string `json:"encoding,omitempty"`
//	Commitment string `json:"commitment,omitempty"`
//	// max version
//	// Legacy = 0, no other version
//	MaxSupportedTransactionVersion string `json:"maxSupportedTransactionVersion,omitempty"`
//}

type GetTransactionResponse struct {
	Jsonrpc string            `json:"jsonrpc"`
	ID      int               `json:"id"`
	Error   *RPCError         `json:"error,omitempty"`
	Result  TransactionResult `json:"result"`
}

type TransactionResult struct {
	Slot        uint64          `json:"slot"`
	Version     any             `json:"version"`
	BlockTime   *int64          `json:"blockTime"`
	Transaction Transaction     `json:"transaction"`
	Meta        TransactionMeta `json:"meta"`
}

type Transaction struct {
	Message    TransactionMessage `json:"message"`
	Signatures []string           `json:"signatures"`
}

type TransactionMeta struct {
	Err               interface{}     `json:"err"`
	Fee               uint64          `json:"fee"`
	PreBalances       []uint64        `json:"preBalances"`
	PostBalances      []uint64        `json:"postBalances"`
	InnerInstructions []interface{}   `json:"innerInstructions"`
	PreTokenBalances  []interface{}   `json:"preTokenBalances"`
	PostTokenBalances []interface{}   `json:"postTokenBalances"`
	LogMessages       []string        `json:"logMessages"`
	LoadedAddresses   LoadedAddresses `json:"loadedAddresses"`
	Status            struct {
		Ok interface{} `json:"Ok"`
	} `json:"status"`
	Rewards              interface{} `json:"rewards"`
	ComputeUnitsConsumed uint64      `json:"computeUnitsConsumed"`
}

type TransactionMessage struct {
	AccountKeys         []string          `json:"accountKeys"`
	AddressTableLookups []interface{}     `json:"addressTableLookups"`
	Header              TransactionHeader `json:"header"`
	Instructions        []Instruction     `json:"instructions"`
	RecentBlockhash     string            `json:"recentBlockhash"`
}

type TransactionHeader struct {
	NumReadonlySignedAccounts   int `json:"numReadonlySignedAccounts"`
	NumReadonlyUnsignedAccounts int `json:"numReadonlyUnsignedAccounts"`
	NumRequiredSignatures       int `json:"numRequiredSignatures"`
}

type Instruction struct {
	Accounts       []int       `json:"accounts"`
	Data           string      `json:"data"`
	ProgramIdIndex int         `json:"programIdIndex"`
	StackHeight    interface{} `json:"stackHeight"`
}

type GetFeeForMessageRequest struct {
	Commitment     string `json:"commitment,omitempty"`
	MinContextSlot uint64 `json:"minContextSlot,omitempty"`
}

type GetFeeForMessageResponse struct {
	Jsonrpc string    `json:"jsonrpc"`
	ID      int       `json:"id"`
	Error   *RPCError `json:"error,omitempty"`
	Result  struct {
		Context struct {
			Slot uint64 `json:"slot"`
		} `json:"context"`
		Value *uint64 `json:"value"`
	} `json:"result"`
}

type getRecentPrioritizationFeesResponse struct {
	Jsonrpc string              `json:"jsonrpc"`
	ID      int                 `json:"id"`
	Error   *RPCError           `json:"error,omitempty"`
	Result  []PrioritizationFee `json:"result"`
}

type PrioritizationFee struct {
	Slot              uint64 `json:"slot"`
	PrioritizationFee uint64 `json:"prioritizationFee"`
}

type GetSignaturesRequest struct {
	Commitment     string `json:"commitment,omitempty"`
	MinContextSlot uint64 `json:"minContextSlot,omitempty"`
	Limit          uint64 `json:"limit,omitempty"`
	Before         string `json:"before,omitempty"`
	Until          string `json:"until,omitempty"`
}

type GetSignaturesResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Error   *RPCError       `json:"error,omitempty"`
	Result  []SignatureInfo `json:"result"`
}

type SignatureInfo struct {
	Signature          string      `json:"signature"`
	Slot               uint64      `json:"slot"`
	Error              interface{} `json:"err"`
	Memo               *string     `json:"memo"`
	BlockTime          *int64      `json:"blockTime"`
	ConfirmationStatus *string     `json:"confirmationStatus"`
}

type SimulateRequest struct {
	Commitment             string        `json:"commitment,omitempty"`
	SigVerify              bool          `json:"sigVerify,omitempty"`
	ReplaceRecentBlockhash bool          `json:"replaceRecentBlockhash,omitempty"`
	MinContextSlot         uint64        `json:"minContextSlot,omitempty"`
	Encoding               string        `json:"encoding,omitempty"`
	Accounts               *AccountsInfo `json:"accounts,omitempty"`
}

type AccountsInfo struct {
	Addresses []string `json:"addresses"`
	Encoding  string   `json:"encoding,omitempty"`
}

type SimulateTransactionResponse struct {
	Jsonrpc string         `json:"jsonrpc"`
	ID      int            `json:"id"`
	Error   *RPCError      `json:"error,omitempty"`
	Result  SimulateResult `json:"result"`
}

type SimulateResult struct {
	Err           interface{} `json:"err"`
	Logs          []string    `json:"logs"`
	UnitsConsumed uint64      `json:"unitsConsumed"`
	Accounts      []struct {
		Executable bool     `json:"executable"`
		Lamports   uint64   `json:"lamports"`
		Owner      string   `json:"owner"`
		RentEpoch  uint64   `json:"rentEpoch"`
		Data       []string `json:"data"`
	} `json:"accounts,omitempty"`
	ReturnData *struct {
		ProgramId string   `json:"programId"`
		Data      []string `json:"data"`
	} `json:"returnData,omitempty"`
	InnerInstructions []struct {
		Index        uint16 `json:"index"`
		Instructions []struct {
			ProgramIdIndex uint8   `json:"programIdIndex"`
			Accounts       []uint8 `json:"accounts"`
			Data           string  `json:"data"`
		} `json:"instructions"`
	} `json:"innerInstructions,omitempty"`
}

type SendTransactionRequest struct {
	Encoding            string `json:"encoding,omitempty"`
	Commitment          string `json:"commitment,omitempty"`
	SkipPreflight       bool   `json:"skipPreflight,omitempty"`
	PreflightCommitment string `json:"preflightCommitment,omitempty"`
	MaxRetries          uint64 `json:"maxRetries,omitempty"`
	MinContextSlot      uint64 `json:"minContextSlot,omitempty"`
}

type SendTransactionResponse struct {
	Jsonrpc string    `json:"jsonrpc"`
	ID      int       `json:"id"`
	Result  string    `json:"result"`
	Error   *RPCError `json:"error,omitempty"`
}

type TxStructure struct {
	Nonce           string `json:"nonce"`
	GasPrice        string `json:"gas_price"`
	GasTipCap       string `json:"gas_tip_cap"`
	GasFeeCap       string `json:"gas_fee_cap"`
	Gas             uint64 `json:"gas"`
	ContractAddress string `json:"contract_address"`
	FromAddress     string `json:"from_address"`
	ToAddress       string `json:"to_address"`
	TokenId         string `json:"token_id"`
	Value           string `json:"value"`
	Signature       string `json:"signature,omitempty"`
}

type GetLatestBlockhashResponse struct {
	JsonRPC string    `json:"jsonrpc"`
	ID      int       `json:"id"`
	Error   *RPCError `json:"error,omitempty"`
	Result  struct {
		Context struct {
			Slot uint64 `json:"slot"`
		} `json:"context"`
		Value struct {
			Blockhash            string `json:"blockhash"`
			LastValidBlockHeight uint64 `json:"lastValidBlockHeight"`
		} `json:"value"`
	} `json:"result"`
}
