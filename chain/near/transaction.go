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

//
//func entrypoint(cctx *cli.Context) (err error) {
//error	networkID := cctx.String("network")
//	senderID := cctx.String("from")
//	recipientID := cctx.String("to")
//	amountValue := cctx.String("amount")
//
//	var amount types.Balance
//
//	amount, err = types.BalanceFromString(amountValue)
//	if err != nil {
//		return fmt.Errorf("failed to parse amount '%s': %w", amountValue, err)
//	}
//
//	network, ok := config.Networks[networkID]
//	if !ok {
//		return fmt.Errorf("unknown network '%s'", networkID)
//	}
//
//	keyPair, err := resolveCredentials(networkID, senderID)
//if err != nil {
//	return fmt.Errorf("failed to load private key: %w", err)
//}
//
//rpc, err := client.NewClient(network.NodeURL)
//if err != nil {
//	return fmt.Errorf("failed to create rpc client: %w", err)
//}
//
//log.Printf("near network: %s", rpc.NetworkAddr())
//
//ctx := client.ContextWithKeyPair(context.Background(), keyPair)
//res, err := rpc.TransactionSendAwait(
//	ctx, senderID, recipientID,
//	[]action.Action{
//		action.NewTransfer(amount),
//	},
//	client.WithLatestBlock(),
//)
//if err != nil {
//	return fmt.Errorf("failed to do txn: %w", err)
//}
//
//log.Printf("tx url: %s/transactions/%s", network.ExplorerURL, res.Transaction.Hash)
//	return
//}
//
//func resolveCredentials(networkName string, id types.AccountID) (kp key.KeyPair, err error) {
//	var creds struct {
//		AccountID  types.AccountID     `json:"account_id"`
//		PublicKey  key.Base58PublicKey `json:"public_key"`
//		PrivateKey key.KeyPair         `json:"private_key"`
//	}
//
//	var home string
//	home, err = os.UserHomeDir()
//	if err != nil {
//		return
//	}
//
//	credsFile := filepath.Join(home, ".near-credentials", networkName, fmt.Sprintf("%s.json", id))
//
//	var cf *os.File
//	if cf, err = os.Open(credsFile); err != nil {
//		return
//	}
//	defer cf.Close()
//
//	if err = json.NewDecoder(cf).Decode(&creds); err != nil {
//		return
//	}
//
//	if creds.PublicKey.String() != creds.PrivateKey.PublicKey.String() {
//		err = fmt.Errorf("inconsistent public key, %s != %s", creds.PublicKey.String(), creds.PrivateKey.PublicKey.String())
//		return
//	}
//	kp = creds.PrivateKey
//
//	return
//}
