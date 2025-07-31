package trongrid

import (
	"context"
	"errors"
	"net/http"
)

var (
	ErrNoDataInResponse = errors.New("no data in response")
)

const (
	mainnetBaseURL        = "https://api.trongrid.io"
	shashtaTestnetBaseURL = "https://api.shasta.trongrid.io"
	nileTestnetBaseURL    = "https://nile.trongrid.io"
)

type Network string

const (
	NetworkMainnet       Network = "mainnet"
	NetworkShastaTestnet Network = "shastatestnet"
	NetworkNileTestnet   Network = "niletestnet"
)

type Client interface {
	GetNowBlock(ctx context.Context) (*Block, error)
	GetAccountBalance(ctx context.Context, address string, blockNumber uint64, blockHash string) (*AccountBalance, error)
	GetBlockByNumber(ctx context.Context, number uint64) (*Block, error)
	GetAccount(ctx context.Context, address string) (*Account, error)
	GetAccountTransactions(ctx context.Context, address string, opts ...GetAccountTransactionsOption) (*GetAccountTransactionsCursor, error)
	BroadcastHex(ctx context.Context, req *BroadcastHexRequest) (*BroadcastHexResponse, error)
	TriggerConstantContract(ctx context.Context, req *TriggerConstantContractRequest) (*TriggerConstantContractResponse, error)
	GetContractTransaction(ctx context.Context, address, contractType string, opts ...GetContractTransactionOption) (*GetContractTransactionCursor, error)
	GetTransactionInfoByID(ctx context.Context, txID string) (*GetTransactionInfoByIDResponse, error)
}

type RateLimiter interface {
	Wait(ctx context.Context) error
}

type clientOptions struct {
	network     Network
	baseURL     string
	httpClient  *http.Client
	apiKey      string
	rateLimiter RateLimiter
}

type ClientOption func(*clientOptions)

func WithNetwork(network Network) ClientOption {
	return func(o *clientOptions) {
		o.network = network

		switch network {
		case NetworkMainnet:
			o.baseURL = mainnetBaseURL
		case NetworkShastaTestnet:
			o.baseURL = shashtaTestnetBaseURL
		case NetworkNileTestnet:
			o.baseURL = nileTestnetBaseURL
		default:
			o.baseURL = mainnetBaseURL
		}
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(o *clientOptions) {
		o.httpClient = httpClient
	}
}

func WithAPIKey(apiKey string) ClientOption {
	return func(o *clientOptions) {
		o.apiKey = apiKey
	}
}

func WithRateLimiter(rateLimiter RateLimiter) ClientOption {
	return func(o *clientOptions) {
		o.rateLimiter = rateLimiter
	}
}

func WithBaseURL(baseURL string) ClientOption {
	return func(o *clientOptions) {
		o.baseURL = baseURL
	}
}
