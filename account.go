package trongrid

type Account struct {
	Address               string `json:"address"`
	Balance               int    `json:"balance"`
	CreateTime            int64  `json:"create_time"`
	LatestOprationTime    int64  `json:"latest_opration_time"`
	FreeNetUsage          int    `json:"free_net_usage"`
	LatestConsumeFreeTime int64  `json:"latest_consume_free_time"`
	NetWindowSize         int    `json:"net_window_size"`
	NetWindowOptimized    bool   `json:"net_window_optimized"`
	AccountResource       struct {
		LatestConsumeTimeForEnergy        int64 `json:"latest_consume_time_for_energy"`
		EnergyWindowSize                  int   `json:"energy_window_size"`
		DelegatedFrozenV2BalanceForEnergy int64 `json:"delegated_frozenV2_balance_for_energy"`
		EnergyWindowOptimized             bool  `json:"energy_window_optimized"`
	} `json:"account_resource"`
	OwnerPermission struct {
		PermissionName string `json:"permission_name"`
		Threshold      int    `json:"threshold"`
		Keys           []struct {
			Address string `json:"address"`
			Weight  int    `json:"weight"`
		} `json:"keys"`
	} `json:"owner_permission"`
	ActivePermission []struct {
		Type           string `json:"type"`
		Id             int    `json:"id"`
		PermissionName string `json:"permission_name"`
		Threshold      int    `json:"threshold"`
		Operations     string `json:"operations"`
		Keys           []struct {
			Address string `json:"address"`
			Weight  int    `json:"weight"`
		} `json:"keys"`
	} `json:"active_permission"`
	FrozenV2 []struct {
		Type string `json:"type,omitempty"`
	} `json:"frozenV2"`
	UnfrozenV2 []struct {
		Type               string `json:"type"`
		UnfreezeAmount     int    `json:"unfreeze_amount"`
		UnfreezeExpireTime int64  `json:"unfreeze_expire_time"`
	} `json:"unfrozenV2"`
	Trc20 []map[string]string `json:"trc20"`
}

type accountResponse struct {
	Data    []*Account `json:"data"`
	Success bool       `json:"success"`
	Meta    struct {
		PageSize int   `json:"page_size"`
		At       int64 `json:"at"`
	} `json:"meta"`
}
