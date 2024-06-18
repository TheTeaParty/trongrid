package trongrid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type client struct {
	options *clientOptions
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

func (c *client) GetContractTransaction(ctx context.Context, address, contractType string) (*GetContractTransactionResponse, error) {

	if c.options.rateLimiter != nil {
		err := c.options.rateLimiter.Wait(ctx)
		if err != nil {
			return nil, err
		}
	}

	fullURL := fmt.Sprintf("%s/v1/accounts/%s/transactions/%s?only_confirmed=true", c.options.baseURL,
		address, contractType)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.options.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	var responseData GetContractTransactionResponse
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, nil
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
