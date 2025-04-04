# cliproject
A Golang command-line application that reads a CSV file containing URLs, downloads their content concurrently, and saves each response as a `.txt` file with a randomly generated filename.

---

## 📦 Features

- Accepts a CSV file containing URLs (1 per line).
- Efficient, memory-safe file reading (streaming, not full-load).
- Concurrent URL downloading (up to 50 goroutines).
- Single-threaded disk writer for safe file output.
- CLI interface using `go-prompt`.
- Built-in test coverage with mocked HTTP server for offline testing.
- Logging support for debug and tracking.

  ## 🛠 Requirements

- Go 1.23 or newer

---

## 📄 CSV Format

The CSV should be structured like this:

```csv
Urls
http://example.com
http://httpbin.org/get
```

- Clone and Run the CLI:
  git clone https://github.com/lakshmi-83659/cliproject.git
  cd cliproject
  go mod tidy
  go run cmd/main.go

- You will see a prompt:
  >>>
- Run with your CSV file
  >>> run samples.csv

- The app will:
Read the file
Download contents
Save responses in the output/ folder as .txt files

- Exit the CLI
At the prompt, type:
>>> exit

- To run Testcases
  cd cliproject/cmd
  go test
  Testcases will be excuted
