// Copyright @larry868 - 2023-2024

package loadfavicon

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/larry868/loadfavicon/v2/pkg/svg"
	"github.com/lolorenzo777/verbose"
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
	"image/vnd.microsoft.icon",
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
// Returns -1 if unable to get the image size or if the image is not loaded yet.
func (icon Favicon) Pixels() int64 {

	if strings.ToLower(icon.Size) == "svg" {
		const MaxInt = int64(^uint64(0) >> 1)
		return MaxInt
	}
	return ToPixels(icon.Size)
}

// IsSVG returns true if the icon file name extension is ".svg"
func (icon Favicon) IsSVG() bool {
	return filepath.Ext(icon.WebIconURL.Path) == ".svg"
}

// IsFaviconIco true if the icon file name is "favicon.ico"
func (icon Favicon) IsFaviconIco() bool {
	return strings.ToLower(path.Base(icon.WebIconURL.Path)) == "favicon.ico"
}

// AbsURL returns the url string of the webicon, making an absolute path with the WebsiteURL if required.
// func (icon Favicon) AbsURL() string {
// 	if icon.WebIconURL.IsAbs() {
// 		return icon.WebIconURL.String()
// 	}
// 	abs := icon.WebsiteURL
// 	abs.Path = icon.WebIconURL.Path
// 	return abs.String()
// }

// Slugify returns a slugified file name based on the WebsiteURL, the Size and WebIconURL.
//
// The DiskFileName pattern is:
//
//	{WebsiteURL.Host}+{size}[+{base WebIconURL}].{ext}
//
// where each part is slugified. Any query and fragment are ignored.
//
// Returns only the file name, not an absolute path.
func (icon Favicon) Slugify(iconpath bool) string {
	webiconurl := icon.WebIconURL
	webiconurl.RawQuery = ""
	webiconurl.Fragment = ""
	fname := ""
	if iconpath {
		fname = path.Base(webiconurl.Path)
	}
	return Slugify(icon.WebsiteURL.Host, icon.Size, fname, path.Ext(webiconurl.String()))
}

// ReadImage reads the image, checks its content type validity, and update the image size according to the image content.
// If WebIconURL returns an error, then try witout the base only.
// If still in error try without subdomain
// func (icon *Favicon) ReadImage(client *http.Client) error {
// 	var errN error
// 	err1 := icon.readImage(client)
// 	if err1 != nil {
// 		// 2nd try: without subdomain
// 		h := strings.Split(icon.WebIconURL.Host, ".")
// 		if len(h) == 3 {
// 			bkp := *icon.WebIconURL
// 			icon.WebIconURL.Host = h[1] + "." + h[2]
// 			errN = icon.readImage(client)
// 			if errN == nil {
// 				err1 = nil
// 			} else {
// 				*icon.WebIconURL = bkp
// 			}
// 		}
// 	}

// 	if errN != nil && icon.WebIconURL.Path != "favicon.ico" {
// 		icon.WebIconURL.Path = "favicon.ico"
// 		errN = icon.readImage(client)
// 		if errN == nil {
// 			err1 = nil
// 		}
// 	}

// 	return err1
// }

func (icon *Favicon) ReadImage(client *http.Client) error {
	siconurl := icon.WebIconURL.String()
	resp, errhttp := doHttpGETRequest(client, siconurl)
	if errhttp != nil {
		return fmt.Errorf("ReadImage %q: %+w", siconurl, errhttp)
	}
	defer resp.Body.Close()

	// unreadable files for whatever reasons
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ReadImage: %q status %s", siconurl, strings.Trim(resp.Status, " "))
	}

	// copy data from HTTP response to Image byte slice
	img, _ := io.ReadAll(resp.Body)
	if len(img) == 0 {
		return fmt.Errorf("ReadImage %q: unable to get the content of the icon", siconurl)
	}

	// check content type
	// warning: DetectContentType detect SVG as text (see https://mimesniff.spec.whatwg.org/#identifying-a-resource-with-an-unknown-mime-type )
	ext := strings.ToLower(filepath.Ext(icon.WebIconURL.Path))
	if ext == ".svg" {
		icon.Size = "svg"
		if !svg.IsValidSVG(img) {
			return fmt.Errorf("ReadImage %q: not a valid svg file", siconurl)
		}
	}

	mimetype := http.DetectContentType(img)
	if ext != ".svg" && find(_FAVICON_MIMETYPE, mimetype) == -1 {
		// DEBUG:
		// os.WriteFile("./.test/log", img, 0755)
		return fmt.Errorf("ReadImage %q: wrong content type %q", siconurl, mimetype)
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
		case "image/x-icon", "image/vnd.microsoft.icon":
			cfg, err = xicon.DecodeConfig(reader)
		default:
			verbose.Printf(verbose.ALERT, "ReadImage %q: unmanaged mimetype %s\n", siconurl, mimetype)
		}
		if err != nil {
			verbose.Printf(verbose.ALERT, "ReadImage %q: %s\n", siconurl, err.Error())
		} else if cfg.Width > 0 && cfg.Height > 0 {
			icon.Size = fmt.Sprintf("%vx%v", cfg.Width, cfg.Height)
		}
	}

	return nil
}
