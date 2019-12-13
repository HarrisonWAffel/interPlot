package networking

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//Webpage is a struct that represents the content of the webpage.
type Webpage struct {
	Title       string
	SpeedLimit  int
	NumberOfIps int
	Content     []byte
}

//loadPage loads our index.html and returns its contents as a data structure.
func loadPage() (*Webpage, error) {
	filename := "templates/index.html"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	html := strings.Split(string(body), "</html>")[0]
	return &Webpage{Title: "InterPlot", SpeedLimit: 0, NumberOfIps: 0, Content: []byte(html)}, nil
}

//  Handlers ///

//viewHandler serves index.html to our http requester
func viewHandler(w http.ResponseWriter, r *http.Request) {

	p, err := loadPage()
	if err != nil {
		log.Fatal("Could not load html")
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println()
		log.Panic("cannot create index template")
	}

	r.Header.Add("Cache-Control", "no-cache, no-store")

	t.Execute(w, p)

}

var scanCtx, cancel = context.WithCancel(context.Background())

//scanHandler is a handler function that returns the index webpage after starting a zmap scan.
//It can be started either by visiting the website endpoint, or by sending the parmeters in its header
func scanHandler(w http.ResponseWriter, r *http.Request) {

	headerScanLimit := r.Header.Get("SpeedLimit")
	headerNum := r.Header.Get("connNum")

	speedLimit := r.FormValue("SpeedLimit")
	connNum := r.FormValue("connNum")

	if headerScanLimit != "" && headerNum != "" {
		//start the
		go ScanInternet(scanCtx, headerScanLimit, headerNum)

	} else if speedLimit != "" && connNum != "" {

		go ScanInternet(scanCtx, speedLimit, connNum)

	} else {

		//This error checking may not be needed, testing required.
		_, e := w.Write([]byte("Failed to parse scanLimit and connNum from request"))
		if e != nil {
			log.Println("Could send error ")
			log.Println(e)

		}
		return
	}
	//request that the client doesn't cache any of the images we send them.
	// this solves the issue of the browser always displaying a cached image
	// and not the image returned by the server.
	r.Header.Add("Cache-Control", "no-cache, no-store")

	viewHandler(w, r)
	//don't leak any child goroutines
	cancel()
}

/// API ///

func stopScan(w http.ResponseWriter, r *http.Request) {

	StopScan()
}

func scanOutput(w http.ResponseWriter, r *http.Request) {
	_, e := w.Write([]byte(ListenToScan()))
	if e != nil {
		log.Println("Could not listen to scan")
		log.Fatal(e)
	}
}

func listFoundIPS(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	o, e := ioutil.ReadFile("results.csv")
	if e != nil {
		w.Write([]byte("No IPS Located"))
	}
	ips := strings.Split(string(o), "\n")
	g := QueryIpLocationsFromAPI(ctx, ips)
	jsonResponse := "[\n"
	for i, e := range g {
		if i == len(g)-1 {
			jsonResponse += ConvertQueryResultToJSON(ctx, e) + "\n"
		} else {
			jsonResponse += ConvertQueryResultToJSON(ctx, e)
			jsonResponse += "," + "\n"
		}
	}

	jsonResponse += "]"

	fmt.Println(jsonResponse)
	w.Write([]byte(jsonResponse))

}

func getImg(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("output.png"))
}

//StartServer listens on localhost:8080 and monitors the handlers required for our server and its API.
func StartServer() {

	mux := http.NewServeMux()
	log.Println("Starting Server...")
	mux.Handle("/", http.FileServer(http.Dir("./templates")))
	mux.HandleFunc("/scan", scanHandler)
	mux.HandleFunc("/stopscan", stopScan)
	mux.HandleFunc("/scanoutput", scanOutput)
	mux.HandleFunc("/listfoundips", listFoundIPS)
	mux.HandleFunc("/img", getImg)
	log.Fatal(http.ListenAndServe(":8080", mux))

}
