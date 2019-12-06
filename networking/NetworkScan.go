package networking

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fiorix/freegeoip"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

var scanner *bufio.Scanner
var cmd *exec.Cmd
var stdin io.WriteCloser
var scanning bool

//ScanInternet runs the zmap command and outputs a csv file
func ScanInternet(ctx context.Context, speedLimit string, n string) {

	if scanning == false {
		scanning = true
		log.Println("Scanning")
		cm := "zmap"

		args := []string{"-B", strings.TrimSpace(speedLimit) + "M", "-p", "80", "-n", strings.TrimSpace(n), "-o", "results.csv"}
		fmt.Println(args)
		cmd = exec.Command(cm, args...)

		//We need to create a reader for the stderr of this script
		cmdReader, err := cmd.StderrPipe()
		if err != nil {
			fmt.Println("Error creating StderrPipe for Cmd", err)
			os.Exit(1)
		}

		defer cmdReader.Close()

		stdn, e := cmd.StdinPipe()
		if e != nil {
			log.Fatal("Cannot get stdin pipe ")
		}
		stdin = stdn

		//A scanner is created to read the stderr of the above command
		scanner = bufio.NewScanner(cmdReader)

		go func() {
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
			return
		}()

		err = cmd.Run()
		if err != nil {
			log.Fatal("Error starting Cmd ", err)
		}
		scanner = nil
		GetIpLocationsFromAPI()

	}
}

//ListenToScan returns the last line of the current scan.
func ListenToScan() string {
	if scanner != nil {
		lastLine := strings.Split(scanner.Text(), "\n")[len(strings.Split(scanner.Text(), "\n"))-1]
		if lastLine != "" {
			return lastLine
		}
	}
	return "No Scan Active"
}

//GetIpLocationsFromAPI is a function that utilizes the ipstack api for ip geolocation. It utilizes a simple curl command to get a json response body containing the desired information.
func GetIpLocationsFromAPI() {

	lines, err := ioutil.ReadFile("results.csv")
	if err != nil {
		fmt.Print("couldn't read results file")
		panic(0)
	}

	ipStrings := strings.Split(string(lines), "\n")

	updateInterval := 24 * time.Hour
	maxRetryInterval := time.Hour
	db, err := freegeoip.OpenURL(freegeoip.MaxMindDB, updateInterval, maxRetryInterval)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	select {
	case <-db.NotifyOpen():

	case err := <-db.NotifyError():
		log.Fatal(err)
	}

	queryResults := make([]freegeoip.DefaultQuery, len(ipStrings))

	for i, e := range ipStrings {

		fmt.Printf("\r%s  %d %s %d %s", "Located", i, "out of ", len(ipStrings), "IPS")

		var result freegeoip.DefaultQuery
		_ = db.Lookup(net.ParseIP(e), &result)

		queryResults[i] = result
	}
	scanning = false
	plotPoints(queryResults)
}

//QueryResult is effectively a touple, with one value being the ip which has been queried, and the result of said query.
type QueryResult struct {
	IP    string
	Query freegeoip.DefaultQuery
}

//QueryLocationsFromAPI utilizes the freegeoip package to produce an array of details for each IP.
func QueryIpLocationsFromAPI(ipStrings []string) []QueryResult {

	updateInterval := 24 * time.Hour
	maxRetryInterval := time.Hour
	db, err := freegeoip.OpenURL(freegeoip.MaxMindDB, updateInterval, maxRetryInterval)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	select {
	case <-db.NotifyOpen():

	case err := <-db.NotifyError():
		log.Fatal(err)
	}

	queryResults := make([]QueryResult, len(ipStrings))

	for i, e := range ipStrings {

		fmt.Printf("\r%s  %d %s %d %s", "Located", i, "out of ", len(ipStrings), "IPS")

		var result freegeoip.DefaultQuery
		_ = db.Lookup(net.ParseIP(e), &result)

		queryResults[i] = QueryResult{e, result}
	}

	return queryResults
}

//ResultJSON is a struct that contains the core attributes of the freegeoip.DefaultQuery struct which are needed for UI display.
//This struct is then marshalled and placed into the response body as JSON format.
type ResultJSON struct {
	IP        string  `json:"IP"`
	Continent string  `json:"Continent"`
	Country   string  `json:"Country"`
	Region    string  `json:"Region"`
	City      string  `json:"City"`
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
	TimeZone  string  `json:"TimeZone"`
}

//ConvertQueryResultToJSON converts the DefaultScan struct from the freegeoip package to a json string.
func ConvertQueryResultToJSON(q QueryResult) string {

	//All fields allow nil, however we have to index the Region.
	//So, to avoid index out of bound errors, we have a simple conditional.
	JSON := ResultJSON{}
	if len(q.Query.Region) == 0 {
		JSON = ResultJSON{
			IP:        q.IP,
			Continent: q.Query.Continent.Names["en"],
			Country:   q.Query.Country.Names["en"],
			Region:    "",
			City:      q.Query.City.Names["en"],
			Latitude:  q.Query.Location.Latitude,
			Longitude: q.Query.Location.Longitude,
			TimeZone:  q.Query.Location.TimeZone,
		}
	} else {
		JSON = ResultJSON{
			IP:        q.IP,
			Continent: q.Query.Continent.Names["en"],
			Country:   q.Query.Country.Names["en"],
			Region:    q.Query.Region[0].Names["en"],
			City:      q.Query.City.Names["en"],
			Latitude:  q.Query.Location.Latitude,
			Longitude: q.Query.Location.Longitude,
			TimeZone:  q.Query.Location.TimeZone,
		}
	}

	ret, err := json.Marshal(JSON)
	if err != nil {
		return "204"
	}

	return string(ret)

}
