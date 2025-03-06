package near

import (
	"context"
	"errors"

	"github.com/eteu-technologies/near-api-go/pkg/client/block"
	"github.com/eteu-technologies/near-api-go/pkg/types"
	"github.com/eteu-technologies/near-api-go/pkg/types/action"
	"github.com/eteu-technologies/near-api-go/pkg/types/key"
	"github.com/eteu-technologies/near-api-go/pkg/types/transaction"
	"github.com/ethereum/go-ethereum/log"
)

type transactionCtx struct {
	txn         transaction.Transaction
	keyPair     *key.KeyPair
	keyNonceSet bool
}

type TransactionOpt func(context.Context, *transactionCtx) error

func (c *NearAdaptor) PrepareTransaction(ctx context.Context, from, to types.AccountID, actions []action.Action, txnOpts ...TransactionOpt) (key.KeyPair, transaction.Transaction, error) {
	var ctx2 = context.WithValue(ctx, clientCtx, c)
	txn := transaction.Transaction{
		SignerID:   from,
		ReceiverID: to,
		Actions:    actions,
	}
	txnCtx := transactionCtx{
		txn:         txn,
		keyPair:     getKeyPair(ctx2),
		keyNonceSet: false,
	}

	var err error
	for _, opt := range txnOpts {
		if err = opt(ctx2, &txnCtx); err != nil {
		}
	}

	if txnCtx.keyPair == nil {
		err = errors.New("no keypair specified")
	}

	txnCtx.txn.PublicKey = txnCtx.keyPair.PublicKey.ToPublicKey()

	// Query the access key nonce, if not specified
	if !txnCtx.keyNonceSet {
		resp, err1 := c.nearClient.AccessKeyView(ctx2, txnCtx.txn.SignerID, txnCtx.keyPair.PublicKey, block.FinalityFinal())
		if err1 != nil {
			log.Error(err.Error())
		}
		nonce := resp.Nonce

		// Increment nonce by 1
		txnCtx.txn.Nonce = nonce + 1
		txnCtx.keyNonceSet = true
	}
	resp, err := c.nearClient.BlockDetails(ctx, block.FinalityFinal())
	if err != nil {
		log.Error(err.Error())
	}
	txnCtx.txn.BlockHash = resp.Header.Hash
	//blob, err = transaction.SignAndSerializeTransaction(*txnCtx.keyPair, txnCtx.txn)
	return *txnCtx.keyPair, txnCtx.txn, err
}

func (c *NearAdaptor) SignTransaction(keyPair key.KeyPair, txn transaction.Transaction) (blob string, err error) {
	blob, err = transaction.SignAndSerializeTransaction(keyPair, txn)
	return blob, err
}
