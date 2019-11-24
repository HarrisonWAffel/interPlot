package main

import (
	"./zan"
	"encoding/json"
	"fmt"
	"github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"image/color"
	"os"

	"github.com/golang/geo/s2"
	"os/exec"
	"strings"
)

type Ip struct {
	IP     string `json:"ip"`
	Geoip2 struct {
		City struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
		} `json:"city"`
		Country struct {
			Name string `json:"name"`
			Code string `json:"code"`
			ID   int    `json:"id"`
		} `json:"country"`
		Continent struct {
			Name string `json:"name"`
			Code string `json:"code"`
			ID   int    `json:"id"`
		} `json:"continent"`
		Postal struct {
			Code string `json:"code"`
		} `json:"postal"`
		Latlong struct {
			AccuracyRadius int     `json:"accuracy_radius"`
			Latitude       float64 `json:"latitude"`
			Longitude      float64 `json:"longitude"`
			MetroCode      int     `json:"metro_code"`
			TimeZone       string  `json:"time_zone"`
		} `json:"latlong"`
		Metadata struct {
		} `json:"metadata"`
	} `json:"geoip2"`
}

func main() {
	//run and wait for zmap
	//write output to csv file
	//run zan api
	//process json into data structures
	//use static plot to plot the long lat
	//print the resulting picture
	//at somepoint we should be able to specify locations to scan

	scanInternet()

}

//scanInternet runs the zmap command and outputs a csv file
func scanInternet() {
	fmt.Println("Scanning")
	cmd := "zmap"
	args := []string{"-B", "10M", "-p", "80", "-n", "10000", "-o", "results.csv"}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Finished Scanning. ")
	zannote()
}

//zan is a command line tool that is used to find the lat long position of each ip.
func zannote() {
	zan.GetLocs()
	plotPoints(createIps())
}

//createIps is a function that reads a zannotated json file and creates an array of Ip structs
func createIps() []Ip {

	jsonout, err := exec.Command("cat", "out.json").Output() //Just for testing, replace with your subProcess
	if err != nil {
		panic(3)
	}

	lines := strings.Split(string(jsonout), "\n")

	ips := make([]Ip, len(lines))

	for i, e := range lines {
		errs := json.Unmarshal([]byte(e), &ips[i])
		if errs != nil {
			fmt.Print(ips)
			fmt.Print(errs)

		}

	}

	fmt.Println(ips)
	return ips
}

//plotPoints uses static maps to create an output image
func plotPoints(ips []Ip) {

	ctx := sm.NewContext()
	ctx.SetSize(4000, 3000)

	for i, _ := range ips {
		ctx.AddMarker(sm.NewMarker(s2.LatLngFromDegrees(ips[i].Geoip2.Latlong.Latitude, ips[i].Geoip2.Latlong.Longitude), color.RGBA{0xff, 0, 0, 0xff}, 16.0))
	}

	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	if err := gg.SavePNG("my-map.png", img); err != nil {
		panic(err)
	}

}
