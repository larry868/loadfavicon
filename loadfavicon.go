// Copyright @larry868 - 2023-2024
/*
loadfavicon module includes a package and a command tool to downloads favicons of a given website.
*/
package loadfavicon

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/larry868/verbose"
	"golang.org/x/net/html"
)

// _FAVICON_REL contains the list of rel id that are usally related to
// define the favicon file/url into the <Head> section of the website's webpage.
// <link rel="{_FAVICON_REL}"...>
var _FAVICON_REL = []string{
	"icon",
	"shortcut icon",
	"apple-touch-icon",
	"apple-touch-icon-precomposed",
	"mask-icon",
}

// _FAVICON_EXT contains the list of allowed icon file extensions.
var _FAVICON_EXT = []string{
	".ico",
	".png",
	".svg",
	".jpg",
	".jpeg",
}

// Download downloads favicons related to a website and save them locally to the 'toDir' directory. The directory is created if it does not exist.
// The written file is suffixed with the icon file name.
// Download does not overwrite existing file on disk if onlymissing parameter is turned on.
//
// You need to provide an http.Client for example with a timeout like this :
//
//	client := &http.Client{Timeout: time.Second * 5}
//
// Returns the number of successfull downloads.
func DownloadAll(client *http.Client, websiteURL string, toDir string, onlymissing bool) (n int, err error) {
	icons, err := Download(client, websiteURL, toDir, "", onlymissing, true)
	if err != nil {
		verbose.Println(verbose.ALERT, err.Error())
	}
	return len(icons), err
}

// DownloadOne downloads the favicon of a website and save it locally to the 'toDir' directory. The directory is created if it does not exist.
// DownloadOne downloads the svg version if it exists otherwize looks for one with the highest resolution.
// The written file is not suffixed with the icon file name.
// DownloadOne does not overwrite existing file on disk if onlymissing parameter is turned on.
//
// You need to provide an http.Client for example with a timeout like this :
//
//	client := &http.Client{Timeout: time.Second * 5}
//
// Returns the filename of the downloaded icon. If none icons are found the returned filename is empty with no error.
func DownloadOne(client *http.Client, websiteURL string, toDir string, onlymissing bool) (iconfilename string, err error) {
	icons, err := Download(client, websiteURL, toDir, "maxres", onlymissing, false)
	if err != nil {
		verbose.Println(verbose.ALERT, err.Error())
		return "", err
	}
	if len(icons) == 0 {
		return "", err
	}
	return icons[0].Slugify(false), nil
}

// Download downloads favicons related to a website and save it locally to the 'toDir' directory. The directory is created if it does not exist.
// Download downloads the svg version if it exists otherwize looks for one in the given size. If no favicon exists in te given size, DownloadOne downloads the closest favicon with an upper number of pixels.
//
// Parameters:
//
//	size = {width}x{height} or maxres // If the size is empty all favicons found are downloaded. If the size is 'maxres' only the favicon with the maximum resolution is downloaded.
//	onlymissing = true // downloads if the file is missing on the disk, do not overwrite.
//	suffix = true // written file name will be suffixed with the icon file name.
//
// You need to provide an http.Client for example with a timeout like this :
//
//	client := &http.Client{Timeout: time.Second * 5}
//
// Returns the slice of downloaded Favicons. If none icons were found the returned slice is empty with no error.
func Download(client *http.Client, websiteURL string, toDir string, size string, onlymissing bool, suffix bool) (favicons []Favicon, err error) {

	toDir = strings.ToLower(strings.Trim(toDir, " "))
	if len(toDir) == 0 {
		return favicons, fmt.Errorf("Download favicons: empty destination directory")
	}

	// create the dest dir
	toDir, err = filepath.Abs(toDir)
	if err != nil {
		return favicons, fmt.Errorf("Download favicons: %+w", err)
	}
	os.MkdirAll(toDir, 0755)

	// get the icons
	icons, errX := Read(client, websiteURL, size)
	if errX != nil {
		return favicons, fmt.Errorf("Download favicons: %+w", errX)
	}

	// save on disk each favicons
	for _, icon := range icons {
		ifn := filepath.Join(toDir, icon.Slugify(suffix))
		if onlymissing {
			_, errX := os.Stat(ifn)
			if errX == nil {
				// file exists
				continue
			}
		}
		outFile, errX := os.Create(ifn)
		if errX != nil {
			return favicons, fmt.Errorf("Download favicons: %+w", errX)
		}
		_, errX = outFile.Write(icon.Image)
		outFile.Close()
		if errX != nil {
			return favicons, fmt.Errorf("Download favicons: %+w", errX)
		}
		favicons = append(favicons, icon)

		verbose.Printf(verbose.INFO, "Download favicons: %q, icon downloaded: %s\n", websiteURL, icon.WebIconURL.String())
	}

	if len(favicons) == 0 {
		verbose.Printf(verbose.INFO, "Download favicons: %q, no icon downloaded\n", websiteURL)
	}
	return favicons, nil
}

