package aptos

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	baseURL   = "https://api.mainnet.aptoslabs.com/"
	apiKey    = "aptoslabs_7Gd8hUMMp85_JxF2SXZCDcmeP4tjuuBXjwFwqyY6nTFup"
	network   = Mainnet
	withDebug = false
)

func TestClient_GetNodeInfo(t *testing.T) {
	// Initialize aptclient
	client, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "aptclient initialization failed")
	assert.NotNil(t, client, "aptclient should not be nil")

	// Test getting node info
	t.Run("Get Node Info", func(t *testing.T) {
		resp, err := client.GetNodeInfo()
		marshal, _ := json.Marshal(resp)
		fmt.Println("GetNodeInfo resp", string(marshal))
		assert.NoError(t, err, "failed to get node info")
		assert.NotNil(t, resp, "response should not be nil")

		// Validate node info data
		t.Run("Node Info Validation", func(t *testing.T) {
			assert.Greater(t, resp.ChainID, uint8(0), "chain id should be greater than 0")
			assert.NotEmpty(t, resp.Epoch, "epoch should not be empty")
			assert.NotEmpty(t, resp.LedgerVersion, "ledger version should not be empty")
			// oldest_ledger_version is "0"
			assert.NotNil(t, resp.OldestLedgerVersion, "oldest ledger version should not be nil")
			assert.NotEmpty(t, resp.LedgerTimestamp, "ledger timestamp should not be empty")
			assert.NotEmpty(t, resp.NodeRole, "node role should not be empty")
			// OldestBlockHeight is "0"
			assert.NotNil(t, resp.OldestBlockHeight, "oldest block height should not be nil")
			assert.NotEmpty(t, resp.BlockHeight, "block height should not be empty")
			assert.NotEmpty(t, resp.GitHash, "git hash should not be empty")
		})

		// Log node info details
		t.Run("Log Details", func(t *testing.T) {
			t.Log("Node Information:")
			t.Logf("Chain ID: %d", resp.ChainID)
			t.Logf("Epoch: %d", resp.Epoch)
			t.Logf("Ledger Version: %d", resp.LedgerVersion)
			t.Logf("Oldest Ledger Version: %d", resp.OldestLedgerVersion)
			t.Logf("Ledger Timestamp: %d", resp.LedgerTimestamp)
			t.Logf("Node Role: %s", resp.NodeRole)
			t.Logf("Oldest Block Height: %d", resp.OldestBlockHeight)
			t.Logf("Block Height: %d", resp.BlockHeight)
			t.Logf("Git Hash: %s", resp.GitHash)
		})
	})
}

func TestRestyClient_GetAccount(t *testing.T) {
	const (
		validAccount   = "0x06671b50c2a5edb709c9e15d7e5a3d6496ae923759a217090e4ba9622720da5e"
		invalidAccount = "0xinvalid_account"
		emptyAccount   = ""
	)

	client, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "failed to initialize aptos aptclient")

	t.Run("Valid Account", func(t *testing.T) {
		accountResponse, err := client.GetAccount(validAccount)

		assert.NoError(t, err)
		assert.NotNil(t, accountResponse)
		//assert.NotZero(t, accountResponse.SequenceNumber)
		assert.NotEmpty(t, accountResponse.AuthenticationKey)

		t.Logf("Sequence Number: %d", accountResponse.SequenceNumber)
		t.Logf("Authentication Key: %s", accountResponse.AuthenticationKey)
	})

	t.Run("Invalid Account", func(t *testing.T) {
		accountResponse, err := client.GetAccount(invalidAccount)

		assert.Error(t, err)
		assert.Nil(t, accountResponse)
		assert.ErrorIs(t, err, errInvalidAddress)

		t.Logf("Error: %v", err)
	})

	t.Run("Empty Account", func(t *testing.T) {
		accountResponse, err := client.GetAccount(emptyAccount)

		assert.Error(t, err)
		assert.Nil(t, accountResponse)
		assert.ErrorIs(t, err, errInvalidAddress)

		t.Logf("Error: %v", err)
	})

	t.Run("Account Not Found", func(t *testing.T) {
		notExistAccount := "0x1234567890123456789012345678901234567890123456789012345678901234"
		accountResponse, err := client.GetAccount(notExistAccount)

		assert.Error(t, err)
		assert.Nil(t, accountResponse)
		assert.ErrorIs(t, err, errHTTPError)

		t.Logf("Error: %v", err)
	})
}

