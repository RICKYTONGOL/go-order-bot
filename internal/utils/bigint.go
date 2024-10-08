// utils package contains utility functions
package utils

import (
	"math/big"
)

func BigIntToFloat64(from *big.Int) float64 {
	one := big.NewInt(1000000000000000000)
	integer := new(big.Int).Set(from)
	decimal := new(big.Int).Set(from)

	integer.Div(integer, one)
	decimal = decimal.Mod(decimal, one)

	return float64(integer.Int64()) + float64(decimal.Int64())/float64(1000000000000000000)
}

func Float64ToBigInt(val float64) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)
	// Set precision if required.
	// bigval.SetPrec(64)

	coin := new(big.Float)
	coin.SetInt(big.NewInt(1000000000000000000))

	bigval.Mul(bigval, coin)

	result := new(big.Int)
	bigval.Int(result) // store converted number in result

	return result
}
