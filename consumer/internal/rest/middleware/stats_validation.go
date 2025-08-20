package middleware

import (
	"regexp"
	"strings"
)

var validPeriods = map[string]bool{
	"5min": true,
	"1h":   true,
	"24h":  true,
}

var validTokens = map[string]bool{
	"USDT": true,
	"BTC":  true,
	"TON":  true,
	"SOL":  true,
	"ETH":  true,
}

// Function to validate the period
func IsValidPeriod(period string) bool {
	_, valid := validPeriods[period]
	return valid
}

// Function to validate the token
func IsValidToken(token string) bool {
	_, valid := validTokens[strings.ToUpper(token)]
	return valid
}

// Function to validate the pair param (format "TOKEN-TOKEN")
func IsValidPair(pair string) bool {
	re := regexp.MustCompile(`^[A-Za-z]{3,4}-[A-Za-z]{3,4}$`)
	return re.MatchString(pair) &&
		IsValidToken(strings.Split(pair, "-")[0]) &&
		IsValidToken(strings.Split(pair, "-")[1])
}
