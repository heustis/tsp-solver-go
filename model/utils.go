package model

// IsBetween returns true if 'val' is between (inclusive) the two values 'testA' and 'testB'
func IsBetween(val float64, testA float64, testB float64) bool {
	if testA < testB {
		return val >= testA && val <= testB
	}
	return val >= testB && val <= testA
}
