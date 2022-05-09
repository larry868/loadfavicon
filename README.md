# loadfavicon

``loadfavicon`` get and download all favicons of given websites, written in go.

## Features

- look up all favicons referenced in the <head><link> of a website, plus the favicon.ico file itself at the root of the website
- returns only valid image files and urls. removes duplicates
- sluggify website names to make a valid disk filename (eg. to put them in cache for example)
- get all favicons or only a single one according to choosen options


* see the [documentation to install and use the loadfavicon command](https://pkg.go.dev/github.com/lolorenzo777/loadfavicon#section-documentation)
* see the [documentation to use the getfavicon package](https://pkg.go.dev/github.com/lolorenzo777/loadfavicon/getfavicon#section-documentation)

## Install

```
go install github.com/lolorenzo777/loadfavicon
```

## Requests or bugs?

https://github.com/lolorenzo777/loadfavicon/issues

# Licence

Copyright @lolorenzo777 - 2022 May

[MIT Licence](LICENCE)

# _References_

- https://github.com/lolorenzo777/favicon-cheat-sheet
- https://en.wikipedia.org/wiki/Favicon

