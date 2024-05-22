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

func (c *client) TriggerConstantContract(ctx context.Context, req *TriggerConstantContractRequest) (*TriggerConstantContractResponse, error) {

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

	endpoint := fmt.Sprintf("%s/wallet/getnowblock", c.options.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return nil, err
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