func TestRestyClient_GetGasPrice(t *testing.T) {
	client, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "failed to initialize aptos aptclient")

	t.Run("Get Gas Price Successfully", func(t *testing.T) {
		gasPrice, err := client.GetGasPrice()

		assert.NoError(t, err)
		assert.NotNil(t, gasPrice)

		assert.Greater(t, gasPrice.GasEstimate, uint64(0))
		assert.Greater(t, gasPrice.PrioritizedGasEstimate, gasPrice.GasEstimate)
		assert.LessOrEqual(t, gasPrice.DeprioritizedGasEstimate, gasPrice.GasEstimate)

		assert.Less(t, gasPrice.DeprioritizedGasEstimate, uint64(10000), "Deprioritized gas price should be reasonable")
		assert.Less(t, gasPrice.GasEstimate, uint64(10000), "Normal gas price should be reasonable")
		assert.Less(t, gasPrice.PrioritizedGasEstimate, uint64(10000), "Prioritized gas price should be reasonable")

		t.Logf("Deprioritized Gas Price: %d", gasPrice.DeprioritizedGasEstimate)
		t.Logf("Normal Gas Price: %d", gasPrice.GasEstimate)
		t.Logf("Prioritized Gas Price: %d", gasPrice.PrioritizedGasEstimate)
	})

	t.Run("Price Relationship Check", func(t *testing.T) {
		gasPrice, err := client.GetGasPrice()

		assert.NoError(t, err)
		assert.NotNil(t, gasPrice)

		assert.Greater(t, gasPrice.PrioritizedGasEstimate, gasPrice.GasEstimate,
			"Prioritized price should be higher than normal price")
		assert.LessOrEqual(t, gasPrice.DeprioritizedGasEstimate, gasPrice.GasEstimate,
			"Deprioritized price should be lower than normal price")

		t.Logf("Price Differences:")
		t.Logf("Priority Premium: %d", gasPrice.PrioritizedGasEstimate-gasPrice.GasEstimate)
		t.Logf("Normal Premium: %d", gasPrice.GasEstimate-gasPrice.DeprioritizedGasEstimate)
	})

	t.Run("Invalid AptClient", func(t *testing.T) {
		invalidClient, err := NewAptosHttpClientAll("https://invalid.url", apiKey, withDebug)
		assert.NoError(t, err)

		gasPrice, err := invalidClient.GetGasPrice()

		assert.Error(t, err)
		assert.Nil(t, gasPrice)
		t.Logf("Expected Error: %v", err)
	})
}

func TestClient_GetTransactionByHash(t *testing.T) {
	const (
		validTxHash = "0x43531969ff8e93de962ea65e5609c2b05de3aa5e78933d8925613e75d3d53772"
	)

	client, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "aptclient initialization failed")
	assert.NotNil(t, client, "aptclient should not be nil")

	t.Run("Valid Transaction", func(t *testing.T) {
		resp, err := client.GetTransactionByHash(validTxHash)
		assert.NoError(t, err, "failed to get transaction")
		assert.NotNil(t, resp, "response should not be nil")

		t.Run("Transaction Metadata", func(t *testing.T) {
			assert.Equal(t, "1878359810", fmt.Sprint(resp.Version))
			assert.Equal(t, validTxHash, resp.Hash)
			assert.Equal(t, "999", fmt.Sprint(resp.GasUsed))
			assert.True(t, resp.Success)
			assert.Equal(t, "Executed successfully", resp.VMStatus)
		})

		t.Run("Transaction Changes", func(t *testing.T) {
			assert.NotEmpty(t, resp.Changes, "changes should not be empty")

			change := resp.Changes[0]
			assert.Equal(t, "write_resource", change.Type)
			assert.NotEmpty(t, change.Address)
			assert.NotEmpty(t, change.StateKeyHash)

			//data := change.Data.Data
			//assert.NotEmpty(t, data, "data should not be empty")
			//
			//rewards := data.Rewards
			//assert.NotEmpty(t, rewards, "rewards should not be empty")

			//lastReward := rewards[len(rewards)-1]
			//assert.Equal(t, "1730609200", lastReward.Eid)
			//assert.Equal(t, "12", lastReward.Share)
			//assert.Equal(t, "1730609324", lastReward.Unlock)
		})

		t.Run("Log Details", func(t *testing.T) {
			t.Log("Transaction Details:")
			t.Logf("Version: %d", resp.Version)
			t.Logf("Hash: %s", resp.Hash)
			t.Logf("Gas Used: %d", resp.GasUsed)
			t.Logf("Success: %v", resp.Success)
			t.Logf("VM Status: %s", resp.VMStatus)

			if len(resp.Changes) > 0 {
				t.Log("\nTransaction Changes:")
				marshal, _ := json.Marshal(resp.Changes)
				t.Logf("resp.Changes: %s", string(marshal))
			}
		})
	})

	t.Run("Error Cases", func(t *testing.T) {
		t.Run("Empty Hash", func(t *testing.T) {
			resp, err := client.GetTransactionByHash("")
			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Contains(t, err.Error(), "transaction hash cannot be empty")
		})

		t.Run("Invalid Hash", func(t *testing.T) {
			invalidHash := "0x1234567890abcdef"
			resp, err := client.GetTransactionByHash(invalidHash)
			assert.Error(t, err)
			assert.Nil(t, resp)
		})
	})
}

