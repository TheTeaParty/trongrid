package trongrid

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
	At       int64 `json:"at"`
	PageSize int   `json:"page_size"`
}
