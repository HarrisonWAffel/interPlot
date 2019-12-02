package Networking

import (
	"context"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

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
	ips := string(o)

	w.Write([]byte(ips))
}

func getImg(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("./templates")
	if err != nil {
		log.Fatal(err)
	}
	x := 0
	for _, file := range files {
		if strings.Contains(file.Name(), "output") {
			x++
		}
		log.Println(file.Name())
	}
	var output_files = make([]string, x)
	k := 0
	for _, file := range files {
		if strings.Contains(file.Name(), "output") {
			output_files[k] = file.Name()
			k++
		}
	}

	sort.Strings(output_files)
	log.Println(output_files[len(output_files)-1])
	w.Write([]byte(output_files[len(output_files)-1]))
}

//startserver launches the handlers required for our server and its API.
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
