package solana

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/chain-explorer-api/common/chain"
	"github.com/dapplink-labs/chain-explorer-api/explorer/solscan"
)

type SolData struct {
	SolDataCli *solscan.ChainExplorerAdaptor
}

func NewSolScanClient(baseUrl, apiKey string, timeout time.Duration) (*SolData, error) {
	solCli, err := solscan.NewChainExplorerAdaptor(apiKey, baseUrl, false, time.Duration(timeout))
	if err != nil {
		log.Error("New solscan client fail", "err", err)
		return nil, err
	}
	return &SolData{SolDataCli: solCli}, err
}

func (ss *SolData) GetTxByAddress(page, pagesize uint64, address string, action account.ActionType) (*account.TransactionResponse[account.AccountTxResponse], error) {
	request := &account.AccountTxRequest{
		PageRequest: chain.PageRequest{
			Page:  page,
			Limit: pagesize,
		},
		Action:  action,
		Address: address,
	}
	fmt.Printf("%#v\n", request)
	txData, err := ss.SolDataCli.GetTxByAddress(request)
	fmt.Printf("%#v\n", txData)
	if err != nil {
		return nil, err
	}
	return txData, nil
}
