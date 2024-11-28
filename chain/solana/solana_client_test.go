package solana

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

const (
	baseURL   = "https://solana-mainnet.g.alchemy.com/v2/dSe_ey3M3YqwXJtyDnkPFpdNlUtQpafS"
	withDebug = true

	tempAddr = "7mSqVJpb8ziMDB7yDAEajeANyDosh1WK5ksS6mCdDHRE"
)

func TestSolClient_GetAccountInfo(t *testing.T) {
	solClient, err := NewSolHttpClientAll(baseURL, withDebug)
	if err != nil {
		return
	}

	t.Run("成功获取账户信息", func(t *testing.T) {
		resp, err := solClient.GetAccountInfo(tempAddr)
		assert.NoError(t, err)
		assert.NotNil(t, resp)

		respJson, err := json.Marshal(resp)
		if err != nil {
			return
		}
		t.Logf("respJson: %s", string(respJson))
	})

	t.Run("无效地址", func(t *testing.T) {
		resp, err := solClient.GetAccountInfo("invalid_address")
		assert.Error(t, err)
		assert.Nil(t, resp)

		respJson, err := json.Marshal(resp)
		if err != nil {
			return
		}
		t.Logf("respJson: %s", string(respJson))
	})
}

func TestSolClient_GetBalance(t *testing.T) {
	solClient, err := NewSolHttpClient(baseURL)
	if err != nil {
		return
	}

	t.Run("成功获取余额", func(t *testing.T) {
		resp, err := solClient.GetBalance(tempAddr)
		assert.NoError(t, err)
		assert.NotNil(t, resp)

		respJson, err := json.Marshal(resp)
		if err != nil {
			return
		}
		t.Logf("respJson: %s", string(respJson))
	})

	t.Run("无效地址", func(t *testing.T) {
		resp, err := solClient.GetBalance("invalid_address")
		assert.Error(t, err)
		assert.Nil(t, resp)

		respJson, err := json.Marshal(resp)
		if err != nil {
			return
		}
		t.Logf("respJson: %s", string(respJson))
	})

}

func TestSolClient_GetBlockHeight(t *testing.T) {
	//solClient, err := NewSolHttpClient(baseURL)
	//if err != nil {
	//	return
	//}

	//t.Run("成功获取区块高度", func(t *testing.T) {
	//	height, err := solClient.GetBlockHeight()
	//	assert.NoError(t, err)
	//	assert.NotZero(t, height)
	//	t.Logf("Current block height: %d", height)
	//})
}

func TestSolClient_GetSlot(t *testing.T) {
	solClient, err := NewSolHttpClient(baseURL)
	if err != nil {
		return
	}

	t.Run("成功获取区块高度", func(t *testing.T) {
		slot, err := solClient.GetSlot(Finalized)
		assert.NoError(t, err)
		assert.NotZero(t, slot)
		t.Logf("Current slot: %d", slot)
	})
}

func TestSolClient_GetBlocksWithLimit(t *testing.T) {
	solClient, err := NewSolHttpClient(baseURL)
	if err != nil {
		return
	}

	t.Run("成功获取区块列表", func(t *testing.T) {
		startSlot, err := solClient.GetSlot(Finalized)
		assert.NoError(t, err)

		limit := uint64(3)
		blocks, err := solClient.GetBlocksWithLimit(startSlot, limit)

		assert.NoError(t, err)
		assert.NotNil(t, blocks)
		assert.LessOrEqual(t, len(blocks), int(limit))

		respJson, err := json.Marshal(blocks)
		if err != nil {
			return
		}
		t.Logf("respJson: %s", string(respJson))
		t.Logf("Blocks from slot %d (limit %d): %v", startSlot, limit, blocks)
	})
}

func TestSolClient_GetBlockBySlot(t *testing.T) {
	solClient, err := NewSolHttpClientAll(baseURL, withDebug)
	if err != nil {
		return
	}

	t.Run("成功获取区块列表", func(t *testing.T) {
		startSlot, err := solClient.GetSlot(Finalized)
		assert.NoError(t, err)
		t.Logf("startSlot: %s", strconv.FormatUint(startSlot, 10))

		blockResponse, err := solClient.GetBlockBySlot(startSlot, Signatures)
		assert.NoError(t, err)
		assert.NotNil(t, blockResponse)

		respJson, err := json.Marshal(blockResponse)
		assert.NoError(t, err)
		t.Logf("respJson: %s", string(respJson))
	})

	t.Run("成功获取区块列表", func(t *testing.T) {
		blockResponse, err := solClient.GetBlockBySlot(300944802, Signatures)
		assert.NoError(t, err)
		assert.NotNil(t, blockResponse)

		respJson, err := json.Marshal(blockResponse)
		assert.NoError(t, err)
		t.Logf("respJson: %s", string(respJson))
	})

	t.Run("成功获取区块列表", func(t *testing.T) {
		blockResponse, err := solClient.GetBlockBySlot(300944802, Full)
		assert.NoError(t, err)
		assert.NotNil(t, blockResponse)

		respJson, err := json.MarshalIndent(blockResponse, "", "    ")
		t.Logf("respJson: %s", "dadadad")

		assert.NoError(t, err)
		outputDir := "test_output"
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			t.Fatalf("创建输出目录失败: %v", err)
		}
		timestamp := time.Now().Format("20060102_150405")
		filename := filepath.Join(outputDir, fmt.Sprintf("block_response_%s.json", timestamp))

		err = os.WriteFile(filename, respJson, 0644)
		assert.NoError(t, err)
		t.Logf("JSON 已写入文件: %s", filename)
	})
}

