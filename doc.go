/*
loadfavicon command download all favicons of given websites, written in go.

	loadfavicon <website> <destination_directory> [--single]

Set the flags [--single] to select and download only one favicon if severals have been found. 

Examples

Some example to run the favicon command

	# https://go.dev has a single favicon file, an .ico, download it into ./myfavicons dir
	# ./myfavicons will be created if does not exists
	$ loadfavicon https://go.dev ./myfavicons

	# https://www.docker.com/ has multiple icons, download all
	$ loadfavicon https://www.docker.com/ ./myfavicons

	# https://github.com has multiple icons, download the best one
	$ loadfavicon https://github.com ./myfavicons --single

Then your ``./myfavicons`` dir should looks like this:

	- ./
	- ../
	- github-com+pinned-octocat.svg
	- go-dev+favicon.ico
	- www-docker-com+cropped-Docker-R-Logo-08-2018-Monochomatic-RGB_Moby-x1-180x180.png
	- www-docker-com+cropped-Docker-R-Logo-08-2018-Monochomatic-RGB_Moby-x1-192x192.png
	- www-docker-com+cropped-Docker-R-Logo-08-2018-Monochomatic-RGB_Moby-x1-32x32.png
	- www-docker-com+favicon.ico

Install

Since Go 1.16 downloading and installing a go programs is very simple:

	go install github.com/lolorenzo777/loadfavicon
	
*/
package main
