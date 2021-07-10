package main


type Result struct {
	BlockData  []BlockInfo
}

type BlockInfo struct {
	BlockNumber int64
	Hash string `json:"hash"`
	Logs string  `json:"logs"`
}