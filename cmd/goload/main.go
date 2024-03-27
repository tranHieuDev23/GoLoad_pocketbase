package main

import (
	"log"

	"github.com/tranHieuDev23/GoLoad_pocketbase/internal/handlers/http"
)

func main() {
	server := http.NewServer()
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
