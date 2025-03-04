package near

// NEAR JSON RPC交互示例（官方接口封装）
import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
	"github.com/eteu-technologies/borsh-go"
	uint128 "github.com/eteu-technologies/golang-uint128"
	nearClient "github.com/eteu-technologies/near-api-go/pkg/client"
	"github.com/eteu-technologies/near-api-go/pkg/client/block"
	"github.com/eteu-technologies/near-api-go/pkg/types"
	"github.com/eteu-technologies/near-api-go/pkg/types/action"
	"github.com/eteu-technologies/near-api-go/pkg/types/key"
	"github.com/eteu-technologies/near-api-go/pkg/types/transaction"
	"github.com/ethereum/go-ethereum/log"
	"strconv"
	"sync"

	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
)

const ChainName = "Near"

// https://pkg.go.dev/github.com/aurora-is-near/near-api-go
type NearAdaptor struct {
	nearClient *nearClient.Client
}

func NewNearAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
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
	// 1.encode params get extract params
	ctx := context.Background()
	nearTransaction, err := n.buildDynamicFeeTx(req.Base64Tx)
	if err != nil {
		return nil, err
	}
	kp, err := key.NewBase58KeyPair(nearTransaction.PrivateKey)

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	ctx = contextWithKeyPair(context.Background(), kp)
	parseUint, err := strconv.ParseUint(nearTransaction.Amount, 10, 64)
	var Balance = types.Balance(uint128.From64(parseUint))
	actions := []action.Action{
		action.NewTransfer(Balance),
	}
	// 2. build transaction and sign
	keypair, transactionData, err := n.PrepareTransaction(ctx, nearTransaction.Sender, nearTransaction.Receiver, actions)

	// 3. sign transcation
	blob, err := n.SignTransaction(keypair, transactionData)
	return &account.SignedTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "build transaction success",
		SignedTx: blob,
	}, nil
}

// ConvertAddress implements chain.IChainAdaptor.
func (n *NearAdaptor) ConvertAddress(req *account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	// 用户可以自定义地址，这个方法无用
	panic("unimplemented")
}

// CreateUnSignTransaction implements chain.IChainAdaptor.
func (n *NearAdaptor) CreateUnSignTransaction(req *account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	// 1.encode params
	ctx := context.Background()
	nearTransaction, err := n.buildDynamicFeeTx(req.Base64Tx)
	if err != nil {
		return nil, err
	}
	kp, err := key.NewBase58KeyPair(nearTransaction.PrivateKey)

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	ctx = contextWithKeyPair(context.Background(), kp)
	parseUint, err := strconv.ParseUint(nearTransaction.Amount, 10, 64)
	var Balance = types.Balance(uint128.From64(parseUint))
	actions := []action.Action{
		action.NewTransfer(Balance),
	}
	// 2. build transaction and sign
	_, transactionData, err := n.PrepareTransaction(ctx, nearTransaction.Sender, nearTransaction.Receiver, actions)

	// 3.Serialize into Borsh
	_, serialized, err := transactionData.Hash()
	return &account.UnSignTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "create un sign transaction success",
		UnSignTx: string(serialized),
	}, nil
}

// DecodeTransaction implements chain.IChainAdaptor.
func (n *NearAdaptor) DecodeTransaction(req *account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	var t transaction.Transaction
	borsh.Deserialize(&t, []byte(req.RawTx))
	marshal, err := json.Marshal(&t)
	if err != nil {
		log.Error(err.Error())
	}
	return &account.DecodeTransactionResponse{
		Code:     common2.ReturnCode_SUCCESS,
		Msg:      "decode transaction success",
		Base64Tx: base64.StdEncoding.EncodeToString(marshal),
	}, nil
}

// GetAccount implements chain.IChainAdaptor.
func (n *NearAdaptor) GetAccount(req *account.AccountRequest) (*account.AccountResponse, error) {
	ctx := context.Background()
	res, err := n.nearClient.AccountView(ctx, req.Address, block.FinalityFinal())
	if err != nil {
		log.Error("GetAccount err", "err", err)
		return nil, err
	}
	var accountData map[string]interface{}
	json.Unmarshal(res.Result, &accountData)
	amount, ok := accountData["amount"].(string)
	if !ok {
		log.Error("amount not found or not a string")
		return nil, err
	}
	return &account.AccountResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "get account success",
		Network: ChainName,
		Balance: amount,
	}, nil
}

