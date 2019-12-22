[![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/github.com/HarrisonWAffel/interPlot)

# InterPlot

A Golang application that utilizes the networkign utility [Zmap](https://zmap.io/) to plot every reachable IPv4 address in the world. The number of points you can map depends on how long you scan the IPv4 address space, Zmap is known to take ~45 minutes to do this on a normal household connection.

But beware, the tool can eat up almost all of your bandwidth if you let it.  You can specify a limit by passing an integer to the respectively named form located on the web GUI. 

Future Goals
---

+ Region specific scans 
+ - [x] A more sustainable API solution which scales 
+ - [x] A web client to act as the server interface
+ - [x] Document the API calls more carefully and undertake best practices for a golang server.


Setup
---
This application relies on the [ZMap](https://github.com/zmap/zmap) network utility, as such you will need to install it before you can use this repository. Once you have zmap installed and placed within your path you will have to go get the Static map repository  

`go get github.com/flopp/go-staticmaps`

After getting the required dependencies you can build all go files and execute main.go. At this point a server will open on localhost:8080. From there you can easily execute scans and see their results.

Resulting Map 
---
Thanks to the awesome [Go-staticmaps](https://github.com/flopp/go-staticmaps) project developed by [flopp](https://github.com/flopp) the program outputs a great looking map of the world. Here is an example of a scan that turned up a few located IP addresses. 


![Alt text](templates/goodexample.png?raw=true "Example Map")
	
	
---
External Repositories / APIs used within this project 

[Go-staticmaps](https://github.com/flopp/go-staticmaps), Developed by [flopp](https://github.com/flopp)

[FreeGeoIP](https://github.com/fiorix/freegeoip)

