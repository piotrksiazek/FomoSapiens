package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

var templates *template.Template

func main() {
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	templates = template.Must(template.ParseGlob("templates/*.html"))

	r.HandleFunc("/", indexHandler).Methods("GET")


	

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
	nids := getPostIds("Bitcoin")
	for index, nid := range nids {
		time.Sleep(2 * time.Second)
		fmt.Println("===================" + string(index) + "================") 
		getComments(nid, "Bitcoin")
	}
}