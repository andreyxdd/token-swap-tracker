package models

type Stats struct {
	Volume  float64 `json:"volume"`
	TxCount int64   `json:"tx_count"`
}
