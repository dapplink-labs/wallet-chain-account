package aptos

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/crypto"
	"github.com/stretchr/testify/assert"
)

func TestCreateDifferentAccounts(t *testing.T) {
	//privKey.ToHex()
	//privKey.Bytes()
	//privKey.PubKey()
	//
	//pubKey.ToHex()
	//pubKey.Bytes()
	//pubKey.AuthKey()
	//
	//account.Address()

	t.Run("Test NewEd25519Account", func(t *testing.T) {
		// Private (hex): 0x482607f684af7b554fb0a5cff46f3e13f55c386a8c2078c85fdd23482cd46d0d
		// pubkey (hex): 0xe053ceb35bdeeb607c404a76ee6a123e2e9527c158a4d283133b9a8053e9599b
		// address: 0xd9a32936401d4c547938ac6af3def73e74d2eea4912e5fc00276d41c04c6f814
		// authKey: 0xd9a32936401d4c547938ac6af3def73e74d2eea4912e5fc00276d41c04c6f814
		// signMessage: 0xc7049e4670af30c404dbd90eaaf36482a06d3dd00f9c8913130862714860961cdfd8d4793d6b74ed2720be72511750989e2c03a67c93d1eaac6e658647b76504
		account, err := aptos.NewEd25519Account()
		if err != nil {
			panic(err)
		}

		privateKey := account.Signer.(*crypto.Ed25519PrivateKey)
		fmt.Printf(" (hex): %s\n", privateKey.ToHex())

		publicKey := account.PubKey()
		fmt.Printf("pubkey (hex): %s\n", publicKey.ToHex())

		address := account.Address
		fmt.Printf("address: %s\n", address.String())

		authKey := account.AuthKey()
		fmt.Printf("authKey: %s\n", authKey.ToHex())

		message := []byte("Hello Aptos!")
		signature, err := account.SignMessage(message)
		if err != nil {
			panic(err)
		}
		fmt.Printf("signMessage: %s\n", signature.ToHex())
	})

	t.Run("Test NewEd25519Account", func(t *testing.T) {
		// Private (hex): 0x437f72e66d2955e6c65d3e97aa6aa7345d63b4afda0c58e78971247bd073843a
		// pubkey (hex): 0xdda59a045b6dc46c378fa84218c0d116c97795c63d8e8b56868fe2492f3bd45b
		// address: 0x34d8e3074323789467ce1e5d2c538312985dcd3b8889f29ce23e08b0d798404d
		// authKey: 0x34d8e3074323789467ce1e5d2c538312985dcd3b8889f29ce23e08b0d798404d
		// signMessage: 0xd18492c1a4f906d88c6842d6b9d3c6ac94a1cce79b2128e53e37ce7dd3e60810d0b88f18c18915175a9b0e9836b9023705080053b40f02a36c9c7bc199cab00c
		account, err := aptos.NewEd25519Account()
		if err != nil {
			panic(err)
		}

		privateKey := account.Signer.(*crypto.Ed25519PrivateKey)
		fmt.Printf("Private (hex): %s\n", privateKey.ToHex())

		publicKey := account.PubKey()
		fmt.Printf("pubkey (hex): %s\n", publicKey.ToHex())

		address := account.Address
		fmt.Printf("address: %s\n", address.String())

		authKey := account.AuthKey()
		fmt.Printf("authKey: %s\n", authKey.ToHex())

		message := []byte("Hello Aptos!")
		signature, err := account.SignMessage(message)
		if err != nil {
			panic(err)
		}
		fmt.Printf("signMessage: %s\n", signature.ToHex())
	})
}

func TestValidateAptosAddress(t *testing.T) {
	t.Run("Validate Aptos Address", func(t *testing.T) {
		address := &aptos.AccountAddress{}

		testAddress := "0x34d8e3074323789467ce1e5d2c538312985dcd3b8889f29ce23e08b0d798404d"

		err := address.ParseStringRelaxed(testAddress)
		if err != nil {
			t.Fatalf("ParseStringRelaxed testAddress fail: %v", err)
		}

		assert.Equal(t, testAddress, address.String(), "testAddress address Equal")

		fullAddress := address.StringLong()
		assert.Equal(t, 66, len(fullAddress), "fullAddress is 66 len")

		addressWithoutPrefix := testAddress[2:]
		err = address.ParseStringRelaxed(addressWithoutPrefix)
		assert.NoError(t, err, "should call ParseStringRelaxed with addressWithoutPrefix")

		fmt.Printf("Validate Aptos Address success:\n")
		fmt.Printf("address normal string: %s\n", address.String())
		fmt.Printf("address full string: %s\n", address.StringLong())
		fmt.Printf("address full string len: %d\n", len(address))
	})
}

