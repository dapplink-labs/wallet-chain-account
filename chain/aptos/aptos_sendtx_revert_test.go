package aptos

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/aptos-labs/aptos-go-sdk/crypto"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	aptosRevertAmount = uint64(100_0000)
)

func Test_v2_revert_CreateUnSignTransaction_and_GenerateSignature(t *testing.T) {
	fromAddr, _ := AddressToAccountAddress("0x84acbdf10e22b9b536d86ad6d017ff2854f7e1a48a4bb1f792ce571ee084fa68")
	toAddr, _ := AddressToAccountAddress("0x06671b50c2a5edb709c9e15d7e5a3d6496ae923759a217090e4ba9622720da5e")
	amount := aptosRevertAmount

	aptosHttpClient, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "failed to initialize aptos aptclient")

	aptosClient, err := NewAptosClient(string(Mainnet))
	assert.NoError(t, err, "failed to initialize aptos aptclient")

	adaptor := ChainAdaptor{
		aptosHttpClient: aptosHttpClient,
		aptosClient:     aptosClient,
	}

	t.Run("normal tx", func(t *testing.T) {
		accountResponse, err := adaptor.aptosHttpClient.GetAccount(fromAddr.String())
		assert.NoError(t, err)
		priceResponse, err := adaptor.aptosHttpClient.GetGasPrice()
		assert.NoError(t, err)

		expirationTimestampSeconds := uint64(time.Now().Unix() + 1800)
		fmt.Println("expirationTimestampSeconds:", expirationTimestampSeconds)

		coinTransferPayload, err := aptos.CoinTransferPayload(nil, toAddr, amount)
		assert.NoError(t, err)

		maxGasAmount := calculateMaxGasAmount(amount)
		options := []any{
			aptos.MaxGasAmount(maxGasAmount),
			aptos.GasUnitPrice(priceResponse.GasEstimate),
			aptos.ExpirationSeconds(int64(expirationTimestampSeconds)),
			aptos.SequenceNumber(accountResponse.SequenceNumber),
			aptos.ChainIdOption(1),
		}
		rawTxn, err := adaptor.aptosClient.BuildTransaction(
			fromAddr,
			aptos.TransactionPayload{Payload: coinTransferPayload},
			options...,
		)
		rawTxnJson, err := json.Marshal(rawTxn)
		fmt.Printf("rawTxnJson: %s\n", string(rawTxnJson))

		assert.NoError(t, err)
		//fmt.Println("=== Raw Transaction Field Details ===")
		//fmt.Printf("Sender: %s\n", rawTxn.Sender.String())
		//fmt.Printf("Sequence Number: %d\n", rawTxn.SequenceNumber)
		//fmt.Printf("Max Gas Amount: %d\n", rawTxn.MaxGasAmount)
		//fmt.Printf("Gas Unit Price: %d\n", rawTxn.GasUnitPrice)
		//fmt.Printf("Expiration Timestamp: %d\n", rawTxn.ExpirationTimestampSeconds)
		//fmt.Printf("Chain ID: %d\n", rawTxn.ChainId)
		//
		//// 打印 payload 信息
		//if payload, ok := rawTxn.Payload.Payload.(*aptos.EntryFunction); ok {
		//	fmt.Printf("Payload Type: Entry Function\n")
		//	fmt.Printf("Module Name: %s\n", payload.Module.Name)
		//	fmt.Printf("Module Address: %s\n", payload.Module.Address.String())
		//	fmt.Printf("Function: %s\n", payload.Function)
		//	fmt.Printf("Type Arguments: %v\n", payload.ArgTypes)
		//	fmt.Printf("Arguments 0: %v\n", hex.EncodeToString(payload.Args[0]))
		//	fmt.Printf("Arguments 1: %v\n", binary.LittleEndian.Uint64(payload.Args[1]))
		//} else {
		//	fmt.Printf("Payload: %+v\n", rawTxn.Payload)
		//}
		//fmt.Println("================================")

		transactionRequest, err := CreateTxReqByAptosCoin(fromAddr.String(), toAddr.String(),
			accountResponse.SequenceNumber, priceResponse.GasEstimate, amount, expirationTimestampSeconds, maxGasAmount, MainnetChainId)
		assert.NoError(t, err)

		transactionRequestJson, err := json.Marshal(transactionRequest)
		assert.NoError(t, err)
		transactionRequestBase64ByteList := base64.StdEncoding.EncodeToString(transactionRequestJson)

		req := &account.UnSignTransactionRequest{
			Chain:    ChainName,
			Network:  string(Mainnet),
			Base64Tx: transactionRequestBase64ByteList,
		}
		unSignTransaction, err := adaptor.CreateUnSignTransaction(req)
		assert.NoError(t, err)
		fmt.Println("unSignTransaction:", unSignTransaction)
		fmt.Println("unSignTransaction UnSignTx:", unSignTransaction.UnSignTx)
	})
}

