package trongrid

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type GetAccountTransactionsOptions struct {
	OnlyConfirmed   *bool
	OnlyUnconfirmed *bool
	OnlyTo          *bool
	OnlyFrom        *bool
	Limit           *int
	Fingerprint     *string
	MinTimestamp    *int64
	MaxTimestamp    *int64
	SearchInternal  *bool
}

type GetAccountTransactionsOption func(*GetAccountTransactionsOptions)

func WithOnlyConfirmed(v bool) GetAccountTransactionsOption {
	return func(o *GetAccountTransactionsOptions) {
		o.OnlyConfirmed = &v
	}
}

func WithOnlyUnconfirmed(v bool) GetAccountTransactionsOption {
	return func(o *GetAccountTransactionsOptions) {
		o.OnlyUnconfirmed = &v
	}
}

func WithOnlyTo(v bool) GetAccountTransactionsOption {
	return func(o *GetAccountTransactionsOptions) {
		o.OnlyTo = &v
	}
}

func WithOnlyFrom(v bool) GetAccountTransactionsOption {
	return func(o *GetAccountTransactionsOptions) {
		o.OnlyFrom = &v
	}
}

func WithLimit(v int) GetAccountTransactionsOption {
	return func(o *GetAccountTransactionsOptions) {
		o.Limit = &v
	}
}

func WithFingerprint(v string) GetAccountTransactionsOption {
	return func(o *GetAccountTransactionsOptions) {
		o.Fingerprint = &v
	}
}

func WithMinTimestamp(v int64) GetAccountTransactionsOption {
	return func(o *GetAccountTransactionsOptions) {
		o.MinTimestamp = &v
	}
}

func WithMaxTimestamp(v int64) GetAccountTransactionsOption {
	return func(o *GetAccountTransactionsOptions) {
		o.MaxTimestamp = &v
	}
}

func WithSearchInternal(v bool) GetAccountTransactionsOption {
	return func(o *GetAccountTransactionsOptions) {
		o.SearchInternal = &v
	}
}

type Transaction struct {
	TxID           string `json:"txID"`
	BlockNumber    int    `json:"blockNumber"`
	BlockTimestamp int64  `json:"block_timestamp"`
	Ret            []struct {
		ContractRet string `json:"contractRet"`
		Fee         int    `json:"fee"`
	} `json:"ret"`
	Signature  []string `json:"signature"`
	RawDataHex string   `json:"raw_data_hex"`
	RawData    struct {
		Contract []struct {
			Parameter struct {
				Value struct {
					OwnerAddress    string `json:"owner_address"`
					ToAddress       string `json:"to_address,omitempty"`
					Amount          int64  `json:"amount,omitempty"`
					UnfreezeBalance int    `json:"unfreeze_balance,omitempty"`
					Resource        string `json:"resource,omitempty"`
					Balance         int64  `json:"balance,omitempty"`
					ReceiverAddress string `json:"receiver_address,omitempty"`
					Lock            bool   `json:"lock,omitempty"`
					LockPeriod      int    `json:"lock_period,omitempty"`
					FrozenBalance   int64  `json:"frozen_balance,omitempty"`
				} `json:"value"`
				TypeUrl string `json:"type_url"`
			} `json:"parameter"`
			Type string `json:"type"`
		} `json:"contract"`
		RefBlockBytes string `json:"ref_block_bytes"`
		RefBlockHash  string `json:"ref_block_hash"`
		Expiration    int64  `json:"expiration"`
		Timestamp     int64  `json:"timestamp"`
	} `json:"raw_data"`
	EnergyFee            int           `json:"energy_fee"`
	EnergyUsage          int           `json:"energy_usage"`
	EnergyUsageTotal     int           `json:"energy_usage_total"`
	NetFee               int           `json:"net_fee"`
	NetUsage             int           `json:"net_usage"`
	InternalTransactions []interface{} `json:"internal_transactions"`
}

type GetAccountTransactionsCursor struct {
	options *GetAccountTransactionsOptions
	client  *client

	currentURL   string
	err          error
	data         []*Transaction
	currentIndex int
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
		options: options,
		client:  c,

		currentURL:   u.String(),
		data:         make([]*Transaction, 0),
		currentIndex: 0,
		err:          nil,
	}

	return cursor, nil
}

func (c *GetAccountTransactionsCursor) Current() (*Transaction, error) {
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

func (c *GetAccountTransactionsCursor) Next(ctx context.Context) bool {

	if c.err != nil {
		return false
	}

	if c.currentIndex < len(c.data) {
		return true
	}

	if c.currentURL == "" {
		return false
	}

	if c.client.options.rateLimiter != nil {
		err := c.client.options.rateLimiter.Wait(ctx)
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

	response, err := c.client.options.httpClient.Do(request)
	if err != nil {
		c.err = err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	var responseData GetAccountTransactionsResponse
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		c.err = err
		return false
	}

	if responseData.Success == false {
		c.err = errors.New("failed to get account transaction")
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

type GetAccountTransactionsResponse struct {
	Data    []*Transaction `json:"data"`
	Success bool           `json:"success"`
	Meta    struct {
		PageSize    int    `json:"page_size"`
		At          int64  `json:"at"`
		Fingerprint string `json:"fingerprint"`
		Links       struct {
			Next string `json:"next"`
		} `json:"links"`
	} `json:"meta"`
}
