// utils package contains utility functions
package utils

func ContainItemOnArray(arrays []interface{}, item interface{}) bool {
	for _, v := range arrays {
		if v == item {
			return true
		}
	}

	return false
}

func RemoveItemsOnArray(original []interface{}, removes []interface{}) []interface{} {
	var ret []interface{}
	for _, v := range original {
		if !ContainItemOnArray(removes, v) {
			ret = append(ret, v)
		}
	}
	return ret
}
