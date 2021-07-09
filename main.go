package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	reddit "github.com/piotrksiazek/fomo-sapiens/reddit"
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
	nids := reddit.GetPostIds("Bitcoin")
	comments := nids.GetCommentsManyPosts("Bitcoin")
	fmt.Println(reddit.GetSentiment(comments))
}