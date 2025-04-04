package main

import (
	"cliproject/pkg/cliprompt"
	"log"
	"os"

	prompt "github.com/c-bata/go-prompt"
)

func main() {
	log.Println("Application started")
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		_ = os.Mkdir("output", os.ModePerm)
	}
	p := prompt.New(
		cliprompt.Executor,
		cliprompt.Completer,
	)
	p.Run()
}
