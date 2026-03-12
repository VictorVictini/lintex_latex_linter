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

// helper function to remove nil arguments from a slice
func RemoveNilArguments(slice []any) []any {
	if slice == nil {
		return slice
	}
	var res []any
	for _, v := range slice {
		if v != nil {
			res = append(res, v)
		}
	}
	return res
}

// helper function to convert any interfaces (any -> []interface{} -> []Component) to Component slices assuming the underlying interfaces conform to Component slices
func AnyInterfaceToComponentSlice(original any) []Component {
	origin := original.([]interface{})
	var res []Component
	for _, v := range origin {
		if v != nil {
			res = append(res, v.(Component))
		} else {
			res = append(res, nil)
		}
	}
	return res
}

// helper function to convert any interfaces (any -> []interface{} -> []T) to a slice containing values of type T, assuming the parameter provided confirms to a T slice
func AnyInterfaceToTSlice[T any](original any) []T {
	origin := original.([]interface{})
	var res []T
	for _, value := range origin {
		res = append(res, value.(T))
	}
	return res
}

// helper function that checks that every value in a list, required, is contained in the provided map
func AllContains(dict map[string]bool, required []string) bool {
	for _, expected := range required {
		if !dict[expected] {
			return false
		}
	}
	return true
}