// GetBlockByHash implements chain.IChainAdaptor.
func (n *NearAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	ctx := context.Background()
	var hashOrFinality block.BlockCharacteristic
	if req.Hash != "" {
		hashOrFinality = block.BlockHashRaw(req.Hash)
	} else {
		hashOrFinality = block.FinalityFinal()
	}
	blockView, err := n.nearClient.BlockDetails(ctx, hashOrFinality)
	if err != nil {
		log.Error("GetBlockByHash err", "err", err)
		return nil, err
	}
	chunks := blockView.Chunks
	var txListRet []*account.BlockInfoTransactionList
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(len(chunks))
	for _, chunk := range chunks {
		go func() {
			defer wg.Done()
			chunkDetailData, err := n.nearClient.ChunkDetails(ctx, chunk.ChunkHash)
			if err != nil {
				log.Error("GetBlockByHash err", "err", err)
				return
			}
			var chunkDetail ChunkDetail
			err = json.Unmarshal(chunkDetailData.Result, &chunkDetail)
			if err != nil {
				log.Error("Unmarshal err", "err", err)
				return
			}
			for _, tx := range chunkDetail.Transactions {
				for _, action := range tx.Actions {
					transfer := action.Transfer
					// only about transfer data
					if transfer != nil {
						bitlItem := &account.BlockInfoTransactionList{
							From:           tx.SignerID,
							To:             tx.ReceiverID,
							TokenAddress:   tx.SignerID,
							ContractWallet: tx.SignerID,
							Hash:           tx.Hash,
							Height:         uint64(chunkDetail.Header.HeightCreated),
							Amount:         transfer.Deposit,
						}
						mu.Lock()
						txListRet = append(txListRet, bitlItem)
						mu.Unlock()
					}
				}
			}
		}()
	}
	wg.Wait()
	return &account.BlockResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "block by hash success",
		Height:       int64(blockView.Header.Height),
		Hash:         blockView.Header.Hash.String(),
		BaseFee:      blockView.Header.GasPrice.String(),
		Transactions: txListRet,
	}, nil
}

// GetBlockByNumber implements chain.IChainAdaptor.
func (n *NearAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	ctx := context.Background()
	var hashOrFinality block.BlockCharacteristic
	if req.Height != 0 {
		hashOrFinality = block.BlockID(uint(req.Height))
	} else {
		hashOrFinality = block.FinalityFinal()
	}
	blockView, err := n.nearClient.BlockDetails(ctx, hashOrFinality)
	if err != nil {
		log.Error("GetBlockByHash err", "err", err)
	}
	chunks := blockView.Chunks
	var txListRet []*account.BlockInfoTransactionList
	for _, chunk := range chunks {
		chunkDetailData, err := n.nearClient.ChunkDetails(ctx, chunk.ChunkHash)
		if err != nil {
			log.Error("ChunkDetails err", "err", err)
		}
		var chunkDetail ChunkDetail
		err = json.Unmarshal(chunkDetailData.Result, &chunkDetail)
		if err != nil {
			log.Error("Unmarshal err", "err", err)
			continue
		}
		for _, tx := range chunkDetail.Transactions {
			for _, action := range tx.Actions {
				transfer := action.Transfer
				// only about transfer data
				if transfer != nil {
					bitlItem := &account.BlockInfoTransactionList{
						From:           tx.SignerID,
						To:             tx.ReceiverID,
						TokenAddress:   tx.SignerID,
						ContractWallet: tx.SignerID,
						Hash:           tx.Hash,
						Height:         uint64(chunkDetail.Header.HeightCreated),
						Amount:         transfer.Deposit,
					}
					txListRet = append(txListRet, bitlItem)
				}
			}
		}
	}
	return &account.BlockResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "block by hash success",
		Height:       int64(blockView.Header.Height),
		Hash:         blockView.Header.Hash.String(),
		BaseFee:      blockView.Header.GasPrice.String(),
		Transactions: txListRet,
	}, nil
}

// GetBlockByRange implements chain.IChainAdaptor.
func (n *NearAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	start := req.Start
	// Convert string to int
	startNum, err := strconv.Atoi(start)
	if err != nil {
		fmt.Println("startNum Error converting string to int:", err)
		return nil, err
	}
	end := req.End
	endNum, err := strconv.Atoi(end)
	if err != nil {
		fmt.Println("endNum Error converting string to int:", err)
		return nil, err
	}
	ctx := context.Background()
	blockHeaderList := make([]*account.BlockHeader, 0, endNum-startNum)
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	wg.Add(endNum - startNum)
	for i := startNum; i < endNum; i++ {
		go func() {
			defer wg.Done()
			var hashOrFinality block.BlockCharacteristic
			hashOrFinality = block.BlockID(uint(i))
			blockView, err := n.nearClient.BlockDetails(ctx, hashOrFinality)
			if err != nil {
				log.Error("GetBlockByRange err", "err", err)
				return
			}
			blockItem := &account.BlockHeader{
				ParentHash:       blockView.Header.PrevHash.String(),
				UncleHash:        "",
				CoinBase:         "",
				Root:             blockView.Header.BlockMerkleRoot.String(),
				TxHash:           blockView.Header.ChunkTxRoot.String(),
				ReceiptHash:      blockView.Header.ChunkReceiptsRoot.String(),
				ParentBeaconRoot: "",
				Difficulty:       "",
				Number:           strconv.FormatUint(blockView.Header.Height, 10),
				GasLimit:         0,
				GasUsed:          0,
				Time:             blockView.Header.Timestamp,
				Extra:            "",
				MixDigest:        "",
				Nonce:            "",
				BaseFee:          "",
				WithdrawalsHash:  blockView.Header.PrevStateRoot.String(),
				BlobGasUsed:      0,
				ExcessBlobGas:    0,
			}
			mutex.Lock()
			blockHeaderList = append(blockHeaderList, blockItem)
			mutex.Unlock()
		}()

	}
	wg.Wait()
	return &account.BlockByRangeResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "get block range success",
		BlockHeader: blockHeaderList,
	}, nil
}

