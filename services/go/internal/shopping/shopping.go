package shopping

import (
	"fmt"
	"math/big"
)

const DefaultTaxRateBPS = 425 // 4.25%
var BPSDivisor *big.Int = big.NewInt(10000)

func CartID(sessionID string) string {
	return fmt.Sprintf("cart_%s", sessionID)
}
func CalculateTaxCents(value, bps int64) int64 {
	rate := big.NewInt(bps)
	sub := big.NewInt(value)
	muled := new(big.Int).Mul(rate, sub)
	return new(big.Int).Div(muled, BPSDivisor).Int64()
}
func CalculateTotalCents(value, bps int64) int64 {
	return CalculateTaxCents(value, bps) + value
}
