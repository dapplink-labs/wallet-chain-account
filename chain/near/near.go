package near

// NEAR JSON RPC交互示例（官方接口封装）
import (
	"github.com/ethereum/go-ethereum/log"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
)

type NearAdaptor struct {
	nearClient *NearClient
}

func NewSuiAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	client, err := NewNearClient(conf)
	if err != nil {
		log.Error("Init Sui Client err", "err", err)
		return nil, err
	}
	return &NearAdaptor{
		nearClient: client,
	}, nil
}

// BuildSignedTransaction implements chain.IChainAdaptor.
func (n *NearAdaptor) BuildSignedTransaction(req *account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	panic("unimplemented")
}

// ConvertAddress implements chain.IChainAdaptor.
func (n *NearAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	panic("unimplemented")
}

// CreateUnSignTransaction implements chain.IChainAdaptor.
func (n *NearAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	panic("unimplemented")
}

// DecodeTransaction implements chain.IChainAdaptor.
func (n *NearAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	panic("unimplemented")
}

// GetAccount implements chain.IChainAdaptor.
func (n *NearAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	panic("unimplemented")
}

// GetBlockByHash implements chain.IChainAdaptor.
func (n *NearAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	panic("unimplemented")
}

// GetBlockByNumber implements chain.IChainAdaptor.
func (n *NearAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	panic("unimplemented")
}

// GetBlockByRange implements chain.IChainAdaptor.
func (n *NearAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	panic("unimplemented")
}

// GetBlockHeaderByHash implements chain.IChainAdaptor.
func (n *NearAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	panic("unimplemented")
}

// GetBlockHeaderByNumber implements chain.IChainAdaptor.
func (n *NearAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	panic("unimplemented")
}

// GetExtraData implements chain.IChainAdaptor.
func (n *NearAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	panic("unimplemented")
}

// GetFee implements chain.IChainAdaptor.
func (n *NearAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	panic("unimplemented")
}

// GetSupportChains implements chain.IChainAdaptor.
func (n *NearAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	panic("unimplemented")
}

// GetTxByAddress implements chain.IChainAdaptor.
func (n *NearAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	panic("unimplemented")
}

// GetTxByHash implements chain.IChainAdaptor.
func (n *NearAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	panic("unimplemented")
}

// SendTx implements chain.IChainAdaptor.
func (n *NearAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	panic("unimplemented")
}

// ValidAddress implements chain.IChainAdaptor.
func (n *NearAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	panic("unimplemented")
}

// VerifySignedTransaction implements chain.IChainAdaptor.
func (n *NearAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	panic("unimplemented")
}
