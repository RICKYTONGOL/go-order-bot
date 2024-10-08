// utils package contains utility functions
package utils

import (
	"errors"
	"math"
)

func GetFloat(unk interface{}) (float64, error) {
	switch i := unk.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	// ...other cases...
	default:
		return math.NaN(), errors.New("getFloat: unknown value is of incompatible type")
	}
}
