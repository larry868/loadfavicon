module github.com/lolorenzo777/loadfavicon

go 1.18

retract (
	// fix bug in SelectSingle
	v1.2.1
	// change signature of Download()
	v1.1.2
	v1.1.1
)

require github.com/PuerkitoBio/goquery v1.8.0

require github.com/gosimple/slug v1.12.0

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
)
