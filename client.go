package trongrid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type client struct {
	options *clientOptions
}

func (c *client) GetNowBlock(ctx context.Context) (*Block, error) {

	endpoint := fmt.Sprintf("%s/wallet/getnowblock", c.options.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
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
