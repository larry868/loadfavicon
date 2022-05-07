package main

import (
    "fmt"
    "strings"
    "os"
    "github.com/lolorenzo777/loadfavicon/getfavicon"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("loadfavicon {website_url} {dest_dir} [--single]")
		return
	}
	website_url := os.Args[1]
	dest_dir := os.Args[2]
    single := false
    if len(os.Args) >=4 {
        single = strings.ToLower(os.Args[3]) == "--single"
    }
    nb, err := getfavicon.Download(website_url, dest_dir, single)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Printf("%d favicons downloaded\n", nb)
}

