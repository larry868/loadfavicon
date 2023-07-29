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
	"https://lolorenzo777.github.io/website4tests-2",
	"https://laurent.lourenco.pro",
	"https://go.dev/",
	"https://brave.com/",
	"https://github.com/",
	"https://www.sncf.com/",
	"https://protonmail.com/",
	"https://getbootstrap.com/",
	"https://www.cloudflare.com/",
	"https://www.docker.com/",
}

var _gwebsitesTricky = []string{
	"https://bitcoin.org/bitcoin.pdf",
	"https://twitter.com/",
	"https://mail.proton.me/u/0/inbox/",
	"https://www.amazon.com/",
	"https://getemoji.com/",
	"https://www.linkedin.com/in/laurentlourenco",
	"https://mail.google.com/mail/u/0/#inbox",
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

	websitesKO := []string{
		"email://www.dummy.com",
		"http:dummy.io",
		"www.dummy.abc",
		"github.com",
	}

	client := &http.Client{Timeout: time.Second * 5}
	for _, v := range websitesKO {
		if _, err := GetFaviconLinks(client, v); err != nil {
			assert.Error(t, err)
		}
	}
}

func ExampleGetFaviconLinks_second() {

	client := &http.Client{Timeout: time.Second * 5}
	for _, w := range _gwebsitesOK {
		fmt.Println("GetFaviconLinks:", w)
		icons, err := GetFaviconLinks(client, w)
		if err != nil {
			fmt.Println(err)
		} else {
			printlist(icons)
		}
		fmt.Println()
	}

	// Output:
	// GetFaviconLinks: https://lolorenzo777.github.io/website4tests-2
	// https://lolorenzo777.github.io/website4tests-2  --  https://lolorenzo777.github.io/website4tests-2/test-32x32.png
	//
	// GetFaviconLinks: https://laurent.lourenco.pro
	// https://laurent.lourenco.pro  --  https://laurent.lourenco.pro/favicon-32x32.png
	// https://laurent.lourenco.pro  --  https://laurent.lourenco.pro/favicon-16x16.png
	// https://laurent.lourenco.pro  --  https://laurent.lourenco.pro/apple-touch-icon.png
	//
	// GetFaviconLinks: https://go.dev/
	// https://go.dev/  --  https://go.dev/images/favicon-gopher.png
	// https://go.dev/  --  https://go.dev/images/favicon-gopher-plain.png
	// https://go.dev/  --  https://go.dev/images/favicon-gopher.svg
	//
	// GetFaviconLinks: https://brave.com/
	// https://brave.com/  --  https://brave.com/static-assets/images/brave-favicon.png
	// https://brave.com/  --  https://brave.com/static-assets/images/cropped-brave_appicon_release-32x32.png
	// https://brave.com/  --  https://brave.com/static-assets/images/cropped-brave_appicon_release-192x192.png
	// https://brave.com/  --  https://brave.com/static-assets/images/cropped-brave_appicon_release-180x180.png
	//
	// GetFaviconLinks: https://github.com/
	// https://github.com/  --  https://github.githubassets.com/pinned-octocat.svg
	// https://github.com/  --  https://github.githubassets.com/favicons/favicon.svg
	//
	// GetFaviconLinks: https://www.sncf.com/
	// https://www.sncf.com/  --  https://www.sncf.com/themes/sncfcom/img/favicon.ico
	// https://www.sncf.com/  --  https://www.sncf.com/themes/sncfcom/img/favicon.png
	// https://www.sncf.com/  --  https://www.sncf.com/themes/sncfcom/img/apple-touch-icon.png
	// https://www.sncf.com/  --  https://www.sncf.com/themes/sncfcom/img/favicon-32x32.png
	// https://www.sncf.com/  --  https://www.sncf.com/themes/sncfcom/img/favicon-16x16.png
	//
	// GetFaviconLinks: https://protonmail.com/
	// https://protonmail.com/  --  https://proton.me/favicons/apple-touch-icon.png
	// https://protonmail.com/  --  https://proton.me/favicons/favicon-32x32.png
	// https://protonmail.com/  --  https://proton.me/favicons/favicon-16x16.png
	// https://protonmail.com/  --  https://proton.me/favicons/safari-pinned-tab.svg
	// https://protonmail.com/  --  https://proton.me/favicons/favicon.ico
	//
	// GetFaviconLinks: https://getbootstrap.com/
	// https://getbootstrap.com/  --  https://getbootstrap.com/docs/5.3/assets/img/favicons/apple-touch-icon.png
	// https://getbootstrap.com/  --  https://getbootstrap.com/docs/5.3/assets/img/favicons/favicon-32x32.png
	// https://getbootstrap.com/  --  https://getbootstrap.com/docs/5.3/assets/img/favicons/favicon-16x16.png
	// https://getbootstrap.com/  --  https://getbootstrap.com/docs/5.3/assets/img/favicons/safari-pinned-tab.svg
	// https://getbootstrap.com/  --  https://getbootstrap.com/docs/5.3/assets/img/favicons/favicon.ico
	//
	// GetFaviconLinks: https://www.cloudflare.com/
	// https://www.cloudflare.com/  --  https://www.cloudflare.com/favicon.ico
	//
	// GetFaviconLinks: https://www.docker.com/
	// https://www.docker.com/  --  https://www.docker.com/wp-content/uploads/2023/04/cropped-Docker-favicon-32x32.png
	// https://www.docker.com/  --  https://www.docker.com/wp-content/uploads/2023/04/cropped-Docker-favicon-192x192.png
	// https://www.docker.com/  --  https://www.docker.com/wp-content/uploads/2023/04/cropped-Docker-favicon-180x180.png
}

