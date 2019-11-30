package main

import "./Networking"

func main() {
	Networking.GetIpLocationsFromAPI()
	Networking.StartServer()
}
