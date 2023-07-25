// Copyright @lolorenzo777 - 2023

package loadfavicon

// find looks for a specific item in a slice and returns the index of the value found.
// Returns -1 if value is not found.
func find(list []string, value string) int {
	for i, v := range list {
		if v == value {
			return i
		}
	}
	return -1
}
