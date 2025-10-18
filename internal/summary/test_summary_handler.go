package summary

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

type Summary struct {
	TestSummary *testSummary
}
type SummaryQuery struct {
	OnlyFailedFiles      bool
	OnlyFailedTests      bool
	OnlyFailedAssertions bool
	OnlyPassedTests      bool
	OnlyPendingTests     bool
}

func (s *Summary) UploadTestSummaryHandler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Upload request received: Method=%s, Content-Type=%s\n", r.Method, r.Header.Get("Content-Type"))

		r.ParseMultipartForm(100 << 20)
		fmt.Printf("Form values: %v\n", r.Form)

		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Error getting form file: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fmt.Printf("File received: %s, Size: %d\n", header.Filename, header.Size)
		jsonData, err := io.ReadAll(file)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(jsonData, &s.TestSummary)
		if err != nil {
			fmt.Printf("Error unmarshalling JSON: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tempFile, err := os.Create("tmp.json")
		if err != nil {
			fmt.Printf("Error creating temp file: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()
		_, err = tempFile.Write(jsonData)
		if err != nil {
			fmt.Printf("Error writing to temp file: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", "/summary")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("File uploaded successfully"))

	}
}

func (s *Summary) GetSummary(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isHTMXRequest := r.Header.Get("HX-Request") == "true"

		if s.TestSummary == nil {
			file, err := os.Open("tmp.json")
			if err != nil {
				fmt.Printf("Error opening temp file: %v\n", err)
				if isHTMXRequest {
					w.Header().Set("HX-Redirect", "/")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("Redirecting to home page"))
				} else {
					http.Redirect(w, r, "/", http.StatusSeeOther)
				}
				return
			}
			defer file.Close()
			err = json.NewDecoder(file).Decode(&s.TestSummary)
			if err != nil {
				fmt.Printf("Error decoding temp file: %v\n", err)
				if isHTMXRequest {
					w.Header().Set("HX-Redirect", "/")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("Redirecting to home page"))
				} else {
					http.Redirect(w, r, "/", http.StatusSeeOther)
				}
				return
			}
		}
		if s.TestSummary == nil {
			if isHTMXRequest {
				w.Header().Set("HX-Redirect", "/")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Redirecting to home page"))
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			return
		}
		query := r.URL.Query()
		summaryQuery := &SummaryQuery{
			OnlyFailedFiles:      query.Get("onlyFailedFiles") == "true",
			OnlyFailedTests:      query.Get("onlyFailedTests") == "true",
			OnlyFailedAssertions: query.Get("onlyFailedAssertions") == "true",
			OnlyPassedTests:      query.Get("onlyPassedTests") == "true",
			OnlyPendingTests:     query.Get("onlyPendingTests") == "true",
		}
		if len(query) > 0 {
			temp := maniplulateTestSummary(*summaryQuery, *s.TestSummary)
			templates.ExecuteTemplate(w, "test_summary.html", temp)

		} else {
			templates.ExecuteTemplate(w, "test_summary.html", s.TestSummary)

		}
	}
}

func maniplulateTestSummary(summaryQuery SummaryQuery, s testSummary) testSummary {
	filteredSummary := s
	filteredTests := make([]TestResult, 0)

	for _, test := range s.TestResults {
		shouldInclude := true

		if summaryQuery.OnlyFailedTests && test.Status != "failed" {
			shouldInclude = false
		}
		if summaryQuery.OnlyPassedTests && test.Status != "passed" {
			shouldInclude = false
		}
		if summaryQuery.OnlyPendingTests && test.Status != "pending" {
			shouldInclude = false
		}

		if summaryQuery.OnlyFailedAssertions {
			hasFailedAssertions := false
			for _, assertion := range test.AssertionResults {
				if assertion.Status == "failed" {
					hasFailedAssertions = true
					break
				}
			}
			if test.Status == "failed" {
				shouldInclude = true
			}
			if !hasFailedAssertions {
				shouldInclude = false
			}
		}

		if shouldInclude {
			filteredTest := test

			if summaryQuery.OnlyFailedAssertions {
				filteredAssertions := make([]Tests, 0)
				for _, assertion := range test.AssertionResults {
					if assertion.Status == "failed" {
						filteredAssertions = append(filteredAssertions, assertion)
					}
				}
				filteredTest.AssertionResults = filteredAssertions
			}

			filteredTests = append(filteredTests, filteredTest)
		}
	}

	filteredSummary.TestResults = filteredTests

	filteredSummary.NumPassedTests = 0
	filteredSummary.NumFailedTests = 0
	filteredSummary.NumPendingTests = 0
	filteredSummary.NumPassedTestSuites = 0
	filteredSummary.NumPendingTestSuites = 0
	filteredSummary.NumTotalTestSuites = len(filteredTests)

	for _, test := range filteredTests {
		switch test.Status {
		case "passed":
			filteredSummary.NumPassedTestSuites++

		case "failed":
			filteredSummary.NumFailedTestSuites = 0
		case "pending":
			filteredSummary.NumPendingTests++
		}
	}
	filteredSummary.NumFailedTests = 0
	filteredSummary.NumPassedTests = 0
	for _, test := range filteredTests {
		for _, assert := range test.AssertionResults {
			if assert.Status == "failed" {
				filteredSummary.NumFailedTests++
			}
			if assert.Status == "passed" {
				filteredSummary.NumPassedTests++
			}
			if assert.Status == "pending" {
				filteredSummary.NumPendingTests++
			}
		}
	}

	return filteredSummary
}