func TestSolClient_GetTransaction(t *testing.T) {
	solClient, err := NewSolHttpClientAll(baseURL, withDebug)
	if err != nil {
		return
	}

	signature := "5myqChdtQuwzHWN6dA74qGvuddqc3oB4RKcn6edXYDJWoMW8ZDaWEeA4mgXjNeP5DFRNAbnXKE8HgKcnKPhg1NAN"

	t.Run("成功获取区块列表", func(t *testing.T) {
		txResponse, err := solClient.GetTransaction(signature)
		assert.NoError(t, err)
		assert.NotNil(t, txResponse)

		respJson, err := json.Marshal(txResponse)
		assert.NoError(t, err)
		t.Logf("respJson: %s", string(respJson))
	})

}

func TestSolClient_GetFeeForMessage(t *testing.T) {
	solClient, err := NewSolHttpClientAll(baseURL, withDebug)
	if err != nil {
		return
	}

	signature := "4mnTvsFsEGY3CKhjZ3X2U9wac8J4sS3uHywrHvUKViN8ta2TD22YghLiw2bQG7CmC47kDBfey1Mzw3HkEDAhwkXS"

	t.Run("成功获取区块列表", func(t *testing.T) {
		baseFee, err := solClient.GetFeeForMessage(signature)
		assert.NoError(t, err)
		assert.NotNil(t, baseFee)

		t.Logf("respJson: %s", strconv.FormatUint(baseFee, 10))
	})
}

func TestSolClient_GetRecentPrioritizationFees(t *testing.T) {
	solClient, err := NewSolHttpClientAll(baseURL, withDebug)
	if err != nil {
		return
	}

	t.Run("成功获取区块列表", func(t *testing.T) {
		feeList, err := solClient.GetRecentPrioritizationFees()
		assert.NoError(t, err)
		assert.NotNil(t, feeList)

		// 打印返回的数据
		fmt.Println("\n=== 优先级费用数据 ===")
		for i, fee := range feeList {
			fmt.Printf("数据 %d:\n", i+1)
			fmt.Printf("  Slot: %d\n", fee.Slot)
			fmt.Printf("  优先级费用: %d lamports (%.9f SOL)\n",
				fee.PrioritizationFee,
				float64(fee.PrioritizationFee)/1e9)
			fmt.Println("-------------------")
		}
		fmt.Printf("总共 %d 条数据\n", len(feeList))
		fmt.Println("===================")

		priorityFee := GetSuggestedPriorityFee(feeList)
		t.Logf("priorityFee: %s", strconv.FormatUint(priorityFee, 10))
	})
}

func TestSolClient_GetTxForAddress(t *testing.T) {
	solClient, err := NewSolHttpClientAll(baseURL, withDebug)
	if err != nil {
		return
	}

	address := "9xQeWvG816bUx9EPjHmaT23yvVM2ZWbrrpZb9PusVFin"
	commitment := Finalized
	limit := uint64(10)
	before := ""
	until := ""

	t.Run("成功获取区块列表", func(t *testing.T) {
		signatureList, err := solClient.GetTxForAddress(
			address,
			commitment,
			limit,
			before,
			until,
		)
		assert.NoError(t, err)
		assert.NotNil(t, signatureList)

		respJson, err := json.Marshal(signatureList)
		assert.NoError(t, err)
		t.Logf("respJson: %s", string(respJson))
	})
}

func TestSolClient_GetLatestBlockhash(t *testing.T) {
	solClient, err := NewSolHttpClientAll(baseURL, withDebug)
	if err != nil {
		t.Fatalf("初始化客户端失败: %v", err)
		return
	}

	t.Run("成功获取最新blockhash", func(t *testing.T) {
		response, err := solClient.GetLatestBlockhash(Finalized)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		respJson, err := json.Marshal(response)
		assert.NoError(t, err)
		t.Logf("respJson: %s", string(respJson))
	})

}
