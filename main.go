package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/joelbladt/go-guestbook/src/guestbook"
)

//go:embed templates/index.html
var tmplFS embed.FS
var tmpl = template.Must(template.ParseFS(tmplFS, "templates/index.html"))

func main() {
	err := guestbook.InitFile()
	if err != nil {
		log.Fatalf("Failed to initialize guestbook file: %v", err)
	}

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
			log.Printf("Error loading entries: %v\n", err)
			http.Error(w, "Error during loading", http.StatusInternalServerError)
			return
		}

		// Handle template execution error
		if err := tmpl.Execute(w, entries); err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			log.Printf("Error rendering template: %v", err)
		}
	})

	// Check return value of ListenAndServe
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
		return
	}

	log.Println("Server running at http://localhost:8080")
}
