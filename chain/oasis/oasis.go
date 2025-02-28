package oasis

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"github.com/ethereum/go-ethereum/log"
	"strconv"
)

const ChainName = "Oasis"

type ChainAdaptor struct {
	apiClient  *OasisRosettaClient
	dataClient *OasisNexusClient
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	apiClient, err := NewOasisRosettaClient(context.Background(), conf.WalletNode.Oasis.RpcUrl, conf.WalletNode.Oasis.TimeOut)
	if err != nil {
		return nil, err
	}
	dataClient, err := NewOasisNexusClient(conf.WalletNode.Oasis.DataApiUrl)
	if err != nil {
		return nil, err
	}

	return &ChainAdaptor{
		apiClient:  apiClient,
		dataClient: dataClient,
	}, nil
}

func (c *ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	return &account.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

func (c *ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	pubKey, _ := hex.DecodeString(req.PublicKey)
	resp, err := c.apiClient.ConstructionDerive(context.Background(), pubKey, types.Edwards25519)
	if err != nil {
		log.Error("ConvertAddress failed", err)
		return &account.ConvertAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "ConvertAddress error",
		}, nil
	}
	return &account.ConvertAddressResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "",
		Address: resp.AccountIdentifier.Address,
	}, nil
}

func (c *ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	address := req.GetAddress()
	if isValidOasisAddress(address) {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_SUCCESS,
			Msg:   "",
			Valid: true,
		}, nil
	}
	return &account.ValidAddressResponse{
		Code:  common2.ReturnCode_ERROR,
		Msg:   "ValidAddress error",
		Valid: false,
	}, nil
}

func isValidOasisAddress(address string) bool {
	if address == "" {
		return false
	}
	hrp, _, err := bech32.Decode(address)
	if err != nil {
		return false
	}
	return hrp == "oasis"
}

func (c *ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	height := req.Height
	//get lastBlock
	if height == 0 {
		lastBlock, _ := c.dataClient.GetLastBlockNumber()
		height = lastBlock.LatestBlock
	}
	block, err := c.dataClient.GetBlockByNumber(height)
	if err != nil {
		log.Error("GetBlockByNumber fail", "err", err)
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "GetBlockByNumber fail",
		}, nil
	}

	var txListRet []*account.BlockInfoTransactionList
	for _, v := range block.Transactions {
		//only do Transfer
		if v.Method != "staking.Transfer" {
			continue
		}
		bitlItem := &account.BlockInfoTransactionList{
			From:           v.Sender,
			To:             v.Body.To,
			TokenAddress:   "",
			ContractWallet: "",
			Hash:           v.Hash,
			Height:         uint64(v.Block),
			Amount:         v.Body.Amount,
		}
		txListRet = append(txListRet, bitlItem)
	}

	return &account.BlockResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "",
		Height:       height,
		Hash:         "string",
		BaseFee:      "string",
		Transactions: txListRet,
	}, nil

}

func (c *ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	accountBalances, err := c.apiClient.FetchAccountBalances(context.Background(), req.Address)
	if err != nil {
		return &account.AccountResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "FetchAccountBalances  error",
		}, nil
	}
	nonce, _ := accountBalances.Metadata["nonce"].(float64)
	nonceStr := strconv.FormatFloat(nonce, 'f', -1, 64)
	return &account.AccountResponse{
		Code:          common2.ReturnCode_SUCCESS,
		Msg:           "",
		Network:       req.Network,
		AccountNumber: req.Address,
		Sequence:      nonceStr,
		Balance:       accountBalances.Balances[0].Value,
	}, nil
}

func (c *ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	submit, err := c.apiClient.ConstructionSubmit(context.Background(), req.GetRawTx())
	if err != nil {
		log.Error("Submit transaction failed", err)
		return &account.SendTxResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "ConstructionSubmit error:" + err.Error(),
		}, nil
	}
	hash := submit.TransactionIdentifier.Hash
	return &account.SendTxResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "SendTx success",
		TxHash: hash,
	}, nil
}

func (c *ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	resp, err := c.dataClient.GetTxByAddress(req.Address)
	if err != nil {
		log.Error("GetTxByAddress fail", "err", err)
		return &account.TxAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "GetTxByAddress fail",
		}, nil
	}
	txs := resp.Transactions
	list := make([]*account.TxMessage, 0, len(txs))
	for i := 0; i < len(txs); i++ {
		//only do Transfer
		if txs[i].Method != "staking.Transfer" {
			continue
		}
		list = append(list, &account.TxMessage{
			Hash:   txs[i].Hash,
			Tos:    []*account.Address{{Address: txs[i].Body.To}},
			Froms:  []*account.Address{{Address: txs[i].Sender}},
			Fee:    txs[i].Fee,
			Status: account.TxStatus_Success,
			Values: []*account.Value{{Value: txs[i].Body.Amount}},
			Type:   1,
			Height: strconv.Itoa(txs[i].Block),
		})
	}
	fmt.Println("resp", resp)
	return &account.TxAddressResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get tx list success",
		Tx:   list,
	}, nil
}

func (c *ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	resp, err := c.dataClient.GetTxByHash(req.Hash)
	if err != nil {
		log.Error("GetTxByHash fail", "err", err)
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "GetTxByHash fail",
		}, nil
	}
	txs := resp.Transactions
	list := make([]*account.TxMessage, 0, len(txs))
	for i := 0; i < len(txs); i++ {
		//only do Transfer
		if txs[i].Method != "staking.Transfer" {
			continue
		}
		list = append(list, &account.TxMessage{
			Hash:   txs[i].Hash,
			Tos:    []*account.Address{{Address: txs[i].Body.To}},
			Froms:  []*account.Address{{Address: txs[i].Sender}},
			Fee:    txs[i].Fee,
			Status: account.TxStatus_Success,
			Values: []*account.Value{{Value: txs[i].Body.Amount}},
			Type:   1,
			Height: strconv.Itoa(txs[i].Block),
		})
	}
	fmt.Println("resp", resp)
	return &account.TxHashResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get tx  success",
		Tx:   list[0],
	}, nil
}

func (c *ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	panic("Does not support loop calls GetBlockByNumber")
}

func (c *ChainAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	panic("not support")
}

func (c *ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	panic("not support")
}

func (c *ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	panic("not support")
}

func (c *ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	panic("not support")
}

func (c *ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	panic("not support")
}

func (c *ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	//can call grpc     estimateGas(nic: client.NodeInternal, signer: Uint8Array): Promise<types.longnum>;
	panic("not support")
}

func (c *ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	panic("not support")
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	panic("not support")
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	panic("not support")
}
