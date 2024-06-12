package trongrid

type GetTransactionInfoByIDResponse struct {
	Id             string                                `json:"id"`
	Fee            int                                   `json:"fee"`
	BlockNumber    int                                   `json:"blockNumber"`
	BlockTimeStamp int64                                 `json:"blockTimeStamp"`
	ContractResult []string                              `json:"contractResult"`
	Receipt        GetTransactionInfoByIDResponseReceipt `json:"receipt"`
}

type GetTransactionInfoByIDResponseReceipt struct {
	EnergyUsage        int    `json:"energy_usage"`
	EnergyFee          int    `json:"energy_fee"`
	OriginEnergyUsage  int    `json:"origin_energy_usage"`
	EnergyUsageTotal   int    `json:"energy_usage_total"`
	NetUsage           int    `json:"net_usage"`
	NetFee             int    `json:"net_fee"`
	Result             string `json:"result"`
	EnergyPenaltyTotal int    `json:"energy_penalty_total"`
}
