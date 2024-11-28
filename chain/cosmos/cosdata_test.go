package cosmos

import (
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dapplink-labs/wallet-chain-account/config"
)

func getCosmosData() (*CosmosData, error) {
	var f = flag.String("c", defaultConfigPath, "config path")
	flag.Parse()
	conf, _ := config.New(*f)
	return NewCosmosData(conf)
}

func TestCosmos_GetThirdNativeBalance(t *testing.T) {
	cosmosData, err := getCosmosData()
	assert.NoError(t, err)

	response, err := cosmosData.GetThirdNativeBalance("cosmos19thxsunl9lzywglsndth5a278wtavawzzpv44q")
	assert.NoError(t, err)
	fmt.Println("response", response)
}

func TestCosmos_GetThirdTxByAddress(t *testing.T) {
	cosmosData, err := getCosmosData()
	assert.NoError(t, err)

	response, err := cosmosData.GetThirdTxByAddress("cosmos19thxsunl9lzywglsndth5a278wtavawzzpv44q", "1", "1")
	assert.NoError(t, err)
	fmt.Println("response", response)
}

func TestCosmos_GetThirdBlockDetail(t *testing.T) {
	cosmosData, err := getCosmosData()
	assert.NoError(t, err)

	response, err := cosmosData.GetThirdBlockDetail("17872234")
	assert.NoError(t, err)
	fmt.Println("response", response)
}
