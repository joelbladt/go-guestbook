package main

import (
	"html/template"
	"log"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})

	log.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
