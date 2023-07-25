// Copyright @lolorenzo777 - 2023

package loadfavicon

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var _gwebsitesOK = []string{
	"https://go.dev/",
	"https://brave.com/",
	"https://github.com/",
	"https://www.amazon.com/",
	"https://www.sncf.com/",
	"https://protonmail.com/",
	"https://getbootstrap.com/",
	"https://www.cloudflare.com/",
	"https://www.docker.com/",
}

const (
	TEST_DIR = ".test/downloadedicons"
)

func init() {
	os.RemoveAll(TEST_DIR)
}

func ExampleFavicon_DiskFileName() {

	type Test struct {
		web  string
		icon string
		sz   string
	}

	data := []Test{
		{web: "https://www.dummy.com"},
		{web: "https://www.dummy.com",
			icon: "favicon.ico",
		},
		{web: "https://www.dummy.com",
			icon: "favicon.ico",
			sz:   "16X16",
		},
		{web: "https://www.dummy.com/website/",
			icon: "/assets/favicon.png",
			sz:   "128X128",
		},
		{web: "https://www.dummy.com/website/",
			icon: "https://cdn.com/dummy/favicon.png",
			sz:   "64x64",
		},
	}

	for _, d := range data {
		icon := new(Favicon)
		icon.WebsiteURL, _ = url.Parse(d.web)
		icon.WebIconURL, _ = url.Parse(d.icon)
		icon.Size = d.sz
		fmt.Println(icon.DiskFileName(false))
	}

	// Output:
	// www-dummy-com+
	// www-dummy-com+.ico
	// www-dummy-com+16x16.ico
	// www-dummy-com+128x128.png
	// www-dummy-com+64x64.png
}

func TestGetFaviconLinks_1(t *testing.T) {

	var _gwebsitesKO = []string{
		"email://www.dummy.com",
		"http:dummy.io",
		"www.dummy.abc",
		"github.com",
	}

	client := &http.Client{Timeout: time.Second * 5}
	for _, v := range _gwebsitesKO {
		if _, err := GetFaviconLinks(client, v); err != nil {
			assert.Error(t, err)
		}
	}
}

func ExampleGetFaviconLinks() {
	printlist := func(icons []Favicon) {
		for _, icon := range icons {
			fmt.Println(icon)
		}
	}

	client := &http.Client{Timeout: time.Second * 5}
	if icons, err := GetFaviconLinks(client, "https://lolorenzo777.github.io/website4tests-2"); err == nil {
		printlist(icons)
	}
	fmt.Println()
	if icons, err := GetFaviconLinks(client, "https://laurent.lourenco.pro"); err == nil {
		printlist(icons)
	}

	// Output:
	// test-32x32.png, Color:"", Size:"32x32", MimeType: Loaded:false
	//
	// /favicon-32x32.png, Color:"", Size:"32x32", MimeType: Loaded:false
	// /favicon-16x16.png, Color:"", Size:"16x16", MimeType: Loaded:false
	// /apple-touch-icon.png, Color:"", Size:"", MimeType: Loaded:false
	// /apple-touch-icon.png, Color:"", Size:"180x180", MimeType: Loaded:false
	// /favicon.ico, Color:"", Size:"", MimeType: Loaded:false
}

func TestGetFaviconLinks_2(t *testing.T) {
	printlist := func(icons []Favicon) {
		for _, icon := range icons {
			fmt.Println(icon)
		}
	}

	client := &http.Client{Timeout: time.Second * 5}
	w := "https://www.cloudflare.com/"
	icons, err := GetFaviconLinks(client, w)
	assert.NoError(t, err, w)
	assert.True(t, err != nil || len(icons) > 0, w)

	fmt.Println(w)
	printlist(icons)
}

func TestGetFaviconLinks_3(t *testing.T) {
	client := &http.Client{Timeout: time.Second * 5}
	for _, v := range _gwebsitesOK {
		icons, err := GetFaviconLinks(client, v)
		assert.NoError(t, err, v)
		assert.True(t, err != nil || len(icons) > 0, v)
	}
}

func TestRead(t *testing.T) {
	w := "https://github.com/"
	client := &http.Client{Timeout: time.Second * 5}
	favicon, err := Read(client, w, false)
	assert.NoError(t, err)
	assert.True(t, len(favicon) > 0)
}

func TestDownloadDummy(t *testing.T) {
	client := &http.Client{Timeout: time.Second * 5}
	_, err := Download(client, "https://www.dummy.dummy", TEST_DIR, false, false, false)
	assert.ErrorContains(t, err, "no such host")
}

func TestDownloadNone(t *testing.T) {
	client := &http.Client{Timeout: time.Second * 5}
	n, err := Download(client, "https://lolorenzo777.github.io/website4tests-1", TEST_DIR, false, false, false)
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
}

func TestDownloadMulti(t *testing.T) {
	client := &http.Client{Timeout: time.Second * 5}
	n, err := Download(client, "https://laurent.lourenco.pro", TEST_DIR, false, false, true)
	assert.NoError(t, err)
	assert.True(t, n > 1)
}

func TestDownloadNoOverwrite(t *testing.T) {
	client := &http.Client{Timeout: time.Second * 5}
	n, err := Download(client, "https://laurent.lourenco.pro", TEST_DIR, false, false, false)
	assert.NoError(t, err)
	assert.True(t, n > 1)

	n, err = Download(client, "https://laurent.lourenco.pro", TEST_DIR, false, true, false)
	assert.NoError(t, err)
	assert.True(t, n == 0)
}

func TestDownloadBatch(t *testing.T) {
	client := &http.Client{Timeout: time.Second * 5}
	for _, v := range _gwebsitesOK {
		n, err := Download(client, v, TEST_DIR, true, false, false)
		assert.NoError(t, err)
		assert.True(t, n > 0)
	}
}

func ExampleDownload() {
	client := &http.Client{Timeout: time.Second * 5}
	Download(client, "https://lolorenzo777.github.io/website4tests-2", "./examples", false, false, false)
	Download(client, "https://laurent.lourenco.pro", "./examples", true, false, false)
	Download(client, "https://www.google.com", "./examples", true, false, true)

	files, _ := ioutil.ReadDir("./examples")
	for _, file := range files {
		fmt.Println(file.Name())
	}

	// Output:
	// laurent-lourenco-pro+180x180.png
	// lolorenzo777-github-io+32x32.png
	// www-google-com+32x32+favicon.ico
}
