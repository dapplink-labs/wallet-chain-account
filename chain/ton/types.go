package ton

type EstimateFeeRequest struct {
	Address      string `json:"address"`
	Body         string `json:"body"`
	IgnoreChksig bool   `json:"ignore_chksig"`
}

type Transactions struct {
	Account       string `json:"account"`
	Hash          string `json:"hash"`
	Lt            string `json:"lt"`
	Now           int    `json:"now"`
	OrigStatus    string `json:"orig_status"`
	EndStatus     string `json:"end_status"`
	TotalFees     string `json:"total_fees"`
	PrevTransHash string `json:"prev_trans_hash"`
	PrevTransLt   string `json:"prev_trans_lt"`
	Description   struct {
		Type   string `json:"type"`
		Action struct {
			Valid       bool `json:"valid"`
			Success     bool `json:"success"`
			NoFunds     bool `json:"no_funds"`
			ResultCode  int  `json:"result_code"`
			TotActions  int  `json:"tot_actions"`
			MsgsCreated int  `json:"msgs_created"`
			SpecActions int  `json:"spec_actions"`
			TotMsgSize  struct {
				Bits  string `json:"bits"`
				Cells string `json:"cells"`
			} `json:"tot_msg_size"`
			StatusChange    string `json:"status_change"`
			TotalFwdFees    string `json:"total_fwd_fees"`
			SkippedActions  int    `json:"skipped_actions"`
			ActionListHash  string `json:"action_list_hash"`
			TotalActionFees string `json:"total_action_fees"`
		} `json:"action"`
		Aborted  bool `json:"aborted"`
		CreditPh struct {
			Credit string `json:"credit"`
		} `json:"credit_ph"`
		Destroyed bool `json:"destroyed"`
		ComputePh struct {
			Mode             int    `json:"mode"`
			Type             string `json:"type"`
			Success          bool   `json:"success"`
			GasFees          string `json:"gas_fees"`
			GasUsed          string `json:"gas_used"`
			VMSteps          int    `json:"vm_steps"`
			ExitCode         int    `json:"exit_code"`
			GasLimit         string `json:"gas_limit"`
			GasCredit        string `json:"gas_credit"`
			MsgStateUsed     bool   `json:"msg_state_used"`
			AccountActivated bool   `json:"account_activated"`
			VMInitStateHash  string `json:"vm_init_state_hash"`
			VMFinalStateHash string `json:"vm_final_state_hash"`
		} `json:"compute_ph"`
		StoragePh struct {
			StatusChange         string `json:"status_change"`
			StorageFeesCollected string `json:"storage_fees_collected"`
		} `json:"storage_ph"`
		CreditFirst bool `json:"credit_first"`
	} `json:"description"`
	BlockRef struct {
		Workchain int    `json:"workchain"`
		Shard     string `json:"shard"`
		Seqno     int    `json:"seqno"`
	} `json:"block_ref"`
	InMsg struct {
		Hash           string      `json:"hash"`
		Source         string      `json:"source"`
		Destination    string      `json:"destination"`
		Value          string      `json:"value"`
		FwdFee         interface{} `json:"fwd_fee"`
		IhrFee         interface{} `json:"ihr_fee"`
		CreatedLt      interface{} `json:"created_lt"`
		CreatedAt      interface{} `json:"created_at"`
		Opcode         string      `json:"opcode"`
		IhrDisabled    interface{} `json:"ihr_disabled"`
		Bounce         interface{} `json:"bounce"`
		Bounced        interface{} `json:"bounced"`
		ImportFee      string      `json:"import_fee"`
		MessageContent struct {
			Hash    string      `json:"hash"`
			Body    string      `json:"body"`
			Decoded interface{} `json:"decoded"`
		} `json:"message_content"`
		InitState interface{} `json:"init_state"`
	} `json:"in_msg"`
	OutMsgs []struct {
		Hash           string      `json:"hash"`
		Source         string      `json:"source"`
		Destination    string      `json:"destination"`
		Value          string      `json:"value"`
		FwdFee         string      `json:"fwd_fee"`
		IhrFee         string      `json:"ihr_fee"`
		CreatedLt      string      `json:"created_lt"`
		CreatedAt      string      `json:"created_at"`
		Opcode         interface{} `json:"opcode"`
		IhrDisabled    bool        `json:"ihr_disabled"`
		Bounce         bool        `json:"bounce"`
		Bounced        bool        `json:"bounced"`
		ImportFee      interface{} `json:"import_fee"`
		MessageContent struct {
			Hash    string      `json:"hash"`
			Body    string      `json:"body"`
			Decoded interface{} `json:"decoded"`
		} `json:"message_content"`
		InitState interface{} `json:"init_state"`
	} `json:"out_msgs"`
	AccountStateBefore struct {
		Hash          string      `json:"hash"`
		Balance       string      `json:"balance"`
		AccountStatus string      `json:"account_status"`
		FrozenHash    interface{} `json:"frozen_hash"`
		CodeHash      string      `json:"code_hash"`
		DataHash      string      `json:"data_hash"`
	} `json:"account_state_before"`
	AccountStateAfter struct {
		Hash          string      `json:"hash"`
		Balance       string      `json:"balance"`
		AccountStatus string      `json:"account_status"`
		FrozenHash    interface{} `json:"frozen_hash"`
		CodeHash      string      `json:"code_hash"`
		DataHash      string      `json:"data_hash"`
	} `json:"account_state_after"`
	McBlockSeqno int `json:"mc_block_seqno"`
}

type Tx struct {
	Transactions []Transactions `json:"transactions"`
	AddressBook  map[string]struct {
		UserFriendly string `json:"user_friendly"`
	} `json:"address_book"`
}

type SendTxResult struct {
	Hash string `json:"message_hash"`
}

type EstimateFeeResult struct {
	InFwdFee   uint64 `json:"in_fwd_fee"`
	StorageFee uint64 `json:"storage_fee"`
	GasFee     uint64 `json:"gas_fee"`
	FwdFee     uint64 `json:"fwd_fee"`
}

type SourceFeesResult struct {
	SourceFees *EstimateFeeResult `json:"source_fees"`
}