func TestAptosTransfer(t *testing.T) {
	// Create aptclient
	client, err := aptos.NewClient(aptos.DevnetConfig)
	assert.NoError(t, err, "Failed to create aptclient")

	// Create test accounts with alice
	alice, err := aptos.NewEd25519Account()
	assert.NoError(t, err, "Failed to create alice account")
	// Create test accounts with bob
	bob, err := aptos.NewEd25519Account()
	assert.NoError(t, err, "Failed to create bob account")

	// Fund alice account
	const fundAmount = 100_000_000
	err = client.Fund(alice.Address, fundAmount)
	assert.NoError(t, err, "Failed to fund alice account")
	t.Logf("Funded account %s with %d APT", alice.Address, fundAmount)

	aliceBalance, err := client.AccountAPTBalance(alice.Address)
	assert.NoError(t, err, "Failed to get alice balance")
	assert.Equal(t, uint64(fundAmount), aliceBalance, "Incorrect initial balance for alice")
	t.Logf("Alice initial balance: %d APT", aliceBalance)

	const transferAmount = 1_000
	transferPayload, err := aptos.CoinTransferPayload(nil, bob.Address, transferAmount)
	assert.NoError(t, err, "build payload fail")

	// BuildTransaction
	rawTxn, err := client.BuildTransaction(alice.AccountAddress(),
		aptos.TransactionPayload{Payload: transferPayload})
	assert.NoError(t, err, "BuildTransaction fail")

	rawTxn1111, _ := json.Marshal(rawTxn)
	fmt.Printf("rawTxn1111: %s\n", rawTxn1111)

	simulationResult, err := client.SimulateTransaction(rawTxn, alice)
	assert.NoError(t, err, "simulationResult fail")
	assert.Equal(t, "Executed successfully", simulationResult[0].VmStatus, "simulationResult fail")
	t.Logf("simulationResult success, emit gas fee: %d", simulationResult[0].GasUsed)

	signedTxn, err := rawTxn.SignedTransaction(alice)
	assert.NoError(t, err, "SignedTransaction fail ")
	signedTxn11, _ := json.Marshal(signedTxn)
	fmt.Printf("signedTxn11: %s\n", signedTxn11)

	flag, msg := VerifySignedTxn(signedTxn)
	fmt.Printf("VerifySignedTxn flag: %v\n", flag)
	fmt.Printf("VerifySignedTxn msg: %s\n", msg)

	submitResult, err := client.SubmitTransaction(signedTxn)
	assert.NoError(t, err, "submitResult fail")

	txn, err := client.WaitForTransaction(submitResult.Hash)
	assert.NoError(t, err, "WaitForTransaction fail")
	assert.True(t, txn.Success, "WaitForTransaction fail")
	t.Logf("WaitForTransaction success, tx hash: %s", submitResult.Hash)

	newAliceBalance, err := client.AccountAPTBalance(alice.Address)
	assert.NoError(t, err, "newAliceBalance fail")

	bobBalance, err := client.AccountAPTBalance(bob.Address)
	assert.NoError(t, err, "bobBalance fail")

	assert.Less(t, newAliceBalance, aliceBalance-transferAmount, "newAliceBalance fail")
	assert.Equal(t, uint64(transferAmount), bobBalance, "bobBalance fail")
	t.Logf("transfer newAliceBalance: %d APT", newAliceBalance)
	t.Logf("transfer Bob fee: %d APT", bobBalance)
	t.Logf("normal gas fee: %d", aliceBalance-transferAmount-newAliceBalance)
}

