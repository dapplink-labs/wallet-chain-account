package solana

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dapplink-labs/wallet-chain-account/chain"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/dapplink-labs/wallet-chain-account/rpc/account"
	common2 "github.com/dapplink-labs/wallet-chain-account/rpc/common"
)

func setup() (chain.IChainAdaptor, error) {
	conf, err := config.New("../../config.yml")
	if err != nil {
		log.Error("load config failed, error:", err)
		return nil, err
	}
	adaptor, err := NewChainAdaptor(conf)
	if err != nil {
		log.Error("create chain adaptor failed, error:", err)
		return nil, err
	}
	return adaptor, nil
}

func Test_GetSupportChains(t *testing.T) {
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
		validPublicKey        = "7e376c64c64e88054b7a2d25dc716f45551d2f796ddc9e7be405e49c522b887c"
		validPublicKeyAddress = "9VhPRjzizPY95TyBrve7heeJTZnofgkQYJpLxRSZGZ3H"
		invalidPublicKey      = "invalid_hex"
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
		validAddress   = "9VhPRjzizPY95TyBrve7heeJTZnofgkQYJpLxRSZGZ3H"
		invalidAddress = "invalid_address"
	)

	adaptor := &ChainAdaptor{}

	t.Run("Valid Address", func(t *testing.T) {
		req := &account.ValidAddressRequest{
			Chain:   ChainName,
			Network: "mainnet",
			Address: validAddress,
		}
		resp, err := adaptor.ValidAddress(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.Equal(t, true, resp.Valid)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Is Valid: %v", resp.Valid)
	})

	t.Run("Invalid Address", func(t *testing.T) {
		req := &account.ValidAddressRequest{
			Chain:   ChainName,
			Network: "mainnet",
			Address: invalidAddress,
		}
		resp, err := adaptor.ValidAddress(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		assert.Equal(t, false, resp.Valid)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Is Valid: %v", resp.Valid)
	})

	t.Run("Empty Address", func(t *testing.T) {
		req := &account.ValidAddressRequest{
			Chain:   ChainName,
			Network: "mainnet",
			Address: "",
		}
		resp, err := adaptor.ValidAddress(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		assert.Equal(t, false, resp.Valid)

		t.Logf("Response Code: %v", resp.Code)
		t.Logf("Response Message: %s", resp.Msg)
		t.Logf("Is Valid: %v", resp.Valid)
	})

}

func TestChainAdaptor_GetBlockByNumber(t *testing.T) {
	solClient, err := NewSolHttpClientAll(baseURL, withDebug)
	assert.NoError(t, err, "failed to initialize sol solclient")

	adaptor := &ChainAdaptor{
		solCli: solClient,
	}

	t.Run("Valid Block Number", func(t *testing.T) {
		req := &account.BlockNumberRequest{
			Chain:  ChainName,
			Height: 300944802,
			ViewTx: true,
		}

		resp, err := adaptor.GetBlockByNumber(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.NotEmpty(t, resp.Hash)
		assert.Equal(t, int64(300944802), resp.Height)

		assert.NotNil(t, resp.Transactions)
		if len(resp.Transactions) > 0 {
			tx := resp.Transactions[0]
			assert.NotEmpty(t, tx.Hash)

			t.Logf("Transaction Hash: %s", tx.Hash)
			t.Logf("From Address: %s", tx.From)
			t.Logf("To Address: %s", tx.To)
			t.Logf("Amount: %s", tx.Amount)
		}

		t.Logf("Block Height: %d", resp.Height)
		t.Logf("Block Hash: %s", resp.Hash)
		t.Logf("Transaction Count: %d", len(resp.Transactions))

		respJson, err := json.Marshal(resp)
		assert.NoError(t, err)
		t.Logf("Block json: %s", string(respJson))
	})

	t.Run("Zero Block Number", func(t *testing.T) {
		req := &account.BlockNumberRequest{
			Chain:  ChainName,
			Height: 0,
			ViewTx: true,
		}

		resp, err := adaptor.GetBlockByNumber(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
		assert.NotEmpty(t, resp.Hash)
		t.Logf("Genesis Block Hash: %s", resp.Hash)

		respJson, err := json.Marshal(resp)
		assert.NoError(t, err)
		t.Logf("Block json: %s", string(respJson))
	})

	t.Run("Invalid Block Number", func(t *testing.T) {
		req := &account.BlockNumberRequest{
			Chain:  ChainName,
			Height: 999999999999,
			ViewTx: true,
		}

		resp, err := adaptor.GetBlockByNumber(req)

		assert.Error(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, common2.ReturnCode_ERROR, resp.Code)
		assert.NotEmpty(t, resp.Msg)

		t.Logf("Error Message: %s", resp.Msg)
	})

}

// tx, err := solana.NewTransaction(
// []solana.Instruction{
// system.NewTransferInstruction(
// value,
// fromPubkey,
// toPubkey,
// ).Build(),
// },
// solana.HashFromBytes(binary.BigEndian.AppendUint64(make([]byte, 24), data.Nonce)),
// solana.TransactionPayer(fromPubkey),
func TestChainAdaptor_CreateUnSignTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.CreateUnSignTransaction(&account.UnSignTransactionRequest{
		Chain:    ChainName,
		Network:  "mainnet",
		Base64Tx: createTestBase64Tx(),
	})
	if err != nil {
		log.Error("CreateUnSignTransaction failed:", err)
		return
	}

	assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.UnSignTx)
}

func TestChainAdaptor_BuildSignedTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.BuildSignedTransaction(&account.SignedTransactionRequest{
		Chain:    ChainName,
		Network:  "mainnet",
		Base64Tx: createTestBase64Tx(),
	})
	if err != nil {
		log.Error("TestChainAdaptor_BuildSignedTransaction failed:", err)
		return
	}

	assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
	fmt.Println(resp.SignedTx)
}

func TestChainAdaptor_VerifySignedTransaction(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}

	resp, err := adaptor.VerifySignedTransaction(&account.VerifyTransactionRequest{
		Chain:     ChainName,
		Network:   "mainnet",
		Signature: "5Gqhswe3hEAEoXd499Pu1KTAz6sWJstDdzxT6Jpf8J1FWuomGEUyZeScH9ekQ1zCZAqpENYmEfEXyyJPsFyfmcDNU1pEAWxKWSKAesFXFjnVsSP3KpBcHcSGMAcfkmorrZD6dj1qjtkQsURGJYbe48t8Q4i6soQw6DYsGQMiVn44MyxGEainBDD64nYfGH7oagXLWQRqpN3M3tK9H1VC9QAP2GqNVcinHWx68Li8NdNhpQSJXP8evbhaGc7CojxS7esHpibJ4DD4f85Riy8oWSYuxAwRo3opz8vJ1dDZCM4UbhSwAndBfH6ZGStg6cjgUrfX7c94pTzoXM",
	})
	if err != nil {
		log.Error("TestChainAdaptor_VerifySignedTransaction failed:", err)
		return
	}

	assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
}

