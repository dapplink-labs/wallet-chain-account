package ton

import (
	"context"
	"fmt"
	"github.com/xssnick/tonutils-go/tlb"

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
)

// TonClient 封装了与 TON 区块链交互的客户端
type TonClient struct {
	api *ton.APIClient
}

// DialTonClient 初始化并返回一个 TonClient 实例
func DialTonClient(ctx context.Context, rpcUrl string) (*TonClient, error) {
	client := liteclient.NewConnectionPool()
	// 从配置 URL 添加连接
	err := client.AddConnectionsFromConfigUrl(ctx, rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to add connections: %w", err)
	}
	api := ton.NewAPIClient(client)
	return &TonClient{api: api}, nil
}

// GetBlockByNumber 根据区块高度获取区块
func (c *TonClient) GetBlockByNumber(number uint32) (*tlb.Block, error) {
	ctx := context.Background()
	// 获取主链最新区块信息
	masterInfo, err := c.api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取主链信息失败: %w", err)
	}
	// 查找指定高度的区块
	blockID, err := c.api.LookupBlock(ctx, masterInfo.Workchain, masterInfo.Shard, number)
	if err != nil {
		return nil, fmt.Errorf("查找区块失败: %w", err)
	}
	// 获取区块详细信息
	block, err := c.api.GetBlockData(ctx, blockID)

	if err != nil {
		return nil, fmt.Errorf("获取区块失败: %w", err)
	}
	return block, nil
}
