package convert

import (
	"fmt"
	"strconv"
)

// StringToInt converts a string to an integer, returns an error if the conversion fails
func StringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		// Provide more context with the error message
		return 0, fmt.Errorf("convert.StringToInt: failed to convert '%s' to int: %w", s, err)
	}
	return i, nil
}

// IntToString converts an integer to a string
func IntToString(value int) string {
	return strconv.Itoa(value)
}
