# loadfavicon

`loadfavicon` gets and downloads favicons of a given website, written in go.

The version 2 has been fully rewritten.

## Features

- look up all favicons referenced in the <head><link> of a website, plus the favicon.ico file itself at the root of the website.
- returns only valid image files and urls. removes duplicates.
- sluggify website names to make a valid disk filename (eg. to put them in cache for example)
- get all favicons or only a single one according to choosen options.

## Usage

You can install it and run it with the `loadfavicon` command line

```
go install github.com/larry/loadfavicon
```

or you can use it in your own code

```
go get github.com/larry868/loadfavicon
```


## Requests or bugs?

https://github.com/larry868/loadfavicon/issues

# Licence

Copyright @larry868 - 2023-2024

[MIT Licence](LICENCE)

# _References_

- https://github.com/larry868/favicon-cheat-sheet
- https://en.wikipedia.org/wiki/Favicon

