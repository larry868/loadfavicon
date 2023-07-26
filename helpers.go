// Copyright @lolorenzo777 - 2023

package loadfavicon

import (
	"strconv"
	"strings"
)

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

// ToPixels converts size string in number of pixels.
// Expected size pattern is {width}x{height}|svg.
// Returns MaxInt for SVG.
// Returns -1 if unable to get the image size or if the image is not loaded yet.
func ToPixels(size string) int64 {
	dim := strings.Split(strings.ToLower(size), "x")
	if len(dim) != 2 {
		return -1
	}
	x, err := strconv.ParseInt(dim[0], 10, 64)
	if err != nil {
		return -1
	}
	y, err := strconv.ParseInt(dim[1], 10, 64)
	if err != nil {
		return -1
	}

	return x * y
}
