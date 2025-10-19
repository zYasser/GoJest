package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/zYasser/GoJest/internal/summary"
)

func main() {
	port := flag.String("port", "", "Port to run the server on")
	help := flag.Bool("help", false, "Show help information")
	flag.Parse()

	if *help {
		fmt.Println("GoJest - Test Summary Server")
		fmt.Println("Usage: gojest [options]")
		fmt.Println("")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

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
	http.HandleFunc("/upload-json-text", summaryHandler.UploadJsonTextHandler(templates))
	http.HandleFunc("/summary", summaryHandler.GetSummary(templates))

	serverPort := *port
	if serverPort == "" {
		serverPort = os.Getenv("PORT")
	}
	if serverPort == "" {
		serverPort = "8080"
	}

	fmt.Printf("🚀 GoJest server starting on port %s\n", serverPort)
	fmt.Printf("📊 Test Summary Dashboard: http://localhost:%s\n", serverPort)
	fmt.Println("Press Ctrl+C to stop the server")

	if err := http.ListenAndServe(":"+serverPort, nil); err != nil {
		panic(err)
	}
}