func ExampleGetFaviconLinks_third() {
	// verbose.IsOn = true
	// verbose.IsDebugging = true

	client := &http.Client{Timeout: time.Second * 5}
	for _, w := range _gwebsitesTricky {
		fmt.Println("GetFaviconLinks:", w)
		icons, err := GetFaviconLinks(client, w)
		if err != nil {
			fmt.Println(err)
		} else {
			printlist(icons)
		}
		fmt.Println()
	}

	// Output:
	// GetFaviconLinks: https://bitcoin.org/bitcoin.pdf
	// https://bitcoin.org  --  https://bitcoin.org/favicon.png?1687792074
	// https://bitcoin.org  --  https://bitcoin.org/img/icons/logo_ios.png?1687792074
	//
	// GetFaviconLinks: https://twitter.com/
	// GetFaviconLinks: Get "/": stopped after 10 redirects
	//
	// GetFaviconLinks: https://mail.proton.me/u/0/inbox/
	// https://mail.proton.me/u/0/inbox/  --  https://mail.proton.me/assets/favicon.ico
	// https://mail.proton.me/u/0/inbox/  --  https://mail.proton.me/assets/favicon.782cda472f79b5eed726.svg
	// https://mail.proton.me/u/0/inbox/  --  https://mail.proton.me/assets/apple-touch-icon-57x57.png
	// https://mail.proton.me/u/0/inbox/  --  https://mail.proton.me/assets/apple-touch-icon-60x60.png
	// https://mail.proton.me/u/0/inbox/  --  https://mail.proton.me/assets/apple-touch-icon-72x72.png
	// https://mail.proton.me/u/0/inbox/  --  https://mail.proton.me/assets/apple-touch-icon-76x76.png
	// https://mail.proton.me/u/0/inbox/  --  https://mail.proton.me/assets/apple-touch-icon-114x114.png
	// https://mail.proton.me/u/0/inbox/  --  https://mail.proton.me/assets/apple-touch-icon-120x120.png
	// https://mail.proton.me/u/0/inbox/  --  https://mail.proton.me/assets/apple-touch-icon-144x144.png
	// https://mail.proton.me/u/0/inbox/  --  https://mail.proton.me/assets/apple-touch-icon-152x152.png
	// https://mail.proton.me/u/0/inbox/  --  https://mail.proton.me/assets/apple-touch-icon-167x167.png
	// https://mail.proton.me/u/0/inbox/  --  https://mail.proton.me/assets/apple-touch-icon-180x180.png
	// https://mail.proton.me/u/0/inbox/  --  https://mail.proton.me/assets/apple-touch-icon-1024x1024.png
	//
	// GetFaviconLinks: https://www.amazon.com/
	// https://www.amazon.com/  --  https://www.amazon.com/favicon.ico
	//
	// GetFaviconLinks: https://getemoji.com/
	// https://getemoji.com/  --  https://getemoji.com/ico/favicon.png
	// https://getemoji.com/  --  https://getemoji.com/ico/apple-touch-icon.png
	//
	// GetFaviconLinks: https://www.linkedin.com/in/laurentlourenco
	// GetFaviconLinks: "https://www.linkedin.com/in/laurentlourenco" returned status 999
	//
	// GetFaviconLinks: https://mail.google.com/mail/u/0/#inbox
	// https://mail.google.com/mail/u/0/#inbox  --  https://mail.google.com/favicon.ico

}

// func TestGetFaviconLinks_2(t *testing.T) {
// 	verbose.IsOn = true
// 	verbose.IsDebugging = true

// 	client := &http.Client{Timeout: time.Second * 5}
// 	icons, err := GetFaviconLinks(client, `https://mail.google.com/mail/u/0/#inbox`)
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		printlist(icons)
// 	}

// }

