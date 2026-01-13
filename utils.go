package main

import "fmt"

// helper function to convert any interfaces (any -> []interface{} -> [][]uint8) to string a string assuming the underlying interfaces conforms to a string format
func AnyInterfaceToString(interfaceArr any) string {
	interArr := interfaceArr.([]interface{})
	var res string
	for _, v := range interArr {
		for _, val := range v.([]uint8) {
			res += fmt.Sprintf("%c", val)
		}
	}
	return res
}
