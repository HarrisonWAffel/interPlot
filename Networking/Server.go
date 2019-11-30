package Networking

import (
	"bufio"
	"html/template"
	"io/ioutil"
	"strings"

	"log"
	"net/http"
)

type Webpage struct {
	Title       string
	SpeedLimit  int
	NumberOfIps int
	Content     []byte
}

var ZmapOutput bufio.Scanner

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

//save page saves the content of index.html as a .txt file. This function will be deprecated soon.
func (p *Webpage) savePage() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Content, 0600)
}

//viewHandler serves index.html to our http requester
func viewHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadPage()
	if err != nil {
		log.Fatal("Could not load html")
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println()
		log.Panic("T is NULL")
	}
	t.Execute(w, p)
}

//scanHandler starts the zmap scan when the user is redirected to /scan. Can be started through the GUI or through the API header values.
func scanHandler(w http.ResponseWriter, r *http.Request) {

	if r.FormValue("SpeedLimit") != "" && r.FormValue("connNum") != "" {

		go scanInternet(r.FormValue("SpeedLimit"), r.FormValue("connNum"))

	} else if r.Header.Get("speedLimit") != "" && r.Header.Get("connNum") != "" {

		go scanInternet(r.Header.Get("speedLimit"), r.Header.Get("connNum"))
		w.Write([]byte("Scan Started."))

	} else {
		log.Println(r.FormValue("SpeedLimit"))
		log.Println(r.FormValue("connNum"))
		log.Println("The form values were empty")
		w.Write([]byte("Form values empty"))

	}
	viewHandler(w, r)
}

func getzmap(w http.ResponseWriter, r *http.Request) {
	o := ListenToScan()
	w.Write([]byte(o))
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	StopScan()
	viewHandler(w, r)
}

//startserver launches the handlers required for our server and its API.
func StartServer() {

	log.Println("Starting Server...")
	http.Handle("/", http.FileServer(http.Dir("./templates")))
	http.HandleFunc("/scan", scanHandler)
	http.HandleFunc("/stop", stopHandler)
	http.HandleFunc("/getzmap", getzmap)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
