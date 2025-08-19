package simulator

import (
	"math/rand"
	"time"
)

type Client struct {
	swapChannel     chan *SwapEvent
	eventsPerSecond float64
	randGen         *rand.Rand
}

func New(swapChannel chan *SwapEvent, eventsPerSecond float64) *Client {
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)
	return &Client{
		swapChannel:     swapChannel,
		eventsPerSecond: eventsPerSecond,
		randGen:         randGen,
	}
}

// simulateSwapEvents simulates the swap events with random values and sends them to the swapChannel
func (c *Client) SimulateSwapEvents() {
	for {
		// Select random tokens for TokenFrom and TokenTo
		tokenFrom := tokens[c.randGen.Intn(len(tokens))]
		tokenTo := tokens[c.randGen.Intn(len(tokens))]
		// Avoid having the same token for both TokenFrom and TokenTo
		for tokenFrom.Name == tokenTo.Name {
			tokenTo = tokens[c.randGen.Intn(len(tokens))]
		}

		// Randomly generate swap event details
		amountFrom := c.randGen.Float64()*999 + 1 // random amount between 1 and 1000
		// Let's assume tokens' price fluctuates compared to USDT in the range of 0.0 to 1.0
		fluctuatingTokenFromUsdPrice := tokenFrom.UsdPrice + c.randGen.Float64()
		if tokenFrom.Name == "USDT" {
			fluctuatingTokenFromUsdPrice = 1.0
		}
		fluctuatingTokenToUsdPrice := tokenTo.UsdPrice + c.randGen.Float64()
		if tokenTo.Name == "USDT" {
			fluctuatingTokenToUsdPrice = 1.0
		}
		tokensExchangeRate := fluctuatingTokenFromUsdPrice / fluctuatingTokenToUsdPrice
		amountTo := tokensExchangeRate * amountFrom
		usdValue := fluctuatingTokenFromUsdPrice * amountFrom

		event := &SwapEvent{
			TokenFrom:  tokenFrom.Name,
			TokenTo:    tokenTo.Name,
			AmountFrom: amountFrom,
			AmountTo:   amountTo,
			UsdValue:   usdValue,
			Timestamp:  time.Now(),
		}

		// Send event to the swap channel and simulate event arrival rate
		c.swapChannel <- event
		sleepFor := time.Duration(1000/c.eventsPerSecond) * time.Millisecond
		time.Sleep(sleepFor)
	}
}
