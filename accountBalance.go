package trongrid

type AccountBalance struct {
	Balance         uint64          `json:"balance"`
	BlockIdentifier BlockIdentifier `json:"block_identifier"`
}

type BlockIdentifier struct {
	Hash   string `json:"hash"`
	Number uint64 `json:"number"`
}
