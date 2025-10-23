# GoJest - Test Summary Viewer

A modern web application for visualizing and analyzing Jest test results. GoJest provides an intuitive interface to upload, view, and filter test summary JSON files with detailed failure analysis and interactive filtering capabilities.

## 🚀 Features

- **File Upload**: Drag-and-drop or click-to-upload Jest test summary JSON files
- **Interactive Dashboard**: Beautiful, responsive UI with real-time statistics
- **Advanced Filtering**: Filter tests by status (passed, failed, pending) and assertions
- **Detailed Failure Analysis**: Modal views with comprehensive failure details and stack traces

## 📋 Prerequisites

- Go 1.22.4 or later
- Modern web browser with JavaScript enabled

## 🛠️ Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/zYasser/GoJest.git
   cd GoJest
   ```

2. **Install dependencies**

   ```bash
   go mod tidy
   ```

3. **Build the application**

   ```bash
   go build -o gojest.exe .
   ```


4. **Access the application**
   Open your browser and navigate to `http://localhost:8080` (or your custom port)

## 📁 Project Structure

```
GoJest/
├── main.go                          # Main application entry point
├── go.mod                           # Go module dependencies
├── internal/
│   └── summary/
│       ├── model_summary.go         # Data models for test results
│       └── test_summary_handler.go  # HTTP handlers and business logic
├── templates/
│   ├── index.html                   # File upload page
│   └── test_summary.html           # Test results dashboard
├── test/
│   └── test.json                    # Sample Jest test results
└── README.md                        # This file
```

## 🎯 Usage

### Uploading Test Results

1. **Generate Jest Test Summary**
   Run your Jest tests with the `--json` flag and output to a file:

   ```bash
   npm test -- --json --outputFile=test-results.json
   ```

2. **Upload to GoJest**
   - Navigate to the home page
   - Drag and drop your JSON file or click to select
   - The application will automatically process and display results

### Viewing and Filtering Results

- **Overview Statistics**: See total passed, failed, pending tests, and test suites
- **Test Results List**: Browse individual test results with detailed information
- **Filter Options**:
  - Show only failed tests
  - Show only failed assertions
  - Show only passed tests
  - Show only pending tests
- **Assertion Details**: Expand individual test assertions with filtering capabilities
- **Failure Analysis**: Click on failed assertions to view detailed error information

## 🔧 API Endpoints

- `GET /` - File upload page
- `POST /upload-test-summary` - Upload and process test summary JSON
- `GET /summary` - View test results dashboard
- `GET /summary?onlyFailedTests=true` - Filtered results (supports multiple filters)

## 📊 Supported Test Frameworks

GoJest is designed to work with Jest test summary JSON output. The application expects JSON files with the following structure:

```json
{
  "numFailedTestSuites": 4,
  "numFailedTests": 2,
  "numPassedTestSuites": 5,
  "numPassedTests": 208,
  "numPendingTestSuites": 0,
  "numPendingTests": 0,
  "numRuntimeErrorTestSuites": 3,
  "numTodoTests": 0,
  "numTotalTestSuites": 9,
  "testResults": [...]
}
```

## 🎨 Technologies Used

- **Backend**: Go 1.22.4
- **Frontend**: HTML5, CSS3, JavaScript (ES6+)
- **Styling**: Tailwind CSS
- **Interactions**: HTMX
- **Templates**: Go HTML templates

## 🚀 Development

### Running in Development Mode

```bash
# Install Air for hot reloading (optional)
go install github.com/cosmtrek/air@latest

# Run with hot reloading
air
```

### Building for Production

```bash
# Build the application
go build -o gojest main.go

# Run the binary
./gojest
```

## 📝 Configuration

The application supports flexible port configuration:

### Command Line Options

- `gojest` - Run on default port 8080
- `gojest -port 3000` - Run on custom port
- `gojest -help` - Show help information
- `gojest -version` - Show version information

### Environment Variables

- `PORT=3000 gojest` - Use environment variable for port

### Priority Order

1. Command line flag (`-port`)
2. Environment variable (`PORT`)
3. Default port (8080)

---


