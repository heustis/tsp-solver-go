package model

// IsBetween returns true if 'val' is between (inclusive) the two values 'testA' and 'testB'
func IsBetween(val float64, testA float64, testB float64) bool {
	if testA < testB {
		return val >= testA && val <= testB
	}
	return val >= testB && val <= testA
}

// DeleteIndexInt removes the int at the specified index in the supplied array, and returns the updated array.
// This may update the supplied array, so it should be updated with the returned array.
func DeleteIndexInt(ints []int, index int) []int {
	if lastIndex := len(ints) - 1; lastIndex < 0 {
		return ints
	} else if index <= 0 {
		return ints[1:]
	} else if index >= lastIndex {
		return ints[:lastIndex]
	} else {
		return append(ints[:index], ints[index+1:]...)
	}
}
