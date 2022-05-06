package main

import (
	"log"
	"net/http"
	"os"

	//	"os"
	"testing"
)

var gwebsitesOK = []string{
	"https://go.dev/",
	"https://brave.com/",
	"https://github.com/",
	"https://twitter.com/",
	"https://www.linkedin.com/",
	"https://protonmail.com/",
	"https://getbootstrap.com/",
	"https://www.cloudflare.com/",
	"https://www.docker.com/",
}

var gwebsitesKO = []string{
	"email://www.dummy.com",
	"http:dummy.io",
	"www.dummy.abc",
	"github.com",
}

var fNeedCleaning bool

const (
	testDIR = ".test/downloadedicons"
)

func init() {
	os.RemoveAll(testDIR)
}

func TestGetFaviconLinks(t *testing.T) {
    // Create the HTTP client, re-usable, with timeout
    client := &http.Client{}
	for _, v := range(gwebsitesOK) {
		log.Printf("---getfaviconLinks: %s\n", v)
		if _, err := getFaviconLinks(client, v); err != nil {
			t.Error(err)
		}
	}

	for _, v := range(gwebsitesKO) {
		if _, err := getFaviconLinks(client, v); err == nil {
			t.Errorf("---getfaviconLinks: %s\n", v)
		}
	}
}

func TestGetfavicons(t *testing.T) {
	Website := "https://github.com/"
	favicons, err := GetFavicons(Website);
	if  err != nil {
		t.Errorf("---getFavicons %q: %v\n", Website, err)
	}
	if len(favicons) == 0 {
		t.Fail()
	}
}

func TestDownloadFaviconsDummyWebsite(t *testing.T) {
	i, err := DownloadFavicons("https://www.dummy.dummy", testDIR, false)
	if err.Error()[len(err.Error())-12:] != "no such host" {
		t.Error(err)
	}
	if i != 0 {
		t.Fail()
	}
}

func TestDownloadFaviconsNone(t *testing.T) {
	i, err := DownloadFavicons("https://lolorenzo777.github.io/website4tests-1", testDIR, false)
	if err != nil {
		t.Error(err)
	}
	if i != 0 {
		t.Fail()
		fNeedCleaning = true
	}
}

func TestDownloadFaviconsSingle(t *testing.T) {
	i, err := DownloadFavicons("https://github.com/", testDIR, true)
	if err != nil {
		t.Error(err)
	}
	if i != 1 {
		t.Fail()
	}
	fNeedCleaning = true
}


func TestDownloadFaviconsMultiple(t *testing.T) {
	i, err := DownloadFavicons("https://www.docker.com", testDIR, false)
	if err != nil {
		t.Error(err)
	}
	if i != 1 {
		t.Fail()
	}
	fNeedCleaning = true
}

func TestDownloadFaviconsBatch(t *testing.T) {
	for _, v := range(gwebsitesOK) {
		log.Printf("---Download favicon from %q\n", v)
		i, err := DownloadFavicons(v, testDIR, false)
		if err != nil {
			t.Error(err)
		}
		if i == 0 {
			t.Fail() 
		}
	}
	fNeedCleaning = true
}

func TestClear(t *testing.T) {
	if fNeedCleaning {
		os.RemoveAll(testDIR)
	}
}