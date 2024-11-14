package main

import (
	"goload/internal/handler/http"
	"log"
)

var (
	version    string
	commitHash string
)

func main() {
	log.Printf("Goload version: %s, commit hash: %s\n", version, commitHash)

	server := http.NewServer()
	if err := server.Start(); err != nil {
		log.Fatalln("failed to run server", err)
	}
}
