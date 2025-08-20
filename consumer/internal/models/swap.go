package models

import "time"

type SwapEvent struct {
	TxHash     string    `json:"tx_hash"`
	TokenFrom  string    `json:"token_from"`
	TokenTo    string    `json:"token_to"`
	AmountFrom float64   `json:"amount_from"`
	AmountTo   float64   `json:"amount_to"`
	UsdValue   float64   `json:"usd_value"`
	Timestamp  time.Time `json:"timestamp"`
}
