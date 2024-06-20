package trongrid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type client struct {
	options *clientOptions
}

func (c *client) GetBlockByNumber(ctx context.Context, number uint64) (*Block, error) {

	if c.options.rateLimiter != nil {
		err := c.options.rateLimiter.Wait(ctx)
		if err != nil {
			return nil, err
		}
	}

	endpoint := fmt.Sprintf("%s/wallet/getblockbynum", c.options.baseURL)

	reqBody := map[string]interface{}{
		"num": number,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if c.options.apiKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", c.options.apiKey)
	}

	resp, err := c.options.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var block Block
	err = json.NewDecoder(resp.Body).Decode(&block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func (c *client) GetAccountBalance(ctx context.Context, address string, blockNumber uint64, blockHash string) (*AccountBalance, error) {

	if c.options.rateLimiter != nil {
		err := c.options.rateLimiter.Wait(ctx)
		if err != nil {
			return nil, err
		}
	}

	endpoint := fmt.Sprintf("%s/wallet/getaccountbalance", c.options.baseURL)

	reqBody := map[string]interface{}{
		"account_identifier": map[string]interface{}{
			"address": address,
		},
		"block_identifier": map[string]interface{}{
			"number": blockNumber,
			"hash":   blockHash,
		},
		"visible": true,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if c.options.apiKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", c.options.apiKey)
	}

	resp, err := c.options.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var accountBalance AccountBalance
	err = json.NewDecoder(resp.Body).Decode(&accountBalance)
	if err != nil {
		return nil, err
	}

	return &accountBalance, nil
}

func (c *client) GetAccount(ctx context.Context, address string) (*Account, error) {

	if c.options.rateLimiter != nil {
		err := c.options.rateLimiter.Wait(ctx)
		if err != nil {
			return nil, err
		}
	}

	endpoint := fmt.Sprintf("%s/v1/accounts/%s", c.options.baseURL, address)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if c.options.apiKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", c.options.apiKey)
	}

	resp, err := c.options.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var accountRsp accountResponse
	err = json.NewDecoder(resp.Body).Decode(&accountRsp)
	if err != nil {
		return nil, err
	}

	if !accountRsp.Success {
		return nil, fmt.Errorf("success false in response")
	}

	if accountRsp.Data == nil || len(accountRsp.Data) == 0 {
		return nil, ErrNoDataInResponse
	}

	return accountRsp.Data[0], nil
}

func (c *client) GetAccountTransactions(ctx context.Context, address string,
	opts ...GetAccountTransactionsOption) (*GetAccountTransactionsCursor, error) {

	options := &GetAccountTransactionsOptions{}

	for _, opt := range opts {
		opt(options)
	}

	urlStr := fmt.Sprintf("%s/v1/accounts/%s/transactions", c.options.baseURL, address)

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	q := u.Query()

	if options.Fingerprint != nil {
		q.Set("fingerprint", *options.Fingerprint)
	}

	if options.MaxTimestamp != nil {
		q.Set("max_timestamp", fmt.Sprintf("%d", *options.MaxTimestamp))
	}

	if options.MinTimestamp != nil {
		q.Set("min_timestamp", fmt.Sprintf("%d", *options.MinTimestamp))
	}

	if options.Limit != nil {
		q.Set("limit", fmt.Sprintf("%d", *options.Limit))
	}

	if options.OnlyConfirmed != nil {
		q.Set("only_confirmed", fmt.Sprintf("%t", *options.OnlyConfirmed))
	}

	if options.OnlyFrom != nil {
		q.Set("only_from", fmt.Sprintf("%t", *options.OnlyFrom))
	}

	if options.OnlyTo != nil {
		q.Set("only_to", fmt.Sprintf("%t", *options.OnlyTo))
	}

	if options.OnlyUnconfirmed != nil {
		q.Set("only_unconfirmed", fmt.Sprintf("%t", *options.OnlyUnconfirmed))
	}

	if options.SearchInternal != nil {
		q.Set("search_internal", fmt.Sprintf("%t", *options.SearchInternal))
	}

	u.RawQuery = q.Encode()

	cursor := &GetAccountTransactionsCursor{
		options:      options,
		client:       c,
		url:          u.String(),
		currentData:  make([]*Transaction, 0),
		currentIndex: 0,
		err:          nil,
	}

	return cursor, nil
}

func (c *client) GetTransactionInfoByID(ctx context.Context, txID string) (*GetTransactionInfoByIDResponse, error) {

	if c.options.rateLimiter != nil {
		err := c.options.rateLimiter.Wait(ctx)
		if err != nil {
			return nil, err
		}
	}

	endpoint := fmt.Sprintf("%s/wallet/gettransactioninfobyid", c.options.baseURL)

	body, err := json.Marshal(map[string]string{"value": txID})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if c.options.apiKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", c.options.apiKey)
	}

	resp, err := c.options.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var getTransactionInfoByIDResponse GetTransactionInfoByIDResponse
	err = json.NewDecoder(resp.Body).Decode(&getTransactionInfoByIDResponse)
	if err != nil {
		return nil, err
	}

	return &getTransactionInfoByIDResponse, nil

}

