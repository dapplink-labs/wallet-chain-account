package aptos

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/aptos-labs/aptos-go-sdk/crypto"
	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

const (
	aptosAmount = uint64(400_0000)
)

func Test_v2_CreateUnSignTransaction_and_GenerateSignature(t *testing.T) {
	fromAddr, _ := AddressToAccountAddress("0x06671b50c2a5edb709c9e15d7e5a3d6496ae923759a217090e4ba9622720da5e")
	//fromPubKey, _ := PubKeyHexToPubKey("0x0caeddcef8612648417b7a07f1634625ab2da615c9178f00fe46fd0ac8d4d0e8")
	toAddr, _ := AddressToAccountAddress("0x84acbdf10e22b9b536d86ad6d017ff2854f7e1a48a4bb1f792ce571ee084fa68")
	amount := aptosAmount

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

		maxGasAmount := calculateMaxGasAmount(amount)

		coinTransferPayload, err := aptos.CoinTransferPayload(nil, toAddr, amount)
		assert.NoError(t, err)
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
		fmt.Println("=== Raw Transaction Field Details ===")
		fmt.Printf("Sender: %s\n", rawTxn.Sender.String())
		fmt.Printf("Sequence Number: %d\n", rawTxn.SequenceNumber)
		fmt.Printf("Max Gas Amount: %d\n", rawTxn.MaxGasAmount)
		fmt.Printf("Gas Unit Price: %d\n", rawTxn.GasUnitPrice)
		fmt.Printf("Expiration Timestamp: %d\n", rawTxn.ExpirationTimestampSeconds)
		fmt.Printf("Chain ID: %d\n", rawTxn.ChainId)

		// 打印 payload 信息
		if payload, ok := rawTxn.Payload.Payload.(*aptos.EntryFunction); ok {
			fmt.Printf("Payload Type: Entry Function\n")
			fmt.Printf("Module Name: %s\n", payload.Module.Name)
			fmt.Printf("Module Address: %s\n", payload.Module.Address.String())
			fmt.Printf("Function: %s\n", payload.Function)
			fmt.Printf("Type Arguments: %v\n", payload.ArgTypes)
			fmt.Printf("Arguments 0: %v\n", hex.EncodeToString(payload.Args[0]))
			fmt.Printf("Arguments 1: %v\n", binary.LittleEndian.Uint64(payload.Args[1]))
		} else {
			fmt.Printf("Payload: %+v\n", rawTxn.Payload)
		}
		fmt.Println("================================")

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

func Test_v2_BuildSignedTransaction(t *testing.T) {
	fromAddr, _ := AddressToAccountAddress("0x06671b50c2a5edb709c9e15d7e5a3d6496ae923759a217090e4ba9622720da5e")
	fromPubKey, _ := PubKeyHexToPubKey("0x0caeddcef8612648417b7a07f1634625ab2da615c9178f00fe46fd0ac8d4d0e8")
	toAddr, _ := AddressToAccountAddress("0x84acbdf10e22b9b536d86ad6d017ff2854f7e1a48a4bb1f792ce571ee084fa68")

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
	expirationTimestampSeconds := uint64(1731307824)
	fmt.Println("expirationTimestampSeconds:", expirationTimestampSeconds)
	amount := aptosAmount
	maxGasAmount := calculateMaxGasAmount(amount)

	transactionRequest, err := CreateTxReqByAptosCoin(fromAddr.String(), toAddr.String(),
		accountResponse.SequenceNumber, priceResponse.GasEstimate, amount, expirationTimestampSeconds, maxGasAmount, MainnetChainId)
	assert.NoError(t, err)

	transactionRequestJson, err := json.Marshal(transactionRequest)
	assert.NoError(t, err)
	transactionRequestBase64ByteList := base64.StdEncoding.EncodeToString(transactionRequestJson)

	signature := "1a11d37cdf6bb68aa13a6cf39520c5af088d0d7adfbc6cdbaeb8349af9de77eae9fe63c73d65980a73f5063151fbac2a9fcb0c3d43829df9db7ac0feb85bc20f"

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

func Test_v2_SendTx(t *testing.T) {

	aptosClient, err := NewAptosClient(string(Mainnet))
	assert.NoError(t, err, "failed to initialize aptos aptclient")

	adaptor := ChainAdaptor{
		aptosHttpClient: nil,
		aptosClient:     aptosClient,
	}

	sendTxRequest := &account.SendTxRequest{
		Chain:   ChainName,
		Network: string(Mainnet),
		RawTx:   "BmcbUMKl7bcJyeFdflo9ZJaukjdZohcJDkupYicg2l4KAAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQ1hcHRvc19hY2NvdW50CHRyYW5zZmVyAAIghKy98Q4iubU22GrW0Bf/KFT34aSKS7H3ks5XHuCE+mgIAAk9AAAAAACIEwAAAAAAAGQAAAAAAAAAMKkxZwAAAAABACAMrt3O+GEmSEF7egfxY0Ylqy2mFckXjwD+Rv0KyNTQ6EAaEdN832u2iqE6bPOVIMWvCI0Net+8bNuuuDSa+d536un+Y8c9ZZgKc/UGMVH7rCqfyww9Q4Kd+dt6wP64W8IP",
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

func CreateTxReqByAptosCoin(from, to string, sequenceNumber, gasPrice, amount, expirationTimestampSeconds, maxGasAmount uint64, chainId ChainId) (*TransactionRequest, error) {

	fromAddr, err := AddressToAccountAddress(from)
	if err != nil {
		return nil, fmt.Errorf("invalid from address: %w", err)
	}
	toAddr, err := AddressToAccountAddress(to)
	if err != nil {
		return nil, fmt.Errorf("invalid to address: %w", err)
	}
	if amount == 0 {
		return nil, errors.New("amount cannot be zero")
	}
	if gasPrice == 0 {
		return nil, errors.New("gas price cannot be zero")
	}
	txRequest := &TransactionRequest{
		Sender:         fromAddr.String(),
		SequenceNumber: sequenceNumber,
		Payload: PayloadWrapper{
			Payload: PayloadFunction{
				Module: ModuleInfo{
					Address: "0x1",
					Name:    "aptos_account",
				},
				Function: "transfer",
				ArgTypes: []string{"address", "u64"},
				Args:     []string{toAddr.String(), strconv.FormatUint(amount, 10)},
			},
		},
		MaxGasAmount: maxGasAmount,
		GasUnitPrice: gasPrice,
		// 1800 Second = 30 min
		ExpirationTimestampSeconds: expirationTimestampSeconds,
		ChainId:                    uint8(chainId),
	}
	if err := validateTransactionRequest(txRequest); err != nil {
		return nil, fmt.Errorf("invalid transaction request: %w", err)
	}

	return txRequest, nil
}

func validateTransactionRequest(tx *TransactionRequest) error {
	if tx == nil {
		return errors.New("transaction request is nil")
	}

	if len(tx.Sender) == 32 {
		return errors.New("sender address cannot be empty")
	}

	if tx.Payload.Payload.Module.Address == "" {
		return errors.New("module address cannot be empty")
	}
	if tx.Payload.Payload.Module.Name == "" {
		return errors.New("module name cannot be empty")
	}
	if tx.Payload.Payload.Function == "" {
		return errors.New("function name cannot be empty")
	}

	if len(tx.Payload.Payload.Args) != 2 {
		return errors.New("invalid number of arguments for transfer")
	}

	if tx.MaxGasAmount == 0 {
		return errors.New("max gas amount cannot be zero")
	}
	if tx.GasUnitPrice == 0 {
		return errors.New("gas unit price cannot be zero")
	}

	if tx.ExpirationTimestampSeconds <= uint64(time.Now().Unix()) {
		return errors.New("transaction already expired")
	}

	if tx.ChainId == 0 {
		return errors.New("chain id cannot be zero")
	}

	return nil
}

func calculateMaxGasAmount(amount uint64) uint64 {
	baseGas := uint64(2000)

	aptAmount := float64(amount) / 100000000.0

	var additionalGas uint64
	switch {
	case aptAmount <= 0.1:
		additionalGas = 3000
	case aptAmount <= 1:
		additionalGas = 5000
	case aptAmount <= 10:
		additionalGas = 10000
	default:
		additionalGas = 20000
	}

	return baseGas + additionalGas
}