func TestRestyClient_SubmitTransaction(t *testing.T) {

}

func TestClient_GetBlockByHeight(t *testing.T) {
	// Test configuration
	const (
		validHeight uint64 = 88012315
	)

	// Initialize aptclient
	client, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "aptclient initialization failed")
	assert.NotNil(t, client, "aptclient should not be nil")

	// Test valid block height
	t.Run("Valid Block Height", func(t *testing.T) {
		resp, err := client.GetBlockByHeight(validHeight)
		assert.NoError(t, err, "failed to get block")
		assert.NotNil(t, resp, "response should not be nil")

		// Validate block data
		t.Run("Block Metadata", func(t *testing.T) {
			assert.NotEmpty(t, resp.BlockHeight, "block height should not be empty")
			assert.NotEmpty(t, resp.BlockHash, "block hash should not be empty")
			assert.NotEmpty(t, resp.BlockTimestamp, "block timestamp should not be empty")
			assert.NotEmpty(t, resp.FirstVersion, "first version should not be empty")
			assert.NotEmpty(t, resp.LastVersion, "last version should not be empty")
		})

		// Log block details
		t.Run("Log Details", func(t *testing.T) {
			t.Log("Block Details:")
			t.Logf("Block Height: %d", resp.BlockHeight)
			t.Logf("Block Hash: %s", resp.BlockHash)
			t.Logf("Timestamp: %d", resp.BlockTimestamp)
			t.Logf("First Version: %d", resp.FirstVersion)
			t.Logf("Last Version: %d", resp.LastVersion)
		})
	})

	// Test error cases
	t.Run("Error Cases", func(t *testing.T) {
		// Test invalid block height
		t.Run("Invalid Height", func(t *testing.T) {
			// Use a very large block height that shouldn't exist yet
			invalidHeight := uint64(999999999999)
			resp, err := client.GetBlockByHeight(invalidHeight)
			assert.Error(t, err, "should return an error")
			assert.Nil(t, resp, "response should be nil")
			t.Logf("Expected error: %v", err)
		})
	})
}

func TestClient_GetBlockByVersion(t *testing.T) {
	// Test configuration
	const (
		validVersion uint64 = 248685535
	)

	// Initialize aptclient
	client, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "aptclient initialization failed")
	assert.NotNil(t, client, "aptclient should not be nil")

	// Test valid version
	t.Run("Valid Version", func(t *testing.T) {
		resp, err := client.GetBlockByVersion(validVersion)
		assert.NoError(t, err, "failed to get block")
		assert.NotNil(t, resp, "response should not be nil")

		// Validate block data
		t.Run("Block Metadata", func(t *testing.T) {
			assert.NotEmpty(t, resp.BlockHeight, "block height should not be empty")
			assert.NotEmpty(t, resp.BlockHash, "block hash should not be empty")
			assert.NotEmpty(t, resp.BlockTimestamp, "block timestamp should not be empty")
			assert.NotEmpty(t, resp.FirstVersion, "first version should not be empty")
			assert.NotEmpty(t, resp.LastVersion, "last version should not be empty")

			// Verify version is within block range
			firstVersion := resp.FirstVersion
			lastVersion := resp.LastVersion
			assert.LessOrEqual(t, firstVersion, validVersion,
				"first version should be less than or equal to requested version")
			assert.GreaterOrEqual(t, lastVersion, validVersion,
				"last version should be greater than or equal to requested version")
		})

		// Log block details
		t.Run("Log Details", func(t *testing.T) {
			t.Log("Block Details:")
			t.Logf("Block Height: %d", resp.BlockHeight)
			t.Logf("Block Hash: %s", resp.BlockHash)
			t.Logf("Timestamp: %d", resp.BlockTimestamp)
			t.Logf("First Version: %d", resp.FirstVersion)
			t.Logf("Last Version: %d", resp.LastVersion)

			//if len(resp.Transactions) > 0 {
			//	t.Log("\nFirst Transaction Details:")
			//	t.Logf("Type: %s", resp.Transactions[0].Type)
			//	t.Logf("Hash: %s", resp.Transactions[0].Hash)
			//	t.Logf("Sender: %s", resp.Transactions[0].Sender)
			//}
		})
	})

	// Test error cases
	t.Run("Error Cases", func(t *testing.T) {
		// Test invalid version
		t.Run("Invalid Version", func(t *testing.T) {
			// Use a very large version that shouldn't exist yet
			invalidVersion := uint64(999999999999)
			resp, err := client.GetBlockByVersion(invalidVersion)
			assert.Error(t, err, "should return an error")
			assert.Nil(t, resp, "response should be nil")
			t.Logf("Expected error: %v", err)
		})

		// Test zero version
		t.Run("Zero Version", func(t *testing.T) {
			resp, err := client.GetBlockByVersion(0)
			assert.NoError(t, err, "should return an error for version 0")
			assert.NotNil(t, resp, "response should be nil")
			t.Logf("Expected error: %v", err)
		})
	})
}