func TestChainAdaptor_SendTx(t *testing.T) {
	adaptor, err := setup()
	if err != nil {
		return
	}
	RawTx := "5AvgiBgcnTZh24LUUGbZ95xT4e2tUwR17AL7XnUZ79fHwcpWRStkMZaiVVv8U93Rh9H9giQbTpF4FJxvcFX766pbeHwXPzzL6EekSmn2qV1nM6SSUQY6Qk18FUpQRFqEQiSj54E468hRj9EZ3oneJhLB7Dn2tQVrJTXZmVSZYm3dMS1cc9NhQErK5Tk28j4VNkiBraHHxmDPxkfhLCcDc5EZ6f7hQEPnGS4s7S5pWg63LLdAwxA525NEsRoQbaTUYFQQYU7YotnbhnHsFWvDoiAdkwQsC4H32Cs9LTz7SuowCFyAgaYSuKYyfqGem2a3Cn6Lwfo7f5P1cb"
	log.Info("1:", RawTx)
	resp, err := adaptor.SendTx(&account.SendTxRequest{
		Chain:   ChainName,
		Network: "mainnet",
		RawTx:   RawTx,
	})
	if err != nil {
		log.Error("TestChainAdaptor_SendTx failed:", err)
		return
	}
	log.Info(resp.TxHash)
	assert.Equal(t, common2.ReturnCode_SUCCESS, resp.Code)
}
func createTestBase64Tx() string {

	testTx := TxStructure{
		Nonce:       "DZW54pmZEFK5hQmbygWmyCMmzh42Yzr8o4bvMqAf9UAh",
		FromAddress: "HhXh35Udy8ZUzVSyhptq51xViyshHyFYkiearNFaVvwE",
		ToAddress:   "EUVrmoaKaSsHNkMFw7mVARR522wwH41BFRMha3WC8gha",
		Value:       "0.001",
		//ContractAddress: "So11111111111111111111111111111111111111112", //5VzPuctbhMdqZBpxgxHCyH41sSckqPEKZ7qxbdgMN29Fbvmnpy3x6GcmUFxFw98oy3LcEEVCxwdr4gyQwcboSW6C
		ContractAddress: "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB", //3L64aQvAmdhbaZJFdWXTSjLgmH1GwBhNE8eezqCFAHRvj9a76bwXoarivTSjzAJLiJ48CxtZ5Zke3djnfhuckKs
		Signature:       "b3fa3fed06877eb99810221b900c79b972f2ba9374506a04cbc2ae8c7e92effe19829035cb949bf7ec60852fa91e021c628420775aa20dbe1159c26d274c900a",
	}

	jsonBytes, err := json.Marshal(testTx)
	if err != nil {
		panic(err)
	}

	base64Str := base64.StdEncoding.EncodeToString(jsonBytes)
	return base64Str
}
