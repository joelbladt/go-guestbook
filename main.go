package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/joelbladt/go-guestbook/src/guestbook"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var name string = r.FormValue("name")
			var message string = r.FormValue("message")
			err := guestbook.Save(name, message)
			if err != nil {
				http.Error(w, "Error during saving", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		entries, err := guestbook.Show()
		if err != nil {
			http.Error(w, "Error during loading", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, entries)
	})

	log.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
