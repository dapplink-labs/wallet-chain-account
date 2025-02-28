package oasis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type OasisNexusClient struct {
	endpoint         string
	oasisNexusClient *http.Client
}

func NewOasisNexusClient(endpoint string) (*OasisNexusClient, error) {
	return &OasisNexusClient{
		endpoint:         endpoint,
		oasisNexusClient: new(http.Client),
	}, nil
}

func (c *OasisNexusClient) GetBlockByNumber(index int64) (*BlockInfo, error) {
	var blockInfo BlockInfo
	url := c.endpoint + "/transactions?block=" + strconv.FormatInt(index, 10)
	resp, err := c.oasisNexusClient.Get(url)
	if err != nil {
		fmt.Println("request error:", err)
		return &blockInfo, nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error:", err)
		return &blockInfo, nil
	}

	bodyStr := string(body)
	errUnmarshal := json.Unmarshal([]byte(bodyStr), &blockInfo)
	if errUnmarshal != nil {
		fmt.Println("JSON unmarshal fail:", err)
		return &blockInfo, nil
	}
	return &blockInfo, nil
}

func (c *OasisNexusClient) GetLastBlockNumber() (*LatestBlock, error) {
	var latestBlock LatestBlock
	url := "https://nexus.oasis.io/v1/"
	resp, err := c.oasisNexusClient.Get(url)
	if err != nil {
		fmt.Println("request error:", err)
		return &latestBlock, nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error:", err)
		return &latestBlock, nil
	}

	bodyStr := string(body)
	errUnmarshal := json.Unmarshal([]byte(bodyStr), &latestBlock)
	if errUnmarshal != nil {
		fmt.Println("JSON unmarshal fail:", err)
		return &latestBlock, nil
	}
	return &latestBlock, nil
}

func (c *OasisNexusClient) GetTxByAddress(address string) (*TxInfo, error) {
	var txByAddress TxInfo
	//https://nexus.oasis.io/v1/consensus/transactions?rel=oasis1qp5fstvek2jep8xvvn67e6kn0u7etce2fsq2w4yr
	url := c.endpoint + "/transactions?rel=" + address
	resp, err := c.oasisNexusClient.Get(url)
	if err != nil {
		fmt.Println("request error:", err)
		return &txByAddress, nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error:", err)
		return &txByAddress, nil
	}

	bodyStr := string(body)
	errUnmarshal := json.Unmarshal([]byte(bodyStr), &txByAddress)
	if errUnmarshal != nil {
		fmt.Println("JSON unmarshal fail:", err)
		return &txByAddress, nil
	}
	return &txByAddress, nil
}

func (c *OasisNexusClient) GetTxByHash(hash string) (*TxInfo, error) {
	var txByHash TxInfo
	//https://nexus.oasis.io/v1/consensus/transactions/ad367bf760255afce848769a9cc601bea7fe16b32a808e24a247b82472123dc3
	url := c.endpoint + "/transactions/" + hash
	resp, err := c.oasisNexusClient.Get(url)
	if err != nil {
		fmt.Println("request error:", err)
		return &txByHash, nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error:", err)
		return &txByHash, nil
	}

	bodyStr := string(body)
	errUnmarshal := json.Unmarshal([]byte(bodyStr), &txByHash)
	if errUnmarshal != nil {
		fmt.Println("JSON unmarshal fail:", err)
		return &txByHash, nil
	}
	return &txByHash, nil
}
