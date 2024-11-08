package solana

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mr-tron/base58"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/dapplink-labs/wallet-chain-account/config"
)

const (
	defaultDialTimeout    = 5 * time.Second
	defaultDialAttempts   = 5
	defaultRequestTimeout = 10 * time.Second
)

type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
type JsonRpcResponse[T any] struct {
	JsonRpc string        `json:"jsonrpc"`
	Id      uint64        `json:"id"`
	Result  T             `json:"result"`
	Error   *JsonRpcError `json:"error,omitempty"`
}

type Block struct {
	Blockhash         string
	BlockTime         *time.Time
	BlockHeight       *int64
	PreviousBlockhash string
	ParentSlot        uint64
	Transactions      []BlockTransaction
	Signatures        []string
}
type BlockTransaction struct {
	// rpc fields
	Meta        *TransactionMeta
	Transaction types.Transaction
}
type TransactionMeta struct {
	Err               any
	Fee               uint64
	PreBalances       []int64
	PostBalances      []int64
	PreTokenBalances  []rpc.TransactionMetaTokenBalance
	PostTokenBalances []rpc.TransactionMetaTokenBalance
	LogMessages       []string

	LoadedAddresses      rpc.TransactionLoadedAddresses
	ComputeUnitsConsumed *uint64
}

type SolClient interface {
	NewSolanaClients(conf *config.Config) (*SolanaClient, error)
	GetLatestBlockHeight() (int64, error)
	GetBalance(address string) (string, error)
	GetTxByHash(hash string) (*TxMessage, error)
	getPreTokenBalance(preTokenBalance []rpc.TransactionMetaTokenBalance, accountIndex uint64) *rpc.TransactionMetaTokenBalance
	RequestAirdrop(address string) (string, error)
	SendTx(rawTx string) (string, error)
	GetNonce(NonceAccount string) (string, error)
	GetMinRent() (string, error)
	GetFee(nonceAccountAddr string) (string, error)
	Close()
}
type SolanaClient struct {
	RpcClient    rpc.RpcClient
	Client       *client.Client
	solanaConfig config.WalletNode
}

type TransactionList struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Hash  string `json:"hash"`
	Value string `json:"value"`
}

type RpcBlock struct {
	Hash         common.Hash       `json:"hash"`
	Height       uint64            `json:"height"`
	Transactions []TransactionList `json:"transactions"`
	BaseFee      string            `json:"baseFeePerGas"`
}

func NewSolanaClients(conf *config.Config) (*SolanaClient, error) {
	endpoint := conf.WalletNode.Sol.RpcUrl
	rpcClient := rpc.NewRpcClient(endpoint)
	clientNew := client.NewClient(endpoint)
	return &SolanaClient{
		RpcClient:    rpcClient,
		Client:       clientNew,
		solanaConfig: conf.WalletNode,
	}, nil
}

func (sol *SolanaClient) GetLatestBlockHeight() (int64, error) {
	res, err := sol.RpcClient.GetBlockHeight(context.Background())
	if err != nil {
		return 0, err
	}
	return int64(res.Result), nil
}

func (sol *SolanaClient) GetBalance(address string) (string, error) {
	balance, err := sol.Client.GetBalanceWithConfig(
		context.TODO(),
		address,
		client.GetBalanceConfig(rpc.GetBalanceConfig{
			Commitment: rpc.CommitmentProcessed,
		}),
	)
	if err != nil {
		return "", err
	}

	var lamportsOnAccount = new(big.Float).SetUint64(balance)

	var solBalance = new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(1000000000))

	return solBalance.String(), nil
}

type GetTxByAddressRes struct {
	Data []GetTxByAddressTx
}

type GetTxByAddressTx struct {
	ID                  string `json:"_id"`
	Src                 string `json:"src"`
	Dst                 string `json:"dst"`
	Lamport             int    `json:"lamport"`
	BlockTime           int    `json:"blockTime"`
	Slot                int    `json:"slot"`
	TxHash              string `json:"txHash"`
	Fee                 int    `json:"fee"`
	Status              string `json:"status"`
	Decimals            int    `json:"decimals"`
	TxNumberSolTransfer int    `json:"txNumberSolTransfer"`
}

func (sol *SolanaClient) GetAccount() (string, string, error) {
	account := types.NewAccount()
	address := account.PublicKey.ToBase58()
	private := base58.Encode(account.PrivateKey)
	return address, private, nil
}

type Header struct {
	NumReadonlySignedAccounts   int `json:"numReadonlySignedAccounts"`
	NumReadonlyUnsignedAccounts int `json:"numReadonlyUnsignedAccounts"`
	NumRequiredSignatures       int `json:"numRequiredSignatures"`
}
type Instructions struct {
	Accounts       []int  `json:"accounts"`
	Data           string `json:"data"`
	ProgramIDIndex int    `json:"programIdIndex"`
}
type Message struct {
	AccountKeys     []string       `json:"accountKeys"`
	Header          Header         `json:"header"`
	Instructions    []Instructions `json:"instructions"`
	RecentBlockhash string         `json:"recentBlockhash"`
}
type Transaction struct {
	Message    Message  `json:"message"`
	Signatures []string `json:"signatures"`
}

type TxMessage struct {
	Hash   string
	From   string
	To     string
	Fee    string
	Status bool
	Value  string
	Type   int32
	Height string
}

