package internet_computer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/dapplink-labs/wallet-chain-account/config"
)

func setup() (*IcpClient, error) {
	conf, err := config.New("../../config.yml")
	if err != nil {
		fmt.Println("load config failed, error:", err)
		return nil, err
	}

	icpClient, err := NewIcpClient(context.Background(), conf.WalletNode.Icp.RpcUrl, conf.WalletNode.Icp.TimeOut)
	if err != nil {
		fmt.Println("New icp client failed:", err)
		return nil, err
	}
	return icpClient, nil
}

func TestIcpClient_ConvertAddress(t *testing.T) {
	icpClient, err := setup()
	if err != nil {
		t.Fatal("set up failed:", err)
	}

	pubKeyHex := "04b3c785bc09c5b0d32799194c40460f99137373bc20876b00d83f5789dd8752797e959635d8e05b33ef1f19bb1f368c45602d24ea97cce8548757fcce40a3aeed"
	pubKey, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		fmt.Println("============================================")
	}
	resp, err := icpClient.ConstructionDerive(context.Background(), pubKey, types.Secp256k1)
	if err != nil {
		t.Fatal("Convert address failed", err)
	}

	bytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Fatal("json marshal failed:", err)
	}
	fmt.Printf(string(bytes))
}

func TestIcpClient_FetchBlocksByIndex(t *testing.T) {
	icpClient, err := setup()
	if err != nil {
		t.Fatal("set up failed:", err)
	}

	resp, err := icpClient.FetchBlocksByIndex(context.Background(), 124)
	if err != nil {
		t.Fatal("Get block by number failed", err)
	}

	bytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Fatal("json marshal failed:", err)
	}
	fmt.Printf(string(bytes))
}

func TestIcpClient_FetchBlocksByHash(t *testing.T) {
	icpClient, err := setup()
	if err != nil {
		t.Fatal("set up failed:", err)
	}

	// This is a test net block hash
	hash := "bcd5f7fbe398496566234bf7b2fa226d9bc85c9e7b72bc671a53b6adccac1fae"
	resp, err := icpClient.FetchBlocksByHash(context.Background(), hash)
	if err != nil {
		t.Fatal("Get block by hash failed", err)
	}

	bytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Fatal("json marshal failed:", err)
	}
	fmt.Printf(string(bytes))
}

func TestIcpClient_FetchAccountBalances(t *testing.T) {
	icpClient, err := setup()
	if err != nil {
		t.Fatal("set up failed:", err)
	}

	address := "cf84cdfd049f5905bb7373ef0d895904e84915d5d968fb820f2f2ff8765de669"
	resp, err := icpClient.FetchAccountBalances(context.Background(), address)
	if err != nil {
		t.Fatal("Get account failed", err)
	}

	bytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		t.Fatal("json marshal failed:", err)
	}
	fmt.Printf(string(bytes))
}
