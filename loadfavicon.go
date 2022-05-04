package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// lookuprel contains the list of <link rel=""...> that are usally related to
// define the favicon file/url
var lookuprel = []string{
    "icon",
    "shortcut icon",
    "apple-touch-icon",
    "apple-touch-icon-precomposed",
    "mask-icon",
    }

// acceptedMIMEtypes contains the list of accepted MIME types for content of icon files
// https://developer.mozilla.org/fr/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
var acceptedMIMEtypes = []string{
    "image/x-icon",
    "image/png",
    "image/svg+xml",
    "image/jpeg",
}

// find is an helper to look for a specific item in a slice
// return the index of the value found, and -1 if value not found
func find(list []string, value string) int {
    for i, v := range(list) {
        if v == value {
            return i
        }
    }
    return -1
}
    
type tFavicon struct {
    url url.URL
    color string
    size string
    buffer []byte
}

// buildHttpURL build an absolute and valid URL to look for an icon file. 
// use host if rawurl does not containt a valid one
// add http schema if missing
// return nil if not http schema or unable to parse rawurl
func buildHttpURL(rawurl string, host url.URL) *url.URL {
    url, err := url.Parse(rawurl)
    if err != nil {
        log.Println(err)
        return nil
    }
    if len(url.Host) == 0 {
        if len(host.Host) == 0 {
            return nil
        }
        url.Host = host.Host
        url.Scheme = host.Scheme
    }
    if len(url.Scheme) == 0 {
        url.Scheme = "http"
    } else if url.Scheme != "http" && url.Scheme != "https"  {
        return nil
    }
    url.User = nil
    url.RawQuery = ""
    url.Fragment = ""
    return url
}

func doHttpGETRequest(client *http.Client, url string) (*http.Response, error) {
    // create and modify http request before sending
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    // make the request
    req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
    //HACK: req.Header.Add("Accept-Encoding", "gzip, deflate, br")
    req.Header.Add("Accept-Language","fr-FR,fr,en-US,en")
    req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
    req.Header.Add("Sec-Fetch-Dest", "document")
    req.Header.Add("sec-fetch-mode", "navigate")
    req.Header.Add("sec-fetch-site", "none")
    req.Header.Add("Sec-Fetch-User", "?1")
    req.Header.Add("Upgrade-Insecure-Requests", "1")
    req.Header.Add("Referrer-Policy", "strict-origin-when-cross-origin")

    resp, err := client.Do(req)
    if err != nil {
        log.Println(err)
        return nil, err
    }

    return resp, nil
}

// getFaviconLinks return a list of urls for the favicon of this host
// extract the list of links declared in the <head> section of the site
// that may corresond to favicon, and add favicon.ico as a valid link
// if the host has responded. Returned url has not been tested
func getFaviconLinks(client *http.Client, host string) (favicons []tFavicon, err error) {
    // ensure hosturl is a host url
    hosturl := buildHttpURL(host, url.URL{})
    if hosturl == nil {
        return nil, fmt.Errorf("fail to understand host url %q", host)
    }

    // sending the request to the host
    resp, err := doHttpGETRequest(client, hosturl.Scheme + "://" + hosturl.Host)
    if err != nil {
        log.Println(err)
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        //log.Println(resp.StatusCode)
        return nil, errors.New(strconv.Itoa(resp.StatusCode))
    }

    // Create a goquery document from the HTTP response
    document, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        log.Println(err)
        return nil, err
    }

    // build the list of favicon url related to this host
    favicons = append(favicons, tFavicon{url:*buildHttpURL("favicon.ico", *hosturl)})

    document.Find("link").Each(func(i int, s *goquery.Selection){
        lrel, hasrel := s.Attr("rel")
        lhref, hasref := s.Attr("href")
        if hasrel && hasref {
            lrel = strings.ToLower(strings.Trim(lrel, " "))
            if find(lookuprel, lrel) >= 0 {
                if href := buildHttpURL(lhref, *hosturl); href != nil {
                    lsize, _ := s.Attr("sizes")
                    lcolor, _ := s.Attr("color")
                    favicons = append(favicons, tFavicon{url:*href, size:lsize, color:lcolor})    
                }
            }
        }
    })
    return favicons, nil
}

// getFavicon load all possible favicon of a host. favicon.ico if it exists
// but alors images declared in the <head> section of the site
// that may corresond to favicon, that has been successuffuly loaded, and 
// that are valid image formats
func getFavicons(host string) (loadedFavicons []tFavicon, err error) {

    // Create the HTTP client, re-usable, with timeout
    client := &http.Client{Timeout:time.Second*5}

    favicons, err := getFaviconLinks(client, host)
    if err != nil {
        return nil, err
    } 

    // scan all links and copy they content
    for _, favicon := range(favicons) {
        fnext := false
        resp, err := doHttpGETRequest(client, favicon.url.String())
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()
        if resp.StatusCode != http.StatusOK {
            //log.Printf("%v status: %v\n", favicon.url, resp.StatusCode)
            continue
        }
        // Copy data from HTTP response to file
        favicon.buffer, err = io.ReadAll(resp.Body)
        if err != nil {
            log.Printf("unable to copy the content of the file: %v\n", err)
            continue
        }
        // Avoid duplicate
        for _, existing := range(loadedFavicons) {
            if bytes.Equal(existing.buffer, favicon.buffer) {
                // log.Printf("Dulicate found %s", &favicon.url )
                fnext = true
                continue
            }
        }
        if fnext {
            continue
        }
        // check content type
        contenttype := http.DetectContentType(favicon.buffer)
        if find(acceptedMIMEtypes, contenttype) == -1 {
            continue
        }
        loadedFavicons = append(loadedFavicons, favicon)
    }

    return loadedFavicons, nil
}

// downloadFavicons loads all favicons files related to a host and store them locally to the destDir.
// files are saved with name prefixed by savePrefix. If dest file already exists, then they're replaced.
// destDir can not be an empty name. savePrefix must be longer than 3 chars
func downloadFavicons(host string, destDir string, savePrefix string) error {

    destDir = strings.ToLower(strings.Trim(destDir, " "))
    if len(destDir) == 0 {
        return fmt.Errorf("destination directory should not be empty")
    }
    savePrefix = strings.ToLower(strings.Trim(savePrefix, " "))
    if len(savePrefix) <= 3 {
        return fmt.Errorf("savePrefix should be longer than 3")
    }

    // create the dest dir
	if filepath.Base(destDir)[0] != '.' && !strings.HasPrefix(destDir, ".") {
		destDir = strings.TrimRight(destDir, "/") 
	}
    os.MkdirAll(destDir, 0755)

    // get the icons
    loadedFavicons, err := getFavicons(host)
    if err != nil {
        fmt.Println(err)
        return err
    }

    // save
    suffix := 1
    for _, favicon := range(loadedFavicons) {
        // Create output file
        strfilename := savePrefix + strconv.Itoa(suffix) + "+" + strings.TrimPrefix(favicon.url.Path, "/") //filepath.Ext
        outFile, err := os.Create(filepath.Join(destDir, strfilename))
        if err != nil {
            fmt.Println(err)
            continue 
        }

        // Copy data from HTTP response to file
        _, err = io.WriteString(outFile, string(favicon.buffer[:]))
        if err != nil {
            fmt.Println(err)
            outFile.Close()
            continue 
        }
        outFile.Close()
        suffix++
    }
    return nil
}

func main() {

}

