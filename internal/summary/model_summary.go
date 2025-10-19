package summary

type testSummary struct {
	NumFailedTestSuites       int          `json:"numFailedTestSuites"`
	NumFailedTests            int          `json:"numFailedTests"`
	NumPassedTestSuites       int          `json:"numPassedTestSuites"`
	NumPassedTests            int          `json:"numPassedTests"`
	NumPendingTestSuites      int          `json:"numPendingTestSuites"`
	NumPendingTests           int          `json:"numPendingTests"`
	NumRuntimeErrorTestSuites int          `json:"numRuntimeErrorTestSuites"`
	NumTodoTests              int          `json:"numTodoTests"`
	NumTotalTestSuites        int          `json:"numTotalTestSuites"`
	TestResults               []TestResult `json:"testResults"`
}

type TestResult struct {
	EndTime          int     `json:"endTime"`
	Message          string  `json:"message"`
	Name             string  `json:"name"`
	StartTime        int     `json:"startTime"`
	Status           string  `json:"status"`
	Summary          string  `json:"summary"`
	AssertionResults []Tests `json:"assertionResults"`
}

type Tests struct {
	FailureDetails    []interface{} `json:"failureDetails"`
	FailureMessages   []string      `json:"failureMessages"`
	FullName          string        `json:"fullName"`
	Invocations       int           `json:"invocations"`
	Location          string        `json:"location"`
	NumPassingAsserts int           `json:"numPassingAsserts"`
	RetryReasons      []interface{} `json:"retryReasons"`
	Status            string        `json:"status"`
	Title             string        `json:"title"`
}

type FailureDetail struct {
	MatcherResult interface{} `json:"matcherResult"`
}

type MatcherResult struct {
	Actual   interface{} `json:"actual"`
	Expected interface{} `json:"expected"`
	Message  string      `json:"message"`
	Name     string      `json:"name"`
	Pass     bool        `json:"pass"`
}