// Read reads favicons of a website and returns them in a slice.
// The returned slice is sorted from the highest icon resolution to the lowest one, starting with SVG ones if any.
//
// The size parameter must follow this pattern:
//
//	size = {width}x{height} or maxres
//
// If the size is empty all favicons found are read. If the size is 'maxres' only the favicon with the maximum resolution is downloaded.
//
// You need to provide an http.Client for example with a timeout like this :
//
//	client := &http.Client{Timeout: time.Second * 5}
//
// The returned slice contains valid images only. (see Favicon.ReadImage)
func Read(client *http.Client, websiteURL string, sz ...string) (favicons []Favicon, err error) {
	size := ""
	if len(sz) > 0 {
		size = strings.ToLower(strings.Trim(sz[0], " "))
	}

	// get Favicon Links from the website header content
	// verbose rror but continue to try with favicon.ico
	icons, errX := getFaviconLinks(client, websiteURL)
	if errX != nil {
		verbose.Printf(verbose.WARNING, "Read [%s]: %s\n", websiteURL, errX.Error())
		err = fmt.Errorf("Read [%s]: %+w", websiteURL, errX)
		//		return nil, fmt.Errorf("Read [%s]: %+w", websiteURL, errX)
	}

	// reduce the scope to the svg file if maxres or svg request, or if unable to understand the requested size
	// to reduce response time
	reqres := ToPixels(size)
	if size == "maxres" || size == "svg" || reqres <= 0 {
		for _, icon := range icons {
			if icon.IsSVG() {
				icons = make([]Favicon, 1)
				icons[0] = icon
				break
			}
		}
	}

	// scan and read all favicon images
	tryfaviconico := false
	for _, icon := range icons {
		if icon.IsFaviconIco() {
			tryfaviconico = true
		}
		errI := icon.ReadImage(client)
		if errI != nil {
			verbose.Printf(verbose.WARNING, "Read [%s]: %s\n", websiteURL, errI.Error())
			continue
		}
		favicons = append(favicons, icon)
	}

	// last try with "favicon.ico"
	if len(favicons) == 0 && !tryfaviconico {
		verbose.Printf(verbose.WARNING, "Read [%s]: reading all images fail. Try with favicon.ico\n", websiteURL)
		dotico := &Favicon{}
		dotico.WebsiteURL, _ = url.Parse(websiteURL)
		dotico.WebIconURL = new(url.URL)
		dotico.WebIconURL.Scheme = dotico.WebsiteURL.Scheme
		dotico.WebIconURL.Host = dotico.WebsiteURL.Host
		dotico.WebIconURL.Path = "favicon.ico"
		dotico.WebIconURL.RawQuery = ""
		dotico.WebIconURL.Fragment = ""

		errI := dotico.ReadImage(client)
		if errI != nil {
			verbose.Printf(verbose.WARNING, "Read [%s]: %s\n", websiteURL, errI.Error())
		} else {
			favicons = append(favicons, *dotico)
		}
	}

	// return empty favicons if none found
	if len(favicons) == 0 {
		return favicons, err
	}

	// sort favicons from high res to low res, starting by svg, finishing by unknown res
	sort.Slice(favicons, func(i, j int) bool {
		return favicons[i].Pixels() > favicons[j].Pixels()
	})

	// return all favicons is no specific size requested
	if size == "" {
		return favicons, nil
	}

	// return the favicon with maxrez if the request is svg, maxres, or if unable to understand the requested size
	oneicon := make([]Favicon, 1)
	if size == "maxres" || size == "svg" || reqres <= 0 {
		oneicon[0] = favicons[0]
		return oneicon, nil
	}

	// scan from low res to high res and returns the first one just above or equal the requested size
	for i := len(favicons) - 1; i >= 0; i-- {
		if favicons[i].Pixels() >= reqres || i == 0 {
			oneicon[0] = favicons[i]
			break
		}
	}
	return oneicon, nil
}

// GetFaviconLinks returns a list of favicons urls for websiteURL.
//
// GetFaviconLinks extracts the list of links declared in the <head> section of the site
// that may correspond to favicon. favicon.ico is added to the list if the list is empty.
//
// Color and Size favicon properties are filled with data found in the <head> items.
// The Size and the MimeType will be updated when loading the Image.
//
// You need to provide an http.Client for example with a timeout like this :
//
//	client := &http.Client{Timeout: time.Second * 5}
func GetFaviconLinks(client *http.Client, websiteURL string) (favicons []Favicon, errX error) {
	return getFaviconLinks(client, websiteURL)
}

