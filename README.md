# loadfavicon

Package ``load-favicon`` look for and get favicons of given websites, written in go.

## Value

- multiple usecase: install and run standalone, embbeded within another go package, download icon files
- extract any kind of favicon for a given websites

## Features

- look up all favicons referenced in the <head><link> of a website, plus the favicon.ico file itself at the root of the website
- returns only valid image files and urls. removes duplicates
- sluggify names to valid filename to enable storage on disk (eg. to put them in cache for example)
- get all favicons or only one according to choosen options

## Examples

```bash
# https://go.dev has a single favicon file, an .ico, download it into ./myfavicons dir
# ./myfavicons will be created if does not exists
$ loadfavicon https://go.dev ./myfavicons
# https://www.docker.com/ has multiple icons, download all
$ loadfavicon https://www.docker.com/ ./myfavicons
# https://github.com has multiple icons, download the best one
$ loadfavicon https://github.com ./myfavicons --single
```

You ``./myfavicons`` dir should looks like this:
```bash
$ ll 
- ./
- ../
- github-com+pinned-octocat.svg
- go-dev+favicon.ico
- www-docker-com+cropped-Docker-R-Logo-08-2018-Monochomatic-RGB_Moby-x1-180x180.png
- www-docker-com+cropped-Docker-R-Logo-08-2018-Monochomatic-RGB_Moby-x1-192x192.png
- www-docker-com+cropped-Docker-R-Logo-08-2018-Monochomatic-RGB_Moby-x1-32x32.png
- www-docker-com+favicon.ico
```


## Requests or bugs?

https://github.com/lolorenzo777/load-favicon/issues


## Install

The program can be installed an used with a command line or the module can be used within another go program.

Go 1.16 introduces a new way to [install Go programs directly](https://play-with-go.dev/installing-go-programs-directly_go116_en/) with the go command: 
```
go install github.com/lolorenzo777/load-favicon
```

and to use the package into your own source code:
```
go get github.com/lolorenzo777/load-favicon
```

# Licence

[MIT Licence](LICENCE)

# _References_

- https://github.com/lolorenzo777/favicon-cheat-sheet
- https://en.wikipedia.org/wiki/Favicon

