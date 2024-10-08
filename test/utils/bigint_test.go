// utils_test package implements the test case on utils.
package utils_test

import (
	"math/big"
	"testing"

	"github.com/dongle/go-order-bot/internal/utils"
)

func TestToFloat64(t *testing.T) {
	one := big.NewInt(1000000000000000000)
	ten := big.NewInt(1000000000000000000)

	ten = ten.Mul(ten, big.NewInt(10))

	oneRet := utils.BigIntToFloat64(one)
	if oneRet != 1.0 {
		t.Error("One is not correct")
	}

	tenRet := utils.BigIntToFloat64(ten)
	if tenRet != 10.0 {
		t.Error("Ten is not correct")
	}

	rand := big.NewInt(1840404949583838384)
	rand = rand.Mul(rand, big.NewInt(100000000))

	randRet := utils.BigIntToFloat64(rand)
	if randRet != 184040494.9583838384 {
		t.Error("Rand is not correct")
	}
}

func TestToBigInt(t *testing.T) {
	rand := 30.0004

	randBigInt := utils.Float64ToBigInt(rand)

	expect := big.NewInt(1000000000000000000)
	expect = expect.Mul(expect, big.NewInt(30))
	expect = expect.Add(expect, big.NewInt(400000000000000))
	if randBigInt.Cmp(expect) != 0 {
		t.Error("Rand is not correct")
	}
}
