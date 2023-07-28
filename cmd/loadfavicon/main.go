// Copyright @lolorenzo777 - 2023

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/lolorenzo777/loadfavicon/v2"
	"github.com/lolorenzo777/verbose"
)

func main() {

	var websiteURL, toDir, size string
	var onlymissing, suffix bool
	flag.StringVar(&websiteURL, "url", "", "url of the website to download the favicons")
	flag.StringVar(&toDir, "to", "", "path to save icon files")
	flag.StringVar(&size, "size", "", "{width}x{height}|maxres|svg. download only one icon. download the icon with closest resolution to the request.")
	flag.BoolVar(&onlymissing, "onlymissing", false, "download the icons file that has not already been downloaded")
	flag.BoolVar(&suffix, "suffix", false, "suffix the written website file with the icon file name")
	flag.BoolVar(&verbose.IsOn, "verbose", false, "verbose output")
	flag.BoolVar(&verbose.IsDebugging, "debug", false, "output debugging informations")
	flag.Parse()

	if websiteURL == "" || toDir == "" {
		fmt.Println("loadfavicon -url={websiteURL} -to={toDir} [--size] [--onlymissing] [--suffix]")
		return
	}

	client := &http.Client{Timeout: time.Second * 5}
	icons, err := loadfavicon.Download(client, websiteURL, toDir, size, onlymissing, suffix)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%d favicons downloaded\n", len(icons))
}