func Test_v2_revert_BuildSignedTransaction(t *testing.T) {
	fromAddr, _ := AddressToAccountAddress("0x84acbdf10e22b9b536d86ad6d017ff2854f7e1a48a4bb1f792ce571ee084fa68")
	fromPubKey, _ := PubKeyHexToPubKey("0xd6e98225ebbb4872dcd5785ee18cd5990fac4482bc8d84f6fcad19e0d001e41f")
	toAddr, _ := AddressToAccountAddress("0x06671b50c2a5edb709c9e15d7e5a3d6496ae923759a217090e4ba9622720da5e")

	aptosHttpClient, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "failed to initialize aptos aptclient")
	aptosClient, err := NewAptosClient(string(Mainnet))
	assert.NoError(t, err, "failed to initialize aptos aptclient")

	adaptor := ChainAdaptor{
		aptosHttpClient: aptosHttpClient,
		aptosClient:     aptosClient,
	}

	accountResponse, err := aptosHttpClient.GetAccount(fromAddr.String())
	assert.NoError(t, err)
	priceResponse, err := aptosHttpClient.GetGasPrice()
	assert.NoError(t, err)
	// 1731249914 = CreateTxReqByAptosCoin input
	// 1731249914 = CreateUnSignTransaction input
	expirationTimestampSeconds := uint64(1731308002)
	fmt.Println("expirationTimestampSeconds:", expirationTimestampSeconds)

	amount := aptosRevertAmount
	maxGasAmount := calculateMaxGasAmount(amount)

	transactionRequest, err := CreateTxReqByAptosCoin(fromAddr.String(), toAddr.String(),
		accountResponse.SequenceNumber, priceResponse.GasEstimate, amount, expirationTimestampSeconds, maxGasAmount, MainnetChainId)
	assert.NoError(t, err)

	transactionRequestJson, err := json.Marshal(transactionRequest)
	assert.NoError(t, err)
	transactionRequestBase64ByteList := base64.StdEncoding.EncodeToString(transactionRequestJson)

	signature := "14ad2594fb13dd18e41cb1e8900430be1bc1f2f92b8984a58ddc8f5deb191e184c0f7fd1d4df791c34557d02f686a5250ac5c7e773019767190317e9088c000b"

	signedTransactionRequest := &account.SignedTransactionRequest{
		Chain:     ChainName,
		Network:   string(Mainnet),
		Base64Tx:  transactionRequestBase64ByteList,
		Signature: signature,
		PublicKey: fromPubKey.ToHex(),
	}
	buildSignedTransaction, err := adaptor.BuildSignedTransaction(signedTransactionRequest)
	assert.NoError(t, err)
	assert.NotEmpty(t, buildSignedTransaction.SignedTx)
	fmt.Println("BuildSignedTransaction resp", buildSignedTransaction)

	signedTxByteList, err := base64.StdEncoding.DecodeString(buildSignedTransaction.SignedTx)
	assert.NoError(t, err)
	fmt.Println("BuildSignedTransaction SignedTx", buildSignedTransaction.SignedTx)

	// 3.1, Deserializer
	signedTxBytesDes := bcs.NewDeserializer(signedTxByteList)
	signedTx := &aptos.SignedTransaction{
		Transaction: &aptos.RawTransaction{},
		Authenticator: &aptos.TransactionAuthenticator{
			Variant: aptos.TransactionAuthenticatorEd25519,
			Auth: &crypto.Ed25519Authenticator{
				PubKey: &crypto.Ed25519PublicKey{},
				Sig:    &crypto.Ed25519Signature{},
			},
		},
	}
	signedTx.UnmarshalBCS(signedTxBytesDes)

	fmt.Printf("\n=== SignedTransaction Details ===\n")
	if rawTx, ok := signedTx.Transaction.(*aptos.RawTransaction); ok {
		fmt.Printf("Raw Transaction:\n")
		fmt.Printf("  Sender: %s\n", rawTx.Sender.String())
		fmt.Printf("  Sequence Number: %d\n", rawTx.SequenceNumber)
		fmt.Printf("  Max Gas Amount: %d\n", rawTx.MaxGasAmount)
		fmt.Printf("  Gas Unit Price: %d\n", rawTx.GasUnitPrice)
		fmt.Printf("  Expiration: %d\n", rawTx.ExpirationTimestampSeconds)
		fmt.Printf("  Chain ID: %d\n", rawTx.ChainId)

		payload := rawTx.Payload
		fmt.Printf("\nPayload:\n")
		fmt.Printf("  Type: %T\n", payload)
		fmt.Printf("  Raw: %+v\n", payload)
	}

	verifyTransactionRequest := &account.VerifyTransactionRequest{
		Chain:     ChainName,
		Network:   string(Mainnet),
		PublicKey: fromPubKey.ToHex(),
		Signature: buildSignedTransaction.SignedTx,
	}
	verifyResp, err := adaptor.VerifySignedTransaction(verifyTransactionRequest)
	if err != nil {
		t.Fatalf("VerifySignedTransaction fail: %v", err)
	}
	fmt.Printf("verifyResp Verify: %v\n", verifyResp.Verify)
	fmt.Printf("verifyResp Code: %v\n", verifyResp.Code)
	fmt.Printf("verifyResp Msg: %v\n", verifyResp.Msg)

}

