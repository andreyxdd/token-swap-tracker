package simulator

import (
	"math/rand"
	"time"
)

const SWAP_EVENTS_ARRIVAL_RATE = 10 * time.Second // 1 * time.Millisecond

// SwapEvent represents a single swap event with basic fields
type SwapEvent struct {
	TokenFrom  string    `json:"token_from"`
	TokenTo    string    `json:"token_to"`
	AmountFrom float64   `json:"amount_from"`
	AmountTo   float64   `json:"amount_to"`
	USDFrom    float64   `json:"usd_from"`
	USDTo      float64   `json:"usd_to"`
	Side       float64   `json:"side"`
	Timestamp  time.Time `json:"timestamp"`
}

// Tokens available for swaps
var tokens = []string{"BTC", "TON", "ETH", "USDT", "SOL"}

// Global random source and generator (called once, then used everywhere)
var randSource = rand.NewSource(time.Now().UnixNano())
var randGen = rand.New(randSource)

// simulateSwapEvents simulates the swap events with random values and sends them to the swapChannel
func SimulateSwapEvents(swapChannel chan *SwapEvent) {
	for {
		// Select random tokens for TokenFrom and TokenTo
		tokenFrom := tokens[randGen.Intn(len(tokens))]
		tokenTo := tokens[randGen.Intn(len(tokens))]
		// Avoid having the same token for both TokenFrom and TokenTo
		for tokenFrom == tokenTo {
			tokenTo = tokens[randGen.Intn(len(tokens))]
		}

		// Generate random AmountFrom and AmountTo within realistic ranges
		amountFrom := randGen.Float64() * 10                    // Random amount between 0 and 10
		amountTo := amountFrom * (0.95 + randGen.Float64()*0.1) // AmountTo is based on amountFrom but with random fluctuation

		// Generate random USDFrom and USDTo, simulate fluctuating USD value over time
		usdFrom := amountFrom * (100 + randGen.Float64()*10) // Approximate USD value for amountFrom
		usdTo := amountTo * (100 + randGen.Float64()*10)     // Approximate USD value for amountTo

		// Generate random "side" (e.g., could represent a buyer/seller side)
		side := randGen.Float64()

		// Create the new swap event
		event := &SwapEvent{
			TokenFrom:  tokenFrom,
			TokenTo:    tokenTo,
			AmountFrom: amountFrom,
			AmountTo:   amountTo,
			USDFrom:    usdFrom,
			USDTo:      usdTo,
			Side:       side,
			Timestamp:  time.Now(),
		}

		// Send event to the swap channel
		swapChannel <- event

		// Simulate event arrival rate
		time.Sleep(SWAP_EVENTS_ARRIVAL_RATE)
	}
}
