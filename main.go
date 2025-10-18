package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"net/http"

	"github.com/zYasser/GoJest/internal/summary"
)

func main() {
	// Parse command line flags
	port := flag.String("port", "", "Port to run the server on")
	flag.Parse()

	funcMap := template.FuncMap{
		"toJson": func(v interface{}) string {
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return "null"
			}
			return string(jsonBytes)
		},
	}

	templates := template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))
	summaryHandler := &summary.Summary{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := templates.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/upload-test-summary", summaryHandler.UploadTestSummaryHandler(templates))
	http.HandleFunc("/summary", summaryHandler.GetSummary(templates))

	// Determine port: command line flag > environment variable > default
	serverPort := *port
	if serverPort == "" {
		serverPort = "8080"
	}

	if err := http.ListenAndServe(":"+serverPort, nil); err != nil {
		panic(err)
	}
}
