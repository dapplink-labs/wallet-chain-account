package ton

import (
	"fmt"

	"github.com/pkg/errors"

	gresty "github.com/go-resty/resty/v2"
)

var errBlockChainHTTPError = errors.New("ton blockchain http error")

type TonDataClient struct {
	client *gresty.Client
}

func NewTonDataClient(url string) (*TonDataClient, error) {
	if url == "" {
		return nil, fmt.Errorf("ton blockchain URL cannot be empty")
	}
	client := gresty.New()
	client.SetBaseURL(url)
	client.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			url := r.Request.URL
			return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errBlockChainHTTPError)
		}
		return nil
	})
	return &TonDataClient{
		client: client,
	}, nil
}

func (tdc *TonDataClient) GetTxByTxHash(txHash string) (*Tx, error) {
	res, err := tdc.client.R().SetQueryParams(map[string]string{
		"hash": txHash,
	}).SetResult(&Tx{}).Get("/transactions")
	if err != nil {
		return nil, errors.New("get transaction by hash fail")
	}
	spt, ok := res.Result().(*Tx)
	if !ok {
		return nil, errors.New("get transaction by hash fail")
	}
	return spt, nil
}

func (tdc *TonDataClient) GetTxByAddr(address string) (*Tx, error) {
	res, err := tdc.client.R().SetQueryParams(map[string]string{
		"account": address,
	}).SetResult(&Tx{}).Get("/transactions")
	if err != nil {
		return nil, errors.New("get transaction by address fail")
	}
	spt, ok := res.Result().(*Tx)
	if !ok {
		return nil, errors.New("get transaction by address fail")
	}
	return spt, nil
}

func (tdc *TonDataClient) PostSendTx(boc string) (string, error) {
	res, err := tdc.client.R().SetBody(map[string]string{
		"boc": boc,
	}).SetResult(&SendTxResult{}).Post("/message")
	if err != nil {
		return "0x00", errors.New("send transaction fail")
	}
	spt, ok := res.Result().(*SendTxResult)
	if !ok {
		return "0x00", errors.New("post transaction fail")
	}
	return spt.Hash, nil
}

func (tdc *TonDataClient) GetEstimateFee(address string, boc string) (*EstimateFeeResult, error) {
	res, err := tdc.client.R().SetBody(map[string]string{
		"address": address,
		"body":    boc,
	}).SetResult(&SendTxResult{}).Post("/estimateFee")
	if err != nil {
		return nil, errors.New("get transaction fee fail")
	}
	spt, ok := res.Result().(*EstimateFeeResult)
	if !ok {
		return nil, errors.New("get transaction fee fail")
	}
	return spt, nil
}
