package Networking

import (
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

func loadPage() (*Webpage, error) {
	filename := "templates/index.html"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	html := strings.Split(string(body), "</html>")[0]
	return &Webpage{Title: "InterPlot", SpeedLimit: 0, NumberOfIps: 0, Content: []byte(html)}, nil
}

func (p *Webpage) savePage() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Content, 0600)
}

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

	err = t.Execute(w, p)
	if err != nil {
		log.Println()
		log.Fatal("Cannot build template")
	}
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("speedLimit") != "" && r.FormValue("connNum") != "" {
		scanInternet(r.FormValue("speedLimit"), r.FormValue("connNum"))
	} else {
		log.Println(r.FormValue("speedLimit"))
		log.Println(r.FormValue("connNum"))
		viewHandler(w, r)
		log.Println("The form values were empty")
	}
	viewHandler(w, r)

}

func StartServer() {
	log.Println("Starting Server...")
	http.Handle("/", http.FileServer(http.Dir("./templates")))
	http.HandleFunc("/scan", scanHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
