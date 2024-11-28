package cosmos

import (
	"context"
	"flag"
	"fmt"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	"strconv"
	"testing"

	authv1beta1 "cosmossdk.io/api/cosmos/auth/v1beta1"
	"github.com/dapplink-labs/wallet-chain-account/config"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
)

const (
	defaultRpcAddress = "https://cosmos-rpc.publicnode.com:443"
)

func getWalletConfig() (*config.Config, error) {
	var f = flag.String("c", defaultConfigPath, "config path")
	flag.Parse()
	return config.New(*f)
}

// success
func TestClient_GetAccount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultDialTimeout)
	defer cancel()

	config, _ := getWalletConfig()
	c, err := DialCosmosClient(ctx, config)
	assert.NoError(t, err)

	response, err := c.GetAccount("cosmos19thxsunl9lzywglsndth5a278wtavawzzpv44q")
	assert.NoError(t, err)

	authAccount := new(authv1beta1.BaseAccount)
	err = ptypes.UnmarshalAny(response.Account, authAccount)
	assert.NoError(t, err)
	fmt.Printf("sequence: %s, account number: %s, address: %s \n",
		strconv.FormatUint(authAccount.GetSequence(), 10),
		strconv.FormatUint(authAccount.GetAccountNumber(), 10),
		authAccount.GetAddress())
}

// success
func TestClient_GetBalance(t *testing.T) {
	config, _ := getWalletConfig()
	c, err := DialCosmosClient(context.Background(), config)
	assert.NoError(t, err)

	balance, err := c.GetBalance("uatom", "cosmos19thxsunl9lzywglsndth5a278wtavawzzpv44q")
	assert.NoError(t, err)
	fmt.Printf("amaount: %s, denom: %s \n", balance.Amount, balance.GetDenom())
}

// success
func TestClient_GetTxByHash(t *testing.T) {
	config, _ := getWalletConfig()
	c, err := DialCosmosClient(context.Background(), config)
	assert.NoError(t, err)

	response, err := c.GetTxByHash("85C84677F466D71C0BB6E744439C3040ABB35B8F2B838CC7B73CD1BFF33D0B88")
	assert.NoError(t, err)
	fmt.Printf("result: %s \n", response.Hash.String())
}

// success
func TestClient_GetBlock(t *testing.T) {
	config, _ := getWalletConfig()
	c, err := DialCosmosClient(context.Background(), config)
	assert.NoError(t, err)

	height := int64(22879895)
	block, err := c.GetBlock(height)
	assert.NoError(t, err)
	fmt.Printf("hash: %s \n", block.BlockID.String())
}

// success
func TestClient_GetBlockByHash(t *testing.T) {
	config, _ := getWalletConfig()
	c, err := DialCosmosClient(context.Background(), config)
	assert.NoError(t, err)

	response, err := c.GetBlockByHash("35290F91317064307B5B1A9A44EEFB1CF3F66F68EAAD4539BCD6A5BA13866E90")
	assert.NoError(t, err)
	fmt.Printf("result: %s \n", response.BlockID.String())
}

// success
func TestClient_GetHeaderByHeight(t *testing.T) {
	config, _ := getWalletConfig()
	c, err := DialCosmosClient(context.Background(), config)
	assert.NoError(t, err)

	height := int64(22879895)
	response, err := c.GetHeaderByHeight(height)
	assert.NoError(t, err)
	fmt.Printf("result: %s \n", response.Header.Hash().String())
}

// success
func TestClient_GetHeaderByHash(t *testing.T) {
	config, _ := getWalletConfig()
	c, err := DialCosmosClient(context.Background(), config)
	assert.NoError(t, err)

	response, err := c.GetHeaderByHash("35290F91317064307B5B1A9A44EEFB1CF3F66F68EAAD4539BCD6A5BA13866E90")
	assert.NoError(t, err)
	fmt.Printf("result: %s \n", response.Header.DataHash.String())
}

// success
func TestClient_BlockchainInfo(t *testing.T) {
	config, _ := getWalletConfig()
	c, err := DialCosmosClient(context.Background(), config)
	assert.NoError(t, err)

	response, err := c.BlockchainInfo(22879895, 22879896)
	assert.NoError(t, err)
	fmt.Printf("result: %v \n", response.LastHeight)
}

func TestClient_BroadcastTx(t *testing.T) {
	config, _ := getWalletConfig()
	c, err := DialCosmosClient(context.Background(), config)
	assert.NoError(t, err)

	var txMap map[string]interface{} = make(map[string]interface{})
	txMsg := "CpoBCo8BChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEm8KLWNvc21vczE5dGh4c3VubDlsenl3Z2xzbmR0aDVhMjc4d3RhdmF3enpwdjQ0cRItY29zbW9zMWw2dnVsMjBxNzRndzU2ZnBlZDhzcmtqcTJ4OGQ5bTMwNWdueHIyGg8KBXVhdG9tEgYxMDAwMDASBjEwMTExMRJoClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDi7ANMzipos1bc45s2u2g+ar2Fu1+Z8lkzFTUtoDjIccSBAoCCAEYChIUCg4KBXVhdG9tEgUxMDAwMBDKswgaQAW5UeOw1oNp6SJQCwbVc10wdBB6lJ1MGVRuTA2i8lUtZbzgeYbU+TJd67iR0UAkjzYvjFI/R18dKlEbmboRykw="
	txMap["tx"] = txMsg
	//decbytes, _ := hex.DecodeString(txMsg)

	//var paramsMap = make(map[string]json.RawMessage, 1)

	t1, err := cmtjson.Marshal(txMap["tx"])
	fmt.Printf("t1: %s , err: %s \n", t1, err)
	//
	//txBytes := t1
	//target := reflect.New(reflect.TypeOf(txMsg)).Interface()
	//resErr := cmtjson.Unmarshal(txBytes, target)
	//
	//actual := reflect.ValueOf(target).Elem().Interface()
	//fmt.Printf("result: %v , err: %s , actual:%s \n", txMsg, resErr.Error(), actual)

	//cmtjson.Marshal()

	resp, err := c.BroadcastTx([]byte(txMsg))
	assert.NoError(t, err)
	fmt.Printf("result: %v \n", resp.TxResponse.TxHash)
}

// success
func TestClient_Tx(t *testing.T) {
	config, _ := getWalletConfig()
	c, err := DialCosmosClient(context.Background(), config)
	assert.NoError(t, err)

	response, err := c.Tx("85C84677F466D71C0BB6E744439C3040ABB35B8F2B838CC7B73CD1BFF33D0B88", true)
	assert.NoError(t, err)
	fmt.Printf("result: %s \n", response.TxResult.String())
	fmt.Printf("result: %s \n", response.Tx.String())
}

// success
func TestClient_GetTx(t *testing.T) {
	config, _ := getWalletConfig()
	c, err := DialCosmosClient(context.Background(), config)
	assert.NoError(t, err)

	response, err := c.GetTx("85C84677F466D71C0BB6E744439C3040ABB35B8F2B838CC7B73CD1BFF33D0B88")
	assert.NoError(t, err)
	fmt.Printf("response TxHash: %s \n", response.GetTxResponse().TxHash)
	fmt.Printf("response: %s \n", response.GetTxResponse().String())
}

func TestClient_GetTxByEvent(t *testing.T) {
	config, _ := getWalletConfig()
	c, err := DialCosmosClient(context.Background(), config)
	assert.NoError(t, err)

	event := []string{"send"}
	ret, err := c.GetTxByEvent(event, 0, 10)
	assert.NoError(t, err)
	fmt.Printf("result: %s \n", ret)
}
