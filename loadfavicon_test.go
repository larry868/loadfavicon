package main

import (
	"log"
	"net/http"
	"os"
	"testing"
)

var ghostsOK = []string{
	"https://twitter.com/",
	"https://www.linkedin.com/",
	"https://protonmail.com/",
	"https://laurent.lourenco.pro",
	"https://go.dev/",
	"https://getbootstrap.com/",
	"https://www.cloudflare.com/",
	"https://brave.com/",
}

var ghostsKO = []string{
	"email://www.dummy.com",
	"http:dummy.io",
	"www.dummy.abc",
	"github.com",
}

func TestGetFaviconLinks(t *testing.T) {
    // Create the HTTP client, re-usable, with timeout
    client := &http.Client{}
	for _, v := range(ghostsOK) {
		log.Printf("---getfaviconLinks: %s\n", v)
		if _, err := getFaviconLinks(client, v); err != nil {
			t.Error(err)
		}
	}

	for _, v := range(ghostsKO) {
		if _, err := getFaviconLinks(client, v); err == nil {
			t.Errorf("---getfaviconLinks: %s\n", v)
		}
	
	}
}

func TestGetfavicons(t *testing.T) {
	if _, err := getFavicons("https://laurent.lourenco.pro/"); err == nil {
		t.Errorf("---getFavicons https://laurent.lourenco.pro/: %v\n", err)
	}
}

func TestDownloadfavicons1(t *testing.T) {
	if err := downloadFavicons("https://laurent.lourenco.pro", ".test/downloadedicons/laurentlourenco", "favicon"); err != nil {
		t.Error(err)
	}
	if err := downloadFavicons("https://blog.lourenco.pro", ".test/downloadedicons/bloglourenco", "favicon"); err != nil {
		t.Error(err)
	}
	if !t.Failed() {
		os.RemoveAll(".downloadedicons")
	}
}

