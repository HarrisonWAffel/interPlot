package main

import "./Networking"

func main() {
	//Networking.StartServer()
	Networking.GetIpLocationsFromAPI()
	Networking.StartServer()
}
