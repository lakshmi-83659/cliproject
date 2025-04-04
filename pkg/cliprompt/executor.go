package cliprompt

import (
	"cliproject/pkg/handlers"
	"cliproject/pkg/models"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func Executor(input string) {
	log.Println("Executing the command...")
	args := strings.Fields(input)
	if len(args) < 1 {
		fmt.Println("Usage: run <csv_file_path> \n exit")
		return
	}
	switch args[0] {
	case "exit":
		os.Exit(0)
		return
	case "run":
		csvDownloader(args)
	}
	fmt.Println("Processing completed.")
}

func csvDownloader(args []string) {
	csvPath := args[1]
	ctx := context.Background()
	jobs := make(chan models.Job)
	results := make(chan models.Result, 50)
	done := make(chan struct{})
	sem := make(chan struct{}, 50)

	go func() {
		if err := handlers.ReadCSVLines(ctx, csvPath, jobs); err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}
	}()

	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go handlers.DownloadWorker(jobs, results, &wg, sem)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go handlers.PersistWorker(results, done)
	<-done
}
