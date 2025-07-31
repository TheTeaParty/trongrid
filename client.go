package trongrid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var accountBalance AccountBalance
	err = json.Unmarshal(b, &accountBalance)
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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var block Block
	err = json.Unmarshal(body, &block)
	if err != nil {
		return nil, err
	}

	if block.BlockHeader == nil || block.BlockHeader.RawData == nil {
		return nil, fmt.Errorf("%w, response: %s", ErrNoDataInResponse, string(body))
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
