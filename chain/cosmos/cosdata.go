package cosmos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/log"

	"github.com/dapplink-labs/wallet-chain-account/config"
)

type CosmosData struct {
	DataApiUrl   string
	DataApiKey   string
	DataApiToken string
	TimeOut      uint64
	Client       *http.Client
}

func NewCosmosData(conf *config.Config) (*CosmosData, error) {
	return &CosmosData{
		DataApiUrl:   conf.WalletNode.Cosmos.DataApiUrl,
		DataApiKey:   conf.WalletNode.Cosmos.DataApiKey,
		DataApiToken: conf.WalletNode.Cosmos.DataApiToken,
		TimeOut:      conf.WalletNode.Cosmos.TimeOut,
		Client:       &http.Client{},
	}, nil
}

func (d *CosmosData) HttpRequestHandle(urlSuffix string, params url.Values) ([]byte, error) {
	apiURL := fmt.Sprintf("%s/%s", d.DataApiUrl, urlSuffix)
	u, err := url.ParseRequestURI(apiURL)
	if err != nil {
		log.Error("HttpRequestHandle NewRequest fail", "err", err)
		return nil, err
	}
	u.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Error("HttpRequestHandle NewRequest fail", "err", err)
		return nil, err
	}

	// set header
	req.Header.Set("Ok-Access-Key", d.DataApiKey)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := d.Client.Do(req)
	if err != nil {
		log.Error("HttpRequestHandle Do fail", "err", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("HttpRequestHandle failed to read response body: %v", "err", err)
		return nil, err
	}
	return body, nil
}

func (d *CosmosData) GetThirdNativeBalance(address string) (*NativeBalanceResponse, error) {
	suffixURL := "api/v5/explorer/address/balance-cosmos"

	params := url.Values{}
	params.Add("address", address)
	params.Add("chainShortName", ChainName)

	body, err := d.HttpRequestHandle(suffixURL, params)
	if err != nil {
		log.Error("GetThirdNativeBalance HttpRequestHandle fail", "err", err)
		return nil, err
	}

	var nativeBalanceResponse NativeBalanceResponse
	err = json.Unmarshal(body, &nativeBalanceResponse)
	if err != nil {
		log.Error("GetThirdNativeBalance failed to unmarshal nativeBalanceResponse: %v", "err", err)
		return nil, err
	}

	return &nativeBalanceResponse, nil
}

func (d *CosmosData) GetThirdTxByAddress(address, page, limit string) (*NormalTransactionResponse, error) {
	suffixURL := "api/v5/explorer/address/normal-transaction-cosmos"
	params := url.Values{}
	params.Add("address", address)
	params.Add("chainShortName", ChainName)
	params.Add("page", page)
	params.Add("limit", limit)
	//params.Add("startBlockHeight", "0")
	//params.Add("endBlockHeight", "0")

	body, err := d.HttpRequestHandle(suffixURL, params)
	if err != nil {
		log.Error("GetThirdNativeBalance HttpRequestHandle fail", "err", err)
		return nil, err
	}

	var nativeNormalTransaction NormalTransactionResponse
	err = json.Unmarshal(body, &nativeNormalTransaction)
	if err != nil {
		log.Error("GetThirdTxByAddress failed to unmarshal nativeNormalTransaction: %v", "err", err)
		return nil, err
	}

	return &nativeNormalTransaction, nil
}

func (d *CosmosData) GetThirdBlockDetail(height string) (*NativeBlockResponse, error) {
	suffixURL := "api/v5/explorer/cosmos/block-fills"
	params := url.Values{}
	params.Add("height", height)
	params.Add("chainShortName", ChainName)

	body, err := d.HttpRequestHandle(suffixURL, params)
	if err != nil {
		log.Error("GetThirdNativeBalance HttpRequestHandle fail", "err", err)
		return nil, err
	}

	var nativeBlockResponse NativeBlockResponse
	err = json.Unmarshal(body, &nativeBlockResponse)
	if err != nil {
		log.Error("GetThirdTxByAddress failed to unmarshal nativeNormalTransaction: %v", "err", err)
		return nil, err
	}

	return &nativeBlockResponse, nil
}