func TestRestyClient_GetTransactionByHash(t *testing.T) {
	const (
		validTxHash   = "0x43531969ff8e93de962ea65e5609c2b05de3aa5e78933d8925613e75d3d53772"
		invalidTxHash = "0xinvalid_hash"
		emptyTxHash   = ""
	)

	client, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "failed to initialize aptos aptclient")

	t.Run("Valid Transaction Hash", func(t *testing.T) {
		resp, err := client.GetTransactionByHash(validTxHash)

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		t.Run("Basic Fields", func(t *testing.T) {
			assert.NotZero(t, resp.Version)
			assert.Equal(t, validTxHash, resp.Hash)
			assert.NotEmpty(t, resp.StateChangeHash)
			assert.NotEmpty(t, resp.EventRootHash)
			assert.NotZero(t, resp.GasUsed)
			assert.True(t, resp.Success)
			assert.Equal(t, "Executed successfully", resp.VMStatus)
			assert.NotEmpty(t, resp.AccumulatorRootHash)
		})

		t.Run("Transaction Details", func(t *testing.T) {
			assert.NotEmpty(t, resp.Sender)
			assert.NotZero(t, resp.SequenceNumber)
			assert.NotZero(t, resp.MaxGasAmount)
			assert.NotZero(t, resp.GasUnitPrice)
			assert.NotZero(t, resp.ExpirationTimestamp)
			assert.NotZero(t, resp.Timestamp)
			assert.NotEmpty(t, resp.Type)
		})

		t.Run("Payload", func(t *testing.T) {
			assert.NotEmpty(t, resp.Payload.Function)
			assert.NotEmpty(t, resp.Payload.Type)
		})

		t.Run("Changes", func(t *testing.T) {
			if len(resp.Changes) > 0 {
				change := resp.Changes[0]
				assert.NotEmpty(t, change.Address)
				assert.NotEmpty(t, change.StateKeyHash)
				assert.NotEmpty(t, change.Type)
			}
		})

		t.Run("Events", func(t *testing.T) {
			if len(resp.Events) > 0 {
				event := resp.Events[0]
				assert.NotEmpty(t, event.Guid.AccountAddress)
				assert.NotEmpty(t, event.SequenceNumber)
				assert.NotEmpty(t, event.Type)
			}
		})

		t.Run("Log Details", func(t *testing.T) {
			t.Log("Transaction Details:")
			t.Logf("Version: %d", resp.Version)
			t.Logf("Hash: %s", resp.Hash)
			t.Logf("Sender: %s", resp.Sender)
			t.Logf("Gas Used: %d", resp.GasUsed)
			t.Logf("Success: %v", resp.Success)
			t.Logf("VM Status: %s", resp.VMStatus)
			t.Logf("Timestamp: %d", resp.Timestamp)

			if len(resp.Changes) > 0 {
				t.Log("\nChanges:")
				for i, change := range resp.Changes {
					t.Logf("Change %d:", i+1)
					t.Logf("  Address: %s", change.Address)
					t.Logf("  Type: %s", change.Type)
				}
			}

			if len(resp.Events) > 0 {
				t.Log("\nEvents:")
				for i, event := range resp.Events {
					t.Logf("Event %d:", i+1)
					t.Logf("  Type: %s", event.Type)
					t.Logf("  Account: %s", event.Guid.AccountAddress)
				}
			}
		})
	})

	t.Run("Error Cases", func(t *testing.T) {
		t.Run("Invalid Hash", func(t *testing.T) {
			resp, err := client.GetTransactionByHash(invalidTxHash)

			assert.Error(t, err)
			assert.Nil(t, resp)
			t.Logf("Expected Error: %v", err)
		})

		t.Run("Empty Hash", func(t *testing.T) {
			resp, err := client.GetTransactionByHash(emptyTxHash)

			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Contains(t, err.Error(), "transaction hash cannot be empty")
			t.Logf("Expected Error: %v", err)
		})

		t.Run("Non-existent Hash", func(t *testing.T) {
			nonExistentHash := "0x1234567890123456789012345678901234567890123456789012345678901234"
			resp, err := client.GetTransactionByHash(nonExistentHash)

			assert.Error(t, err)
			assert.Nil(t, resp)
			t.Logf("Expected Error: %v", err)
		})
	})
}

