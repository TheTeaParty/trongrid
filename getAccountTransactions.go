package trongrid

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
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

	url          string
	currentData  []*Transaction
	currentIndex int
	err          error
}

func (c *GetAccountTransactionsCursor) Current() (*Transaction, error) {
	if c.err != nil {
		return nil, c.err
	}
	if c.currentIndex >= len(c.currentData) {
		return nil, nil
	}
	return c.currentData[c.currentIndex], nil
}

func (c *GetAccountTransactionsCursor) Next(ctx context.Context) bool {

	if c.currentIndex < len(c.currentData)-1 {
		c.currentIndex++
		return true
	}

	if c.url == "" {
		return false
	}

	if c.client.options.rateLimiter != nil {
		err := c.client.options.rateLimiter.Wait(ctx)
		if err != nil {
			c.err = err
			return false
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		c.err = err
		return false
	}

	resp, err := c.client.options.httpClient.Do(req)
	if err != nil {
		c.err = err
		return false
	}

	defer resp.Body.Close()

	var responseData GetAccountTransactionsResponse
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		c.err = err
		return false
	}

	if !responseData.Success {
		c.err = errors.New("API returned an error response")
		return false
	}

	if responseData.Data == nil {
		return false
	}

	if len(responseData.Data) == 0 {
		return false
	}

	c.currentData = responseData.Data
	c.currentIndex = 0
	c.url = responseData.Meta.Links.Next

	return len(c.currentData) > 0
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