func (sol *SolanaClient) GetTxByHash(hash string) (*TxMessage, error) {
	out, err := sol.RpcClient.GetTransaction(
		context.Background(),
		hash,
	)
	fmt.Printf("%#v\n", out)
	if err != nil {
		log.Fatalf("failed to request airdrop, err: %v", err)
		return nil, err
	}
	if out.Result == nil {
		return nil, fmt.Errorf("out or out.Result is nil")
	}

	transaction, ok := out.Result.Transaction.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Transaction is not of type map[string]interface{}")
	}

	message, exists := transaction["message"]
	if !exists {
		return nil, fmt.Errorf("key 'message' does not exist in transaction")
	}
	accountKeys := message.((map[string]interface{}))["accountKeys"].([]interface{})
	signatures := out.Result.Transaction.(map[string]interface{})["signatures"].([]interface{})
	log.Println("out:", out.Result)
	_hash := signatures[0]
	if out.Result.Meta.Err != nil || len(out.Result.Meta.LogMessages) == 0 || _hash == "" {
		//log.Fatalf("not found tx, err: %v", err)

		return nil, err
	}
	var txMessage []*TxMessage
	for i := 0; i < len(accountKeys); i++ {
		to := accountKeys[i].(string)
		amount := out.Result.Meta.PostBalances[i] - out.Result.Meta.PreBalances[i]

		if to != "" && amount > 0 {
			txMessage = append(txMessage, &TxMessage{
				Hash:   hash,
				From:   "",
				To:     to,
				Fee:    strconv.FormatUint(out.Result.Meta.Fee, 10),
				Status: true,
				Value:  strconv.FormatInt(amount, 10),
				Type:   1,
				Height: strconv.FormatUint(out.Result.Slot, 10),
			})
		}
	}

	for i := 0; i < len(out.Result.Meta.PostTokenBalances); i++ {
		postToken := out.Result.Meta.PostTokenBalances[i]

		preTokenBalance := getPreTokenBalance(out.Result.Meta.PreTokenBalances, postToken.AccountIndex)
		if preTokenBalance == nil {
			continue
		}
		postAmount, _ := strconv.ParseFloat(postToken.UITokenAmount.Amount, 64)
		preAmount, _ := strconv.ParseFloat(preTokenBalance.UITokenAmount.Amount, 64)
		amount := postAmount - preAmount
		if amount > 0 {
			txMessage = append(txMessage, &TxMessage{
				Hash:   hash,
				From:   "",
				To:     postToken.Owner,
				Fee:    strconv.FormatUint(out.Result.Meta.Fee, 10),
				Status: true,
				Value:  strconv.FormatFloat(amount, 'E', -1, 10),
				Type:   1,
				Height: strconv.FormatUint(out.Result.Slot, 10),
			})
		}
	}
	if len(txMessage) > 0 {
		return txMessage[0], nil
	}

	return &TxMessage{
		Hash:   hash,
		From:   "",
		To:     "",
		Fee:    "",
		Status: false,
		Value:  "",
		Type:   0,
		Height: "",
	}, nil
}

func getPreTokenBalance(preTokenBalance []rpc.TransactionMetaTokenBalance, accountIndex uint64) *rpc.TransactionMetaTokenBalance {
	for j := 0; j < len(preTokenBalance); j++ {
		preToken := preTokenBalance[j]
		if preToken.AccountIndex == accountIndex {
			return &preTokenBalance[j]
		}
	}
	return nil
}
func (sol *SolanaClient) RequestAirdrop(address string) (string, error) {
	c := client.NewClient(rpc.DevnetRPCEndpoint)
	sig, err := c.RequestAirdrop(
		context.TODO(),
		address,
		1e9, // lamports (1 SOL = 10^9 lamports)
	)
	if err != nil {
		log.Fatalf("failed to request airdrop, err: %v", err)
		return "", err
	}
	return sig, nil
}

func (sol *SolanaClient) SendTx(rawTx string) (string, error) {
	res, err := sol.RpcClient.SendTransactionWithConfig(
		context.Background(),
		base64.StdEncoding.EncodeToString([]byte(rawTx)),
		rpc.SendTransactionConfig{
			Encoding: rpc.SendTransactionConfigEncodingBase64,
		},
	)
	if err != nil {
		return "", err
	}
	if res.Error != nil {
		return "", res.Error
	}
	return res.Result, nil
}

func (sol *SolanaClient) GetNonce(NonceAccount string) (string, error) {
	//NonceAccount, err := sol.Client.GetNonceAccount(context.Background(), AccountAddr)

	nonce, err := sol.Client.GetNonceFromNonceAccount(context.Background(), NonceAccount)
	println("nonce:", nonce)
	if err != nil {
		log.Fatalf("failed to get nonce account, err: %v", err)
		return "", err
	}
	return nonce, nil
}

func (sol *SolanaClient) Getfee(AccountAddr string, NonceAccount string) (string, error) {

	panic("implement me")
}
func (sol *SolanaClient) GetMinRent() (string, error) {
	bal, err := sol.RpcClient.GetMinimumBalanceForRentExemption(context.Background(), 100)
	if err != nil {
		log.Fatalf("failed to get GetMinimumBalanceForRentExemption , err: %v", err)
		return "", err
	}
	return strconv.FormatUint(bal.Result, 10), nil
}

func (sol *SolanaClient) BlockByNumber(hash string) (*JsonRpcResponse[*Block], error) {

	return nil, nil
}