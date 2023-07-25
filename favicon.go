// Copyright @lolorenzo777 - 2023

package loadfavicon

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gosimple/slug"
	xicon "github.com/mat/besticon/ico"
	"golang.org/x/image/webp"
)

// _FAVICON_MIMETYPE contains the list of accepted MIME types for content of icon files
// https://developer.mozilla.org/fr/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
var _FAVICON_MIMETYPE = []string{
	"image/x-icon",
	"image/png",
	"image/svg+xml",
	"image/jpeg",
	"image/webp",
}

type Favicon struct {
	WebsiteURL *url.URL // The absolute URL of the favicon's host website
	WebIconURL *url.URL // The URL of the favicon's file

	Image    []byte // The loaded raw image
	Color    string // Color specfications if any specified in the <link> node
	Size     string // Size specfications if any specified in the <link> node
	MimeType string // The mimetype detected when loading the image.
}

func (icon Favicon) String() string {
	return fmt.Sprintf("%s, Color:%q, Size:%q, MimeType:%v Loaded:%v", icon.WebIconURL.String(), icon.Color, icon.Size, icon.MimeType, len(icon.Image) > 0)
}

// Pixels returns the total number of pixels.
// Returns MaxInt for SVG.
// Returns -1 if unable to for get the image size or if the image is not loaded yet.
func (icon Favicon) Pixels() int64 {

	if icon.Size == "SVG" {
		const MaxInt = int64(^uint64(0) >> 1)
		return MaxInt
	}

	dim := strings.Split(strings.ToLower(icon.Size), "x")
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

// IsSVG returns true if the icon file name extension is ".svg"
func (icon Favicon) IsSVG() bool {
	return filepath.Ext(icon.WebIconURL.Path) == ".svg"
}

// AbsURL returns the url string of the webicon, making an absolute path with the WebsiteURL if required.
func (icon Favicon) AbsURL() string {
	if icon.WebIconURL.IsAbs() {
		return icon.WebIconURL.String()
	}
	j, err := url.JoinPath(icon.WebsiteURL.String(), icon.WebIconURL.String())
	if err != nil {
		return ""
	}
	return j
}

// DiskFileName returns a slugified file name based on the WebsiteURL, the Size and WebIconURL.
//
// The DiskFileName pattern is:
//
//	{WebsiteURL.Host}+{size}[+{base WebIconURL}].{ext}
//
// where each part is slugified.
//
// Returns only the file name and not an absolute path.
func (icon Favicon) DiskFileName(suffix bool) string {

	fn := slug.Make(icon.WebsiteURL.Host)
	fn += "+"
	if icon.Size != "" {
		fn += strings.ToLower(strings.Trim(icon.Size, " "))
	}

	ext := path.Ext(icon.WebIconURL.String())
	if suffix {
		fn += "+"
		base := path.Base(icon.WebIconURL.Path)
		if base != "" && base != "." && base != "/" {
			base, _ = strings.CutSuffix(base, ext)
			fn += slug.Make(base)
		}
	}

	if ext != "" {
		fn += ext
	}
	return fn
}

// ReadImage reads the image, checks its content type validity, and update the image size according to the image content.
func (icon *Favicon) ReadImage(client *http.Client) error {
	req := icon.AbsURL()
	resp, errhttp := doHttpGETRequest(client, req)
	if errhttp != nil {
		return fmt.Errorf("ReadImage %q: %+w", icon.WebIconURL.String(), errhttp)
	}
	defer resp.Body.Close()

	// ignore unreadable files, for whatever reasons
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ReadImage: %q status %s", icon.WebIconURL.String(), resp.Status)
	}

	// copy data from HTTP response to Image byte slice
	img, _ := io.ReadAll(resp.Body)
	if len(img) == 0 {
		return fmt.Errorf("ReadImage %q: unable to get the content of the icon", icon.WebIconURL.String())
	}

	// check content type
	// warning: DetectContentType detect SVG as text (see https://mimesniff.spec.whatwg.org/#identifying-a-resource-with-an-unknown-mime-type )
	ext := filepath.Ext(icon.WebIconURL.Path)
	if ext == ".svg" {
		icon.Size = "SVG"
		if !isValidSVG(img) {
			return fmt.Errorf("ReadImage %q: not a valid svg file", icon.WebIconURL.String())
		}
	}
	mimetype := http.DetectContentType(img)
	if ext != ".svg" && find(_FAVICON_MIMETYPE, mimetype) == -1 {
		return fmt.Errorf("ReadImage %q: wrong content type %q", icon.WebIconURL.String(), mimetype)
	}
	icon.MimeType = mimetype
	icon.Image = img

	// fulfill the image size for non SVG image type
	if ext != ".svg" {
		var err error
		var cfg image.Config

		reader := bytes.NewReader(img)
		switch mimetype {
		case "image/png":
			cfg, err = png.DecodeConfig(reader)
		case "image/jpeg":
			cfg, err = jpeg.DecodeConfig(reader)
		case "image/webp":
			cfg, err = webp.DecodeConfig(reader)
		case "image/x-icon":
			cfg, err = xicon.DecodeConfig(reader)
		default:
			log.Printf("ReadImage %q: unmanaged mimetype %s\n", icon.WebIconURL.String(), mimetype)
		}
		if err != nil {
			log.Printf("ReadImage %q: %s\n", icon.WebIconURL.String(), err.Error())
		} else if cfg.Width > 0 && cfg.Height > 0 {
			icon.Size = fmt.Sprintf("%vx%v", cfg.Width, cfg.Height)
		}
	}

	return nil
}

// Returns true if the given buffer is a valid SVG.
// Algo based on https://github.com/h2non/go-is-svg/blob/master/svg.go
func isValidSVG(buf []byte) bool {
	var (
		htmlCommentRegex = regexp.MustCompile(`(?i)\<\!\-\-(?:.|\n|\r)*?-->`)
		svgRegex         = regexp.MustCompile(`(?i)^\s*(?:<\?xml[^>]*>\s*)?(?:<!doctype svg[^>]*>\s*)?<svg[^>]*>[^*]*<\/svg>\s*$`)
	)
	return !isBinary(buf) && svgRegex.Match(htmlCommentRegex.ReplaceAll(buf, []byte{}))
}

// isBinary checks if the given buffer is a binary file.
func isBinary(buf []byte) bool {
	if len(buf) < 24 {
		return false
	}
	for i := 0; i < 24; i++ {
		charCode, _ := utf8.DecodeRuneInString(string(buf[i]))
		if charCode == 65533 || charCode <= 8 {
			return true
		}
	}
	return false
}
