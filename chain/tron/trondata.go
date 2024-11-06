package tron

import (
	"github.com/dapplink-labs/chain-explorer-api/common/gas_fee"
	"github.com/ethereum/go-ethereum/log"
	"time"

	"github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/chain-explorer-api/common/chain"
	"github.com/dapplink-labs/chain-explorer-api/explorer/oklink"
)

type TronData struct {
	TronDataCli *oklink.ChainExplorerAdaptor
}

func NewTronDataClient(baseUrl, apiKey string, timeout time.Duration) (*TronData, error) {
	tronDataCli, err := oklink.NewChainExplorerAdaptor(apiKey, baseUrl+"/", false, time.Duration(timeout))
	if err != nil {
		log.Error("New tronerscan client fail", "err", err)
		return nil, err
	}
	return &TronData{TronDataCli: tronDataCli}, err
}

// GetTxByAddress Transaction for obtaining address
func (td *TronData) GetTxByAddress(page, pagesize uint64, address string, action account.ActionType) (*account.TransactionResponse[account.AccountTxResponse], error) {
	request := &account.AccountTxRequest{
		ChainShortName: ChainName,
		ExplorerName:   oklink.ChainExplorerName,
		Action:         action,
		Address:        address,
		PageRequest: chain.PageRequest{
			Page:  page,
			Limit: pagesize,
		},
	}
	txData, err := td.TronDataCli.GetTxByAddress(request)
	if err != nil {
		return nil, err
	}
	return txData, nil
}

// GetEstimateGasFee Obtain estimated transaction fees
func (td *TronData) GetEstimateGasFee() (*gas_fee.GasEstimateFeeResponse, error) {
	request := &gas_fee.GasEstimateFeeRequest{
		ChainShortName: ChainName,
		ExplorerName:   oklink.ChainExplorerName,
	}
	gasFee, err := td.TronDataCli.GetEstimateGasFee(request)
	if err != nil {
		return nil, err
	}
	return gasFee, nil
}
