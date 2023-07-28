// Copyright @lolorenzo777 - 2023

package loadfavicon

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/lolorenzo777/verbose"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _gwebsitesOK = []string{
	"https://go.dev/",
	"https://brave.com/",
	"https://github.com/",
	"https://www.amazon.com/",
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

func ExampleFavicon_Slugify() {

	type Test struct {
		web  string
		icon string
		sz   string
	}

	data := []Test{
		{web: "https://www.dummy.com",
			icon: "favicon.ico",
		},
		{web: "https://www.dummy.com",
			icon: "favicon.ico",
			sz:   "16X16",
		},
		{web: "https://www.dummy.com/website/",
			icon: "https://cdn.com/dummy/favicon.png",
			sz:   "64x64",
		},
		{web: "https://www.dummy.com/website/",
			icon: "/assets/favicon.png",
			sz:   "128X128",
		},
		{web: "https://www.dummy.com/website/",
			icon: "./assets/favicon.png?1234567",
			sz:   "256x256",
		},
	}

	for _, d := range data {
		for f := 0; f <= 1; f++ {
			icon := new(Favicon)
			icon.WebsiteURL, _ = url.Parse(d.web)
			icon.WebIconURL, _ = url.Parse(d.icon)
			icon.Size = d.sz
			fmt.Println(icon.Slugify(f == 1))
		}
	}

	// Output:
	// www-dummy-com+.ico
	// www-dummy-com+favicon.ico
	// www-dummy-com+16x16+.ico
	// www-dummy-com+16x16+favicon.ico
	// www-dummy-com+64x64+.png
	// www-dummy-com+64x64+favicon.png
	// www-dummy-com+128x128+.png
	// www-dummy-com+128x128+favicon.png
	// www-dummy-com+256x256+.png
	// www-dummy-com+256x256+favicon.png
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
			fmt.Println(icon.WebsiteURL, " -- ", icon.WebIconURL)
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
	// https://lolorenzo777.github.io/website4tests-2  --  https://lolorenzo777.github.io/website4tests-2/test-32x32.png
	//
	// https://laurent.lourenco.pro  --  https://laurent.lourenco.pro/favicon-32x32.png
	// https://laurent.lourenco.pro  --  https://laurent.lourenco.pro/favicon-16x16.png
	// https://laurent.lourenco.pro  --  https://laurent.lourenco.pro/apple-touch-icon.png
	// https://laurent.lourenco.pro  --  https://laurent.lourenco.pro/favicon.ico
}

func TestGetFaviconLinks_3(t *testing.T) {
	client := &http.Client{Timeout: time.Second * 5}
	for _, v := range _gwebsitesOK {
		icons, err := GetFaviconLinks(client, v)
		require.NoError(t, err, v)
		assert.True(t, err != nil || len(icons) > 0, v)
	}
}

func TestRead_1(t *testing.T) {
	w := "https://laurent.lourenco.pro"
	client := &http.Client{Timeout: time.Second * 5}
	favicon, err := Read(client, w)
	require.NoError(t, err)
	assert.True(t, len(favicon) == 3)
}

func TestRead_2(t *testing.T) {
	verbose.IsOn = false
	// special cases
	w := "https://bitcoin.org/bitcoin.pdf"
	// w := "https://twitter.com/"
	// w := "https://mail.proton.me/u/0/inbox/"
	client := &http.Client{Timeout: time.Second * 5}
	favicon, err := Read(client, w)
	require.NoError(t, err)
	assert.True(t, len(favicon) > 0)
}

func TestDownloadDummy(t *testing.T) {
	client := &http.Client{Timeout: time.Second * 5}
	_, err := Download(client, "https://www.dummy.dummy", TEST_DIR, "", false, false)
	assert.ErrorContains(t, err, "no such host")
}

func TestDownloadNone(t *testing.T) {
	client := &http.Client{Timeout: time.Second * 5}
	icons, err := Download(client, "https://lolorenzo777.github.io/website4tests-1", TEST_DIR, "", false, false)
	require.NoError(t, err)
	assert.Equal(t, 0, len(icons))
}

func TestDownloadMulti(t *testing.T) {
	client := &http.Client{Timeout: time.Second * 5}
	icons, err := Download(client, "https://laurent.lourenco.pro", TEST_DIR, "", false, false)
	require.NoError(t, err)
	assert.True(t, len(icons) == 3)

	icons, err = Download(client, "https://laurent.lourenco.pro", TEST_DIR, "", true, false)
	require.NoError(t, err)
	assert.True(t, len(icons) == 0)
}

func TestDownloadBatch(t *testing.T) {
	client := &http.Client{Timeout: time.Second * 5}
	for _, v := range _gwebsitesOK {
		one, err := DownloadOne(client, v, TEST_DIR, false)
		require.NoError(t, err)
		assert.True(t, one != "", v)
	}
}

func ExampleDownload() {
	verbose.IsOn = false
	os.RemoveAll("./examples")

	client := &http.Client{Timeout: time.Second * 5}
	DownloadOne(client, "https://lolorenzo777.github.io/website4tests-2", "./examples", false)
	DownloadAll(client, "https://laurent.lourenco.pro", "./examples", false)
	DownloadOne(client, "https://web.archive.org/", "./examples", false)
	DownloadOne(client, "https://wikipedia.org", "./examples", false)

	files, _ := os.ReadDir("./examples")
	for _, file := range files {
		fmt.Println(file.Name())
	}

	// Output:
	// laurent-lourenco-pro+16x16+favicon-16x16.png
	// laurent-lourenco-pro+180x180+apple-touch-icon.png
	// laurent-lourenco-pro+32x32+favicon-32x32.png
	// lolorenzo777-github-io+32x32+.png
	// web-archive-org+32x32+.ico
	// wikipedia-org+160x160+.png
}
