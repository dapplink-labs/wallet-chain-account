package tezos

type Operations struct {
	Type      string `json:"type"`
	Id        uint64 `json:"id"`
	Level     uint64 `json:"level"`
	TimeStamp string `json:"timestamp"`
	Block     string `json:"block"`
	Hash      string `json:"hash"`
	Counter   uint64 `json:"counter"`
	Sender    struct {
		Address string `json:"address"`
		Alias   string `json:"alias"`
	} `json:"sender"`
	GasLimit      uint64 `json:"gasLimit"`
	GasUsed       uint64 `json:"gasUsed"`
	StorageLimit  uint64 `json:"storageLimit"`
	StorageUsed   uint64 `json:"storageUsed"`
	BakerFee      uint64 `json:"bakerFee"`
	StorageFee    uint64 `json:"storageFee"`
	AllocationFee uint64 `json:"allocationFee"`
	Target        struct {
		Address string `json:"address"`
		Alias   string `json:"alias"`
	} `json:"target"`
	Amount       uint64 `json:"amount"`
	Status       string `json:"status"`
	HasInternals bool   `json:"hasInternals"`
}

type Block struct {
	Cycle              uint64 `json:"cycle"`
	Level              uint64 `json:"level"`
	Hash               string `json:"hash"`
	TimeStamp          string `json:"timestamp"`
	Proto              uint64 `json:"proto"`
	PayloadRound       uint64 `json:"payloadRound"`
	BlockRound         uint64 `json:"blockRound"`
	Validations        uint64 `json:"validations"`
	Deposit            uint64 `json:"deposit"`
	RewardDelegated    uint64 `json:"rewardDelegated"`
	RewardStakedOwn    uint64 `json:"rewardStakedOwn"`
	RewardStakedEdge   uint64 `json:"rewardStakedEdge"`
	RewardStakedShared uint64 `json:"rewardStakedShared"`
	Fees               uint64 `json:"fees"`
	NonceRevealed      bool   `json:"nonceRevealed"`
	Proposer           struct {
		Address string `json:"address"`
	} `json:"proposer"`
	Producer struct {
		Address string `json:"address"`
	} `json:"producer"`
	Software struct {
		Version string `json:"version"`
		Date    string `json:"date"`
	} `json:"software"`
	LbToggleEma  uint64 `json:"lbToggleEma"`
	AiToggleEma  uint64 `json:"aiToggleEma"`
	RewardLiquid uint64 `json:"rewardLiquid"`
	BonusLiquid  uint64 `json:"bonusLiquid"`
	Reward       uint64 `json:"reward"`
	Bonus        uint64 `json:"bonus"`
	Priority     uint64 `json:"priority"`
	Baker        struct {
		Address string `json:"address"`
	} `json:"baker"`
	LbEscapeVote bool   `json:"lbEscapeVote"`
	LbEscapeEma  uint64 `json:"lbEscapeEma"`
}

type Account struct {
	Id                uint64 `json:"id"`
	Type              string `json:"type"`
	Address           string `json:"address"`
	PublicKey         string `json:"publicKey"`
	Revealed          bool   `json:"revealed"`
	Balance           uint64 `json:"balance"`
	Counter           uint64 `json:"counter"`
	FirstActivity     uint64 `json:"firstActivity"`
	FirstActivityTime string `json:"firstActivityTime"`
	LastActivity      uint64 `json:"lastActivity"`
	LastActivityTime  string `json:"lastActivityTime"`
	LostBalance       uint64 `json:"lostBalance"`
}

type SendTransaction struct {
	Branch   string `json:"branch"`
	Protocol string `json:"protocol"`
	Contents []struct {
		Kind         string `json:"kind"`
		Source       string `json:"source"`
		Fee          string `json:"fee"`
		Counter      string `json:"counter"`
		GasLimit     string `json:"gas_limit"`
		StorageLimit string `json:"storage_limit"`
		Destination  string `json:"destination"`
		Amount       string `json:"amount"`
	} `json:"contents"`
	Signature string `json:"signature"`
}

type PreApplyTransaction struct {
	Branch   string `json:"branch"`
	Protocol string `json:"protocol"`
	Contents []struct {
		Kind         string `json:"kind"`
		Source       string `json:"source"`
		Fee          string `json:"fee"`
		Counter      string `json:"counter"`
		GasLimit     string `json:"gas_limit"`
		StorageLimit string `json:"storage_limit"`
		Destination  string `json:"destination"`
		Amount       string `json:"amount"`
		Metadata     struct {
			BalanceUpdates []struct {
				Kind     string `json:"kind"`
				Contract string `json:"contract"`
				Category string `json:"category"`
				Change   string `json:"change"`
				Origin   string `json:"origin"`
			} `json:"balance_updates"`
			OperationResult struct {
				Status         string `json:"status"`
				BalanceUpdates []struct {
					Kind     string `json:"kind"`
					Contract string `json:"contract"`
					Category string `json:"category"`
					Change   string `json:"change"`
					Origin   string `json:"origin"`
				} `json:"balance_updates"`
				ConsumedMilligas string `json:"consumed_milligas"`
			} `json:"operation_result"`
		}
	} `json:"contents"`
	Signature string `json:"signature"`
}
