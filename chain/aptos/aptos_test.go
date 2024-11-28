package aptos

import (
	"encoding/json"
	"fmt"
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/crypto"
	"github.com/stretchr/testify/assert"
	"log"
	"math"
	"strconv"
	"testing"

	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

func Test_GenerateAptosAccount(t *testing.T) {
	t.Run("Test NewEd25519Account", func(t *testing.T) {
		accountTemp, err := aptos.NewEd25519Account()
		if err != nil {
			panic(err)
		}

		privateKey := accountTemp.Signer.(*crypto.Ed25519PrivateKey)
		fmt.Printf("Private (hex): %s\n", privateKey.ToHex())

		publicKey := accountTemp.PubKey()
		fmt.Printf("pubkey (hex): %s\n", publicKey.ToHex())

		address := accountTemp.Address
		fmt.Printf("address: %s\n", address.String())

		authKey := accountTemp.AuthKey()
		fmt.Printf("authKey: %s\n", authKey.ToHex())

		message := []byte("Hello Aptos!")
		signature, err := accountTemp.SignMessage(message)
		if err != nil {
			panic(err)
		}
		fmt.Printf("signMessage: %s\n", signature.ToHex())
	})
}

func TestChainAdaptor_GetSupportChains(t *testing.T) {
	adaptor := ChainAdaptor{}

	req := &account.SupportChainsRequest{
		Chain:   ChainName,
		Network: "mainnet",
	}

	resp, err := adaptor.GetSupportChains(req)

	if err != nil {
		t.Errorf("GetSupportChains failed with error: %v", err)
	}
	fmt.Printf("resp: %s\n", resp)

	if resp.Code != common2.ReturnCode_SUCCESS {
		t.Errorf("Expected success code, got %v", resp.Code)
	}

	if !resp.Support {
		t.Error("Expected Support to be true")
	}
}

func TestChainAdaptor_ConvertAddress(t *testing.T) {
	const (
		validPublicKey        = "0x032a5b5d8a16bcdf70ab23e9e717693672b969366c0a61b7fbddd723c986bf2925"
		validPublicKeyAddress = "0x3a8eef8a52bc873f5416e835e7ec7da6dd978e5f6a8a12d278df0c42ef01d131"
		emptyKeyAddress       = "0xa7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"
		invalidPublicKey      = "0x06b331be57280b9050ea43fbcaf8f0e6f194987c9d9b00a42c1b5367742a2d05"
	)

	adaptor := &ChainAdaptor{}

	t.Run("Valid Public Key", func(t *testing.T) {
		req := &account.ConvertAddressRequest{
			Chain:     ChainName,
			Network:   "mainnet",
			PublicKey: validPublicKey,
		}
		resp, err := adaptor.ConvertAddress(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.Equal(t, validPublicKeyAddress, resp.Address)
		assert.Equal(t, "convert address success", resp.Msg)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Converted Address: %s", resp.Address)
	})

	t.Run("Empty Public Key", func(t *testing.T) {
		req := &account.ConvertAddressRequest{
			Chain:     ChainName,
			Network:   "mainnet",
			PublicKey: "",
		}
		resp, err := adaptor.ConvertAddress(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Invalid Public Key Format", func(t *testing.T) {
		req := &account.ConvertAddressRequest{
			Chain:     ChainName,
			Network:   "mainnet",
			PublicKey: invalidPublicKey,
		}
		resp, err := adaptor.ConvertAddress(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestChainAdaptor_ValidAddress(t *testing.T) {
	const (
		validAddress   = "0x3a8eef8a52bc873f5416e835e7ec7da6dd978e5f6a8a12d278df0c42ef01d131"
		allZeroAddress = "0x0000000000000000000000000000000000000000000000000000000000000000"
		invalidChars   = "0xfc38d27af874e409de8056d11cc8e10b8f8449e6f723a59251f04e62a24d74zz"
		shortAddress   = "0xfc38d27af874e409de8056d11cc8e10b8f8449e6"
	)

	adaptor := &ChainAdaptor{}

	t.Run("Valid Address", func(t *testing.T) {
		req := &account.ValidAddressRequest{
			Address: validAddress,
			Chain:   ChainName,
		}
		resp, err := adaptor.ValidAddress(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.True(t, resp.Valid)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Is Valid: %v", resp.Valid)
	})

	t.Run("Missing 0x Prefix", func(t *testing.T) {
		req := &account.ValidAddressRequest{
			Chain:   ChainName,
			Address: validAddress[2:],
		}
		resp, err := adaptor.ValidAddress(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.True(t, resp.Valid)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Is Valid: %v", resp.Valid)
	})

	t.Run("All Zeros Address", func(t *testing.T) {
		req := &account.ValidAddressRequest{
			Chain:   ChainName,
			Address: allZeroAddress,
		}
		resp, err := adaptor.ValidAddress(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Is Valid: %v", resp.Valid)
	})

	t.Run("Invalid Characters", func(t *testing.T) {
		req := &account.ValidAddressRequest{
			Chain:   ChainName,
			Address: invalidChars,
		}
		resp, err := adaptor.ValidAddress(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		//assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		//assert.False(t, resp.Valid)

		//t.Logf("Response Code: %v", resp.Code)
		//t.Logf("Response Message: %s", resp.Msg)
		//t.Logf("Is Valid: %v", resp.Valid)
	})

	t.Run("Short Address", func(t *testing.T) {
		req := &account.ValidAddressRequest{
			Chain:   ChainName,
			Address: shortAddress,
		}
		resp, err := adaptor.ValidAddress(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.True(t, resp.Valid)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Is Valid: %v", resp.Valid)
	})
}

func TestChainAdaptor_GetBlockByNumber(t *testing.T) {
	const (
		latestBlock    = int64(0)
		specificHeight = int64(247764636)
		invalidHeight  = int64(-1)
		withTxHeight   = int64(1000)
	)
	aptosClient, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "failed to initialize aptos aptclient")

	adaptor := ChainAdaptor{
		aptosHttpClient: aptosClient,
	}

	t.Run("Get Latest Block", func(t *testing.T) {
		req := &account.BlockNumberRequest{
			Chain:  ChainName,
			Height: latestBlock,
			ViewTx: false,
		}
		resp, err := adaptor.GetBlockByNumber(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.Greater(t, resp.Height, latestBlock)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Block Height: %d", resp.Height)
	})

	t.Run("Get Block By Specific Height", func(t *testing.T) {
		req := &account.BlockNumberRequest{
			Chain:  ChainName,
			Height: specificHeight,
			ViewTx: false,
		}
		resp, err := adaptor.GetBlockByNumber(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.Equal(t, specificHeight, resp.Height)
		assert.NotEmpty(t, resp.Hash)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Block Height: %d", resp.Height)
		t.Logf("Block Hash: %s", resp.Hash)
	})

	t.Run("Get Block With Invalid Height", func(t *testing.T) {
		req := &account.BlockNumberRequest{
			Chain:  ChainName,
			Height: invalidHeight,
			ViewTx: false,
		}
		resp, err := adaptor.GetBlockByNumber(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		//assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		//assert.Contains(t, resp.Msg, "invalid block height")

		//t.Logf("Response Code: %v", resp.Code)
		//t.Logf("Response Message: %s", resp.Msg)
	})

	t.Run("Get Block With Transactions", func(t *testing.T) {
		req := &account.BlockNumberRequest{
			Chain:  ChainName,
			Height: withTxHeight,
			ViewTx: true,
		}
		resp, err := adaptor.GetBlockByNumber(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.Equal(t, withTxHeight, resp.Height)
		assert.NotEmpty(t, resp.Hash)
		// Note: Currently transactions are not implemented in the adapter
		assert.Nil(t, resp.Transactions)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Block Height: %d", resp.Height)
		t.Logf("Block Hash: %s", resp.Hash)
	})
}

func TestChainAdaptor_GetBlockHeaderByNumber(t *testing.T) {
	const (
		latestBlock    = int64(0)
		specificHeight = int64(247764636)
		invalidHeight  = int64(-1)
	)

	aptosClient, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "failed to initialize aptos aptclient")

	adaptor := ChainAdaptor{
		aptosHttpClient: aptosClient,
	}

	t.Run("Get Latest Block Header", func(t *testing.T) {
		req := &account.BlockHeaderNumberRequest{
			Height: latestBlock,
			Chain:  ChainName,
		}
		resp, err := adaptor.GetBlockHeaderByNumber(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.NotNil(t, resp.BlockHeader)
		//assert.NotEmpty(t, resp.BlockHeader.Hash)
		//assert.NotEmpty(t, resp.BlockHeader.ParentHash)
		assert.NotEmpty(t, resp.BlockHeader.Number)
		assert.Greater(t, resp.BlockHeader.Time, uint64(0))

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Block Header Hash: %s", resp.BlockHeader.Hash)
		t.Logf("Block Parent Hash: %s", resp.BlockHeader.ParentHash)
		t.Logf("Block Number: %s", resp.BlockHeader.Number)
		t.Logf("Block Time: %d", resp.BlockHeader.Time)
	})

	t.Run("Get Block Header By Specific Height", func(t *testing.T) {
		req := &account.BlockHeaderNumberRequest{
			Height: specificHeight,
			Chain:  ChainName,
		}
		resp, err := adaptor.GetBlockHeaderByNumber(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.NotNil(t, resp.BlockHeader)
		//assert.NotEmpty(t, resp.BlockHeader.Hash)
		//assert.NotEmpty(t, resp.BlockHeader.ParentHash)
		assert.Equal(t, fmt.Sprintf("%d", specificHeight), resp.BlockHeader.Number)
		assert.Greater(t, resp.BlockHeader.Time, uint64(0))

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Block Header Hash: %s", resp.BlockHeader.Hash)
		t.Logf("Block Parent Hash: %s", resp.BlockHeader.ParentHash)
		t.Logf("Block Number: %s", resp.BlockHeader.Number)
		t.Logf("Block Time: %d", resp.BlockHeader.Time)
	})

	t.Run("Get Block Header With Invalid Height", func(t *testing.T) {
		req := &account.BlockHeaderNumberRequest{
			Height: invalidHeight,
			Chain:  ChainName,
		}
		resp, err := adaptor.GetBlockHeaderByNumber(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		//assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		//assert.Nil(t, resp.BlockHeader)

		//t.Logf("Response Code: %v", resp.Code)
		//t.Logf("Response Message: %s", resp.Msg)
	})

	t.Run("Get Block Header With Future Height", func(t *testing.T) {
		req := &account.BlockHeaderNumberRequest{
			Height: math.MaxInt64,
			Chain:  ChainName,
		}
		resp, err := adaptor.GetBlockHeaderByNumber(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		//assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		//assert.Nil(t, resp.BlockHeader)

		//t.Logf("Response Code: %v", resp.Code)
		//t.Logf("Response Message: %s", resp.Msg)
	})
}

func TestChainAdaptor_GetAccount(t *testing.T) {
	const (
		validAccount   = "0x8d2d7bcde13b2513617df3f98cdd5d0e4b9f714c6308b9204fe18ad900d92609"
		expectedAPT    = 0.68374979
		invalidAccount = "0xinvalid_account"
		emptyAccount   = ""
	)

	aptosClient, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "failed to initialize aptos aptclient")

	adaptor := ChainAdaptor{
		aptosHttpClient: aptosClient,
	}

	t.Run("Valid Account", func(t *testing.T) {
		req := &account.AccountRequest{
			ConsumerToken:    "test_token",
			Chain:            ChainName,
			Coin:             "APT",
			Network:          "mainnet",
			Address:          validAccount,
			ContractAddress:  "",
			ProposerKeyIndex: 0,
		}

		resp, err := adaptor.GetAccount(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.Equal(t, req.Network, resp.Network)
		assert.NotEmpty(t, resp.Sequence)
		assert.NotEmpty(t, resp.Balance)

		aptValue, err := strconv.ParseInt(resp.Balance, 10, 64)
		if err != nil {
			log.Printf("convert err: %v", err)
		}
		t.Logf("Account %s APT balance: %d", validAccount, aptValue)
		aptValueFloat64 := float64(aptValue) / 100000000
		t.Logf("Account %s APT balance: %v", validAccount, aptValueFloat64)
		assert.InDelta(t, expectedAPT, aptValueFloat64, 0.00000001, "APT balance should match expected amount")

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Network: %s", resp.Network)
		t.Logf("Account Number: %s", resp.AccountNumber)
		t.Logf("Sequence: %s", resp.Sequence)
		t.Logf("Balance: %s", resp.Balance)
	})

	t.Run("Invalid Account", func(t *testing.T) {
		req := &account.AccountRequest{
			ConsumerToken:    "test_token",
			Chain:            ChainName,
			Coin:             "APT",
			Network:          "mainnet",
			Address:          invalidAccount,
			ContractAddress:  "",
			ProposerKeyIndex: 0,
		}

		resp, err := adaptor.GetAccount(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		//assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		//assert.Equal(t, "GetAccount fail", resp.Msg)
		//assert.Empty(t, resp.Sequence)
		//assert.Empty(t, resp.AccountNumber)
		//assert.Empty(t, resp.Balance)

		//t.Logf("Response Code: %v", resp.Code)
		//t.Logf("Response Message: %s", resp.Msg)
	})

	t.Run("Empty Account", func(t *testing.T) {
		req := &account.AccountRequest{
			ConsumerToken:    "test_token",
			Chain:            ChainName,
			Coin:             "APT",
			Network:          "mainnet",
			Address:          emptyAccount,
			ContractAddress:  "",
			ProposerKeyIndex: 0,
		}

		resp, err := adaptor.GetAccount(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		//assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		//assert.Equal(t, "GetAccount fail", resp.Msg)
		//assert.Empty(t, resp.Sequence)
		//assert.Empty(t, resp.AccountNumber)
		//assert.Empty(t, resp.Balance)

		//t.Logf("Response Code: %v", resp.Code)
		//t.Logf("Response Message: %s", resp.Msg)
	})
}

func TestChainAdaptor_GetFee(t *testing.T) {
	aptosClient, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "failed to initialize aptos aptclient")

	adaptor := ChainAdaptor{
		aptosHttpClient: aptosClient,
	}

	t.Run("Valid Fee Request", func(t *testing.T) {
		req := &account.FeeRequest{
			ConsumerToken: "test_token",
			Chain:         ChainName,
			Coin:          "APT",
			Network:       "mainnet",
			RawTx:         "",
			Address:       "0xfc38d27af874e409de8056d11cc8e10b8f8449e6f723a59251f04e62a24d7475",
		}

		resp, err := adaptor.GetFee(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.Equal(t, "GetFee success", resp.Msg)

		assert.NotEmpty(t, resp.SlowFee)
		assert.NotEmpty(t, resp.NormalFee)
		assert.NotEmpty(t, resp.FastFee)

		slowFee, _ := strconv.ParseUint(resp.SlowFee, 10, 64)
		normalFee, _ := strconv.ParseUint(resp.NormalFee, 10, 64)
		fastFee, _ := strconv.ParseUint(resp.FastFee, 10, 64)

		assert.LessOrEqual(t, slowFee, normalFee, "Slow fee should be less than normal fee")
		assert.Less(t, normalFee, fastFee, "Normal fee should be less than fast fee")

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Slow Fee: %s", resp.SlowFee)
		t.Logf("Normal Fee: %s", resp.NormalFee)
		t.Logf("Fast Fee: %s", resp.FastFee)
	})

	t.Run("Invalid Chain", func(t *testing.T) {
		req := &account.FeeRequest{
			ConsumerToken: "test_token",
			Chain:         "InvalidChain",
			Coin:          "APT",
			Network:       "mainnet",
			RawTx:         "",
			Address:       "0xfc38d27af874e409de8056d11cc8e10b8f8449e6f723a59251f04e62a24d7475",
		}

		resp, err := adaptor.GetFee(req)

		assert.Error(t, err)
		assert.NotEmpty(t, resp)
		//assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		//assert.Equal(t, "GetFee fail", resp.Msg)
		//assert.Empty(t, resp.SlowFee)
		//assert.Empty(t, resp.NormalFee)
		//assert.Empty(t, resp.FastFee)

		//t.Logf("Response Code: %v", resp.Code)
		//t.Logf("Response Message: %s", resp.Msg)
	})

	t.Run("Empty Address", func(t *testing.T) {
		req := &account.FeeRequest{
			ConsumerToken: "test_token",
			Chain:         ChainName,
			Coin:          "APT",
			Network:       "mainnet",
			RawTx:         "",
			Address:       "",
		}

		resp, err := adaptor.GetFee(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.NotEmpty(t, resp.SlowFee)
		assert.NotEmpty(t, resp.NormalFee)
		assert.NotEmpty(t, resp.FastFee)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
	})

	t.Run("Multiple Consecutive Requests", func(t *testing.T) {
		req := &account.FeeRequest{
			ConsumerToken: "test_token",
			Chain:         ChainName,
			Coin:          "APT",
			Network:       "mainnet",
			RawTx:         "",
			Address:       "0xfc38d27af874e409de8056d11cc8e10b8f8449e6f723a59251f04e62a24d7475",
		}

		var fees []struct {
			slow   uint64
			normal uint64
			fast   uint64
		}

		for i := 0; i < 3; i++ {
			resp, err := adaptor.GetFee(req)
			assert.NoError(t, err)
			assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)

			slow, _ := strconv.ParseUint(resp.SlowFee, 10, 64)
			normal, _ := strconv.ParseUint(resp.NormalFee, 10, 64)
			fast, _ := strconv.ParseUint(resp.FastFee, 10, 64)

			fees = append(fees, struct {
				slow   uint64
				normal uint64
				fast   uint64
			}{slow, normal, fast})

			t.Logf("Request %d - Slow: %d, Normal: %d, Fast: %d", i+1, slow, normal, fast)
		}

		for i := 1; i < len(fees); i++ {
			assert.InDelta(t, fees[i-1].slow, fees[i].slow, float64(fees[i-1].slow)*0.5)
			assert.InDelta(t, fees[i-1].normal, fees[i].normal, float64(fees[i-1].normal)*0.5)
			assert.InDelta(t, fees[i-1].fast, fees[i].fast, float64(fees[i-1].fast)*0.5)
		}
	})
}

func TestChainAdaptor_GetTxByAddress(t *testing.T) {
	const (
		validAddress   = "0xb5e1cc180e603037887c9e9eb4a8a06774ebcddafac37ceea9e33f3b6552bb25"
		invalidAddress = "0xinvalid_address"
		emptyAddress   = ""
	)

	aptosClient, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "failed to initialize aptos aptclient")

	adaptor := ChainAdaptor{
		aptosHttpClient: aptosClient,
	}

	t.Run("Valid Address", func(t *testing.T) {
		req := &account.TxAddressRequest{
			ConsumerToken: "test_token",
			Chain:         ChainName,
			Coin:          "APT",
			Network:       "mainnet",
			Address:       validAddress,
		}

		resp, err := adaptor.GetTxByAddress(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.Equal(t, "GetTxByAddress success", resp.Msg)

		if len(resp.Tx) > 0 {
			for _, tx := range resp.Tx {
				assert.NotEmpty(t, tx.Hash, "Transaction hash should not be empty")
				assert.NotEmpty(t, tx.Height, "Transaction height should not be empty")
				assert.NotEmpty(t, tx.Fee, "Transaction fee should not be empty")
				assert.NotEmpty(t, tx.Datetime, "Transaction datetime should not be empty")

				assert.NotEmpty(t, tx.Froms, "From addresses should not be empty")
				assert.Equal(t, validAddress, tx.Froms[0].Address)

				assert.Contains(t, []account.TxStatus{
					account.TxStatus_Success,
					account.TxStatus_Failed,
				}, tx.Status)

				fee, err := strconv.ParseUint(tx.Fee, 10, 64)
				assert.NoError(t, err, "Fee should be a valid number")
				assert.Greater(t, fee, uint64(0), "Fee should be greater than 0")

				height, err := strconv.ParseUint(tx.Height, 10, 64)
				assert.NoError(t, err, "Height should be a valid number")
				assert.Greater(t, height, uint64(0), "Height should be greater than 0")

				t.Logf("Transaction Hash: %s", tx.Hash)
				t.Logf("Transaction Height: %s", tx.Height)
				t.Logf("Transaction Fee: %s", tx.Fee)
				t.Logf("Transaction Status: %v", tx.Status)
			}
		}
	})

	t.Run("Invalid Address", func(t *testing.T) {
		req := &account.TxAddressRequest{
			ConsumerToken: "test_token",
			Chain:         "dasdas",
			Coin:          "APT",
			Network:       "mainnet",
			Address:       invalidAddress,
		}

		resp, err := adaptor.GetTxByAddress(req)

		assert.Error(t, err)
		assert.NotNil(t, resp)
		//assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		//assert.Equal(t, "GetTxByAddress GetTransactionByAddress fail", resp.Msg)
		//assert.Empty(t, resp.Tx)

		//t.Logf("Response Code: %v", resp.Code)
		//t.Logf("Response Message: %s", resp.Msg)
	})

	t.Run("Empty Address", func(t *testing.T) {
		req := &account.TxAddressRequest{
			ConsumerToken: "test_token",
			Chain:         ChainName,
			Coin:          "APT",
			Network:       "mainnet",
			Address:       emptyAddress,
		}

		resp, err := adaptor.GetTxByAddress(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		//assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		//assert.Equal(t, "GetTxByAddress GetTransactionByAddress fail", resp.Msg)
		//assert.Empty(t, resp.Tx)

		//t.Logf("Response Code: %v", resp.Code)
		//t.Logf("Response Message: %s", resp.Msg)
	})

	t.Run("Invalid Chain", func(t *testing.T) {
		req := &account.TxAddressRequest{
			ConsumerToken: "test_token",
			Chain:         "InvalidChain",
			Coin:          "APT",
			Network:       "mainnet",
			Address:       validAddress,
		}

		resp, err := adaptor.GetTxByAddress(req)

		assert.Error(t, err)
		assert.NotNil(t, resp)
		//assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		//assert.Equal(t, "invalid chain", resp.Msg)
		//assert.Empty(t, resp.Tx)

		//t.Logf("Response Code: %v", resp.Code)
		//t.Logf("Response Message: %s", resp.Msg)
	})
}

func TestChainAdaptor_GetTxByHash(t *testing.T) {
	const (
		validTxHash   = "0x43531969ff8e93de962ea65e5609c2b05de3aa5e78933d8925613e75d3d53772"
		invalidTxHash = "0xinvalid_hash"
		emptyTxHash   = ""
	)

	aptosClient, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "failed to initialize aptos aptclient")

	adaptor := ChainAdaptor{
		aptosHttpClient: aptosClient,
	}

	t.Run("Valid Transaction Hash", func(t *testing.T) {
		req := &account.TxHashRequest{
			Chain: ChainName,
			Hash:  validTxHash,
		}
		resp, err := adaptor.GetTxByHash(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.NotNil(t, resp.Tx)
		assert.Equal(t, validTxHash, resp.Tx.Hash)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Transaction Hash: %s", resp.Tx.Hash)
		txJson, _ := json.Marshal(resp.Tx)
		t.Logf("Transaction txJson: %s", txJson)
	})

	t.Run("Invalid Transaction Hash", func(t *testing.T) {
		req := &account.TxHashRequest{
			Hash:  invalidTxHash,
			Chain: ChainName,
		}
		resp, err := adaptor.GetTxByHash(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		//assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		//assert.Nil(t, resp.Tx)

		//t.Logf("Response Code: %v", resp.Code)
		//t.Logf("Response Message: %s", resp.Msg)
	})

	t.Run("Empty Transaction Hash", func(t *testing.T) {
		req := &account.TxHashRequest{
			Chain: ChainName,
			Hash:  emptyTxHash,
		}
		resp, err := adaptor.GetTxByHash(req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		//assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		//assert.Nil(t, resp.Tx)

		//t.Logf("Response Code: %v", resp.Code)
		//t.Logf("Response Message: %s", resp.Msg)
	})

}

func TestChainAdaptor_GetBlockByRange(t *testing.T) {
	aptosClient, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "failed to initialize aptos aptclient")

	adaptor := ChainAdaptor{
		aptosHttpClient: aptosClient,
	}

	t.Run("Valid Block Range", func(t *testing.T) {
		req := &account.BlockByRangeRequest{
			ConsumerToken: "test_token",
			Chain:         ChainName,
			Network:       "mainnet",
			Start:         "1000",
			End:           "1000",
		}

		resp, err := adaptor.GetBlockByRange(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.Equal(t, "GetBlockByRange success", resp.Msg)
		assert.NotEmpty(t, resp.BlockHeader)

		for _, block := range resp.BlockHeader {
			blockJson, _ := json.Marshal(block)
			t.Logf("resp.BlockHeader blockJson: %s", blockJson)
			assert.NotEmpty(t, block.Hash, "Block hash should not be empty")
			//assert.NotEmpty(t, block.Number, "Block number should not be empty")
			assert.NotEmpty(t, block.TxHash, "Transaction hash should not be empty")
			//assert.NotEmpty(t, block.ReceiptHash, "Receipt hash should not be empty")
			assert.NotZero(t, block.Time, "Block timestamp should not be zero")
			assert.NotNil(t, block.GasLimit, "Gas limit should not be zero")
			assert.NotEmpty(t, block.Extra, "Extra data should not be empty")

			blockNum, err := strconv.ParseUint(block.Number, 10, 64)
			assert.NoError(t, err, "Block number should be valid")
			startNum, _ := strconv.ParseUint(req.Start, 10, 64)
			endNum, _ := strconv.ParseUint(req.End, 10, 64)
			assert.GreaterOrEqual(t, blockNum, startNum)
			assert.LessOrEqual(t, blockNum, endNum)

			t.Logf("Block Hash: %s", block.Hash)
			t.Logf("Block Number: %s", block.Number)
			t.Logf("Block Time: %d", block.Time)
			t.Logf("Gas Used: %d", block.GasUsed)
		}
	})

	t.Run("Invalid Range - Start Greater Than End", func(t *testing.T) {
		req := &account.BlockByRangeRequest{
			ConsumerToken: "test_token",
			Chain:         ChainName,
			Network:       "mainnet",
			Start:         "1000",
			End:           "999",
		}

		_, err := adaptor.GetBlockByRange(req)

		assert.Error(t, err)
		//assert.Contains(t, err.Error(), "start version cannot be greater than end version")
	})

	t.Run("Invalid Block Number Format", func(t *testing.T) {
		req := &account.BlockByRangeRequest{
			ConsumerToken: "test_token",
			Chain:         ChainName,
			Network:       "mainnet",
			Start:         "invalid",
			End:           "1000",
		}

		_, err := adaptor.GetBlockByRange(req)

		assert.Error(t, err)
		//assert.Contains(t, err.Error(), "invalid start version")
	})

	t.Run("Empty Block Range", func(t *testing.T) {
		req := &account.BlockByRangeRequest{
			ConsumerToken: "test_token",
			Chain:         ChainName,
			Network:       "mainnet",
			Start:         "",
			End:           "",
		}

		_, err := adaptor.GetBlockByRange(req)

		assert.Error(t, err)
		//assert.Contains(t, err.Error(), "invalid start version")
	})

	t.Run("Invalid Chain", func(t *testing.T) {
		req := &account.BlockByRangeRequest{
			ConsumerToken: "test_token",
			Chain:         "InvalidChain",
			Network:       "mainnet",
			Start:         "1000",
			End:           "1005",
		}

		resp, err := adaptor.GetBlockByRange(req)

		assert.Error(t, err)
		assert.NotNil(t, resp)
		//assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		//assert.Equal(t, "invalid chain", resp.Msg)
		//assert.Empty(t, resp.BlockHeader)
	})

	t.Run("Latest Blocks", func(t *testing.T) {
		latestReq := &account.BlockNumberRequest{
			Chain:  ChainName,
			Height: 0,
		}
		latestResp, err := adaptor.GetBlockByNumber(latestReq)
		assert.NoError(t, err)

		latestHeight := latestResp.Height
		startHeight := latestHeight - 5

		req := &account.BlockByRangeRequest{
			ConsumerToken: "test_token",
			Chain:         ChainName,
			Network:       "mainnet",
			Start:         strconv.FormatInt(startHeight, 10),
			End:           strconv.FormatInt(latestHeight, 10),
		}

		resp, err := adaptor.GetBlockByRange(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.NotEmpty(t, resp.BlockHeader)

		lastBlock := resp.BlockHeader[len(resp.BlockHeader)-1]
		lastBlockNum, _ := strconv.ParseInt(lastBlock.Number, 10, 64)
		assert.Equal(t, latestHeight, lastBlockNum)
	})
}

func TestChainAdaptor_DecodeTransaction(t *testing.T) {

}

func TestChainAdaptor_VerifySignedTransaction(t *testing.T) {

}

func TestChainAdaptor_GetExtraData(t *testing.T) {

}
