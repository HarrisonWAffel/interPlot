package Networking

import (
	"bufio"
	"fmt"
	"github.com/fiorix/freegeoip"

	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const ApiKey = "959279e06d5ca3430cf26fdbca17ea7c"

var scanner *bufio.Scanner
var cmd *exec.Cmd
var stdin io.WriteCloser

//scanInternet runs the zmap command and outputs a csv file
func scanInternet(speedLimit string, n string) {

	if scanner == nil {
		fmt.Println("Scanning")
		cm := "zmap"

		args := []string{"-B", strings.TrimSpace(speedLimit) + "M", "-p", "21", "-n", strings.TrimSpace(n), "-o", "results.csv"}
		fmt.Println(args)
		cmd = exec.Command(cm, args...)

		//We need to create a reader for the stderr of this script
		cmdReader, err := cmd.StderrPipe()
		if err != nil {
			fmt.Println("Error creating StdoutPipe for Cmd", err)
			os.Exit(1)
		}

		defer cmdReader.Close()

		stdn, e := cmd.StdinPipe()
		if e != nil {
			log.Fatal("Cannot get stdin pipe ")
		}
		stdin = stdn
		//A scanner is created to read the stdout of the above command
		scanner = bufio.NewScanner(cmdReader)

		//A new go thread is created to handle the output
		//fmt not log
		go func() {
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
			return
		}()

		log.Println("Started scan")

		//We need to start our goroutine from the main thread
		err = cmd.Run()
		if err != nil {
			log.Fatal("Error starting Cmd ", err)
		}
		scanner = nil
		GetIpLocationsFromAPI()
	}
}

func StopScan() {

	e := cmd.Process.Signal(syscall.SIGINT)
	if e != nil {
		log.Fatal("Cannot sigint ")
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

//getIpLocationsFromAPI is a function that utilizes the ipstack api for ip geolocation. It utilizes a simple curl command to get a json response body containing the desired information.
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

	plotPoints(queryResults)
}
