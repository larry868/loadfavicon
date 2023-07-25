// Copyright @lolorenzo777 - 2023

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/lolorenzo777/loadfavicon/v2"
)

func main() {

	var websiteURL, toDir string
	var maxres, missing, suffix bool
	flag.StringVar(&websiteURL, "url", "", "url of the website to download the favicons")
	flag.StringVar(&toDir, "to", "", "path to save icon files")
	flag.BoolVar(&maxres, "maxres", false, "download only the favicon with the maximum resolution")
	flag.BoolVar(&missing, "missing", false, "download only if the icons file has not already been downloaded")
	flag.BoolVar(&suffix, "suffix", false, "suffix the written website file with the icon file name")
	flag.Parse()

	if websiteURL == "" || toDir == "" {
		fmt.Println("loadfavicon -url={websiteURL} -to={toDir} [--maxres] [--missing] [--prefix]")
		return
	}

	client := &http.Client{Timeout: time.Second * 5}
	n, err := loadfavicon.Download(client, websiteURL, toDir, maxres, missing, suffix)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%d favicons downloaded\n", n)
}