func Test_v2_revert_SendTx(t *testing.T) {

	aptosClient, err := NewAptosClient(string(Mainnet))
	assert.NoError(t, err, "failed to initialize aptos aptclient")

	adaptor := ChainAdaptor{
		aptosHttpClient: nil,
		aptosClient:     aptosClient,
	}

	sendTxRequest := &account.SendTxRequest{
		Chain:   ChainName,
		Network: string(Mainnet),
		RawTx:   "hKy98Q4iubU22GrW0Bf/KFT34aSKS7H3ks5XHuCE+mgAAAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQ1hcHRvc19hY2NvdW50CHRyYW5zZmVyAAIgBmcbUMKl7bcJyeFdflo9ZJaukjdZohcJDkupYicg2l4IQEIPAAAAAACIEwAAAAAAAGQAAAAAAAAA4qkxZwAAAAABACDW6YIl67tIctzVeF7hjNWZD6xEgryNhPb8rRng0AHkH0AUrSWU+xPdGOQcseiQBDC+G8Hy+SuJhKWN3I9d6xkeGEwPf9HU33kcNFV9AvaGpSUKxcfncwGXZxkDF+kIjAAL",
	}
	sendTxResponse, err := adaptor.SendTx(sendTxRequest)
	if err != nil {
		t.Fatalf("Failed to SendTx signedTx: %v", err)
	}
	sendTxResponseJson, _ := json.Marshal(sendTxResponse)
	fmt.Printf("sendTxResponseJson: %s\n", sendTxResponseJson)

	// 5, wait tx
	txn, err := aptosClient.WaitForTransaction(sendTxResponse.TxHash)
	assert.NoError(t, err, "WaitForTransaction fail")
	txnJson, _ := json.Marshal(txn)
	fmt.Printf("txnJson: %s\n", txnJson)
	assert.True(t, txn.Success, "WaitForTransaction fail")
	t.Logf("tx success, tx hash: %s", sendTxResponse.TxHash)
}