// GetBlockHeaderByHash implements chain.IChainAdaptor.
func (n *NearAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	ctx := context.Background()
	var hashOrFinality block.BlockCharacteristic
	if req.Hash != "" {
		hashOrFinality = block.BlockHashRaw(req.Hash)
	} else {
		hashOrFinality = block.FinalityFinal()
	}
	blockView, err := n.nearClient.BlockDetails(ctx, hashOrFinality)
	if err != nil {
		log.Error("BlockDetails err", "err", err)
	}
	blockItem := &account.BlockHeader{
		ParentHash:       blockView.Header.PrevHash.String(),
		UncleHash:        "",
		CoinBase:         "",
		Root:             blockView.Header.BlockMerkleRoot.String(),
		TxHash:           blockView.Header.ChunkTxRoot.String(),
		ReceiptHash:      blockView.Header.ChunkReceiptsRoot.String(),
		ParentBeaconRoot: "",
		Difficulty:       "",
		Number:           strconv.FormatUint(blockView.Header.Height, 10),
		GasLimit:         0,
		GasUsed:          0,
		Time:             blockView.Header.Timestamp,
		Extra:            "",
		MixDigest:        "",
		Nonce:            "",
		BaseFee:          "",
		WithdrawalsHash:  blockView.Header.PrevStateRoot.String(),
		BlobGasUsed:      0,
		ExcessBlobGas:    0,
	}
	return &account.BlockHeaderResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "get block header by hash success",
		BlockHeader: blockItem,
	}, nil
}

// GetBlockHeaderByNumber implements chain.IChainAdaptor.
func (n *NearAdaptor) GetBlockHeaderByNumber(req *account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	ctx := context.Background()
	var hashOrFinality block.BlockCharacteristic
	if req.Height != 0 {
		hashOrFinality = block.BlockID(uint(req.Height))
	} else {
		hashOrFinality = block.FinalityFinal()
	}
	blockView, err := n.nearClient.BlockDetails(ctx, hashOrFinality)
	if err != nil {
		log.Error("GetBlockByRange err", "err", err)
	}
	blockItem := &account.BlockHeader{
		ParentHash:       blockView.Header.PrevHash.String(),
		UncleHash:        "",
		CoinBase:         "",
		Root:             blockView.Header.BlockMerkleRoot.String(),
		TxHash:           blockView.Header.ChunkTxRoot.String(),
		ReceiptHash:      blockView.Header.ChunkReceiptsRoot.String(),
		ParentBeaconRoot: "",
		Difficulty:       "",
		Number:           strconv.FormatUint(blockView.Header.Height, 10),
		GasLimit:         0,
		GasUsed:          0,
		Time:             blockView.Header.Timestamp,
		Extra:            "",
		MixDigest:        "",
		Nonce:            "",
		BaseFee:          "",
		WithdrawalsHash:  blockView.Header.PrevStateRoot.String(),
		BlobGasUsed:      0,
		ExcessBlobGas:    0,
	}
	return &account.BlockHeaderResponse{
		Code:        common2.ReturnCode_SUCCESS,
		Msg:         "get block header by number success",
		BlockHeader: blockItem,
	}, nil
}

// GetExtraData implements chain.IChainAdaptor.
func (n *NearAdaptor) GetExtraData(req *account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	return &account.ExtraDataResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "get extra data success",
		Value: "not data",
	}, nil
}

// GetFee implements chain.IChainAdaptor.
func (n *NearAdaptor) GetFee(req *account.FeeRequest) (*account.FeeResponse, error) {
	ctx := context.Background()
	gasData, err := n.nearClient.GasPriceView(ctx, block.FinalityFinal())
	if err != nil {
		log.Error("GetFee err", "err", err)
	}
	var gasPrice map[string]interface{}
	err = json.Unmarshal(gasData.Result, &gasPrice)
	return &account.FeeResponse{
		Code:      common2.ReturnCode_SUCCESS,
		Msg:       "get fee success",
		NormalFee: gasPrice["gas_price"].(string),
	}, nil

}

