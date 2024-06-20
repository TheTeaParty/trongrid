package trongrid

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type GetContractTransactionOptions struct {
	onlyConfirmed   *bool
	onlyUnconfirmed *bool
	limit           *int
	fingerprint     *string
	orderBy         *string
	minTimestamp    *int64
	maxTimestamp    *int64
	contractAddress *string
	onlyTo          *bool
	onlyFrom        *bool
}

type GetContractTransactionOption func(*GetContractTransactionOptions)

func WithContractTransactionOnlyConfirmed(onlyConfirmed bool) GetContractTransactionOption {
	return func(o *GetContractTransactionOptions) {
		o.onlyConfirmed = &onlyConfirmed
	}
}

func WithContractTransactionOnlyUnconfirmed(onlyUnconfirmed bool) GetContractTransactionOption {
	return func(o *GetContractTransactionOptions) {
		o.onlyUnconfirmed = &onlyUnconfirmed
	}
}

func WithContractTransactionLimit(limit int) GetContractTransactionOption {
	return func(o *GetContractTransactionOptions) {
		o.limit = &limit
	}
}

func WithContractTransactionFingerprint(fingerprint string) GetContractTransactionOption {
	return func(o *GetContractTransactionOptions) {
		o.fingerprint = &fingerprint
	}
}

func WithContractTransactionOrderBy(orderBy string) GetContractTransactionOption {
	return func(o *GetContractTransactionOptions) {
		o.orderBy = &orderBy
	}
}

func WithContractTransactionMinTimestamp(minTimestamp int64) GetContractTransactionOption {
	return func(o *GetContractTransactionOptions) {
		o.minTimestamp = &minTimestamp
	}
}

func WithContractTransactionMaxTimestamp(maxTimestamp int64) GetContractTransactionOption {
	return func(o *GetContractTransactionOptions) {
		o.maxTimestamp = &maxTimestamp
	}
}

func WithContractTransactionContractAddress(contractAddress string) GetContractTransactionOption {
	return func(o *GetContractTransactionOptions) {
		o.contractAddress = &contractAddress
	}
}

func WithContractTransactionOnlyTo(onlyTo bool) GetContractTransactionOption {
	return func(o *GetContractTransactionOptions) {
		o.onlyTo = &onlyTo
	}
}

func WithContractTransactionOnlyFrom(onlyFrom bool) GetContractTransactionOption {
	return func(o *GetContractTransactionOptions) {
		o.onlyFrom = &onlyFrom
	}
}

type GetContractTransactionCursor struct {
	contractType string
	address      string
	options      *clientOptions
	currentURL   string
	nextURL      string
	err          error
	data         []*ContractTransaction
	currentIndex int
}

func (c *GetContractTransactionCursor) Next(ctx context.Context) bool {

	if c.err != nil {
		return false
	}

	if c.currentURL == "" {
		return false
	}

	if c.options.rateLimiter != nil {
		err := c.options.rateLimiter.Wait(ctx)
		if err != nil {
			c.err = err
			return false
		}
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.currentURL, nil)
	if err != nil {
		c.err = err
		return false
	}

	response, err := c.options.httpClient.Do(request)
	if err != nil {
		c.err = err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	var responseData GetContractTransactionResponse
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		c.err = err
		return false
	}

	if responseData.Success == false {
		c.err = errors.New("failed to get contract transaction")
		return false
	}

	if len(responseData.Data) == 0 {
		return false
	}

	c.data = responseData.Data
	c.currentURL = responseData.Meta.Links.Next
	c.currentIndex = 0

	return true
}

func (c *GetContractTransactionCursor) Current() (*ContractTransaction, error) {
	if c.err != nil {
		return nil, c.err
	}

	if c.currentIndex >= len(c.data) {
		return nil, nil
	}

	return c.data[c.currentIndex], nil

}

type GetContractTransactionResponse struct {
	Data    []*ContractTransaction `json:"data"`
	Meta    Meta                   `json:"meta"`
	Success bool                   `json:"success"`
}

type ContractTransaction struct {
	TransactionId  string                    `json:"transaction_id"`
	TokenInfo      *ContractTransactionToken `json:"token_info"`
	BlockTimestamp int64                     `json:"block_timestamp"`
	From           string                    `json:"from"`
	To             string                    `json:"to"`
	Type           string                    `json:"type"`
	Value          string                    `json:"value"`
}

type ContractTransactionToken struct {
	Symbol   string `json:"symbol"`
	Address  string `json:"address"`
	Decimals int    `json:"decimals"`
	Name     string `json:"name"`
}

type Meta struct {
	At          int64  `json:"at"`
	Fingerprint string `json:"fingerprint"`
	Links       struct {
		Next string `json:"next"`
	} `json:"links"`
	PageSize int `json:"page_size"`
}