// VerifySignedTxn
func VerifySignedTxn(signedTxn *aptos.SignedTransaction) (bool, string) {
	var messages []string
	isValid := true

	rawTxn, ok := signedTxn.Transaction.(*aptos.RawTransaction)
	if !ok {
		return false, "signedTxn Transaction error"
	}

	signingMessage, err := rawTxn.SigningMessage()
	if err != nil {
		return false, "SigningMessage err"
	}

	auth := signedTxn.Authenticator.Auth
	authJson, _ := json.Marshal(auth)
	fmt.Printf("authJson Json: %s\n", authJson)
	fmt.Printf("Type of auth: %T\n", auth)

	if senderAuth, ok := auth.(*aptos.Ed25519TransactionAuthenticator); ok {
		if ed25519Auth, ok := senderAuth.Sender.Auth.(*crypto.Ed25519Authenticator); ok {
			// Step 1: Derive account address from the public key in signedTxn.Authenticator.Auth
			accountAddress, err := PubKeyToAccountAddress(ed25519Auth.PubKey)
			if err != nil {
				return false, "failed to derive address from public key"
			}
			// Step 2: Verify if the derived address matches the sender address in raw transaction
			// This ensures the transaction is truly from the claimed sender
			// rawTxn.Sender from rawTxn, this is req.Sender.address
			// accountAddress from signedTxn.Authenticator.Auth, this is signedTxn.address
			if *accountAddress != rawTxn.Sender {
				isValid = false
				messages = append(messages, fmt.Sprintf("sender address mismatch\nexpected: %s\nactual: %s",
					accountAddress.String(), rawTxn.Sender.String()))
			}

			// Verify if the signature is valid for this transaction
			// ed25519Auth.PubKey: the public key of the signer
			// Parameters:
			// - signingMessage: the original transaction data to be signed
			// - ed25519Auth.Sig: the signature created by the private key
			if !ed25519Auth.PubKey.Verify(signingMessage, ed25519Auth.Sig) {
				isValid = false
				messages = append(messages, "invalid signature")
			}
		} else {
			return false, "invalid ed25519 authenticator type"
		}
	} else {
		return false, "invalid sender authenticator type"
	}

	// Step 4: Verify transaction basic parameters
	messages = append(messages, fmt.Sprintf("\n=== Transaction Parameters ==="+
		"\nSender Address: %s"+
		"\nSequence Number: %d"+
		"\nGas Limit: %d"+
		"\nGas Unit Price: %d"+
		"\nExpiration Time: %d"+
		"\nChain ID: %d",
		rawTxn.Sender.String(),
		rawTxn.SequenceNumber,
		rawTxn.MaxGasAmount,
		rawTxn.GasUnitPrice,
		rawTxn.ExpirationTimestampSeconds,
		rawTxn.ChainId))

	payloadJson, _ := json.Marshal(rawTxn.Payload)
	fmt.Printf("Payload Json: %s\n", payloadJson)
	fmt.Printf("Type of rawTxn.Payload: %T\n", rawTxn.Payload)
	fmt.Printf("Type of rawTxn.Payload.Payload: %T\n", rawTxn.Payload.Payload)

	// Step 5: Verify transfer parameters if it's a transfer transaction
	if entryFunction, ok := rawTxn.Payload.Payload.(*aptos.EntryFunction); ok {
		messages = append(messages, "\n=== Transfer Parameters ===")

		// 1. Verify module and function name
		if entryFunction.Module.Name != "aptos_account" {
			isValid = false
			messages = append(messages, "invalid module name")
		}
		if entryFunction.Function != "transfer" {
			isValid = false
			messages = append(messages, "invalid function name")
		}

		// 2. Verify arguments length
		if len(entryFunction.Args) < 2 {
			isValid = false
			messages = append(messages, "insufficient transfer arguments")
			return isValid, strings.Join(messages, "\n")
		}

		// 3. Verify recipient address
		toAddrBytes := entryFunction.Args[0]
		if len(toAddrBytes) != 32 { // Aptos address length
			isValid = false
			messages = append(messages, "invalid recipient address length")
		}
		messages = append(messages, fmt.Sprintf("Recipient Address (bytes): %x", toAddrBytes))

		// 4. Verify transfer amount
		amountBytes := entryFunction.Args[1]
		if len(amountBytes) != 8 { // uint64 length
			isValid = false
			messages = append(messages, "invalid amount format")
		}
		amount := binary.LittleEndian.Uint64(amountBytes)

		// Check amount constraints
		if amount == 0 {
			isValid = false
			messages = append(messages, "transfer amount cannot be zero")
		}
		const MAX_TRANSFER_AMOUNT = uint64(1000000000000)
		if amount > MAX_TRANSFER_AMOUNT {
			isValid = false
			messages = append(messages, fmt.Sprintf("transfer amount exceeds maximum limit: %d", MAX_TRANSFER_AMOUNT))
		}
		messages = append(messages, fmt.Sprintf("Transfer Amount: %d", amount))

		// 5. Verify sender has sufficient balance
		// totalRequired := amount + (rawTxn.MaxGasAmount * rawTxn.GasUnitPrice)
		// if senderBalance < totalRequired {
		//     isValid = false
		//     messages = append(messages, "insufficient balance for transfer and gas")
		// }

		// 6. Verify recipient address is valid
		if bytes.Equal(toAddrBytes, rawTxn.Sender[:]) {
			isValid = false
			messages = append(messages, "cannot transfer to self")
		}
	}

	return isValid, strings.Join(messages, "\n")
}
