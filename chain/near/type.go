package near

import (
	"encoding/json"
	"fmt"
)

type CongestionInfo struct {
	AllowedShard        int    `json:"allowed_shard"`
	BufferedReceiptsGas string `json:"buffered_receipts_gas"`
	DelayedReceiptsGas  string `json:"delayed_receipts_gas"`
	ReceiptBytes        int    `json:"receipt_bytes"`
}

type Header struct {
	BalanceBurnt         string         `json:"balance_burnt"`
	BandwidthRequests    interface{}    `json:"bandwidth_requests"`
	ChunkHash            string         `json:"chunk_hash"`
	CongestionInfo       CongestionInfo `json:"congestion_info"`
	EncodedLength        int            `json:"encoded_length"`
	EncodedMerkleRoot    string         `json:"encoded_merkle_root"`
	GasLimit             int64          `json:"gas_limit"`
	GasUsed              int64          `json:"gas_used"`
	HeightCreated        int64          `json:"height_created"`
	HeightIncluded       int64          `json:"height_included"`
	OutcomeRoot          string         `json:"outcome_root"`
	OutgoingReceiptsRoot string         `json:"outgoing_receipts_root"`
	PrevBlockHash        string         `json:"prev_block_hash"`
	PrevStateRoot        string         `json:"prev_state_root"`
	RentPaid             string         `json:"rent_paid"`
	ShardID              int            `json:"shard_id"`
	Signature            string         `json:"signature"`
	TxRoot               string         `json:"tx_root"`
	ValidatorProposals   []interface{}  `json:"validator_proposals"`
	ValidatorReward      string         `json:"validator_reward"`
}

type DelegateAction struct {
	Actions        []Action `json:"actions"`
	MaxBlockHeight int64    `json:"max_block_height"`
	Nonce          int64    `json:"nonce"`
	PublicKey      string   `json:"public_key"`
	ReceiverID     string   `json:"receiver_id"`
	SenderID       string   `json:"sender_id"`
	Signature      string   `json:"signature"`
}

type Action struct {
	Type          string         `json:"type"`
	Delegate      *Delegate      `json:"Delegate,omitempty"`
	FunctionCall  *FunctionCall  `json:"FunctionCall,omitempty"`
	Transfer      *Transfer      `json:"Transfer,omitempty"`
	AddKey        *AddKey        `json:"AddKey,omitempty"`
	CreateAccount *CreateAccount `json:"CreateAccount,omitempty"`
}

func (a *Action) UnmarshalJSON(data []byte) error {
	// Check if data is a string and handle it
	if data[0] == '"' {
		var str string
		if err := json.Unmarshal(data, &str); err != nil {
			return err
		}
		// Handle the plain string case
		a.Type = "PlainString"
		return nil
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for key, value := range raw {
		a.Type = key
		switch key {
		case "Delegate":
			a.Delegate = &Delegate{}
			return json.Unmarshal(value, a.Delegate)
		case "FunctionCall":
			a.FunctionCall = &FunctionCall{}
			return json.Unmarshal(value, a.FunctionCall)
		case "Transfer":
			a.Transfer = &Transfer{}
			return json.Unmarshal(value, a.Transfer)
		case "AddKey":
			a.AddKey = &AddKey{}
			return json.Unmarshal(value, a.AddKey)
		case "CreateAccount":
			a.CreateAccount = &CreateAccount{}
			return nil
		default:
			return fmt.Errorf("unknown action type: %s", key)
		}
	}

	return nil
}

type Delegate struct {
	DelegateAction DelegateAction `json:"delegate_action"`
	Signature      string         `json:"signature"`
}

type FunctionCall struct {
	Args       string `json:"args"`
	Deposit    string `json:"deposit"`
	Gas        int64  `json:"gas"`
	MethodName string `json:"method_name"`
}

type Transfer struct {
	Deposit string `json:"deposit"`
}

type AddKey struct {
	AccessKey struct {
		Nonce      int64  `json:"nonce"`
		Permission string `json:"permission"`
	} `json:"access_key"`
	PublicKey string `json:"public_key"`
}

type CreateAccount struct{}

type Receipt struct {
	PredecessorID string `json:"predecessor_id"`
	Priority      int    `json:"priority"`
	Receipt       struct {
		Action struct {
			Actions             []Action      `json:"actions"`
			GasPrice            string        `json:"gas_price"`
			InputDataIDs        []interface{} `json:"input_data_ids"`
			IsPromiseYield      bool          `json:"is_promise_yield"`
			OutputDataReceivers []interface{} `json:"output_data_receivers"`
			SignerID            string        `json:"signer_id"`
			SignerPublicKey     string        `json:"signer_public_key"`
		} `json:"Action"`
	} `json:"receipt"`
	ReceiptID  string `json:"receipt_id"`
	ReceiverID string `json:"receiver_id"`
}

type Transaction struct {
	Actions     []Action `json:"actions"`
	Hash        string   `json:"hash"`
	Nonce       int64    `json:"nonce"`
	PriorityFee int64    `json:"priority_fee"`
	PublicKey   string   `json:"public_key"`
	ReceiverID  string   `json:"receiver_id"`
	Signature   string   `json:"signature"`
	SignerID    string   `json:"signer_id"`
}

type ChunkDetail struct {
	Author       string        `json:"author"`
	Header       Header        `json:"header"`
	Receipts     []Receipt     `json:"receipts"`
	Transactions []Transaction `json:"transactions"`
}

type TransactionStatus struct {
	FinalExecutionStatus string           `json:"final_execution_status"`
	Receipts             []Receipt        `json:"receipts"`
	ReceiptsOutcome      []ReceiptOutcome `json:"receipts_outcome"`
	Status               struct {
		SuccessValue string `json:"SuccessValue"`
	} `json:"status"`
	Transaction        Transaction        `json:"transaction"`
	TransactionOutcome TransactionOutcome `json:"transaction_outcome"`
}

type TransactionOutcome struct {
	BlockHash string `json:"block_hash"`
	ID        string `json:"id"`
	Outcome   struct {
		ExecutorID string        `json:"executor_id"`
		GasBurnt   int64         `json:"gas_burnt"`
		Logs       []interface{} `json:"logs"`
		Metadata   struct {
			GasProfile interface{} `json:"gas_profile"`
			Version    int         `json:"version"`
		} `json:"metadata"`
		ReceiptIDs []string `json:"receipt_ids"`
		Status     struct {
			SuccessReceiptID string `json:"SuccessReceiptId"`
		} `json:"status"`
		TokensBurnt string `json:"tokens_burnt"`
	} `json:"outcome"`
	Proof []struct {
		Direction string `json:"direction"`
		Hash      string `json:"hash"`
	} `json:"proof"`
}

type ReceiptOutcome struct {
	BlockHash string `json:"block_hash"`
	ID        string `json:"id"`
	Outcome   struct {
		ExecutorID string        `json:"executor_id"`
		GasBurnt   int64         `json:"gas_burnt"`
		Logs       []interface{} `json:"logs"`
		Metadata   struct {
			GasProfile []interface{} `json:"gas_profile"`
			Version    int           `json:"version"`
		} `json:"metadata"`
		ReceiptIDs []string `json:"receipt_ids"`
		Status     struct {
			SuccessValue string `json:"SuccessValue"`
		} `json:"status"`
		TokensBurnt string `json:"tokens_burnt"`
	} `json:"outcome"`
	Proof []struct {
		Direction string `json:"direction"`
		Hash      string `json:"hash"`
	} `json:"proof"`
}