// GetSupportChains implements chain.IChainAdaptor.
func (n *NearAdaptor) GetSupportChains(req *account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	return &account.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

// GetTxByAddress implements chain.IChainAdaptor.
func (n *NearAdaptor) GetTxByAddress(req *account.TxAddressRequest) (*account.TxAddressResponse, error) {
	transactions, err := getAccountTransactions(req.Address)
	if err != nil {
		log.Error("GetTxByAddress err", "err", err)
	}
	var txMessageData []*account.TxMessage
	for _, t := range transactions {
		for _, ta := range t.Actions {
			if ta.Action == "TRANSFER" {
				message := &account.TxMessage{
					Hash:  t.TransactionHash,
					Index: 0,
					Froms: []*account.Address{{
						Address: t.SignerAccountID,
					}},
					Tos: []*account.Address{{
						Address: t.ReceiverAccountID,
					}},
					Fee:      strconv.FormatFloat(ta.Fee, 'f', -1, 64),
					Height:   strconv.FormatInt(int64(t.Block.BlockHeight), 10),
					Datetime: t.BlockTimestamp,
				}
				txMessageData = append(txMessageData, message)
			}
		}

	}
	return &account.TxAddressResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get GetTxByAddress success",
		Tx:   txMessageData,
	}, nil
}

// GetTxByHash implements chain.IChainAdaptor.
func (n *NearAdaptor) GetTxByHash(req *account.TxHashRequest) (*account.TxHashResponse, error) {
	res, err := n.nearClient.RPCClient.CallRPC(context.Background(), "EXPERIMENTAL_tx_status", map[string]interface{}{
		"tx_hash":           req.Hash,
		"sender_account_id": req.Coin,
	})
	if err != nil {
		log.Error("GetTxByHash err", "err", err)
	}
	result := res.Result
	var tramscatopmsStatus TransactionStatus
	err = json.Unmarshal(result, &tramscatopmsStatus)
	var totalGasPrice uint64
	for _, r := range tramscatopmsStatus.Receipts {
		parseUint, err := strconv.ParseUint(r.Receipt.Action.GasPrice, 10, 64)
		if err != nil {
			log.Error("GetTxByHash err", "err", err)
			continue
		}
		totalGasPrice = totalGasPrice + parseUint
	}
	message := &account.TxMessage{
		Hash:  req.Hash,
		Index: 0,
		Froms: []*account.Address{{
			Address: tramscatopmsStatus.Transaction.SignerID,
		}},
		Tos: []*account.Address{{
			Address: tramscatopmsStatus.Transaction.ReceiverID,
		}},
		Fee:  strconv.FormatUint(totalGasPrice, 10), // 交易并未
		Data: string(res.Result),
	}
	return &account.TxHashResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get transaction success",
		Tx:   message,
	}, nil
}

// SendTx implements chain.IChainAdaptor.
func (n *NearAdaptor) SendTx(req *account.SendTxRequest) (*account.SendTxResponse, error) {
	ctx := context.Background()
	name, err := n.nearClient.RPCTransactionSend(ctx, req.RawTx)
	if err != nil {
		return &account.SendTxResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Send tx error" + err.Error(),
		}, err
	}
	return &account.SendTxResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "send tx success",
		TxHash: name.String(),
	}, nil
}

// ValidAddress implements chain.IChainAdaptor.
func (n *NearAdaptor) ValidAddress(req *account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	// Near chain 可以自定义地址 所以无法验证地址 理论上说任何地址都是有效的 之后发送交易到链上才能知道地址是否真的有效
	panic("unimplemented")
}

// VerifySignedTransaction implements chain.IChainAdaptor.
func (n *NearAdaptor) VerifySignedTransaction(req *account.VerifyTransactionRequest) (*account.VerifyTransactionResponse, error) {
	return &account.VerifyTransactionResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "verify tx success",
		Verify: true,
	}, nil
}

// buildDynamicFeeTx
func (c *NearAdaptor) buildDynamicFeeTx(base64Tx string) (*NearTransaction, error) {
	// Decode base64 string
	txReqJsonByte, err := base64.StdEncoding.DecodeString(base64Tx)
	if err != nil {
		log.Error("decode string fail", "err", err)
		return nil, err
	}
	var nearTransaction NearTransaction
	if err := json.Unmarshal(txReqJsonByte, &nearTransaction); err != nil {
		log.Error("parse json fail", "err", err)
		return nil, err
	}
	return &nearTransaction, nil
}
