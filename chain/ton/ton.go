package ton

import (
	"context"
	"github.com/ethereum/go-ethereum/log"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

type ChainAdaptor struct {
	tonClient *TonClient
}

// GetSupportChains 返回是否支持链
func (c *ChainAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	return &account.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

// ConvertAddress 将公钥转换为地址
func (c *ChainAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	address, err := wallet.AddressFromPubKey(req.PublicKey, wallet.V4R2, 0)
	if err != nil {
		return &account.ConvertAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	} else {
		return &account.ConvertAddressResponse{
			Code:    common2.ReturnCode_SUCCESS,
			Msg:     "convert address successs",
			Address: address.String(),
		}, nil
	}
}

// ValidAddress 验证地址
func (c *ChainAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	_, err := address.ParseAddr(req.Address)
	return &account.ValidAddressResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "convert address successs",
		Valid: err == nil,
	}, nil
}

// GetBlockByNumber 根据区块高度获取区块
func (c *ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	block, err := c.tonClient.GetBlockByNumber(uint32(req.Height))
	if err != nil {
		log.Error("block by number error", err)
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	} else {
		// 提取交易列表
		return &account.BlockResponse{
			Code: common2.ReturnCode_SUCCESS,
			Msg:  "get block success",
			//Hash:         block.BlockInfo.,
			//BaseFee:      block.BlockInfo.,
			//Transactions: block.BlockInfo.,
		}, err
	}
	//c.tonClient.rw.Lock()
	//defer c.tonClient.rw.Unlock()
	//return c.tonClient.GetBlockByNumber(req.Number)
}

func (c *ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ChainAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	tonClient, err := DialTonClient(context.Background(), conf.WalletNode.Eth.RPCs[0].RPCURL)
	if err != nil {
		return nil, err
	}
	return &ChainAdaptor{
		tonClient: tonClient,
	}, nil
}
