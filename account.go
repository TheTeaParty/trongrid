package trongrid

type Account struct {
	Address               string `json:"address"`
	Balance               int64  `json:"balance"`
	CreateTime            int64  `json:"create_time"`
	LatestOprationTime    int64  `json:"latest_opration_time"`
	LatestConsumeFreeTime int64  `json:"latest_consume_free_time"`
	NetWindowSize         int    `json:"net_window_size"`
	NetWindowOptimized    bool   `json:"net_window_optimized"`
	AccountResource       struct {
		LatestConsumeTimeForEnergy int64 `json:"latest_consume_time_for_energy"`
		EnergyWindowSize           int   `json:"energy_window_size"`
		EnergyWindowOptimized      bool  `json:"energy_window_optimized"`
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
	AssetV2 []struct {
		Key   string `json:"key"`
		Value int64  `json:"value"`
	} `json:"assetV2"`
	FreeAssetNetUsageV2 []struct {
		Key   string `json:"key"`
		Value int    `json:"value"`
	} `json:"free_asset_net_usageV2"`
	AssetOptimized bool `json:"asset_optimized"`
}
