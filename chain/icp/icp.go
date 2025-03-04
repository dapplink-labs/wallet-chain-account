package icp

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"github.com/ethereum/go-ethereum/log"
	"strconv"
	"strings"

	"github.com/aviate-labs/agent-go/principal"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

const (
	ChainName = "ICP"
)

var DefaultSubaccount = [32]byte{}

type ChainAdaptor struct {
	icpClient *Client
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	icpClient, err := NewIcpClient(context.Background(), conf.WalletNode.Icp.RpcUrl, conf.WalletNode.Icp.TimeOut)
	if err != nil {
		log.Error("new icp client failed:", err)
		return nil, err
	}

	return &ChainAdaptor{
		icpClient: icpClient,
	}, nil
}

func (c ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	return &account.SupportChainsResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

func (c ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	publicKeyBytes, _ := hex.DecodeString(req.PublicKey)
	p := principal.NewSelfAuthenticating(publicKeyBytes)
	a := principal.NewAccountID(p, DefaultSubaccount)
	return &account.ConvertAddressResponse{
		Code: common2.ReturnCode_SUCCESS, Msg: "convert address success",
		Address: a.String(),
	}, nil
}

func (c ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	_, err := principal.DecodeAccountID(req.Address)
	if err != nil {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_ERROR,
			Msg:   "invalid address",
			Valid: false,
		}, nil
	} else {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_SUCCESS,
			Msg:   "valid address",
			Valid: true,
		}, nil
	}
}

func (c ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	block, err := c.icpClient.GetBlockByNumber(req.Height)
	if err != nil {
		log.Error("get balance err", "err", err)
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Get blockByNumber err",
		}, err
	}
	height := block.Block.BlockIdentifier.Index
	hash := block.Block.BlockIdentifier.Hash
	var blockTxList []*account.BlockInfoTransactionList
	for _, tx := range block.Block.Transactions {
		blockTransactionList, _ := convertBlockTransactionToBlockTransactionList(tx)
		blockTxList = append(blockTxList, blockTransactionList)
	}
	return &account.BlockResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "get block success",
		Height:  height,
		Hash:    hash,
		BaseFee: "10000",
	}, nil
}

func (c ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	block, err := c.icpClient.GetBlockByHash(req.Hash)
	if err != nil {
		log.Error("get balance err", "err", err)
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Get blockByNumber err",
		}, err
	}
	height := block.Block.BlockIdentifier.Index
	hash := block.Block.BlockIdentifier.Hash
	var blockTxList []*account.BlockInfoTransactionList
	for _, tx := range block.Block.Transactions {
		blockTransactionList, _ := convertBlockTransactionToBlockTransactionList(tx)
		blockTxList = append(blockTxList, blockTransactionList)
	}
	return &account.BlockResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "get block success",
		Height:  height,
		Hash:    hash,
		BaseFee: "10000",
	}, nil
}

func (c ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	balance, err := c.icpClient.GetAccountBalance(req.Address)
	if err != nil {
		log.Error("get balance err", "err", err)
		return &account.AccountResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Get balance err",
		}, err
	}
	var stringData = ""
	if req.Coin == "" {
		bytes, _ := json.Marshal(&balance.Balances)
		stringData = string(bytes)
	} else {
		for _, b := range balance.Balances {
			if strings.EqualFold(req.Coin, b.Currency.Symbol) {
				stringData = b.Value
			}
		}
	}
	log.Info("balance result", "balance=", balance.Balances)
	return &account.AccountResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "get account success",
		Network: ChainName,
		Balance: stringData,
	}, nil
}

func (c ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	return &account.FeeResponse{
		Code:      common2.ReturnCode_SUCCESS,
		Msg:       "get fee success",
		SlowFee:   "10000",
		NormalFee: "10000",
		FastFee:   "10000",
	}, nil
}

func (c ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	txList, err := c.icpClient.GetTxByAddress(req.Address)
	if err != nil {
		log.Error("get tx err", "err", err)
		return &account.TxAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Get tx err",
		}, err
	}
	var txs []*account.TxMessage
	for _, t := range txList.Transactions {
		txMessage, _ := convertTransactionToTxMessage(t)
		txs = append(txs, txMessage)
	}
	return &account.TxAddressResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get txByAddress success",
		Tx:   txs,
	}, nil
}

func (c ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	transaction, err := c.icpClient.GetTxByHash(req.Hash)
	if err != nil {
		log.Error("get tx err", "err", err)
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Get tx err",
		}, err
	}
	if transaction.TotalCount != 1 {
		log.Error("get tx err", "err", err)
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Get tx err",
		}, nil
	}
	tx, err := convertTransactionToTxMessage(transaction.Transactions[0])
	if err != nil {
		log.Error("get tx err", "err", err)
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Get tx err",
		}, err
	}
	return &account.TxHashResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get txByHash success",
		Tx:   tx,
	}, nil
}

func (c ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	//TODO implement me
	panic("implement me")
}

func convertTransactionToTxMessage(transaction *types.BlockTransaction) (*account.TxMessage, error) {
	height := transaction.BlockIdentifier.Index
	hash := transaction.BlockIdentifier.Hash
	var fee string
	var fromAddrs []*account.Address
	var toAddrs []*account.Address
	var valueList []*account.Value
	for _, operation := range transaction.Transaction.Operations {
		if strings.EqualFold(operation.Type, "FEE") {
			fee = operation.Amount.Value
		} else if strings.EqualFold(operation.Type, "TRANSACTION") {
			amount, _ := strconv.Atoi(operation.Amount.Value)
			if amount > 0 {
				toAddrs = append(toAddrs, &account.Address{Address: operation.Account.Address})
				valueList = append(valueList, &account.Value{Value: operation.Amount.Value})
			} else if amount < 0 {
				fromAddrs = append(fromAddrs, &account.Address{Address: operation.Account.Address})
			}
		}
	}
	return &account.TxMessage{
		Hash:   hash,
		Index:  uint32(height),
		Froms:  fromAddrs,
		Tos:    toAddrs,
		Values: valueList,
		Fee:    fee,
		Status: account.TxStatus_Success,
		Type:   0,
		Height: strconv.FormatInt(height, 10),
	}, nil
}

func convertBlockTransactionToBlockTransactionList(transaction *types.Transaction) (*account.BlockInfoTransactionList, error) {
	return nil, nil
}
