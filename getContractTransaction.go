package trongrid

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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
	err          error
	data         []*ContractTransaction
	currentIndex int
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
		err:          nil,
		currentIndex: 0,
		data:         make([]*ContractTransaction, 0),
	}

	return cursor, nil
}

func (c *GetContractTransactionCursor) Next(ctx context.Context) bool {

	if c.err != nil {
		return false
	}

	if c.currentIndex < len(c.data) {
		return true
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

	data := c.data[c.currentIndex]
	c.currentIndex++

	return data, nil

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
