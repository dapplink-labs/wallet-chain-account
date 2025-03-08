package xlm

type ResponseGetTransactionForLedgers struct {
	Embedded ResultDataEmbedded `json:"_embedded"`
}

type ResultDataEmbedded struct {
	Records []ResultDataTransaction `json:"records"`
}

type ResultDataTransaction struct {
	Id                      string `json:"id"`
	Successful              bool   `json:"successful"`
	Hash                    string `json:"hash"`
	Ledger                  int64  `json:"ledger"`
	Created_at              string `json:"created_at"`
	Source_account          string `json:"source_account"`
	Source_account_Sequence string `json:"source_account_sequence"`
	Fee_account             string `json:"fee_account"`
	Fee_charged             string `json:"fee_charged"`
	Max_fee                 string `json:"max_fee"`
	Memo_type               string `json:"memo_type"`
}

////////////////////////////////////////////////////////////////////

type ResponseAccountInfo struct {
	ID                   string              `json:"id"`
	AccountID            string              `json:"account_id"`
	Sequence             string              `json:"sequence"`
	SequenceLedger       int                 `json:"sequence_ledger"`
	SequenceTime         string              `json:"sequence_time"`
	SubentryCount        int                 `json:"subentry_count"`
	InflationDestination string              `json:"inflation_destination"`
	HomeDomain           string              `json:"home_domain"`
	LastModifiedLedger   int                 `json:"last_modified_ledger"`
	LastModifiedTime     string              `json:"last_modified_time"`
	Balances             []ResultDataBalance `json:"balances"`
}

type ResultDataBalance struct {
	Balance            string `json:"balance"`
	BuyingLiabilities  string `json:"buying_liabilities"`
	SellingLiabilities string `json:"selling_liabilities"`
	AssetType          string `json:"asset_type"`
}

////////////////////////////////////////////////////////////////////

type RequestGetFeeStats struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
}

type ResponseGetFeeStats struct {
	Jsonrpc string           `json:"jsonrpc"`
	ID      int              `json:"id"`
	Result  ResultDataFeeRes `json:"result"`
}

type ResultDataFeeRes struct {
	SorobanInclusionFee ResultDataInclusionFee `json:"sorobanInclusionFee"`
	InclusionFee        ResultDataInclusionFee `json:"inclusionFee"`
	LatestLedger        int                    `json:"latestLedger"`
}

type ResultDataInclusionFee struct {
	Max              string `json:"max"`
	Min              string `json:"min"`
	Mode             string `json:"mode"`
	P10              string `json:"p10"`
	P20              string `json:"p20"`
	P30              string `json:"p30"`
	P40              string `json:"p40"`
	P50              string `json:"p50"`
	P60              string `json:"p60"`
	P70              string `json:"p70"`
	P80              string `json:"p80"`
	P90              string `json:"p90"`
	P95              string `json:"p95"`
	P99              string `json:"p99"`
	TransactionCount string `json:"transactionCount"`
	LedgerCount      int    `json:"ledgerCount"`
}

////////////////////////////////////////////////////////////////////

type ResponseGetTransactionEffect struct {
	Embedded struct {
		Records []struct {
			ID          string `json:"id"`
			PagingToken string `json:"paging_token"`
			Account     string `json:"account"`
			Type        string `json:"type"`
			TypeI       int    `json:"type_i"`
			CreatedAt   string `json:"created_at"`
			AssetType   string `json:"asset_type"`
			Amount      string `json:"amount"`
		} `json:"records"`
	} `json:"_embedded"`
}

type ResponseGetTransactionPart01 struct {
	ID                    string   `json:"id"`
	PagingToken           string   `json:"paging_token"`
	Successful            bool     `json:"successful"`
	Hash                  string   `json:"hash"`
	Ledger                uint32   `json:"ledger"`
	CreatedAt             string   `json:"created_at"`
	SourceAccount         string   `json:"source_account"`
	SourceAccountSequence string   `json:"source_account_sequence"`
	FeeAccount            string   `json:"fee_account"`
	FeeCharged            string   `json:"fee_charged"`
	MaxFee                string   `json:"max_fee"`
	OperationCount        uint32   `json:"operation_count"`
	EnvelopeXDR           string   `json:"envelope_xdr"`
	ResultXDR             string   `json:"result_xdr"`
	FeeMetaXDR            string   `json:"fee_meta_xdr"`
	MemoType              string   `json:"memo_type"`
	Signatures            []string `json:"signatures"`
}

////////////////////////////////////////////////////////////////////

type ResponseGetBlockHeader struct {
	ID                         string `json:"id"`
	PagingToken                string `json:"paging_token"`
	Hash                       string `json:"hash"`
	PrevHash                   string `json:"prev_hash"`
	Sequence                   int    `json:"sequence"`
	SuccessfulTransactionCount int    `json:"successful_transaction_count"`
	FailedTransactionCount     int    `json:"failed_transaction_count"`
	OperationCount             int    `json:"operation_count"`
	TxSetOperationCount        int    `json:"tx_set_operation_count"`
	ClosedAt                   string `json:"closed_at"`
	TotalCoins                 string `json:"total_coins"`
	FeePool                    string `json:"fee_pool"`
	BaseFeeInStroops           int    `json:"base_fee_in_stroops"`
	BaseReserveInStroops       int    `json:"base_reserve_in_stroops"`
	MaxTxSetSize               int    `json:"max_tx_set_size"`
	ProtocolVersion            int    `json:"protocol_version"`
	HeaderXDR                  string `json:"header_xdr"`
}

////////////////////////////////////////////////////////////////////

type RequestSendTransaction struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  struct {
		Transaction string `json:"transaction"`
		XdrFormat   string `json:"xdrFormat"`
	} `json:"params"`
}

type ResponseSendTransaction struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Status                string `json:"status"`
		Hash                  string `json:"hash"`
		LatestLedger          int    `json:"latestLedger"`
		LatestLedgerCloseTime string `json:"latestLedgerCloseTime"`
	} `json:"result"`
}

////////////////////////////////////////////////////////////////////

type RequestCreateUnsignTransaction struct {
	AddrFrom     string `json:"addrFrom"`
	AddrTo       string `json:"addrTo"`
	SequenceFrom int64  `json:"sequenceFrom"`
	Amount       string `json:"amount"`
}