func ExampleRead() {
	client := &http.Client{Timeout: time.Second * 5}
	for _, w := range _gwebsitesOK {
		fmt.Println("Read:", w)
		icons, err := Read(client, w)
		if err != nil {
			fmt.Println(err)
		} else {
			printlist(icons)
		}
		fmt.Println()
	}

	// Output:
	// Read: https://lolorenzo777.github.io/website4tests-2
	// https://lolorenzo777.github.io/website4tests-2  --  https://lolorenzo777.github.io/website4tests-2/test-32x32.png
	//
	// Read: https://laurent.lourenco.pro
	// https://laurent.lourenco.pro  --  https://laurent.lourenco.pro/apple-touch-icon.png
	// https://laurent.lourenco.pro  --  https://laurent.lourenco.pro/favicon-32x32.png
	// https://laurent.lourenco.pro  --  https://laurent.lourenco.pro/favicon-16x16.png
	//
	// Read: https://go.dev/
	// https://go.dev/  --  https://go.dev/images/favicon-gopher.svg
	//
	// Read: https://brave.com/
	// https://brave.com/  --  https://brave.com/static-assets/images/cropped-brave_appicon_release-192x192.png
	// https://brave.com/  --  https://brave.com/static-assets/images/cropped-brave_appicon_release-180x180.png
	// https://brave.com/  --  https://brave.com/static-assets/images/brave-favicon.png
	// https://brave.com/  --  https://brave.com/static-assets/images/cropped-brave_appicon_release-32x32.png
	//
	// Read: https://github.com/
	// https://github.com/  --  https://github.githubassets.com/pinned-octocat.svg
	//
	// Read: https://www.sncf.com/
	// https://www.sncf.com/  --  https://www.sncf.com/themes/sncfcom/img/favicon.png
	// https://www.sncf.com/  --  https://www.sncf.com/themes/sncfcom/img/apple-touch-icon.png
	// https://www.sncf.com/  --  https://www.sncf.com/themes/sncfcom/img/favicon.ico
	// https://www.sncf.com/  --  https://www.sncf.com/themes/sncfcom/img/favicon-32x32.png
	// https://www.sncf.com/  --  https://www.sncf.com/themes/sncfcom/img/favicon-16x16.png
	//
	// Read: https://protonmail.com/
	// https://protonmail.com/  --  https://proton.me/favicons/safari-pinned-tab.svg
	//
	// Read: https://getbootstrap.com/
	// https://getbootstrap.com/  --  https://getbootstrap.com/docs/5.3/assets/img/favicons/safari-pinned-tab.svg
	//
	// Read: https://www.cloudflare.com/
	// https://www.cloudflare.com/  --  https://www.cloudflare.com/favicon.ico
	//
	// Read: https://www.docker.com/
	// https://www.docker.com/  --  https://www.docker.com/wp-content/uploads/2023/04/cropped-Docker-favicon-192x192.png
	// https://www.docker.com/  --  https://www.docker.com/wp-content/uploads/2023/04/cropped-Docker-favicon-180x180.png
	// https://www.docker.com/  --  https://www.docker.com/wp-content/uploads/2023/04/cropped-Docker-favicon-32x32.png
}

func ExampleRead_second() {
	// verbose.IsOn = true
	// verbose.IsDebugging = true

	client := &http.Client{Timeout: time.Second * 5}
	for _, w := range _gwebsitesTricky {
		fmt.Println("Read:", w)
		icons, err := Read(client, w)
		if err != nil {
			fmt.Println(err)
		} else {
			printlist(icons)
		}
		fmt.Println()
	}

	// Output:
	// Read: https://bitcoin.org/bitcoin.pdf
	// https://bitcoin.org  --  https://bitcoin.org/img/icons/logo_ios.png?1687792074
	// https://bitcoin.org  --  https://bitcoin.org/favicon.png?1687792074
	//
	// Read: https://twitter.com/
	// https://twitter.com/  --  https://twitter.com/favicon.ico
	//
	// Read: https://mail.proton.me/u/0/inbox/
	// https://mail.proton.me/u/0/inbox/  --  https://mail.proton.me/assets/favicon.782cda472f79b5eed726.svg
	//
	// Read: https://www.amazon.com/
	// https://www.amazon.com/  --  https://www.amazon.com/favicon.ico
	//
	// Read: https://getemoji.com/
	// https://getemoji.com/  --  https://getemoji.com/favicon.ico
	//
	// Read: https://www.linkedin.com/in/laurentlourenco
	// https://www.linkedin.com/in/laurentlourenco  --  https://www.linkedin.com/favicon.ico
	//
	// Read: https://mail.google.com/mail/u/0/#inbox
	// https://mail.google.com/mail/u/0/#inbox  --  https://mail.google.com/favicon.ico
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
		assert.NoError(t, err)
		if err == nil {
			assert.True(t, one != "", v)
		}
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

func printlist(icons []Favicon) {
	for _, icon := range icons {
		fmt.Println(icon.WebsiteURL, " -- ", icon.WebIconURL)
	}
}
