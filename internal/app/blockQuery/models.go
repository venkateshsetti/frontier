package blockQuery

// BlockInfo block data model
type BlockInfo  struct {
	BlockNumber string   `json:"block_number"`
	Hash 		string   `json:"hash"`
	Logs 		string   `json:"logs"`
}

// Response http response
type Response struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Result  []BlockInfo `json:"result"`
}