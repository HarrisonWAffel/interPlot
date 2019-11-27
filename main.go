package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"

	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Ip struct {
	IP            string  `json:"ip"`
	Type          string  `json:"type"`
	ContinentCode string  `json:"continent_code"`
	ContinentName string  `json:"continent_name"`
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	RegionCode    string  `json:"region_code"`
	RegionName    string  `json:"region_name"`
	City          string  `json:"city"`
	Zip           string  `json:"zip"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Location      struct {
		GeonameID int    `json:"geoname_id"`
		Capital   string `json:"capital"`
		Languages []struct {
			Code   string `json:"code"`
			Name   string `json:"name"`
			Native string `json:"native"`
		} `json:"languages"`
		CountryFlag             string `json:"country_flag"`
		CountryFlagEmoji        string `json:"country_flag_emoji"`
		CountryFlagEmojiUnicode string `json:"country_flag_emoji_unicode"`
		CallingCode             string `json:"calling_code"`
		IsEu                    bool   `json:"is_eu"`
	} `json:"location"`
}

const ApiKey = "959279e06d5ca3430cf26fdbca17ea7c"

func main() {
	//run and wait for zmap
	//write output to csv file
	//run zan api
	//process json into data structures
	//use static plot to plot the long lat
	//print the resulting picture
	//at somepoint we should be able to specify locations to scan
	getIpLocations()
	//scanInternet()

}

//scanInternet runs the zmap command and outputs a csv file
func scanInternet() {
	fmt.Println("Scanning")
	cm := "zmap"
	args := []string{"-B", "2M", "-p", "21", "-n", "700", "-o", "test.csv"}

	cmd := exec.Command(cm, args...)
	//We need to create a reader for the stdout of this script
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	//A scanner is created to read the stdout of the above command
	scanner := bufio.NewScanner(cmdReader)
	defer cmdReader.Close()

	//A new go thread is created to handle the output
	go func() {
		for {

			fmt.Print("ok")
			for scanner.Scan() {
				//todo; make this output zmap
				fmt.Println(scanner.Text())

			}
		}
	}()

	//We need to start our goroutine from the main thread
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error starting Cmd", err)
		os.Exit(1)
	}
	fmt.Print(scanner.Text())

}

//getIpLocations is a function that utilizes the ipstack api for ip geolocation. It utilizes a simple curl command to get a json response body containing the desired information.
func getIpLocations() {

	lines, err := ioutil.ReadFile("results.csv")
	if err != nil {
		fmt.Print("couldn't read results file")
		panic(0)
	}

	ipStrings := strings.Split(string(lines), "\n")

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	ips := make([]Ip, len(ipStrings))
	fmt.Println("Getting Location Of IPS.")
	total := len(ipStrings)

	for i, e := range ipStrings {
		fmt.Print("Located ")
		fmt.Print(i)
		fmt.Print(" out of ")
		fmt.Print(total)
		fmt.Println(" Ip addresses")

		//go func() {
		response, err := netClient.Get("http://api.ipstack.com/" + e + "?access_key=" + ApiKey + "&output=json")
		if err != nil {
			fmt.Println("bad APi connection reported")
		}
		defer response.Body.Close()
		if response.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}

			ip := Ip{}
			e := json.Unmarshal(bodyBytes, &ip)
			if e != nil {
				fmt.Println(e)

			}

			ips = append(ips, ip)
		}
		//	return
		//}()

	}

	//https://ipstack.com/quickstart

	plotPoints(ips)
}

//plotPoints uses static maps to create an output image, it takes  an array of Ip structs as its only parameter
func plotPoints(ips []Ip) {

	ctx := sm.NewContext()
	ctx.SetSize(4000, 3000)

	for i, _ := range ips {
		ctx.AddMarker(sm.NewMarker(s2.LatLngFromDegrees(ips[i].Latitude, ips[i].Longitude), color.RGBA{0xff, 0, 0, 0xff}, 16.0))
	}

	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	if err := gg.SavePNG("my-map.png", img); err != nil {
		panic(err)
	}

}