func getFaviconLinks(client *http.Client, websiteURL string) (favicons []Favicon, errX error) {

	weburl, err := url.Parse(websiteURL)
	if weburl == nil {
		return nil, fmt.Errorf("GetFaviconLinks: %+w", err)
	}
	sweburl := weburl.String()

	// load the website page
	resp, err := doHttpGETRequest(client, sweburl)
	if err != nil {
		return favicons, fmt.Errorf("GetFaviconLinks: %+w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return favicons, fmt.Errorf("GetFaviconLinks: %q returned status %s", sweburl, strings.Trim(resp.Status, " "))
	}

	// if the page is not a text content then try with the host without the path
	// if does not succeed yet then return an error
	ctyp := resp.Header.Get("Content-Type")
	if len(ctyp) < 4 || ctyp[:4] != "text" {
		if weburl.Path != "" {
			weburl.Path = ""
			verbose.Printf(verbose.INFO, "GetFaviconLinks: %q, not an html page: %s, try %q\n", sweburl, ctyp, weburl.String())
			return getFaviconLinks(client, weburl.String())
		}
		return favicons, fmt.Errorf("GetFaviconLinks: %q, unable to find favicon", sweburl)
	}

	// redirected response ?
	finalurl := resp.Request.URL
	sfinalurl := finalurl.String()
	if sfinalurl != sweburl {
		maxchar := len(sfinalurl)
		if maxchar > 50 {
			maxchar = 50
		}
		verbose.Printf(verbose.INFO, "GetFaviconLinks: %q, redirected to: %s\n", sweburl, sfinalurl[:maxchar]+" ...")
	}

	// extract <link rel={_FAVICON_REL} href={url{_FAVICON_EXT}}>
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return favicons, fmt.Errorf("GetFaviconLinks: %+w", err)
	}

	var scan func(n *html.Node)
	scan = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "link" {
			var rel, href, sizes, color string
			for _, a := range n.Attr {
				switch a.Key {
				case "rel":
					rel = strings.ToLower(strings.Trim(a.Val, " "))
				case "href":
					href = a.Val
				case "sizes":
					sizes = a.Val
				case "color":
					color = a.Val
				}
			}

			// appends favicon for elements with a valid rel and a valid href file extension only
			if rel != "" && href != "" {
				if find(_FAVICON_REL, rel) >= 0 {
					if iconurl, err := url.Parse(href); iconurl != nil && err == nil {

						iconurlext := filepath.Ext(iconurl.Path)

						// FAVIRON_REL but with a bad file extension
						if iconurlext != "" && find(_FAVICON_EXT, iconurlext) == -1 {
							verbose.Debug("GetFaviconLinks: %q, unmanaged rel type. rel=%q, ext=%q", sweburl, rel, iconurlext)
							return
						}

						// make iconurl and absolute url with the finalurl host
						if !iconurl.IsAbs() {
							iconurl.Scheme = finalurl.Scheme
							iconurl.Host = finalurl.Host
						}

						// avoid duplicate, rare case !
						for _, already := range favicons {
							if already.WebIconURL.String() == iconurl.String() {
								return
							}
						}

						favicons = append(favicons, Favicon{
							WebsiteURL: weburl,
							WebIconURL: iconurl,
							Size:       sizes,
							Color:      color})

					} else {
						verbose.Printf(verbose.WARNING, "GetFaviconLinks: %q, %s\n", sweburl, err.Error())
					}
				}
			}
		}

		// traverses the HTML of the webpage from the first child node
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			scan(c)
		}
	}
	scan(doc)

	verbose.Printf(verbose.INFO, "GetFaviconLinks: %q, %v favicons links found\n", sweburl, len(favicons))

	// append the favicon.ico to the list as the default file to lookup
	if len(favicons) == 0 {
		dotico := Favicon{}
		dotico.WebsiteURL = weburl
		dotico.WebIconURL = new(url.URL)
		dotico.WebIconURL.Scheme = weburl.Scheme
		dotico.WebIconURL.Host = weburl.Host
		dotico.WebIconURL.Path = "favicon.ico"
		dotico.WebIconURL.RawQuery = ""
		dotico.WebIconURL.Fragment = ""

		favicons = append(favicons, dotico)
		verbose.Printf(verbose.INFO, "GetFaviconLinks: => website: %q, favicons.ico single candidate\n", sweburl)
	}

	return favicons, nil
}

// doHttpGETRequest creates, setup, and sends a http GET request.
// Returns the http response. client is not closed and can be reused.
func doHttpGETRequest(client *http.Client, getrequest string) (*http.Response, error) {

	req, err := http.NewRequest("GET", getrequest, nil)
	if err != nil {
		return nil, err
	}

	//HACK: req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("Sec-Ch-Ua", `"Not A;Brand";v="99", "Brave";v="115", "Chromium";v="115"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	// req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
