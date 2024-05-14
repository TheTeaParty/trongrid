package trongrid

type BroadcastHexRequest struct {
	Transaction string `json:"transaction"`
}

type BroadcastHexResponse struct {
	Result      bool   `json:"result"`
	Code        string `json:"code"`
	Txid        string `json:"txid"`
	Message     string `json:"message"`
	Transaction string `json:"transaction"`
}
