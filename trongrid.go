package trongrid

import (
	"context"
	"net/http"
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
	BroadcastHex(ctx context.Context, req *BroadcastHexRequest) (*BroadcastHexResponse, error)
	TriggerConstantContract(ctx context.Context, req *TriggerConstantContractRequest) (*TriggerConstantContractResponse, error)
}

type clientOptions struct {
	network    Network
	baseURL    string
	httpClient *http.Client
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
