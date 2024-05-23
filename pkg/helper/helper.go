package helper

import (
	"errors"
	"reflect"
	"strconv"
)

// Check if a number or a string number is in a spesific range.
// num can have string or number type
func NumberInRange[V string | int](num V, low int, high int) (bool, error) {
	var (
		value int
		err   error
	)

	realValue := reflect.ValueOf(num)

	// Convert to int if the value is string
	if reflect.TypeOf(num).Kind() == reflect.String {
		value, err = strconv.Atoi(realValue.String())
	} else {
		value = int(realValue.Int())
	}

	if err != nil {
		return false, err
	}

	if value < low || value > high {
		return false, errors.New("value out of bound")
	}

	return true, nil
}
