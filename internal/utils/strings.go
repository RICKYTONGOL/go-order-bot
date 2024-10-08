// utils package contains utility functions
package utils

import (
	"strconv"
	"strings"
)

// Contains returns the result that the string array contains the string
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func SplitPairString(strPair string) (string, string) {
	pairArray := strings.Split(strPair, "/")

	symbol1 := pairArray[0]
	symbol2 := pairArray[1]

	return symbol1, symbol2
}

func StringToFloat64(str string) *float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return nil
	}

	return &f
}

func StringToUint64(str string) *uint64 {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return nil
	}

	return &i
}
