package conversion

import "strconv"

// BoolToString converts a boolean to its string representation ("true" or "false")
func BoolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// StringToInt converts a string to an integer (int).
// It returns the resulting integer and an error if the conversion fails.
func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// IntToString converts an integer (int) to its string representation.
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// StringToFloat64 converts a string to a float64.
// It returns the resulting float64 and an error if the conversion fails.
func StringToFloat64(s string) (float64, error) {
	// Use 64 to specify float64 type
	return strconv.ParseFloat(s, 64)
}
