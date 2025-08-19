package simulator

import "time"

// SwapEvent represents a single swap event with base fields
type SwapEvent struct {
	TokenFrom  string    `json:"token_from"`
	TokenTo    string    `json:"token_to"`
	AmountFrom float64   `json:"amount_from"`
	AmountTo   float64   `json:"amount_to"`
	UsdValue   float64   `json:"usd_value"`
	Timestamp  time.Time `json:"timestamp"`
}

// Token info to simulate swaps
type TokenInfo struct {
	Name     string  `json:"name"`
	UsdPrice float64 `json:"usdt_price"`
}

var tokens = []TokenInfo{
	{"BTC", 114500},
	{"SOL", 180},
	{"TON", 3.4},
	{"ETH", 4200},
	{"USDT", 1},
}
