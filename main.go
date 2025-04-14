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
			name := r.FormValue("name")
			message := r.FormValue("message")

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

		// Handle template execution error
		if err := tmpl.Execute(w, entries); err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			log.Printf("template execution failed: %v", err)
		}
	})

	// Check return value of ListenAndServe
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
		return
	}

	log.Println("Server running at http://localhost:8080")
}
