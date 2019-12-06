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

var ctx context.Context

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

//scanHandler is a handler function that returns the index webpage after starting a zmap scan.
//It can be started either by visiting the website endpoint, or by sending the parmeters in its header
func scanHandler(w http.ResponseWriter, r *http.Request) {

	headerScanLimit := w.Header().Get("SpeedLimit")
	headerNum := w.Header().Get("connNum")

	speedLimit := r.FormValue("SpeedLimit")
	connNum := r.FormValue("connNum")

	if headerScanLimit != "" && headerNum != "" {
		ScanInternet(nil, headerScanLimit, headerNum)
	} else if speedLimit != "" && connNum != "" {
		ScanInternet(nil, speedLimit, connNum)
	}
	r.Header.Add("Cache-Control", "no-cache, no-store")
	viewHandler(w, r)
}

/// API ///

func scanOutput(w http.ResponseWriter, r *http.Request) {
	_, e := w.Write([]byte(ListenToScan()))
	if e != nil {
		log.Println("Could not listen to scan")
		log.Fatal(e)
	}
}

func listFoundIPS(w http.ResponseWriter, r *http.Request) {

	o, e := ioutil.ReadFile("results.csv")
	if e != nil {
		w.Write([]byte("No IPS Located"))
	}
	ips := strings.Split(string(o), "\n")
	g := QueryIpLocationsFromAPI(ips)
	jsonResponse := "[\n"
	for i, e := range g {
		if i == len(g)-1 {
			jsonResponse += ConvertQueryResultToJSON(e) + "\n"
		} else {
			jsonResponse += ConvertQueryResultToJSON(e)
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

//StartServer launches the handlers required for our server and its API.
func StartServer() {
	ctx = context.Background()
	log.Println("Starting Server...")
	http.Handle("/", http.FileServer(http.Dir("./templates")))
	http.HandleFunc("/scan", scanHandler)
	http.HandleFunc("/scanoutput", scanOutput)
	http.HandleFunc("/listfoundips", listFoundIPS)
	http.HandleFunc("/img", getImg)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
