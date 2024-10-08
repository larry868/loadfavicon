// Copyright @larry868 - 2023-2024

package loadfavicon

import (
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/gosimple/slug"
)

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

// Slugify returns the slugified name base on a simple filename with an extension, and a prefix.
// The returned string is ready to be used as a filename.
//
// The returned string pattern is:
//
//	[{prefix1}+][{prefix2}+][{filename}].{ext}
//
// where each part is slugified. Any query and fragment are ignored.
//
// Returns only the file name, not an absolute path.
func Slugify(prefix1 string, prefix2 string, filename string, ext string) string {

	var s string

	if prefix1 != "" {
		s = slug.Make(prefix1)
		s += "+"
	}

	if prefix2 != "" {
		s += strings.ToLower(prefix2)
		s += "+"
	}

	if filename != "" && filename != "." && filename != "/" {
		filename, _ = strings.CutSuffix(filename, ext)
		s += slug.Make(filename)
	}

	if ext != "" {
		s += ext
	}
	return s
}

// SlugHost slugify the website URL
func SlugHost(website string) string {
	hosturl, err := url.Parse(website)
	if err != nil {
		log.Println(err)
		return ""
	}
	if len(hosturl.Scheme) == 0 {
		hosturl.Scheme = "http"
	}
	hosturl.User = nil
	hosturl.RawQuery = ""
	hosturl.Fragment = ""
	hosturl.Path = strings.TrimPrefix(hosturl.Path, "/")
	return slug.Make(hosturl.Hostname() + hosturl.Path)
}

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
