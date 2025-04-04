package handlers

import (
	"bufio"
	"cliproject/pkg/models"
	"context"
	"crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func ReadCSVLines(ctx context.Context, path string, jobs chan<- models.Job) error {
	log.Println("Reading CSV file")
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("CSV File not present")
		return err
	}
	defer file.Close()

	r := csv.NewReader(bufio.NewReader(file))
	head, err := r.Read()
	if err != nil || len(head) != 1 || strings.ToLower(head[0]) != "urls" {
		log.Fatal("Invalid csv file")
		return errors.New("invalid CSV format: missing or incorrect header")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			record, err := r.Read()
			//fmt.Println("Full file conent:", record)
			if err == io.EOF {
				close(jobs)
				return nil
			}
			if err != nil {
				return err
			}
			if len(record) != 1 {
				continue
			}
			//fmt.Println("File content:", record[0])
			jobs <- models.Job{URL: record[0]}
		}
	}
}

func DownloadWorker(jobs <-chan models.Job, results chan<- models.Result, wg *sync.WaitGroup, sem chan struct{}) {
	log.Println("Downloading csv file content")
	defer wg.Done()
	for job := range jobs {
		sem <- struct{}{}
		resp, err := http.Get(job.URL)
		if err == nil {
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil {
				results <- models.Result{Content: body}
			}
		}
		<-sem
	}
}

func PersistWorker(results <-chan models.Result, done chan<- struct{}) {
	log.Println("Storing the data...")
	for result := range results {
		fileName, err := generateRandomFileName()
		if err != nil {
			continue
		}
		_ = os.WriteFile(filepath.Join("output", fileName), result.Content, 0644)
	}
	done <- struct{}{}
}

func generateRandomFileName() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.txt", hex.EncodeToString(b)), nil
}