func TestClient_GetTransactionByVersion(t *testing.T) {
	// Test configuration
	const (
		// Use a known transaction version for testing
		validVersion = "1878359810"
	)

	// Initialize aptclient
	client, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "aptclient initialization failed")
	assert.NotNil(t, client, "aptclient should not be nil")

	t.Run("Valid Transaction Version", func(t *testing.T) {
		resp, err := client.GetTransactionByVersion(validVersion)
		assert.NoError(t, err, "failed to get transaction")
		assert.NotNil(t, resp, "response should not be nil")
		json1, _ := json.Marshal(resp)
		fmt.Println("GetTransactionByVersion", string(json1))

		t.Run("Transaction Metadata", func(t *testing.T) {
			assert.Equal(t, validVersion, fmt.Sprint(resp.Version), "version should match")
			assert.NotEmpty(t, resp.Hash)
			assert.NotNil(t, resp.GasUsed)
			assert.True(t, resp.Success)
			assert.Equal(t, "Executed successfully", resp.VMStatus)
		})

		t.Run("Transaction Changes", func(t *testing.T) {
			assert.NotEmpty(t, resp.Changes, "changes should not be empty")

			change := resp.Changes[0]
			assert.Equal(t, "write_resource", change.Type)
			assert.NotEmpty(t, change.Address)
			assert.NotEmpty(t, change.StateKeyHash)
		})

		t.Run("Log Details", func(t *testing.T) {
			t.Log("Transaction Details:")
			t.Logf("Version: %d", resp.Version)
			t.Logf("Hash: %s", resp.Hash)
			t.Logf("Gas Used: %d", resp.GasUsed)
			t.Logf("Success: %v", resp.Success)
			t.Logf("VM Status: %s", resp.VMStatus)

			if len(resp.Changes) > 0 {
				t.Log("\nTransaction Changes:")
				marshal, _ := json.Marshal(resp.Changes)
				t.Logf("resp.Changes: %s", string(marshal))
			}
		})
	})

	t.Run("Error Cases", func(t *testing.T) {
		t.Run("Empty Version", func(t *testing.T) {
			resp, err := client.GetTransactionByVersion("")
			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Contains(t, err.Error(), "version cannot be empty")
		})

		t.Run("Invalid Version Format", func(t *testing.T) {
			resp, err := client.GetTransactionByVersion("invalid-version")
			assert.Error(t, err)
			assert.Nil(t, resp)
		})

		t.Run("Non-existent Version", func(t *testing.T) {
			resp, err := client.GetTransactionByVersion("999999999999999")
			assert.Error(t, err)
			assert.Nil(t, resp)
		})
	})

	//t.Run("Version Zero", func(t *testing.T) {
	//	resp, err := client.GetTransactionByVersion("0")
	//	assert.NoError(t, err)
	//	assert.NotNil(t, resp)
	//	assert.Equal(t, "0", fmt.Sprint(resp.Version))
	//
	//	t.Log("Genesis Transaction Details:")
	//	t.Logf("Version: %d", resp.Version)
	//	t.Logf("Type: %s", resp.Type)
	//	t.Logf("Hash: %s", resp.Hash)
	//	if len(resp.Changes) > 0 {
	//		t.Logf("Number of Changes: %d", len(resp.Changes))
	//	}
	//})
}