func (c *client) GetContractTransaction(ctx context.Context, address, contractType string, opts ...GetContractTransactionOption) (*GetContractTransactionCursor, error) {

	if c.options.rateLimiter != nil {
		err := c.options.rateLimiter.Wait(ctx)
		if err != nil {
			return nil, err
		}
	}

	options := &GetContractTransactionOptions{}

	for _, opt := range opts {
		opt(options)
	}

	fullURLStr := fmt.Sprintf("%s/v1/accounts/%s/transactions/%s", c.options.baseURL,
		address, contractType)

	u, err := url.Parse(fullURLStr)
	if err != nil {
		return nil, err
	}

	q := u.Query()

	if options.onlyConfirmed != nil {
		q.Set("only_confirmed", fmt.Sprintf("%t", *options.onlyConfirmed))
	}

	if options.onlyUnconfirmed != nil {
		q.Set("only_unconfirmed", fmt.Sprintf("%t", *options.onlyUnconfirmed))
	}

	if options.limit != nil {
		q.Set("limit", fmt.Sprintf("%d", *options.limit))
	}

	if options.fingerprint != nil {
		q.Set("fingerprint", *options.fingerprint)
	}

	if options.orderBy != nil {
		q.Set("order_by", *options.orderBy)
	}

	if options.minTimestamp != nil {
		q.Set("min_timestamp", fmt.Sprintf("%d", *options.minTimestamp))
	}

	if options.maxTimestamp != nil {
		q.Set("max_timestamp", fmt.Sprintf("%d", *options.maxTimestamp))
	}

	if options.contractAddress != nil {
		q.Set("contract_address", *options.contractAddress)
	}

	if options.onlyTo != nil {
		q.Set("only_to", fmt.Sprintf("%t", *options.onlyTo))
	}

	if options.onlyFrom != nil {
		q.Set("only_from", fmt.Sprintf("%t", *options.onlyFrom))
	}

	u.RawQuery = q.Encode()

	cursor := &GetContractTransactionCursor{
		contractType: contractType,
		address:      address,
		options:      c.options,
		currentURL:   u.String(),
		nextURL:      "",
		err:          nil,
		currentIndex: 0,
		data:         make([]*ContractTransaction, 0),
	}

	return cursor, nil
}

func (c *client) TriggerConstantContract(ctx context.Context, req *TriggerConstantContractRequest) (*TriggerConstantContractResponse, error) {

	if c.options.rateLimiter != nil {
		err := c.options.rateLimiter.Wait(ctx)
		if err != nil {
			return nil, err
		}
	}

	endpoint := fmt.Sprintf("%s/wallet/triggerconstantcontract", c.options.baseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	if c.options.apiKey != "" {
		httpReq.Header.Set("TRON-PRO-API-KEY", c.options.apiKey)
	}

	resp, err := c.options.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var triggerConstantContractResponse TriggerConstantContractResponse
	err = json.NewDecoder(resp.Body).Decode(&triggerConstantContractResponse)
	if err != nil {
		return nil, err
	}

	return &triggerConstantContractResponse, nil
}

func (c *client) BroadcastHex(ctx context.Context, broadcastHexRequest *BroadcastHexRequest) (*BroadcastHexResponse, error) {

	if c.options.rateLimiter != nil {
		err := c.options.rateLimiter.Wait(ctx)
		if err != nil {
			return nil, err
		}
	}

	endpoint := fmt.Sprintf("%s/wallet/broadcasthex", c.options.baseURL)

	body, err := json.Marshal(broadcastHexRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if c.options.apiKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", c.options.apiKey)
	}

	resp, err := c.options.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var broadcastHexResponse BroadcastHexResponse
	err = json.NewDecoder(resp.Body).Decode(&broadcastHexResponse)
	if err != nil {
		return nil, err
	}

	return &broadcastHexResponse, nil
}

func (c *client) GetNowBlock(ctx context.Context) (*Block, error) {

	if c.options.rateLimiter != nil {
		err := c.options.rateLimiter.Wait(ctx)
		if err != nil {
			return nil, err
		}
	}

	endpoint := fmt.Sprintf("%s/wallet/getnowblock", c.options.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return nil, err
	}

	if c.options.apiKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", c.options.apiKey)
	}

	resp, err := c.options.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var block Block
	err = json.NewDecoder(resp.Body).Decode(&block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

// New returns a new TronGrid API client.
func New(opts ...ClientOption) Client {

	options := &clientOptions{
		network:    NetworkMainnet,
		baseURL:    mainnetBaseURL,
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(options)
	}

	return &client{
		options: options,
	}
}
