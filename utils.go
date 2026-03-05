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

// helper function to convert any interfaces (any -> []interface{} -> []Argument) to Argument slices assuming the underlying interfaces conform to Argument slices
func AnyInterfaceToArgumentSlice(original any) []Argument {
	origin := original.([]interface{})
	var res []Argument
	for _, v := range origin {
		res = append(res, v.(Argument))
	}
	return res
}

// helper function to convert any interfaces (any -> []interface{} -> []Line) to Line slices assuming the underlying interfaces conform to Line slices
func AnyInterfaceToLineSlice(original any) []Line {
	origin := original.([]interface{})
	var res []Line
	for _, v := range origin {
		res = append(res, v.(Line))
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

// helper function that checks that every value in a list, required, is contained in the provided map
func AllContains(dict map[string]bool, required []string) bool {
	for _, expected := range required {
		if !dict[expected] {
			return false
		}
	}
	return true
}
