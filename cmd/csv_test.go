package main

import (
	"cliproject/pkg/handlers"
	"cliproject/pkg/models"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestReadCSVLines(t *testing.T) {
	log.Println("Running TestReadCSVLines")
	csvContent := "Urls\nexample.com\nexample.org"
	tempFile, err := os.CreateTemp("", "test-*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, _ = tempFile.WriteString(csvContent)
	_ = tempFile.Close()

	jobs := make(chan models.Job)
	ctx := context.Background()

	go func() {
		err = handlers.ReadCSVLines(ctx, tempFile.Name(), jobs)
		if err != nil {
			t.Errorf("readCSVLines failed: %v", err)
		}
	}()

	expected := []string{"example.com", "example.org"}
	for i, url := range expected {
		job, ok := <-jobs
		if !ok {
			t.Fatalf("expected job but got closed channel at index %d", i)
		}
		if job.URL != url {
			t.Errorf("expected URL %s, got %s", url, job.URL)
		}
	}
}

func TestDownloadWorkerWithMock(t *testing.T) {
	log.Println("Running TestDownloadWorkerWithMock")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("mocked response"))
	}))
	defer server.Close()

	mockURL := strings.TrimPrefix(server.URL, "http://")
	jobs := make(chan models.Job, 1)
	results := make(chan models.Result, 1)
	sem := make(chan struct{}, 1)
	var wg sync.WaitGroup

	jobs <- models.Job{URL: mockURL}
	close(jobs)

	wg.Add(1)
	go handlers.DownloadWorker(jobs, results, &wg, sem)

	wg.Wait()
	close(results)

	res, ok := <-results
	if !ok || string(res.Content) != "mocked response" {
		t.Errorf("downloadWorker did not return expected content, got: %s", string(res.Content))
	}
}

func TestPersistWorker(t *testing.T) {
	log.Println("Running TestPersistWorker")
	results := make(chan models.Result, 1)
	done := make(chan struct{})

	_ = os.Mkdir("output", os.ModePerm)
	defer os.RemoveAll("output")

	results <- models.Result{Content: []byte("test content")}
	close(results)

	go handlers.PersistWorker(results, done)
	<-done

	files, err := ioutil.ReadDir("output")
	if err != nil {
		t.Fatalf("failed to read output directory: %v", err)
	}

	if len(files) != 1 || !strings.HasSuffix(files[0].Name(), ".txt") {
		t.Errorf("expected one .txt file in output, found: %v", files)
	}
}

func TestReadCSVLines_FileNotFound(t *testing.T) {
	log.Println("Running TestReadCSVLines_FileNotFound")
	jobs := make(chan models.Job)
	err := handlers.ReadCSVLines(context.Background(), "nonexistent.csv", jobs)
	if err == nil {
		t.Errorf("expected error for missing file, got nil")
	}
}

func TestDownloadWorker_InvalidURL(t *testing.T) {
	log.Println("Running TestDownloadWorker_InvalidURL")
	jobs := make(chan models.Job, 1)
	results := make(chan models.Result, 1)
	sem := make(chan struct{}, 1)
	var wg sync.WaitGroup

	jobs <- models.Job{URL: "http://invalid_url"} // badly formatted or unreachable
	close(jobs)

	wg.Add(1)
	go handlers.DownloadWorker(jobs, results, &wg, sem)

	wg.Wait()
	close(results)

	if len(results) != 0 {
		for r := range results {
			t.Errorf("expected no result, but got: %+v", r)
		}
	}
}
