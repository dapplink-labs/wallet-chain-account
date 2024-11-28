package sui

import (
	"encoding/hex"
	"fmt"

	"math/big"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/crypto/blake2b"

	"github.com/block-vision/sui-go-sdk/models"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

const (
	ChainName   = "Sui"
	SuiCoinType = "0x2::sui::SUI"

	PublicKeySize    = 32
	SuiAddressLength = 32
)

var SIGNATURE_SCHEME_TO_FLAG = map[string]byte{
	"ED25519": 0x00,
}

type SuiAdaptor struct {
	suiClient *SuiClient
}

func NewSuiAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	client, err := NewSuiClient(conf)
	if err != nil {
		log.Error("Init Sui Client err", "err", err)
		return nil, err
	}
	return &SuiAdaptor{
		suiClient: client,
	}, nil
}

func (s *SuiAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	return &account.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

func (s *SuiAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	publicKey, err := hex.DecodeString(req.PublicKey)
	if err != nil {
		log.Error("hex decode err", "err", err)
		return &account.ConvertAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "req decode error ",
		}, nil
	}
	newPubkey := []byte{byte(0x00)}
	newPubkey = append(newPubkey, publicKey...)
	addrBytes := blake2b.Sum256(newPubkey)
	address := fmt.Sprintf("0x%s", hex.EncodeToString(addrBytes[:])[:64])
	return &account.ConvertAddressResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "create address success",
		Address: address,
	}, nil
}

func (s *SuiAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	if len(req.Address) != 66 || !strings.HasPrefix(req.Address, "0x") {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_ERROR,
			Msg:   "invalid address",
			Valid: false,
		}, nil
	}
	ok := regexp.MustCompile("^[0-9a-fA-F]{64}$").MatchString(req.Address[2:])
	if ok {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_SUCCESS,
			Msg:   "valid address",
			Valid: true,
		}, nil
	} else {
		return &account.ValidAddressResponse{
			Code:  common2.ReturnCode_SUCCESS,
			Msg:   "invalid address",
			Valid: false,
		}, nil
	}
}

// not nessary
func (s *SuiAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	//TODO implement me
	panic("implement me")
}

// not nessary
func (s *SuiAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SuiAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SuiAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SuiAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	balanceRes, err := s.suiClient.GetAccountBalance(SuiCoinType, req.Address)
	if err != nil {
		log.Error("get balance err", "err", err)
		return &account.AccountResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Get balance err",
		}, err
	}
	//nonceResult, err := c.ethClient.TxCountByAddress(common.HexToAddress(req.Address))
	//if err != nil {
	//	log.Error("get nonce by address fail", "err", err)
	//	return &account.AccountResponse{
	//		Code: common2.ReturnCode_ERROR,
	//		Msg:  "get nonce by address fail",
	//	}, nil
	//}

	log.Info("balance result", "balance=", balanceRes.TotalBalance)
	return &account.AccountResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "get account success",
		Network: ChainName,
		Balance: balanceRes.TotalBalance,
	}, nil
}

func (s *SuiAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	price, err := s.suiClient.GetGasPrice()
	if err != nil {
		return &account.FeeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get gas price err",
		}, nil
	}
	return &account.FeeResponse{
		Code:      common2.ReturnCode_SUCCESS,
		Msg:       "get gas price success",
		SlowFee:   fmt.Sprintf("%d", price) + "|",
		NormalFee: fmt.Sprintf("%d", price) + "|" + "*2",
		FastFee:   fmt.Sprintf("%d", price) + "|" + "*3",
	}, nil

}

func (s *SuiAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	_, err := s.suiClient.SendTx(req.RawTx)
	if err != nil {
		return &account.SendTxResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	return &account.SendTxResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "send tx success",
		TxHash: "",
	}, nil
}

func (s *SuiAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	cursor := req.Cursor
	txList, err := s.suiClient.GetTxListByAddress(req.Address, cursor, req.Pagesize)
	if err != nil {
		return &account.TxAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get transactions fail",
		}, err
	}
	// todo sui 专有的交易结构，直接放到value中，前端自定义获取解析
	var tx_list []*account.TxMessage
	for _, tx := range txList.Data {
		message, _ := s.getTxMessage(tx)
		tx_list = append(tx_list, message)
	}
	return &account.TxAddressResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get transactions success",
		Tx:   tx_list,
	}, nil
}

func (s *SuiAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	txDetail, err := s.suiClient.GetTxDetailByDigest(req.Hash)
	if err != nil {
		return &account.TxHashResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get transaction fail",
		}, err
	}

	message, _ := s.getTxMessage(txDetail)

	return &account.TxHashResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get transaction success",
		Tx:   message,
	}, nil
}

func (a *SuiAdaptor) getTxMessage(suiTransaction models.SuiTransactionBlockResponse) (*account.TxMessage, error) {
	var from_addrs []*account.Address
	var to_addrs []*account.Address
	var value_list []*account.Value
	totalAmount := big.NewInt(0)
	toAmount := big.NewInt(0)
	for _, bc := range suiTransaction.BalanceChanges {
		if bc.Owner.AddressOwner != "" {
			from_addrs = append(from_addrs, &account.Address{Address: bc.Owner.AddressOwner})
			totalAmount = new(big.Int).Add(totalAmount, stringToInt(bc.Amount))
		} else {
			to_addrs = append(to_addrs, &account.Address{Address: bc.Owner.ObjectOwner})
			toAmount = new(big.Int).Add(toAmount, stringToInt(bc.Amount))
			value_list = append(value_list, &account.Value{Value: bc.Amount})
		}
	}
	totalAmount = new(big.Int).Abs(totalAmount)
	fee := new(big.Int).Sub(totalAmount, toAmount).String()
	return &account.TxMessage{
		Hash:     suiTransaction.Digest,
		Height:   suiTransaction.Checkpoint,
		Status:   account.TxStatus_Success,
		Type:     0,
		Datetime: suiTransaction.TimestampMs,
		Froms:    from_addrs,
		Tos:      to_addrs,
		Values:   value_list,
		Fee:      fee,
	}, nil
}

func (s *SuiAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SuiAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SuiAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SuiAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SuiAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SuiAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	//TODO implement me
	panic("implement me")
}

func stringToInt(amount string) *big.Int {
	log.Info("string to Int", "amount", amount)
	intAmount, success := big.NewInt(0).SetString(amount, 0)
	if !success {
		return nil
	}
	return intAmount
}