func TestClient_GetTransactionByVersionRange(t *testing.T) {
	// Test configuration
	const (
		validVersion = 1881899111
		endVersion   = 1881811111
	)

	// Initialize aptclient
	client, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "aptclient initialization failed")
	assert.NotNil(t, client, "aptclient should not be nil")

	t.Run("Single Version", func(t *testing.T) {
		startVersion := uint64(validVersion)
		endVersion := startVersion

		txs, err := client.GetTransactionByVersionRange(startVersion, endVersion)
		assert.NoError(t, err)
		assert.NotNil(t, txs)
		assert.Len(t, txs, 1)

		// Convert validVersion to string for comparison
		assert.Equal(t, fmt.Sprint(validVersion), fmt.Sprint(txs[0].Version), "version should match")
		t.Logf("Successfully retrieved single transaction at version %d", startVersion)
	})

	t.Run("Small Range", func(t *testing.T) {
		startVersion := uint64(validVersion)
		endVersion := startVersion + 5

		txs, err := client.GetTransactionByVersionRange(startVersion, endVersion)
		assert.NoError(t, err)
		assert.NotNil(t, txs)
		assert.Len(t, txs, 6) // inclusive range

		// Verify transactions are in order
		for i, tx := range txs {
			expectedVersion := fmt.Sprintf("%d", startVersion+uint64(i))
			assert.Equal(t, expectedVersion, fmt.Sprint(tx.Version))
		}

		t.Logf("Successfully retrieved %d transactions", len(txs))
	})

	t.Run("Medium Range", func(t *testing.T) {
		startVersion := uint64(validVersion)
		endVersion := startVersion + 99 // Test with groupSize

		txs, err := client.GetTransactionByVersionRange(startVersion, endVersion)
		assert.NoError(t, err)
		assert.NotNil(t, txs)
		assert.Len(t, txs, 100)

		t.Logf("Successfully retrieved %d transactions", len(txs))
	})

	t.Run("Error Cases", func(t *testing.T) {
		t.Run("Invalid Range - Start > End", func(t *testing.T) {
			startVersion := uint64(validVersion)
			endVersion := uint64(endVersion)

			txs, err := client.GetTransactionByVersionRange(startVersion, endVersion)
			assert.Error(t, err)
			assert.Nil(t, txs)
			assert.Contains(t, err.Error(), "start version")
		})

		t.Run("Non-existent Version", func(t *testing.T) {
			startVersion := uint64(999999999999)
			endVersion := startVersion + 5

			txs, err := client.GetTransactionByVersionRange(startVersion, endVersion)
			assert.Error(t, err)
			assert.Nil(t, txs)
		})
	})

	t.Run("Log Transaction Details", func(t *testing.T) {
		startVersion := uint64(validVersion)
		endVersion := startVersion + 2

		txs, err := client.GetTransactionByVersionRange(startVersion, endVersion)
		assert.NoError(t, err)

		for i, tx := range txs {
			t.Logf("\nTransaction %d Details:", i+1)
			t.Logf("Version: %v", tx.Version)
			t.Logf("Hash: %s", tx.Hash)
			t.Logf("Success: %v", tx.Success)
			t.Logf("VM Status: %s", tx.VMStatus)
			if len(tx.Changes) > 0 {
				t.Logf("Number of Changes: %d", len(tx.Changes))
			}
		}
	})
}

func TestRestyClient_GetAccountBalance(t *testing.T) {
	const (
		validAddress    = "0x06671b50c2a5edb709c9e15d7e5a3d6496ae923759a217090e4ba9622720da5e"
		expectedAPT     = 0.68374979
		ResourceTypeAPT = "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>"
	)
	client, err := NewAptosHttpClientAll(baseURL, apiKey, withDebug)
	assert.NoError(t, err, "aptclient initialization failed")
	assert.NotNil(t, client, "aptclient should not be nil")

	t.Run("Get APT Balance", func(t *testing.T) {
		address := validAddress
		resourceType := ResourceTypeAPT

		balance, err := client.GetAccountBalance(address, resourceType)
		t.Logf("Account %s APT balance: %d", address, balance)

		assert.NoError(t, err, "should not return error")
		assert.NotZero(t, balance, "balance should not be zero")

		aptValue := float64(balance) / 100000000
		t.Logf("Account %s APT balance: %.8f", address, aptValue)

		assert.InDelta(t, expectedAPT, aptValue, 0.00000001, "APT balance should match expected amount")
	})
}
