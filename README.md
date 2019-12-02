# InterPlot

A Golang application that utilizes the networkign utility [Zmap](https://zmap.io/) to plot every reachable IPv4 address in the world. The number of points you can map depends on how long you scan the IPv4 address space, Zmap is known to take ~45 minutes to do this on a normal household connection.

But beware, the tool can eat up almost all of your bandwidth if you let it.  The current implementation limits the scan speed to 20 Mbps, which for some households can be aggresive. Feel free to modify the value to best fit your network setup. 

Future Goals
---

+ Region specific scans 
+ Server functionality to allow for continuous mapping 
+ A more sustainable API solution which scales more readily
+ A web client to act as the server interface
	


Setup
---
This application relies on the [ZMap](https://github.com/zmap/zmap) network utility, as such you will need to install it before you can use this repository. Once you have zmap installed and placed within your path you will have to go get the Static map repository  

`go get github.com/flopp/go-staticmaps`

You will also need to create your own [IP Stack API Key](https://ipstack.com/plan). It's free and doesn't require a credit card to use. Once you have your key, open `main.go` and paste it into the const variable `ApiKey`  within `main.go`. 

Resulting Map 
---
Thanks to the awesome [Go-staticmaps](https://github.com/flopp/go-staticmaps) project developed by [flopp](https://github.com/flopp) the program outputs a great looking map of the world. Here is an example of a scan that turned up a few located IP addresses. 


![Alt text](goodexample.png?raw=true "Example Map")


#### A Word Of Caution
As mentioned before, the networking tool used to undertake these IP scans can destroy bandwidth speeds and cripple a network. This project has the scan speed limited to 20Mbps, however this may be aggresive for some areas. To change this value modify the second argument of the zmap exec process call, to whatever value you desire. 

For example, 

A scan limited to 2 Mbps would look like such

    args := []string{"-B", "2M", "-p", "21", "-n", "700", "-o", "test.csv"}

If we wanted to increase that limit to 20 Mbps it would be changed to 

	args := []string{"-B", "20M", "-p", "21", "-n", "700", "-o", "test.csv"}
	
<br/>

It should also be noted that the free tier on the IP Stack API will not support a full scan of the IPv4 address space. 	
	
	
---
External Repositories / APIs used within this project 

[Go-staticmaps](https://github.com/flopp/go-staticmaps), Developed by [flopp](https://github.com/flopp)

[IP Stack ](https://ipstack.com/) IP Geolocation API
