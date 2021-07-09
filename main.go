package main

import (
	"fmt"
	"net/http"
	"text/template"

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
	// c := make(chan string)
	// nids.getCommentsManyPosts(c, "Bitcoin")
	comments := nids.getCommentsManyPosts("Bitcoin")
	fmt.Println(getSentiment(comments))
	// for msg := range c {
	// 	fmt.Println(msg)
	// }
}