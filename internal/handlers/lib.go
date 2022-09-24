package handlers

import (
	"fmt"
	"regexp"
	"strconv"
)

// add support for commas. only returns part before comma
func SanitizeAmount(i interface{}) (amount interface{}) {
	switch v := i.(type) {
	case float64:
		n := fmt.Sprintf("$%v", v)
		return n
	case string:
		r, _ := regexp.Compile("([0-9.]+)")
		n := r.FindString(v)
		amount, err := strconv.ParseFloat(n, 64)
		if err == nil {
			return amount
		}
	}
	return
}
